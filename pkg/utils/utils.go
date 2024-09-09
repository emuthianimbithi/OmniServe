package utils

import (
	"embed"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"os"
	"path/filepath"
)

var (
	//go:embed templates/*
	templatesFS embed.FS

	Verbose bool
)

func IsValidLanguage(lang string) bool {
	return variables.SupportedLanguages[lang]
}

func GetDefaultEntryPoint(language string) string {
	return variables.EntryPointTemplate[language]
}

// VerboseLog used to give more detailed output
func VerboseLog(message string) {
	if Verbose {
		fmt.Fprintln(os.Stderr, "VERBOSE:", message)
	}
}

func CreateEntryPointFile(projectPath, entryPoint, language string) error {
	fullPath := filepath.Join(projectPath, entryPoint)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	template, err := getTemplate(language)
	if err != nil {
		return err
	}

	_, err = file.WriteString(template)
	return err
}

func getTemplate(language string) (string, error) {
	templateFile := fmt.Sprintf("templates/%s.tmpl", language)
	templateContent, err := templatesFS.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("error reading template file: %v", err)
	}
	return string(templateContent), nil
}
