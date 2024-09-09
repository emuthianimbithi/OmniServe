package template

import (
	"embed"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
)

var (
	//go:embed templates/*
	templatesFS embed.FS
)

func CreateEntryPointFile(projectPath, entryPoint, language string) error {
	fullPath := filepath.Join(projectPath, entryPoint)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	template, err := GetTemplate(language)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fullPath, []byte(template), 0644)
}

func GetTemplate(language string) (string, error) {
	// First, try to get a custom template
	customTemplate, err := getCustomTemplate(language)
	if err == nil {
		utils.VerboseLog(fmt.Sprintf("Using custom template for %s", language))
		return customTemplate, nil
	}
	utils.VerboseLog(fmt.Sprintf("No custom template found for %s, using built-in template", language))

	// If no custom template, fall back to the built-in template
	return getBuiltInTemplate(language)
}

func getCustomTemplate(language string) (string, error) {
	templateDir := cliconfig.Config.Paths.Templates
	if templateDir == "" {
		return "", fmt.Errorf("template path not set in configuration")
	}

	// Expand the ~ to the home directory if present
	if templateDir[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %v", err)
		}
		templateDir = filepath.Join(home, templateDir[2:])
	}

	customTemplatePath := filepath.Join(templateDir, fmt.Sprintf("%s.tmpl", language))
	content, err := ioutil.ReadFile(customTemplatePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getBuiltInTemplate(language string) (string, error) {
	templateFile := fmt.Sprintf("templates/%s.tmpl", language)
	templateContent, err := templatesFS.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("error reading built-in template file: %v", err)
	}
	return string(templateContent), nil
}
