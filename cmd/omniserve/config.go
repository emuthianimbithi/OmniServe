package main

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	forceOverwrite bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage OmniServe configuration",
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default configuration file",
	Run:   runInitConfig,
}

var deleteConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the configuration file",
	Run:   runDeleteConfig,
}

func init() {
	initConfigCmd.Flags().BoolVar(&forceOverwrite, "force", false, "Force overwrite of existing config file")
	configCmd.AddCommand(initConfigCmd)
	configCmd.AddCommand(deleteConfigCmd)

	// add the config commands to the root command
	rootCmd.AddCommand(configCmd)
}

func runInitConfig(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	if configPath == "" {
		configPath = filepath.Join(home, ".omniserve.yaml")
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		if !forceOverwrite {
			fmt.Println("Configuration file already exists. Use --force to overwrite back to default configuration.")
			return
		}
	}

	// Create default configuration
	defaultConfig := variables.DefaultConfig

	err = os.WriteFile(configPath, []byte(defaultConfig), 0644)
	if err != nil {
		fmt.Println("Error writing configuration file:", err)
		return
	}

	fmt.Println("Default configuration file created at:", configPath)
}

func runDeleteConfig(cmd *cobra.Command, args []string) {
	configPath := cliconfig.GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file does not exist.")
		return
	}

	fmt.Printf("Are you sure you want to delete the config file at %s? (y/N): ", configPath)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil || (response != "y" && response != "Y") {
		fmt.Println("Operation cancelled.")
		return
	}

	err = os.Remove(configPath)
	if err != nil {
		fmt.Printf("Error deleting config file: %v\n", err)
		return
	}

	fmt.Println("Config file deleted successfully.")
}
