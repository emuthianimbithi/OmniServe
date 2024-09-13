package utils

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"github.com/sabhiram/go-gitignore"
	"os"
	"path/filepath"
)

var Verbose bool

var ignoreList *ignore.GitIgnore

func InitIgnoreList() error {
	var err error
	ignoreList, err = ignore.CompileIgnoreFile(".omniserve-ignore")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func ShouldIgnore(path string) bool {
	if ignoreList == nil {
		return false
	}
	return ignoreList.MatchesPath(path)
}

func IsValidLanguage(lang string) bool {
	if langConfig, exists := cliconfig.CliConfig.Languages[lang]; exists {
		return langConfig.EntryPoint != ""
	}
	return variables.DefaultSupportedLanguages[lang]
}

func GetDefaultEntryPoint(language string) string {
	if langConfig, exists := cliconfig.CliConfig.Languages[language]; exists {
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
