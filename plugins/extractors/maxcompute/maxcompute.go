package maxcompute

import (
	"context"
	_ "embed" // used to print the embedded assets
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"
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

	client *odps.Odps
	eg     *errgroup.Group
}

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

func New(logger log.Logger) *Extractor {
	e := &Extractor{
		logger: logger,
	}
	e.BaseExtractor = plugins.NewBaseExtractor(info, &e.config)
	e.ScopeNotRequired = true

	return e
}

func (e *Extractor) Init(ctx context.Context, config plugins.Config) error {
	if err := e.BaseExtractor.Init(ctx, config); err != nil {
		return err
	}

	aliAccount := account.NewAliyunAccount(e.config.AccessKey.ID, e.config.AccessKey.Secret)
	e.client = odps.NewOdps(aliAccount, e.config.EndpointProject)
	e.client.SetDefaultProjectName(e.config.ProjectName)

	e.eg = &errgroup.Group{}
	e.eg.SetLimit(e.config.Concurrency)

	return nil
}

func (e *Extractor) Extract(_ context.Context, emit plugins.Emit) error {
	schemas := e.client.Project(e.config.ProjectName).Schemas()

	err := schemas.List(func(schema *odps.Schema, err error) {
		if err != nil {
			e.logger.Error("failed to list schemas", "error", err)
			return
		}
		if e.config.SchemaName != "" && schema.Name() != e.config.SchemaName {
			return
		}
		if contains(e.config.Exclude.Schemas, schema.Name()) {
			return
		}

		tables := odps.NewTables(e.client, e.config.ProjectName, schema.Name())

		e.fetchTablesFromSchema(tables, emit)
	})
	if err != nil {
		return err
	}
	return e.eg.Wait()
}

func (e *Extractor) fetchTablesFromSchema(tables *odps.Tables, emit plugins.Emit) {
	tables.List(
		func(t *odps.Table, err error) {
			if err != nil {
				e.logger.Error("table list", err)
				return
			}

			if contains(e.config.Exclude.Tables, fmt.Sprintf("%s.%s", t.SchemaName(), t.Name())) {
				return
			}

			e.eg.Go(func() error {
				return e.processTable(t, emit)
			})
		},
	)
}

func (e *Extractor) processTable(table *odps.Table, emit plugins.Emit) error {
	err := table.Load()
	if err != nil {
		isView := table.Schema().IsVirtualView || table.Schema().IsMaterializedView
		isLoaded := table.IsLoaded()
		if !isView || (isView && !isLoaded) {
			e.logger.Error("failed to get table info", "schema", table.SchemaName(), "table", table.Name(), "error", err)
			return err
		}
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
		return New(plugins.GetLog())
	}); err != nil {
		panic(err)
	}
}
