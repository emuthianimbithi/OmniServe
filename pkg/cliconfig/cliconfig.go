package cliconfig

import (
	"errors"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/models"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var defaultConfig = models.CLIConfig{
	Defaults: struct {
		Language  string `mapstructure:"language"`
		License   string `mapstructure:"license"`
		Version   string `mapstructure:"version"`
		Author    string `mapstructure:"author"`
		GitInit   bool   `mapstructure:"git_init"`
		Dockerize bool   `mapstructure:"dockerize"`
	}{
		Language:  "go",
		License:   "default",
		Version:   "0.1.0",
		Author:    "default",
		GitInit:   false,
		Dockerize: true,
	},
	Paths: struct {
		Templates string `mapstructure:"templates"`
	}{
		Templates: "~/.omniserve/templates",
	},
	Languages: map[string]struct {
		EntryPoint   string `mapstructure:"entry_point"`
		BuildCommand string `mapstructure:"build_command"`
	}{
		"go": {
			EntryPoint:   "main.go",
			BuildCommand: "go build",
		},
		"python": {
			EntryPoint:   "main.py",
			BuildCommand: "python -m compileall",
		},
		"javascript": {
			EntryPoint:   "index.js",
			BuildCommand: "npm run build",
		},
		"c": {
			EntryPoint:   "main.c",
			BuildCommand: "gcc -o main main.c",
		},
	},
	CLI: struct {
		Verbose     bool `mapstructure:"verbose"`
		ColorOutput bool `mapstructure:"color_output"`
	}{
		Verbose:     false,
		ColorOutput: true,
	},
	Server: struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}{
		Host: "0.0.0.0:8765",
		Port: "8765",
	},
}

var CliConfig models.CLIConfig

func LoadConfig(cfgFile string) error {

	// get the file from the specified path
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		// get the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".omniserve")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		// Check if the error is a config file not found error
		// If it is, use the default config
		if errors.As(err, &configFileNotFoundError) {
			// CliConfig file not found; use defaults
			CliConfig = defaultConfig
			return nil
		}
		// CliConfig file was found but another error was produced
		return err
	}

	// CliConfig file found and successfully parsed
	err := viper.Unmarshal(&CliConfig)
	if err != nil {
		return err
	}

	variables.Verbose = CliConfig.CLI.Verbose

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	if viper.ConfigFileUsed() != "" {
		return viper.ConfigFileUsed()
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return ""
	}

	return filepath.Join(home, ".omniserve.yaml")
}
