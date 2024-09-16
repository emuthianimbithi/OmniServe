package main

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/template"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage project templates",
}

var addTemplateCmd = &cobra.Command{
	Use:   "add [language] [file]",
	Short: "Add a custom template for a language",
	Args:  cobra.ExactArgs(2),
	Run:   runAddTemplate,
}

var listTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all custom templates",
	Run:   runListTemplates,
}

func init() {
	templateCmd.AddCommand(addTemplateCmd)
	templateCmd.AddCommand(listTemplatesCmd)
	rootCmd.AddCommand(templateCmd)
}

func runAddTemplate(cmd *cobra.Command, args []string) {
	utils.VerboseLog("Starting runAddTemplate function")
	utils.VerboseLog(fmt.Sprintf("Arguments: language=%s, filePath=%s", args[0], args[1]))

	language := args[0]
	filePath := args[1]

	utils.VerboseLog(fmt.Sprintf("Reading template file: %s", filePath))
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error reading template file: %v", err))
		fmt.Printf("Error reading template file: %v\n", err)
		return
	}

	templateDir := cliconfig.CliConfig.Paths.Templates
	utils.VerboseLog(fmt.Sprintf("Initial template directory: %s", templateDir))

	// Expand the ~ to the home directory if present
	if templateDir[:2] == "~/" {
		utils.VerboseLog("Expanding ~ in template directory path")
		home, err := os.UserHomeDir()
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error getting home directory: %v", err))
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}
		templateDir = filepath.Join(home, templateDir[2:])
		utils.VerboseLog(fmt.Sprintf("Expanded template directory: %s", templateDir))
	}

	fmt.Printf("Creating template in: %s\n", templateDir)

	utils.VerboseLog(fmt.Sprintf("Creating template directory: %s", templateDir))
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		utils.VerboseLog(fmt.Sprintf("Error creating template directory: %v", err))
		fmt.Printf("Error creating template directory: %v\n", err)
		return
	}

	templatePath := filepath.Join(templateDir, fmt.Sprintf("%s.tmpl", language))
	utils.VerboseLog(fmt.Sprintf("Writing template to: %s", templatePath))
	if err := ioutil.WriteFile(templatePath, content, 0644); err != nil {
		utils.VerboseLog(fmt.Sprintf("Error saving template: %v", err))
		fmt.Printf("Error saving template: %v\n", err)
		return
	}

	utils.VerboseLog(fmt.Sprintf("Template for %s added successfully at %s", language, templatePath))
	fmt.Printf("Template for %s added successfully at %s\n", language, templatePath)
}

func runListTemplates(cmd *cobra.Command, args []string) {
	utils.VerboseLog("Starting runListTemplates function")

	utils.VerboseLog("Listing built-in templates")
	fmt.Println("Built-in templates:")
	listBuiltInTemplates()

	utils.VerboseLog("Listing custom templates")
	fmt.Println("\nCustom templates:")
	listCustomTemplates()
}

func listBuiltInTemplates() {
	utils.VerboseLog("Starting listBuiltInTemplates function")

	utils.VerboseLog("Reading built-in templates directory")
	entries, err := template.TemplatesFS.ReadDir("templates")
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error reading built-in templates: %v", err))
		fmt.Printf("Error reading built-in templates: %v\n", err)
		return
	}

	templateCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			utils.VerboseLog(fmt.Sprintf("Found built-in template: %s", entry.Name()))
			fmt.Printf("- %s\n", entry.Name())
			templateCount++
		}
	}

	utils.VerboseLog(fmt.Sprintf("Total built-in templates found: %d", templateCount))
}

func listCustomTemplates() {
	utils.VerboseLog("Starting listCustomTemplates function")

	templateDir := cliconfig.CliConfig.Paths.Templates
	utils.VerboseLog(fmt.Sprintf("Initial custom template directory: %s", templateDir))

	// Expand the ~ to the home directory if present
	if templateDir[:2] == "~/" {
		utils.VerboseLog("Expanding ~ in custom template directory path")
		home, err := os.UserHomeDir()
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error getting home directory: %v", err))
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}
		templateDir = filepath.Join(home, templateDir[2:])
		utils.VerboseLog(fmt.Sprintf("Expanded custom template directory: %s", templateDir))
	}

	utils.VerboseLog(fmt.Sprintf("Looking for custom templates in: %s", templateDir))

	// Check if the directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		utils.VerboseLog(fmt.Sprintf("Custom template directory does not exist: %s", templateDir))
		fmt.Printf("Custom template directory does not exist: %s\n", templateDir)
		return
	}

	utils.VerboseLog("Reading custom template directory")
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error reading custom template directory: %v", err))
		fmt.Printf("Error reading custom template directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		utils.VerboseLog("No files found in custom template directory")
		fmt.Println("No custom templates found.")
		return
	}

	templateCount := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tmpl" {
			utils.VerboseLog(fmt.Sprintf("Found custom template: %s", file.Name()))
			fmt.Printf("- %s\n", file.Name())
			templateCount++
		}
	}

	if templateCount == 0 {
		utils.VerboseLog("No .tmpl files found in the custom template directory")
		fmt.Println("No custom .tmpl files found in the template directory.")
	} else {
		utils.VerboseLog(fmt.Sprintf("Total custom templates found: %d", templateCount))
	}
}
