package main

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/config"
	"github.com/emuthianimbithi/OmniServe/pkg/docker"
	"github.com/emuthianimbithi/OmniServe/pkg/template"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	projectName string
	language    string
	entryPoint  string
	version     string
	author      string
	description string
	licenseName string
	gitInit     bool
	dockerize   bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new serverless project",
	Run:   runInit,
}

func init() {
	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project (required)")
	initCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language (go, c, python, javascript) (required)")
	initCmd.Flags().StringVarP(&entryPoint, "entry-point", "e", "", "Path to the entry point file (optional)")
	initCmd.Flags().StringVarP(&version, "version", "", "0.1.0", "Initial version of the project")
	initCmd.Flags().StringVarP(&author, "author", "a", "", "Author of the project")
	initCmd.Flags().StringVarP(&description, "description", "d", "", "Short description of the project")
	initCmd.Flags().StringVarP(&licenseName, "license", "", "MIT", "License for the project")
	initCmd.Flags().BoolVarP(&gitInit, "git-init", "g", false, "Initialize a git repository")
	initCmd.Flags().BoolVarP(&dockerize, "dockerize", "D", false, "Add a Dockerfile to the project")

	err := initCmd.MarkFlagRequired("name")
	if err != nil {
		return
	}
	err = initCmd.MarkFlagRequired("language")
	if err != nil {
		return
	}

	// Add the init command to the root command
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	utils.VerboseLog("Starting runInit function")
	utils.VerboseLog(fmt.Sprintf("Project Name: %s, Language: %s, Entry Point: %s", projectName, language, entryPoint))

	if !utils.IsValidLanguage(language) {
		utils.VerboseLog(fmt.Sprintf("Unsupported language: %s", language))
		fmt.Printf("Unsupported language: %s. Supported languages are: go, c, python, javascript\n", language)
		return
	}
	utils.VerboseLog("Language validation passed")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error getting current directory: %v", err))
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}
	utils.VerboseLog(fmt.Sprintf("Current working directory: %s", cwd))

	// Create project directory
	projectPath := filepath.Join(cwd, projectName)
	utils.VerboseLog(fmt.Sprintf("Creating project directory: %s", projectPath))
	err = os.MkdirAll(projectPath, 0755)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error creating project directory: %v", err))
		fmt.Printf("Error creating project directory: %v\n", err)
		return
	}
	utils.VerboseLog("Project directory created successfully")

	// Create configuration file
	utils.VerboseLog("Creating project configuration")
	cfg, err := config.NewProjectConfig(projectPath, projectName, language, entryPoint)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error creating project configuration: %v", err))
		fmt.Printf("Error creating project configuration: %v\n", err)
		return
	}
	utils.VerboseLog("Project configuration created successfully")

	// Create entry point file
	utils.VerboseLog(fmt.Sprintf("Creating entry point file: %s", cfg.EntryPoint))
	err = template.CreateEntryPointFile(projectPath, cfg.EntryPoint, language)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error creating entry point file: %v", err))
		fmt.Printf("Error creating entry point file: %v\n", err)
		return
	}
	utils.VerboseLog("Entry point file created successfully")

	// Create docker file
	if cliconfig.CliConfig.Defaults.Dockerize {
		utils.VerboseLog("Dockerize flag is set, creating Dockerfile")
		fmt.Println("Creating Dockerfile")
		err = docker.CreateDockerfile(projectPath, language, entryPoint)
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error creating Dockerfile: %v", err))
			fmt.Printf("Error creating Dockerfile: %v\n", err)
		} else {
			utils.VerboseLog("Dockerfile created successfully")
		}
	} else {
		utils.VerboseLog("Dockerize flag is not set, skipping Dockerfile creation")
	}

	utils.VerboseLog("Project initialization completed successfully")
	fmt.Printf("Successfully initialized %s project: %s\n", language, projectName)
	fmt.Println("Next steps:")
	fmt.Printf("1. cd %s\n", projectName)
	fmt.Println("2. Edit omniserve.json to add dependencies or modify configuration")
	fmt.Printf("3. Implement your serverless function in %s\n", cfg.EntryPoint)
}
