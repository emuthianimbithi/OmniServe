package variables

// DefaultSupportedLanguages  map of supported languages
var DefaultSupportedLanguages = map[string]bool{
	"go":         true,
	"c":          true,
	"python":     true,
	"javascript": true,
}

// DefaultEntryPointTemplate  map of entry point templates
var DefaultEntryPointTemplate = map[string]string{
	"go":         "main.go",
	"c":          "main.c",
	"python":     "main.py",
	"javascript": "main.js",
}

var Verbose bool

const DefaultConfig = `
defaults:
  language: go
  license: default
  version: 0.1.0
  author: default
  git_init: false
  dockerize: true

paths:
  templates: ~/.omniserve/templates

languages:
  go:
    entry_point: main.go
    build_command: go build
  python:
    entry_point: main.py
    build_command: python -m compileall
  javascript:
    entry_point: index.js
    build_command: npm run build
  c:
    entry_point: main.c
    build_command: gcc -o main main.c

cli:
  verbose: false
  color_output: true

server:
  host: 0.0.0.0:8765
  port: 8765
`

var DefaultIgnorePatterns = []string{
	"# OS generated files",
	".DS_Store",
	"Thumbs.db",
	"",
	"# Editor and IDE files",
	"*.swp",
	"*.git/",
	"*.github/",
	"*.swo",
	"*.idea/",
	"*.vscode/",
	"",
	"# Compiled files",
	"*.class",
	"*.pyc",
	"*.pyo",
	"__pycache__/",
	"",
	"# Dependency directories",
	"node_modules/",
	"vendor/",
	"",
	"# Log files",
	"*.log",
	"",
	"# Build directories",
	"build/",
	"dist/",
	"",
	"",
	"# Package files",
	"*.tar.gz",
	"*.zip",
}
