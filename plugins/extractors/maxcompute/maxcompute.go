package maxcompute

import (
	"context"
	_ "embed" // used to print the embedded assets
	"fmt"

	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/registry"
	"github.com/goto/salt/log"

	client2 "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	maxcomputeclient "github.com/alibabacloud-go/maxcompute-20220104/client"
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
	MaxPageSize          int      `mapstructure:"max_page_size"`
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

// Init initializes the extractor
func (e *Extractor) Init(ctx context.Context, config plugins.Config) error {
	if err := e.BaseExtractor.Init(ctx, config); err != nil {
		return err
	}

	return nil
}

// Extract checks if the table is valid and extracts the table schema
func (e *Extractor) Extract(ctx context.Context, emit plugins.Emit) error {
	apiClient, err := maxcomputeclient.NewClient(&client2.Config{
		AccessKeyId:     &e.config.AccessKey.ID,
		AccessKeySecret: &e.config.AccessKey.Secret,
		Endpoint:        &e.config.EndpointProject,
	})
	if err != nil {
		panic(err)
	}

	resp, err := apiClient.ListTables(&e.config.ProjectName, &maxcomputeclient.ListTablesRequest{})
	if err != nil {
		panic(err)
	}

	for _, table := range resp.Body.Data.Tables {
		fmt.Printf("Table: %s\n", *table.Name)
		resp, err := apiClient.GetTableInfo(&e.config.ProjectName, table.Name, &maxcomputeclient.GetTableInfoRequest{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Table Info: %v\n", resp.Body.Data)
	}
	return nil
}

// Register the extractor to catalog
func init() {
	if err := registry.Extractors.Register("maxcompute", func() plugins.Extractor {
		return New(plugins.GetLog())
	}); err != nil {
		panic(err)
	}
}
