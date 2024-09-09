# OmniServe: Multi-Language Serverless CLI Tool

OmniServe is a powerful, flexible, and user-friendly command-line interface (CLI) tool designed to streamline the development of serverless applications across multiple programming languages. It provides a unified workflow for initializing, building, and managing serverless projects, making it easier for developers to work with various cloud platforms and programming languages.

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Getting Started](#getting-started)
4. [Command Reference](#command-reference)
5. [Configuration](#configuration)
6. [Templates](#templates)
7. [Supported Languages](#supported-languages)
8. [Extending OmniServe](#extending-omniserve)
9. [Troubleshooting](#troubleshooting)
10. [Contributing](#contributing)

## Features

OmniServe offers a rich set of features to enhance your serverless development experience:

- **Multi-Language Support**: Initialize and manage projects in various programming languages, including Go, Python, JavaScript, C, Ruby, and more.
- **Custom Templates**: Create, manage, and use custom project templates for each supported language, allowing you to standardize project structures across your team or organization.
- **Flexible Configuration**: Easily configure project defaults, CLI behavior, and language-specific settings through a YAML configuration file.
- **Extensible Language Support**: Add support for new programming languages without modifying the CLI's source code, simply by updating the configuration and providing appropriate templates.
- **Project Initialization**: Quickly set up new serverless projects with proper structure, boilerplate code, and configuration files.
- **Template Management**: Add, list, and use custom templates for project initialization, enabling you to tailor the initial project setup to your specific needs.
- **Verbose Mode**: Get detailed output about CLI operations for better debugging and understanding of the tool's processes.
- **Cross-Platform Compatibility**: OmniServe works seamlessly on Windows, macOS, and Linux.

## Installation

### Prerequisites

- Go 1.23 or later (for installation from source)
- Git (optional, for version control integration)

### Using Go Install

If you have Go installed on your system, you can install OmniServe directly from the source:

```bash
go install github.com/emuthianimbithi/OmniServe/cmd/omniserve@latest
```

Make sure your Go bin directory is in your PATH.

### Pre-built Binaries

1. Navigate to the [Releases page](https://github.com/emuthianimbithi/OmniServe/releases) of the OmniServe repository.
2. Download the appropriate binary for your operating system and architecture.
3. Rename the binary to `omniserve` (or `omniserve.exe` on Windows).
4. Move the binary to a directory in your system's PATH.

#### Linux and macOS

```bash
chmod +x ./omniserve
sudo mv ./omniserve /usr/local/bin/omniserve
```

#### Windows

Move the `omniserve.exe` file to a directory in your PATH, or add the directory containing the executable to your PATH environment variable.

### Verifying the Installation

To verify that OmniServe is installed correctly, open a new terminal window and run:

```bash
omniserve --version
```

This should display the version number of OmniServe.

## Getting Started

### Initializing Your First Project

To create a new serverless project, use the `init` command:

```bash
omniserve init --name myproject --language go
```

This command will:
1. Create a new directory named `myproject`.
2. Generate a project structure based on the Go template.
3. Create a `omniserve.json` file with project configuration.

### Using Custom Templates

OmniServe automatically uses custom templates if they exist. When you initialize a project, it first looks for a custom template for the specified language. If found, it uses that template; otherwise, it falls back to the built-in template.

```bash
omniserve init --name myproject --language go 
```

Make sure you've added the custom template first (see [Templates](#templates) section).

### Building Your Project

(Note: Implementation of the build command in future versions)

To build your serverless project:

```bash
omniserve build
```

This command will compile your code and prepare it for deployment.

### Deploying Your Project

(Note: Implementation of the deployment command in future versions)

To deploy your serverless project:

```bash
omniserve deploy
```

This command will package your project and deploy it to the configured serverless platform.

## Command Reference

### Global Flags

- `--verbose, -v`: Enable verbose output
- `--config`: Specify a custom config file (default is $HOME/.omniserve.yaml)

### `omniserve init`

Initialize a new serverless project.

Flags:
- `--name, -n`: Name of the project (required)
- `--language, -l`: Programming language (go, python, javascript, c, ruby, etc.) (required)
- `--template, -t`: Name of the custom template to use (optional)
- `--entry-point, -e`: Path to the entry point file (optional)
- `--version`: Initial version of the project (default: "0.1.0")
- `--author, -a`: Author of the project
- `--description, -d`: Short description of the project
- `--license`: License for the project (default: EMM)
- `--git-init, -g`: Initialize a git repository
- `--dockerize, -D`: Add a Dockerfile to the project

### `omniserve template`

Manage project templates.

Subcommands:
- `add`: Add a custom template
- `list`: List all custom templates

#### `omniserve template add`

Add a custom template for a language.

Usage:
```bash
omniserve template add [language] [file]
```

#### `omniserve template list`

List all custom templates.

Usage:
```bash
omniserve template list
```

## Configuration

OmniServe uses a YAML configuration file to manage default settings and language-specific configurations. The default location for this file is `$HOME/.omniserve.yaml`.

### Sample Configuration

```yaml
defaults:
  language: go
  license: EMM
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
  ruby:
    entry_point: app.rb
    build_command: bundle install

cli:
  verbose: false
  color_output: true
```

### Configuration Options

- `defaults`: Set default values for project initialization
- `paths`: Configure paths for templates and other resources
- `languages`: Define supported languages and their configurations
- `cli`: Set CLI behavior options

## Templates

Templates are used to generate the initial structure and code for new projects. OmniServe supports both built-in and custom templates.

### Template Location

By default, custom templates are stored in `~/.omniserve/templates/`. You can change this location in the configuration file.

### Custom Templates

Custom templates should be text files, typically with the extension `.tmpl`. These templates are used as-is when initializing a new project. The content of the template will be copied directly into the new project's entry point file.

Example Go template (`go.tmpl`):

```go
package main

import "fmt"

func main() {
    fmt.Println("Welcome to Omniserve!")
}
```

## Supported Languages

OmniServe comes with built-in support for:

- Go
- Python
- JavaScript
- C

You can add support for additional languages by updating your configuration file and providing appropriate templates.

## Extending OmniServe

### Adding Support for a New Language

1. Update your `~/.omniserve.yaml` configuration file:

```yaml
languages:
  newlang:
    entry_point: main.newlang
    build_command: newlang build
```

2. Create a template file for the new language:

```bash
omniserve template add newlang path/to/newlang_template.txt
```

3. Use the new language when initializing projects:

```bash
omniserve init --name mynewproject --language newlang
```

## Troubleshooting

### Common Issues and Solutions

1. **Template Not Found**
   - Ensure the template exists in the configured templates directory
   - Check the spelling of the language name
   - Verify that you have read permissions for the template file

2. **Configuration File Not Loaded**
   - Check if `~/.omniserve.yaml` exists
   - Use the `--config` flag to specify a custom configuration file location
   - Ensure the YAML syntax in your configuration file is correct

3. **Permission Denied Errors**
   - Check file and directory permissions
   - Ensure you have write access to the project directory
   - Run the command with elevated privileges if necessary (use with caution)

### Debugging

Use the `--verbose` flag to get more detailed output:

```bash
omniserve --verbose init --name myproject --language go
```

This will provide additional information about each step of the process, which can be helpful in identifying issues.

## Contributing

We welcome contributions to OmniServe! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request