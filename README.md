```markdown
# OmniServe

OmniServe is a powerful multi-language serverless platform CLI tool that helps developers initialize, build, and manage serverless projects across various programming languages.

## Features

- Initialize new serverless projects
- Support for multiple programming languages (Go, C, Python, JavaScript)
- Customizable entry points
- JSON-based project configuration
- CLI configuration management
- Git and Docker integration options

## Installation

### Prerequisites

- Go 1.23 or later

### Using Go Install

If you have Go installed on your system, you can install OmniServe directly from the source:
```
```bash
go install github.com/emuthianimbithi/OmniServe/cmd/omniserve@latest
```

Make sure your Go bin directory is in your PATH.

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/emuthianimbithi/OmniServe.git
   cd OmniServe
   ```

2. Build the project:
   ```bash
   make build
   ```

3. (Optional) Install globally:
   ```bash
   make install
   ```

## Usage

After installation, you can use OmniServe from anywhere in your terminal:

```bash
omniserve [command] [flags]
```

### Available Commands

- `init`: Initialize a new serverless project
- `config`: Manage OmniServe configuration
- `version`: Print the version number of OmniServe
- `info`: Print information about OmniServe

### Global Flags

- `--verbose, -v`: Enable verbose output
- `--config`: Specify a custom config file (default is $HOME/.omniserve.yaml)

### Init Command

Initialize a new serverless project:

```bash
omniserve init --name myproject --language go
```

#### Init Command Options

- `--name, -n`: Name of the project (required)
- `--language, -l`: Programming language (go, c, python, javascript) (required)
- `--entry-point, -e`: Path to the entry point file (optional)
- `--version`: Initial version of the project (default: "0.1.0")
- `--author, -a`: Author of the project
- `--description, -d`: Short description of the project
- `--license`: License for the project (default: MIT)
- `--git-init, -g`: Initialize a git repository
- `--dockerize, -D`: Add a Dockerfile to the project

### Config Command

Manage OmniServe configuration:

```bash
omniserve config init  # Initialize default configuration file
omniserve config delete  # Delete the configuration file
```

## Project Structure

When you initialize a new project, OmniServe creates the following structure:

```
myproject/
├── omniserve.json
└── [entry-point file]
```

- `omniserve.json`: Contains project configuration
- Entry-point file: The main file for your serverless function (e.g., `main.go` for Go projects)

## Configuration

OmniServe can be configured using a YAML file. By default, it looks for `~/.omniserve.yaml`.

To generate a default configuration file:

```bash
omniserve config init
```

This will create a configuration file with default values. You can then edit this file to customize OmniServe's behavior.

If no configuration file is found, OmniServe will use built-in default values.

## Development

To contribute to OmniServe, follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Cleaning up

To remove the built binary:

```bash
make clean
```

## Future Enhancements

- Implement `build` command for compiling projects
- Add `deploy` functionality to upload projects to serverless platforms
- Introduce `run` command for local testing
- Add support for more languages and serverless platforms

## License

[Add your chosen license here]

## Contact

Emmanuel Muthiani Mbithi - [Your email or contact information]

Project Link: [https://github.com/emuthianimbithi/OmniServe](https://github.com/emuthianimbithi/OmniServe)

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
