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
	"github.com/goto/meteor/registry"
	"github.com/goto/meteor/utils"
	"github.com/goto/salt/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Extractor struct {
	plugins.BaseExtractor
	logger log.Logger
	config config.Config

	client    Client
	newClient NewClientFunc
	eg        *errgroup.Group
}

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

	asset, err := e.buildAsset(schema, table, tableType, tableSchema)
	if err != nil {
		e.logger.Error("failed to build asset", "table", table.Name(), "error", err)
		return err
	}

	emit(models.NewRecord(asset))
	return nil
}

func (e *Extractor) buildAsset(schema *odps.Schema, _ *odps.Table, tableType string, tableSchema *tableschema.TableSchema) (*v1beta2.Asset, error) {
	defaultSchema := "default"
	schemaName := schema.Name()
	if schemaName == "" {
		schemaName = defaultSchema
	}

	tableURN := plugins.MaxComputeURN(e.config.ProjectName, schemaName, tableSchema.TableName)

	asset := &v1beta2.Asset{
		Urn:         tableURN,
		Name:        tableSchema.TableName,
		Type:        tableType,
		Description: tableSchema.Comment,
		CreateTime:  timestamppb.New(time.Time(tableSchema.CreateTime)),
		UpdateTime:  timestamppb.New(time.Time(tableSchema.LastModifiedTime)),
		Service:     "maxcompute",
	}

	tableAttributesData := e.buildTableAttributesData(schemaName, tableSchema)

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

	tbl, err := anypb.New(tableData)
	if err != nil {
		e.logger.Warn("error creating Any struct", "error", err)
	}
	asset.Data = tbl

	return asset, nil
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

func (e *Extractor) buildTableAttributesData(schemaName string, tableInfo *tableschema.TableSchema) map[string]interface{} {
	attributesData := map[string]interface{}{}

	attributesData["project_name"] = e.config.ProjectName
	attributesData["schema"] = schemaName

	rb := common.ResourceBuilder{ProjectName: e.config.ProjectName}
	attributesData["resource_url"] = rb.Table(tableInfo.TableName)

	if tableInfo.ViewText != "" {
		attributesData["sql"] = tableInfo.ViewText
	}

	if tableInfo.PartitionColumns != nil && len(tableInfo.PartitionColumns) > 0 {
		partitionNames := make([]string, len(tableInfo.PartitionColumns))
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
	if err := registry.Extractors.Register("maxcompute", func() plugins.Extractor {
		return New(plugins.GetLog(), CreateClient)
	}); err != nil {
		panic(err)
	}
}

func CreateClient(_ context.Context, _ log.Logger, conf config.Config) (Client, error) {
	return client.New(conf), nil
}
