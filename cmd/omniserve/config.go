package main

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
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
	utils.VerboseLog("Starting runInitConfig function")

	home, err := os.UserHomeDir()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error getting home directory: %v", err))
		fmt.Println("Error getting home directory:", err)
		return
	}
	utils.VerboseLog(fmt.Sprintf("Home directory: %s", home))

	if configPath == "" {
		configPath = filepath.Join(home, ".omniserve.yaml")
		utils.VerboseLog(fmt.Sprintf("No config path provided, using default: %s", configPath))
	} else {
		utils.VerboseLog(fmt.Sprintf("Using provided config path: %s", configPath))
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		utils.VerboseLog("Configuration file already exists")
		if !forceOverwrite {
			utils.VerboseLog("Force overwrite not enabled, exiting")
			fmt.Println("Configuration file already exists. Use --force to overwrite back to default configuration.")
			return
		}
		utils.VerboseLog("Force overwrite enabled, proceeding with file creation")
	} else {
		utils.VerboseLog("Configuration file does not exist, proceeding with creation")
	}

	// Create default configuration
	defaultConfig := variables.DefaultConfig
	utils.VerboseLog("Created default configuration")

	utils.VerboseLog(fmt.Sprintf("Attempting to write configuration file to: %s", configPath))
	err = os.WriteFile(configPath, []byte(defaultConfig), 0644)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error writing configuration file: %v", err))
		fmt.Println("Error writing configuration file:", err)
		return
	}

	utils.VerboseLog("Configuration file successfully written")
	fmt.Println("Default configuration file created at:", configPath)
}

func runDeleteConfig(cmd *cobra.Command, args []string) {
	utils.VerboseLog("Starting runDeleteConfig function")

	configPath := cliconfig.GetConfigPath()
	utils.VerboseLog(fmt.Sprintf("Config path: %s", configPath))

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		utils.VerboseLog("Config file does not exist")
		fmt.Println("Config file does not exist.")
		return
	}
	utils.VerboseLog("Config file exists")

	fmt.Printf("Are you sure you want to delete the config file at %s? (y/N): ", configPath)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error reading user input: %v", err))
	}
	utils.VerboseLog(fmt.Sprintf("User response: %s", response))

	if err != nil || (response != "y" && response != "Y") {
		utils.VerboseLog("Operation cancelled by user or due to input error")
		fmt.Println("Operation cancelled.")
		return
	}

	utils.VerboseLog("Attempting to delete config file")
	err = os.Remove(configPath)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error deleting config file: %v", err))
		fmt.Printf("Error deleting config file: %v\n", err)
		return
	}

	utils.VerboseLog("Config file deleted successfully")
	fmt.Println("Config file deleted successfully.")
}
