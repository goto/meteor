package maxcompute

import (
	"context"
	_ "embed" // used to print the embedded assets
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/datatype"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tableschema"
	"github.com/goto/meteor/models"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/registry"
	"github.com/goto/meteor/utils"
	"github.com/goto/salt/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Config struct {
	ProjectName     string `mapstructure:"project_name"`
	EndpointProject string `mapstructure:"endpoint_project"`
	AccessKey       struct {
		ID     string `mapstructure:"id"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"access_key"`
	SchemaName string `mapstructure:"schema_name,omitempty"`
	Exclude    struct {
		Schemas []string `mapstructure:"schemas"`
		Tables  []string `mapstructure:"tables"`
	} `mapstructure:"exclude,omitempty"`
	Concurrency int `mapstructure:"concurrency,omitempty"`
}

type Extractor struct {
	plugins.BaseExtractor
	logger log.Logger
	config Config

	client    Client
	newClient NewClientFunc
	eg        *errgroup.Group
}

type NewClientFunc func(ctx context.Context, logger log.Logger, config Config) (Client, error)

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
	GetTable(ctx context.Context, table *odps.Table) (*odps.Table, error)
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

func (e *Extractor) Init(ctx context.Context, config plugins.Config) error {
	if err := e.BaseExtractor.Init(ctx, config); err != nil {
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
	if err != nil {
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
	if err != nil {
		return err
	}

	for _, table := range tables {
		if contains(e.config.Exclude.Tables, fmt.Sprintf("%s.%s", table.SchemaName(), table.Name())) {
			continue
		}

		tbl := table
		e.eg.Go(func() error {
			return e.processTable(ctx, tbl, emit)
		})
	}

	return nil
}

func (e *Extractor) processTable(ctx context.Context, table *odps.Table, emit plugins.Emit) error {
	table, err := e.client.GetTable(ctx, table)
	if err != nil {
		return err
	}

	asset, err := e.buildAsset(table)
	if err != nil {
		e.logger.Error("failed to build asset", "table", table.Name(), "error", err)
		return err
	}

	emit(models.NewRecord(asset))
	return nil
}

func (e *Extractor) buildAsset(tableInfo *odps.Table) (*v1beta2.Asset, error) {
	defaultSchema := "default"
	schemaName := tableInfo.SchemaName()
	if schemaName == "" {
		schemaName = defaultSchema
	}

	tableURN := plugins.MaxComputeURN(tableInfo.ProjectName(), schemaName, tableInfo.Name())

	schema := tableInfo.Schema()
	asset := &v1beta2.Asset{
		Urn:         tableURN,
		Name:        schema.TableName,
		Type:        tableInfo.Type().String(),
		Description: schema.Comment,
		CreateTime:  timestamppb.New(tableInfo.CreatedTime()),
		UpdateTime:  timestamppb.New(tableInfo.LastModifiedTime()),
		Service:     "maxcompute",
	}

	tableAttributesData := buildTableAttributesData(tableInfo)

	var columns []*v1beta2.Column
	for i, col := range schema.Columns {
		columnData := &v1beta2.Column{
			Name:        col.Name,
			DataType:    dataTypeToString(col.Type),
			Description: col.Comment,
			IsNullable:  col.IsNullable,
			Attributes:  utils.TryParseMapToProto(buildColumnAttributesData(&schema.Columns[i])),
			Columns:     buildColumns(col.Type),
		}
		columns = append(columns, columnData)
	}

	tableData := &v1beta2.Table{
		Attributes: utils.TryParseMapToProto(tableAttributesData),
		Columns:    columns,
		CreateTime: timestamppb.New(tableInfo.CreatedTime()),
		UpdateTime: timestamppb.New(tableInfo.LastModifiedTime()),
	}

	table, err := anypb.New(tableData)
	if err != nil {
		e.logger.Warn("error creating Any struct", "error", err)
	}
	asset.Data = table

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

func buildTableAttributesData(tableInfo *odps.Table) map[string]interface{} {
	attributesData := map[string]interface{}{}

	if tableInfo == nil {
		return attributesData
	}

	if tableInfo.ProjectName() != "" {
		attributesData["project_name"] = tableInfo.ProjectName()
	}

	if tableInfo.SchemaName() != "" {
		attributesData["schema"] = tableInfo.SchemaName()
	}

	if tableInfo.ViewText() != "" {
		attributesData["sql"] = tableInfo.ViewText()
	}

	if tableInfo.ResourceUrl() != "" {
		attributesData["resource_url"] = tableInfo.ResourceUrl()
	}

	// TODO: map the column and add the column name in the field
	// if tableInfo.PartitionColumns() != nil {
	// 	attributesData["partition_field"] = tableInfo.PartitionColumns()
	// }

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

func CreateClient(_ context.Context, _ log.Logger, config Config) (Client, error) {
	return NewMaxComputeClient(config), nil
}
