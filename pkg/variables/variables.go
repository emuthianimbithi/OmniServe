package variables

// SupportedLanguages map of supported languages
var SupportedLanguages = map[string]bool{
	"go":         true,
	"c":          true,
	"python":     true,
	"javascript": true,
}

// EntryPointTemplate map of entry point templates
var EntryPointTemplate = map[string]string{
	"go":         "main.go",
	"c":          "main.c",
	"python":     "main.py",
	"javascript": "main.js",
}
