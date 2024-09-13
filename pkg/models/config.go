package models

type ProjectConfig struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	Language   string `json:"language"`
	Version    string `json:"version"`
	EntryPoint string `json:"entryPoint"`
	BootLoader string `json:"runTime"`
}
