package config

const TableTypeView = "VIRTUAL_VIEW"

type Config struct {
	ProjectName     string `mapstructure:"project_name"`
	EndpointProject string `mapstructure:"endpoint_project"`
	AccessKey       struct {
		ID     string `mapstructure:"id"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"access_key"`
	SchemaName string `mapstructure:"schema_name,omitempty"`
	Exclude    struct {
		Schemas           []string `mapstructure:"schemas"`
		Tables            []string `mapstructure:"tables"`
		MinTableLifecycle int      `mapstructure:"min_table_lifecycle"`
	} `mapstructure:"exclude,omitempty"`
	MaxPreviewRows   int  `mapstructure:"max_preview_rows,omitempty"`
	MixValues        bool `mapstructure:"mix_values,omitempty"`
	Concurrency      int  `mapstructure:"concurrency,omitempty"`
	BuildViewLineage bool `mapstructure:"build_view_lineage,omitempty"`
}
