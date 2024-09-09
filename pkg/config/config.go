package config

import (
	"encoding/json"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/models"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"os"
	"path/filepath"
)

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
