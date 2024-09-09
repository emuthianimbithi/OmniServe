package models

type ProjectConfig struct {
	Name         string   `json:"name"`
	Language     string   `json:"language"`
	Version      string   `json:"version"`
	EntryPoint   string   `json:"entryPoint"`
	Dependencies []string `json:"dependencies,omitempty"`
}
