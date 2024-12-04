package maxcompute

import (
	"context"
	_ "embed" // used to print the embedded assets
	"fmt"
	"time"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/common"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/datatype"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tableschema"
	"github.com/goto/meteor/models"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/plugins/extractors/maxcompute/client"
	"github.com/goto/meteor/plugins/extractors/maxcompute/config"
	"github.com/goto/meteor/plugins/internal/upstream"
	"github.com/goto/meteor/registry"
	"github.com/goto/meteor/utils"
	"github.com/goto/salt/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	maxcomputeService = "maxcompute"
	typeTable         = "table"
)

type Extractor struct {
	plugins.BaseExtractor
	logger log.Logger
	config config.Config
	randFn randFn

	client    Client
	newClient NewClientFunc
	eg        *errgroup.Group
}

type randFn func(rndSeed int64) func(int64) int64

type NewClientFunc func(ctx context.Context, logger log.Logger, conf config.Config) (Client, error)

//go:embed README.md
var summary string

var sampleConfig = `
project_name: goto_test
endpoint_project: https://service.ap-southeast-5.maxcompute.aliyun.com/api
access_key:
    id: access_key_id
    secret: access_key_secret
schema_name: default
exclude:
	schemas:
    	- schema_a
    	- schema_b
    tables:
    	- schema_c.table_a
concurrency: 10
max_preview_rows: 3
mix_values: false
build_view_lineage: true,
`

var info = plugins.Info{
	Description:  "MaxCompute metadata and metrics",
	SampleConfig: sampleConfig,
	Tags:         []string{"alicloud", "table"},
	Summary:      summary,
}

//go:generate mockery --name=Client -r --case underscore --with-expecter --structname MaxComputeClient --filename maxcompute_client_mock.go --output=./mocks
type Client interface {
	ListSchema(ctx context.Context) ([]*odps.Schema, error)
	ListTable(ctx context.Context, schemaName string) ([]*odps.Table, error)
	GetTableSchema(ctx context.Context, table *odps.Table) (string, *tableschema.TableSchema, error)
	GetTablePreview(ctx context.Context, partitionValue string, table *odps.Table, maxRows int) ([]string, *structpb.ListValue, error)
}

func New(logger log.Logger, clientFunc NewClientFunc) *Extractor {
	e := &Extractor{
		logger:    logger,
		newClient: clientFunc,
	}
	e.BaseExtractor = plugins.NewBaseExtractor(info, &e.config)
	e.ScopeNotRequired = true

	return e
}

func (e *Extractor) Init(ctx context.Context, conf plugins.Config) error {
	if err := e.BaseExtractor.Init(ctx, conf); err != nil {
		return err
	}

	if e.config.ProjectName == "" {
		return fmt.Errorf("project_name is required")
	}
	if e.config.AccessKey.ID == "" || e.config.AccessKey.Secret == "" {
		return fmt.Errorf("access_key is required")
	}
	if e.config.EndpointProject == "" {
		return fmt.Errorf("endpoint_project is required")
	}
	if e.config.Concurrency == 0 {
		e.config.Concurrency = 1
	}

	var err error
	e.client, err = e.newClient(ctx, e.logger, e.config)
	if err != nil {
		return err
	}

	e.eg = &errgroup.Group{}
	e.eg.SetLimit(e.config.Concurrency)

	return nil
}

func (e *Extractor) Extract(ctx context.Context, emit plugins.Emit) error {
	schemas, err := e.client.ListSchema(ctx)
	if err != nil && len(schemas) == 0 {
		return err
	}

	for _, schema := range schemas {
		if e.config.SchemaName != "" && schema.Name() != e.config.SchemaName {
			continue
		}
		if contains(e.config.Exclude.Schemas, schema.Name()) {
			continue
		}

		err := e.fetchTablesFromSchema(ctx, schema, emit)
		if err != nil {
			return err
		}
	}

	return e.eg.Wait()
}

func (e *Extractor) fetchTablesFromSchema(ctx context.Context, schema *odps.Schema, emit plugins.Emit) error {
	tables, err := e.client.ListTable(ctx, schema.Name())
	if err != nil && len(tables) == 0 {
		return err
	}

	for _, table := range tables {
		if contains(e.config.Exclude.Tables, fmt.Sprintf("%s.%s", table.SchemaName(), table.Name())) {
			continue
		}

		tbl := table
		e.eg.Go(func() error {
			return e.processTable(ctx, schema, tbl, emit)
		})
	}

	return nil
}

