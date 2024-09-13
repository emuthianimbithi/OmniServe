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
	language := args[0]
	filePath := args[1]

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		return
	}

	templateDir := cliconfig.CliConfig.Paths.Templates

	// Expand the ~ to the home directory if present
	if templateDir[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}
		templateDir = filepath.Join(home, templateDir[2:])
	}

	fmt.Printf("Creating template in: %s\n", templateDir)

	if err := os.MkdirAll(templateDir, 0755); err != nil {
		fmt.Printf("Error creating template directory: %v\n", err)
		return
	}

	templatePath := filepath.Join(templateDir, fmt.Sprintf("%s.tmpl", language))
	if err := ioutil.WriteFile(templatePath, content, 0644); err != nil {
		fmt.Printf("Error saving template: %v\n", err)
		return
	}

	fmt.Printf("Template for %s added successfully at %s\n", language, templatePath)
}

func runListTemplates(cmd *cobra.Command, args []string) {
	fmt.Println("Built-in templates:")
	listBuiltInTemplates()

	fmt.Println("\nCustom templates:")
	listCustomTemplates()
}

func listBuiltInTemplates() {
	entries, err := template.TemplatesFS.ReadDir("templates")
	if err != nil {
		fmt.Printf("Error reading built-in templates: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			fmt.Printf("- %s\n", entry.Name())
		}
	}
}

func listCustomTemplates() {
	templateDir := cliconfig.CliConfig.Paths.Templates

	// Expand the ~ to the home directory if present
	if templateDir[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}
		templateDir = filepath.Join(home, templateDir[2:])
	}

	utils.VerboseLog(fmt.Sprintf("Looking for custom templates in: %s\n", templateDir))

	// Check if the directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		fmt.Printf("Custom template directory does not exist: %s\n", templateDir)
		return
	}

	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		fmt.Printf("Error reading custom template directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No custom templates found.")
		return
	}

	templateFound := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tmpl" {
			fmt.Printf("- %s\n", file.Name())
			templateFound = true
		}
	}

	if !templateFound {
		fmt.Println("No custom .tmpl files found in the template directory.")
	}
}
