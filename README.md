```markdown
# OmniServe

OminiServe is a multi-language serverless platform CLI tool that helps developers initialize, build, and manage serverless projects across various programming languages.

## Features

- Initialize new serverless projects
- Support for multiple programming languages (Go, C, Python, JavaScript)
- Customizable entry points
- JSON-based project configuration

## Installation

### Using Go Install

If you have Go installed on your system, you can install OminServe directly from the source:
```

```bash
go install github.com/emuthianimbithi/OmniServe/cmd/ominiserve@latest
```

Make sure your Go bin directory is in your PATH.

### Pre-built Binaries

(Note: Add links to pre-built binaries when available)

## Usage

After installation, you can use OminServe from anywhere in your terminal:

```bash
ominserve init --name myproject --language go
```

### Commands

Currently, OminServe supports the following command:

- `init`: Initialize a new serverless project

#### Init Command Options

- `--name, -n`: Name of the project (required)
- `--language, -l`: Programming language (go, c, python, javascript) (required)
- `--entry-point, -e`: Path to the entry point file (optional)

## Project Structure

When you initialize a new project, OminServe creates the following structure:

```
myproject/
├── ominserve.json
└── [entry-point file]
```

- `ominserve.json`: Contains project configuration
- Entry-point file: The main file for your serverless function (e.g., `main.go` for Go projects)

## Development

To contribute to OminServe, follow these steps:

1. Clone the repository
2. Make your changes
3. Run tests (when implemented)
4. Submit a pull request

## Future Enhancements

- Implement `build` command for compiling projects
- Add `deploy` functionality to upload projects to serverless platforms
- Introduce `run` command for local testing
- Implement unit tests
- Add support for more languages and serverless platforms