func (e *Extractor) processTable(ctx context.Context, schema *odps.Schema, table *odps.Table, emit plugins.Emit) error {
	tableType, tableSchema, err := e.client.GetTableSchema(ctx, table)
	if err != nil {
		return err
	}

	fmt.Println("********************" + tableType + "********************")
	asset, err := e.buildAsset(ctx, schema, table, tableType, tableSchema)
	if err != nil {
		e.logger.Error("failed to build asset", "table", table.Name(), "error", err)
		return err
	}

	emit(models.NewRecord(asset))
	return nil
}

func (e *Extractor) buildAsset(ctx context.Context, schema *odps.Schema,
	table *odps.Table, tableType string, tableSchema *tableschema.TableSchema,
) (*v1beta2.Asset, error) {
	defaultSchema := "default"
	schemaName := schema.Name()
	if schemaName == "" {
		schemaName = defaultSchema
	}

	tableURN := plugins.MaxComputeURN(e.config.ProjectName, schemaName, tableSchema.TableName)

	var previewFields []string
	var previewRows *structpb.ListValue
	var err error
	previewFields, previewRows, err = e.buildPreview(ctx, table, tableSchema)
	if err != nil {
		e.logger.Warn("error building preview", "err", err, "table", tableSchema.TableName)
	}

	asset := &v1beta2.Asset{
		Urn:         tableURN,
		Name:        tableSchema.TableName,
		Type:        typeTable,
		Description: tableSchema.Comment,
		CreateTime:  timestamppb.New(time.Time(tableSchema.CreateTime)),
		UpdateTime:  timestamppb.New(time.Time(tableSchema.LastModifiedTime)),
		Service:     maxcomputeService,
	}

	tableAttributesData := e.buildTableAttributesData(schemaName, tableType, tableSchema)

	if tableType == config.TableTypeView {
		query := tableSchema.ViewText
		tableAttributesData["sql"] = query
		if e.config.BuildViewLineage {
			upstreamResources := getUpstreamResources(query)
			asset.Lineage = &v1beta2.Lineage{
				Upstreams: upstreamResources,
			}
		}
	}

	var columns []*v1beta2.Column
	for i, col := range tableSchema.Columns {
		columnData := &v1beta2.Column{
			Name:        col.Name,
			DataType:    dataTypeToString(col.Type),
			Description: col.Comment,
			IsNullable:  col.IsNullable,
			Attributes:  utils.TryParseMapToProto(buildColumnAttributesData(&tableSchema.Columns[i])),
			Columns:     buildColumns(col.Type),
		}
		columns = append(columns, columnData)
	}

	tableData := &v1beta2.Table{
		Attributes: utils.TryParseMapToProto(tableAttributesData),
		Columns:    columns,
		CreateTime: timestamppb.New(time.Time(tableSchema.CreateTime)),
		UpdateTime: timestamppb.New(time.Time(tableSchema.LastModifiedTime)),
	}

	maxPreviewRows := e.config.MaxPreviewRows
	if maxPreviewRows > 0 {
		tableData.PreviewFields = previewFields
		tableData.PreviewRows = previewRows
	}

	tbl, err := anypb.New(tableData)
	if err != nil {
		e.logger.Warn("error creating Any struct", "error", err)
	}
	asset.Data = tbl

	return asset, nil
}

func getUpstreamResources(query string) []*v1beta2.Resource {
	upstreamDependencies := upstream.ParseTopLevelUpstreamsFromQuery(query)
	uniqueUpstreamDependencies := upstream.UniqueFilterResources(upstreamDependencies)
	var upstreams []*v1beta2.Resource
	for _, dependency := range uniqueUpstreamDependencies {
		urn := plugins.MaxComputeURN(dependency.Project, dependency.Dataset, dependency.Name)
		upstreams = append(upstreams, &v1beta2.Resource{
			Urn:     urn,
			Name:    dependency.Name,
			Type:    "table",
			Service: maxcomputeService,
		})
	}
	return upstreams
}

