package stagedfiles

import (
	"encoding/json"
	"os"
)

const stagedFilesPath = ".omniserve_staged"

func SaveStagedFiles(files []string) error {
	data, err := json.Marshal(files)
	if err != nil {
		return err
	}
	return os.WriteFile(stagedFilesPath, data, 0644)
}

func LoadStagedFiles() ([]string, error) {
	data, err := os.ReadFile(stagedFilesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	var files []string
	err = json.Unmarshal(data, &files)
	return files, err
}

func ClearStagedFiles() error {
	return os.Remove(stagedFilesPath)
}
