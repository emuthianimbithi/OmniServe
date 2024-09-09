package models

type CLIConfig struct {
	Defaults struct {
		Language  string `mapstructure:"language"`
		License   string `mapstructure:"license"`
		Version   string `mapstructure:"version"`
		Author    string `mapstructure:"author"`
		GitInit   bool   `mapstructure:"git_init"`
		Dockerize bool   `mapstructure:"dockerize"`
	} `mapstructure:"defaults"`

	Paths struct {
		Templates string `mapstructure:"templates"`
	} `mapstructure:"paths"`

	Languages map[string]struct {
		EntryPoint   string `mapstructure:"entry_point"`
		BuildCommand string `mapstructure:"build_command"`
	} `mapstructure:"languages"`

	Aliases map[string]string `mapstructure:"aliases"`

	CLI struct {
		Verbose     bool `mapstructure:"verbose"`
		ColorOutput bool `mapstructure:"color_output"`
	} `mapstructure:"cli"`

	CloudProviders map[string]map[string]string `mapstructure:"cloud_providers"`
}