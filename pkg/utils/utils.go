package utils

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"os"
	"path/filepath"
)

var Verbose bool

func IsValidLanguage(lang string) bool {
	if langConfig, exists := cliconfig.Config.Languages[lang]; exists {
		return langConfig.EntryPoint != ""
	}
	return variables.DefaultSupportedLanguages[lang]
}

func GetDefaultEntryPoint(language string) string {
	if langConfig, exists := cliconfig.Config.Languages[language]; exists {
		return langConfig.EntryPoint
	}
	return variables.DefaultEntryPointTemplate[language]
}

// VerboseLog used to give more detailed output
func VerboseLog(message string) {
	if Verbose {
		fmt.Fprintln(os.Stderr, "VERBOSE:", message)
	}
}
func GetAllFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func GetFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{path}, nil
	}
	return GetAllFiles(path)
}
