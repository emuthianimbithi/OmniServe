package main

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/config"
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

	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("language")

	// Add the init command to the root command
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	if !utils.IsValidLanguage(language) {
		fmt.Printf("Unsupported language: %s. Supported languages are: go, c, python, javascript\n", language)
		return
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Create project directory
	projectPath := filepath.Join(cwd, projectName)
	err = os.MkdirAll(projectPath, 0755)
	if err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		return
	}

	// Create configuration file
	cfg, err := config.NewProjectConfig(projectPath, projectName, language, entryPoint)
	if err != nil {
		fmt.Printf("Error creating project configuration: %v\n", err)
		return
	}

	// Create entry point file
	err = utils.CreateEntryPointFile(projectPath, cfg.EntryPoint, language)
	if err != nil {
		fmt.Printf("Error creating entry point file: %v\n", err)
		return
	}

	fmt.Printf("Successfully initialized %s project: %s\n", language, projectName)
	fmt.Println("Next steps:")
	fmt.Printf("1. cd %s\n", projectName)
	fmt.Println("2. Edit omniserve.json to add dependencies or modify configuration")
	fmt.Printf("3. Implement your serverless function in %s\n", cfg.EntryPoint)
}
