package maxcompute

import (
	"context"
	_ "embed" // used to print the embedded assets
	"fmt"
	"math/rand"
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
	maxcomputeService             = "maxcompute"
	attributesDataSQL             = "sql"
	attributesDataMaskingPolicy   = "masking_policy"
	attributesDataProjectName     = "project_name"
	attributesDataSchema          = "schema"
	attributesDataType            = "type"
	attributesDataResourceURL     = "resource_url"
	attributesDataPartitionFields = "partition_fields"
	attributesDataLabel           = "label"
	attributesDataLifecycle       = "lifecycle"
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
	GetMaskingPolicies(table *odps.Table) (map[client.Column][]client.Policy, error)
}

func New(logger log.Logger, clientFunc NewClientFunc, randFn randFn) *Extractor {
	e := &Extractor{
		logger:    logger,
		newClient: clientFunc,
		randFn:    randFn,
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
			e.logger.Info("skipping schema as it is in the exclude list", "schema", schema.Name())
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
		tableName := fmt.Sprintf("%s.%s", schema.Name(), table.Name())
		if contains(e.config.Exclude.Tables, tableName) {
			e.logger.Info("skipping table as it is in the exclude list", "table", tableName)
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

	// If lifecycle is less than the minimum lifecycle (days), skip the table
	if e.config.Exclude.MinTableLifecycle > 1 {
		lifecyclePermanent := tableSchema.Lifecycle == -1
		lifecycleNotConfigured := tableSchema.Lifecycle == 0
		if !lifecyclePermanent && !lifecycleNotConfigured && tableSchema.Lifecycle < e.config.Exclude.MinTableLifecycle {
			tableName := fmt.Sprintf("%s.%s", schema.Name(), table.Name())
			e.logger.Info("skipping table due to lifecycle less than minimum configured lifecycle",
				"table", tableName, "lifecycle", tableSchema.Lifecycle)
			return nil
		}
	}

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
		Type:        "table",
		Description: tableSchema.Comment,
		CreateTime:  timestamppb.New(time.Time(tableSchema.CreateTime)),
		UpdateTime:  timestamppb.New(time.Time(tableSchema.LastModifiedTime)),
		Service:     maxcomputeService,
	}

	tableAttributesData := e.buildTableAttributesData(schemaName, tableType, tableSchema)

	if tableType == config.TableTypeView {
		query := tableSchema.ViewText
		tableAttributesData[attributesDataSQL] = query
		if e.config.BuildViewLineage {
			upstreamResources := getUpstreamResources(query)
			asset.Lineage = &v1beta2.Lineage{
				Upstreams: upstreamResources,
			}
		}
	}

	maskingPolicy, err := e.client.GetMaskingPolicies(table)
	if err != nil {
		e.logger.Warn("error getting masking policy", "error", err)
	}

	var columns []*v1beta2.Column
	for i, col := range tableSchema.Columns {
		columnData := &v1beta2.Column{
			Name:        col.Name,
			DataType:    dataTypeToString(col.Type),
			Description: col.Comment,
			IsNullable:  !col.NotNull,
			Attributes:  utils.TryParseMapToProto(buildColumnAttributesData(&tableSchema.Columns[i])),
			Columns:     buildColumns(col.Type),
		}

		if policies, found := maskingPolicy[col.Name]; found {
			policyValues := make([]*structpb.Value, 0, len(policies))
			for _, policy := range policies {
				policyValues = append(policyValues, structpb.NewStringValue(policy))
			}
			columnData.Attributes.Fields[attributesDataMaskingPolicy] = &structpb.Value{
				Kind: &structpb.Value_ListValue{
					ListValue: &structpb.ListValue{Values: policyValues},
				},
			}
		}

		columns = append(columns, columnData)
	}

	tableData := &v1beta2.Table{
		Profile:    &v1beta2.TableProfile{TotalRows: int64(tableSchema.RecordNum)},
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

	attributesData[attributesDataProjectName] = e.config.ProjectName
	attributesData[attributesDataSchema] = schemaName
	attributesData[attributesDataType] = tableType

	rb := common.ResourceBuilder{ProjectName: e.config.ProjectName}
	attributesData[attributesDataResourceURL] = rb.Table(schemaName, tableInfo.TableName)

	if tableInfo.ViewText != "" {
		attributesData[attributesDataSQL] = tableInfo.ViewText
	}

	if tableInfo.Lifecycle != 0 {
		attributesData[attributesDataLifecycle] = tableInfo.Lifecycle
	}

	var partitionNames []interface{}
	if tableInfo.PartitionColumns != nil && len(tableInfo.PartitionColumns) > 0 {
		partitionNames = make([]interface{}, len(tableInfo.PartitionColumns))
		for i, column := range tableInfo.PartitionColumns {
			partitionNames[i] = column.Name
		}
		attributesData[attributesDataPartitionFields] = partitionNames
	}

	return attributesData
}

func buildColumnAttributesData(column *tableschema.Column) map[string]interface{} {
	attributesData := map[string]interface{}{}

	if column == nil {
		return attributesData
	}

	if column.Label != "" {
		attributesData[attributesDataLabel] = column.Label
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
		return New(plugins.GetLog(), CreateClient, seededRandom)
	}); err != nil {
		panic(err)
	}
}

func seededRandom(seed int64) func(max int64) int64 {
	rnd := rand.New(rand.NewSource(seed)) //nolint:gosec
	return func(max int64) int64 {
		return rnd.Int63n(max)
	}
}

func CreateClient(_ context.Context, _ log.Logger, conf config.Config) (Client, error) {
	return client.New(conf)
}
