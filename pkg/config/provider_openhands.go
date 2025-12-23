package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// --- OpenHands Provider ---

type OpenHandsProvider struct{}

func (p *OpenHandsProvider) Name() string {
	return "OpenHands"
}

func (p *OpenHandsProvider) Create(scope Scope, projectPath string) (string, error) {
	if scope != ScopeProject {
		return "", fmt.Errorf("OpenHands primarily uses project-level config.toml")
	}
	if projectPath == "" {
		return "", fmt.Errorf("project path is required")
	}
	path := filepath.Join(projectPath, "config.toml")
	if FileExists(path) {
		return "", fmt.Errorf("file exists: %s", path)
	}
	if err := os.WriteFile(path, []byte("[core]\n"), 0o600); err != nil {
		return "", fmt.Errorf("failed to write: %w", err)
	}
	return path, nil
}

func (p *OpenHandsProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	if projectPath != "" {
		path := filepath.Join(projectPath, "config.toml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "config.toml",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatTOML,
				Exists:   true,
			})
		}
	}
	return items, nil
}

func (p *OpenHandsProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "OpenHands status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
