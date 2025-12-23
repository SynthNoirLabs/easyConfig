package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Jules Provider ---

type JulesProvider struct{}

func (p *JulesProvider) Name() string {
	return "Jules"
}

func (p *JulesProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".jules-mcp", "data.json")
	default:
		return "", fmt.Errorf("unsupported scope")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
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

func (p *JulesProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global Data
	if home != "" {
		path := filepath.Join(home, ".jules-mcp", "data.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "data.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Context (AGENTS.md)
	if projectPath != "" {
		pathProj := filepath.Join(projectPath, "AGENTS.md")
		if FileExists(pathProj) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Agents Context",
				FileName: "AGENTS.md",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}
	}

	return items, nil
}

func (p *JulesProvider) CheckStatus() ProviderStatus {
	const (
		msgHomeMissing   = "Home directory not found."
		msgConfigMissing = "Global config not found. Create one to get started."
		msgConfigOK      = "Configuration file found. (Authentication not yet verified)."
	)

	status := ProviderStatus{
		ProviderName: p.Name(),
		LastChecked:  time.Now().Format(time.RFC3339),
	}

	home := paths.GetHomeDir()
	if home == "" {
		status.Health = StatusUnhealthy
		status.StatusMessage = msgHomeMissing
		return status
	}

	configPath := filepath.Join(home, ".jules-mcp", "data.json")
	files, _ := p.Discover("")
	status.DiscoveredFiles = files

	if !FileExists(configPath) {
		status.Health = StatusUnhealthy
		status.StatusMessage = msgConfigMissing
	} else {
		status.Health = StatusHealthy
		status.StatusMessage = msgConfigOK
	}

	return status
}
