package maxcompute

import (
	"context"
	_ "embed" // used to print the embedded assets
	"sync"
	"time"

	client2 "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	maxcomputeclient "github.com/alibabacloud-go/maxcompute-20220104/client"
	"github.com/goto/meteor/models"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/registry"
	"github.com/goto/meteor/utils"
	"github.com/goto/salt/log"
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
	SchemaName     string `mapstructure:"schema_name"`
	MaxPreviewRows int    `mapstructure:"max_preview_rows"`
	Exclude        struct {
		Schemas []string `mapstructure:"schemas"`
		Tables  []string `mapstructure:"tables"`
	} `mapstructure:"exclude"`
	MaxPageSize          int32    `mapstructure:"max_page_size"`
	Concurrency          int      `mapstructure:"concurrency"`
	MixValues            bool     `mapstructure:"mix_values"`
	IncludeColumnProfile bool     `mapstructure:"include_column_profile"`
	BuildViewLineage     bool     `mapstructure:"build_view_lineage"`
	IsCollectTableUsage  bool     `mapstructure:"collect_table_usage"`
	UsagePeriodInDay     int      `mapstructure:"usage_period_in_day"`
	UsageProjectNames    []string `mapstructure:"usage_project_names"`
}

type Extractor struct {
	plugins.BaseExtractor
	logger log.Logger
	config Config
}

//go:embed README.md
var summary string

var sampleConfig = `
project_name: goto_test
endpoint_project: http://goto_test-maxcompute.com
access_key:
    id: access_key_id
    secret: access_key_secret
schema_name: DEFAULT
max_preview_rows: 3
exclude:
	schemas:
	- schema_a
	- schema_b
	tables:
	- schema_c.table_a
max_page_size: 100
concurrency: 10
mix_values: false
include_column_profile: true
build_view_lineage: true
collect_table_usage: false
usage_period_in_day: 7
usage_project_names:
	- maxcompute-project-name
	- other-maxcompute-project-name`

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
	return e.BaseExtractor.Init(ctx, config)
}

func (e *Extractor) Extract(_ context.Context, emit plugins.Emit) error {
	apiClient, err := maxcomputeclient.NewClient(&client2.Config{
		AccessKeyId:     &e.config.AccessKey.ID,
		AccessKeySecret: &e.config.AccessKey.Secret,
		Endpoint:        &e.config.EndpointProject,
	})
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	var marker string
	var counter int
	for {
		e.logger.Info("fetching tables", "marker", marker)
		resp, err := apiClient.ListTables(&e.config.ProjectName, &maxcomputeclient.ListTablesRequest{
			MaxItem: &e.config.MaxPageSize,
			Marker:  &marker,
		})
		if err != nil {
			panic(err)
		}

		counter += len(resp.Body.Data.Tables)
		e.logger.Info("fetched tables", "count", counter)
		wg.Add(len(resp.Body.Data.Tables))

		for _, table := range resp.Body.Data.Tables {
			go func(table *maxcomputeclient.ListTablesResponseBodyDataTables) {
				defer wg.Done()
				tableInfo, err := apiClient.GetTableInfo(&e.config.ProjectName, table.Name, &maxcomputeclient.GetTableInfoRequest{})
				if err != nil {
					panic(err)
				}

				asset, err := e.buildAsset(tableInfo)
				if err != nil {
					e.logger.Error("failed to build asset", "table", *table.Name, "error", err)
					return
				}

				emit(models.NewRecord(asset))
			}(table)
		}
		wg.Wait()
		if len(resp.Body.Data.Tables) == 0 || len(resp.Body.Data.Tables) < int(e.config.MaxPageSize) || *resp.Body.Data.Marker == "" {
			break
		}
		marker = *resp.Body.Data.Marker
	}

	return nil
}

func (e *Extractor) buildAsset(tableInfo *maxcomputeclient.GetTableInfoResponse) (*v1beta2.Asset, error) {
	defaultSchema := "default"
	if tableInfo.Body.Data.Schema == nil {
		tableInfo.Body.Data.Schema = &defaultSchema
	}

	tableURN := plugins.MaxComputeURN(*tableInfo.Body.Data.ProjectName, *tableInfo.Body.Data.Schema, *tableInfo.Body.Data.Name)

	// TODO(mayur): Add all the P0 metadata
	// Table Name ✅
	// Project Name (Data Layer) ✅
	// Table Type ✅
	// Created_At ✅
	// Last_Updated_At (for both DML & DDL) ✅ ✅
	// Partition Field
	// Table Schema ✅
	// Table Description ✅
	// Resource URL
	// Table Refresh Frequency → previously get from Optimus, but we need to provide it as well with tengo script
	// Table SQL
	// User Query References

	asset := &v1beta2.Asset{
		Urn:         tableURN,
		Name:        *tableInfo.Body.Data.Name,
		Type:        *tableInfo.Body.Data.Type,
		Description: *tableInfo.Body.Data.Comment,
		CreateTime:  convertInt64ToTimestamp(tableInfo.Body.Data.CreationTime),
		UpdateTime:  convertInt64ToTimestamp(tableInfo.Body.Data.LastModifiedTime),
		Service:     "maxcompute",
	}

	attributesData := map[string]interface{}{
		"project_name": *tableInfo.Body.Data.ProjectName,
		"schema":       *tableInfo.Body.Data.Schema,
	}
	tableData := &v1beta2.Table{
		Attributes: utils.TryParseMapToProto(attributesData),
		UpdateTime: convertInt64ToTimestamp(tableInfo.Body.Data.LastDDLTime),
	}

	table, err := anypb.New(tableData)
	if err != nil {
		e.logger.Warn("error creating Any struct", "error", err)
	}
	asset.Data = table

	return asset, nil
}

func convertInt64ToTimestamp(unixTimeMs *int64) *timestamppb.Timestamp {
	if unixTimeMs == nil {
		return nil
	}

	seconds := *unixTimeMs / 1000
	nanoseconds := (*unixTimeMs % 1000) * 1_000_000

	t := time.Unix(seconds, nanoseconds).UTC()
	return timestamppb.New(t)
}

func init() {
	if err := registry.Extractors.Register("maxcompute", func() plugins.Extractor {
		return New(plugins.GetLog())
	}); err != nil {
		panic(err)
	}
}
