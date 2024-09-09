package cliconfig

import (
	"errors"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/models"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
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
		License:   "",
		Version:   "0.1.0",
		Author:    "",
		GitInit:   true,
		Dockerize: true,
	},
}

var Config models.CLIConfig

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
			// Config file not found; use defaults
			Config = defaultConfig
			return nil
		}
		// Config file was found but another error was produced
		return err
	}

	// Config file found and successfully parsed
	err := viper.Unmarshal(&Config)
	if err != nil {
		return err
	}

	utils.VerboseLog(fmt.Sprintf("Config loaded from %s", viper.ConfigFileUsed()))
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
