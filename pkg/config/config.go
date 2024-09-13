package config

import (
	"encoding/json"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/models"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"os"
	"path/filepath"
)

var ProjectConfig models.ProjectConfig

func NewProjectConfig(projectPath, name, language, entryPoint string) (*models.ProjectConfig, error) {
	if entryPoint == "" {
		entryPoint = utils.GetDefaultEntryPoint(language)
	}

	config := &models.ProjectConfig{
		Name:       name,
		Language:   language,
		Version:    "0.1.0",
		EntryPoint: entryPoint,
	}

	configFile, err := os.Create(filepath.Join(projectPath, "omniserve.json"))
	if err != nil {
		return nil, fmt.Errorf("error creating config file: %v", err)
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return nil, fmt.Errorf("error writing config file: %v", err)
	}

	return config, nil
}

func LoadProjectConfig(projectPath string) error {
	configFile := filepath.Join(projectPath, "omniserve.json")

	// Check if the config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("project config file not found: %s", configFile)
	}

	// Open the config file
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	// Decode the JSON into a ProjectConfig struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ProjectConfig); err != nil {
		return fmt.Errorf("error decoding config file: %v", err)
	}

	return nil
}
