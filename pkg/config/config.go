package config

import (
	"encoding/json"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/models"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

var ProjectConfig models.ProjectConfig

func NewProjectConfig(projectPath, name, language, entryPoint string) (*models.ProjectConfig, error) {
	utils.VerboseLog("Starting NewProjectConfig function")
	utils.VerboseLog(fmt.Sprintf("Input parameters: projectPath=%s, name=%s, language=%s, entryPoint=%s", projectPath, name, language, entryPoint))

	if entryPoint == "" {
		utils.VerboseLog("Entry point not provided, getting default for language")
		entryPoint = utils.GetDefaultEntryPoint(language)
		utils.VerboseLog(fmt.Sprintf("Default entry point: %s", entryPoint))
	}

	config := &models.ProjectConfig{
		Name:       name,
		Code:       uuid.New().String(),
		Language:   language,
		Version:    "0.1.0",
		EntryPoint: entryPoint,
	}
	utils.VerboseLog(fmt.Sprintf("Created project config: %+v", config))

	configFilePath := filepath.Join(projectPath, "omniserve.json")
	utils.VerboseLog(fmt.Sprintf("Creating config file at: %s", configFilePath))
	configFile, err := os.Create(configFilePath)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error creating config file: %v", err))
		return nil, fmt.Errorf("error creating config file: %v", err)
	}
	defer configFile.Close()

	utils.VerboseLog("Encoding config to JSON")
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		utils.VerboseLog(fmt.Sprintf("Error writing config file: %v", err))
		return nil, fmt.Errorf("error writing config file: %v", err)
	}

	utils.VerboseLog("Project config created and written successfully")
	return config, nil
}

func LoadProjectConfig(projectPath string) error {
	utils.VerboseLog("Starting LoadProjectConfig function")
	utils.VerboseLog(fmt.Sprintf("Project path: %s", projectPath))

	configFile := filepath.Join(projectPath, "omniserve.json")
	utils.VerboseLog(fmt.Sprintf("Config file path: %s", configFile))

	// Check if the config file exists
	utils.VerboseLog("Checking if config file exists")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		utils.VerboseLog(fmt.Sprintf("Project config file not found: %s", configFile))
		return fmt.Errorf("project config file not found: %s", configFile)
	}

	// Open the config file
	utils.VerboseLog("Opening config file")
	file, err := os.Open(configFile)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error opening config file: %v", err))
		return fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	// Decode the JSON into a ProjectConfig struct
	utils.VerboseLog("Decoding JSON into ProjectConfig struct")
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ProjectConfig); err != nil {
		utils.VerboseLog(fmt.Sprintf("Error decoding config file: %v", err))
		return fmt.Errorf("error decoding config file: %v", err)
	}

	utils.VerboseLog("Project config loaded successfully")
	utils.VerboseLog(fmt.Sprintf("Loaded config: %+v", ProjectConfig))
	return nil
}