func buildColumns(dataType datatype.DataType) []*v1beta2.Column {
	if dataType.ID() != datatype.STRUCT {
		return nil
	}

	structType, ok := dataType.(datatype.StructType)
	if !ok {
		return nil
	}

	var columns []*v1beta2.Column
	for _, field := range structType.Fields {
		column := &v1beta2.Column{
			Name:     field.Name,
			DataType: dataTypeToString(field.Type),
			Columns:  buildColumns(field.Type),
		}
		columns = append(columns, column)
	}
	return columns
}

func (e *Extractor) buildTableAttributesData(schemaName, tableType string, tableInfo *tableschema.TableSchema) map[string]interface{} {
	attributesData := map[string]interface{}{}

	attributesData["project_name"] = e.config.ProjectName
	attributesData["schema"] = schemaName
	attributesData["type"] = tableType

	rb := common.ResourceBuilder{ProjectName: e.config.ProjectName}
	attributesData["resource_url"] = rb.Table(tableInfo.TableName)

	if tableInfo.ViewText != "" {
		attributesData["sql"] = tableInfo.ViewText
	}

	var partitionNames []interface{}
	if tableInfo.PartitionColumns != nil && len(tableInfo.PartitionColumns) > 0 {
		partitionNames = make([]interface{}, len(tableInfo.PartitionColumns))
		for i, column := range tableInfo.PartitionColumns {
			partitionNames[i] = column.Name
		}
		attributesData["partition_fields"] = partitionNames
	}

	return attributesData
}

func buildColumnAttributesData(column *tableschema.Column) map[string]interface{} {
	attributesData := map[string]interface{}{}

	if column == nil {
		return attributesData
	}

	if column.Label != "" {
		attributesData["label"] = column.Label
	}

	return attributesData
}

func (e *Extractor) buildPreview(ctx context.Context, t *odps.Table, tSchema *tableschema.TableSchema) ([]string, *structpb.ListValue, error) {
	maxPreviewRows := e.config.MaxPreviewRows
	if maxPreviewRows <= 0 {
		return nil, nil, nil
	}

	previewFields, previewRows, err := e.client.GetTablePreview(ctx, "", t, maxPreviewRows)
	if err != nil {
		e.logger.Error("failed to preview table", "table", t.Name(), "error", err)
		return nil, nil, err
	}

	if e.config.MixValues {
		tempRows := make([]interface{}, len(previewRows.GetValues()))
		for i, val := range previewRows.GetValues() {
			tempRows[i] = val.AsInterface()
		}

		tempRows, err = e.mixValuesIfNeeded(tempRows, time.Time(tSchema.LastModifiedTime).Unix())
		if err != nil {
			return nil, nil, fmt.Errorf("mix values: %w", err)
		}

		previewRows, err = structpb.NewList(tempRows)
		if err != nil {
			return nil, nil, fmt.Errorf("create preview list: %w", err)
		}
	}

	return previewFields, previewRows, nil
}

func (e *Extractor) mixValuesIfNeeded(rows []interface{}, rndSeed int64) ([]interface{}, error) {
	if !e.config.MixValues || len(rows) < 2 {
		return rows, nil
	}

	var table [][]any
	for _, row := range rows {
		arr, ok := row.([]any)
		if !ok {
			return nil, fmt.Errorf("row %d is not a slice", row)
		}
		table = append(table, arr)
	}

	numRows := len(table)
	numColumns := len(table[0])

	if e.randFn == nil {
		return nil, fmt.Errorf("randFn is not initialized")
	}

	rndGen := e.randFn(rndSeed)
	for col := 0; col < numColumns; col++ {
		for row := 0; row < numRows; row++ {
			randomRow := rndGen(int64(numRows))

			table[row][col], table[randomRow][col] = table[randomRow][col], table[row][col]
		}
	}

	mixedRows := make([]any, numRows)
	for i, row := range table {
		mixedRows[i] = row
	}
	return mixedRows, nil
}

func dataTypeToString(dataType datatype.DataType) string {
	if dataType.ID() == datatype.MAP {
		return dataType.Name()
	}
	if dataType.ID() == datatype.ARRAY {
		return dataType.Name()
	}
	return dataType.ID().String()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	if err := registry.Extractors.Register(maxcomputeService, func() plugins.Extractor {
		return New(plugins.GetLog(), CreateClient)
	}); err != nil {
		panic(err)
	}
}

func CreateClient(_ context.Context, _ log.Logger, conf config.Config) (Client, error) {
	return client.New(conf)
}
