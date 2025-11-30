package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- OpenAI Provider ---

type OpenAIProvider struct{}

func (p *OpenAIProvider) Name() string {
	return "OpenAI"
}

func (p *OpenAIProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "version: 1\n"
	var path string

	switch scope {
	case ScopeGlobal:
		configDir := paths.GetConfigDir("openai")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "config.yaml")
	default:
		return "", fmt.Errorf("unsupported scope")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("failed to create dir: %w", err)
	}
	if FileExists(path) {
		return "", fmt.Errorf("file exists: %s", path)
	}
	if err := os.WriteFile(path, []byte(defaultContent), 0o600); err != nil {
		return "", fmt.Errorf("failed to write: %w", err)
	}
	return path, nil
}

func (p *OpenAIProvider) Discover(_ string) ([]Item, error) {
	var items []Item

	configDir := paths.GetConfigDir("openai")
	if configDir != "" {
		path := filepath.Join(configDir, "config.yaml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "config.yaml",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
	}
	return items, nil
}

func (p *OpenAIProvider) CheckStatus() ProviderStatus {
	status := ProviderStatus{
		ProviderName: p.Name(),
		LastChecked:  time.Now().Format(time.RFC3339),
	}

	configDir := paths.GetConfigDir("openai")
	if configDir == "" {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Config directory not found."
		return status
	}

	configPath := filepath.Join(configDir, "config.yaml")
	files, _ := p.Discover("")
	status.DiscoveredFiles = files

	if !FileExists(configPath) {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Global config not found. Create one to get started."
	} else {
		status.Health = StatusHealthy
		status.StatusMessage = "Configuration file found. (Authentication not yet verified)."
	}

	return status
}
