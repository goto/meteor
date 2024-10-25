package maxcompute

import (
	"context"
	"log"

	"github.com/goto/meteor/plugins"
)

type Config struct {
	ProjectName     string `mapstructure:"project_name"`
	EndpointProject string `mapstructure:"endpoint_project"`
	AccessKeyJSON   struct {
		ID           string `mapstructure:"id"`
		SecretBase64 string `mapstructure:"secret_base64"`
	} `mapstructure:"access_key_json"`
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
access_key_json: {
	id: xyz
	secret_base64: __base64__
}
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
	return nil
}
