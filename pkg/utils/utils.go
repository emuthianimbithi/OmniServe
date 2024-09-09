package utils

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"os"
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
