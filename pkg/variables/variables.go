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

// Version of the application
const Version = "0.1.0"

const DefaultConfig = `# OmniServe CLI Configuration

defaults:
  language: go
  license: MIT
  version: 0.1.0
  author: Your Name
  git_init: true
  dockerize: false

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
`
