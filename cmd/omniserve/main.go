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
)

var rootCmd = &cobra.Command{
	Use:   "omniserve",
	Short: "OmniServe - A multi-language serverless platform CLI",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new serverless project",
	Run:   runInit,
}

func init() {
	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project (required)")
	initCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language (go, c, python, javascript) (required)")
	initCmd.Flags().StringVarP(&entryPoint, "entry-point", "e", "", "Path to the entry point file (optional)")
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("language")
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
