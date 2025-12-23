package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Aider Provider ---

type AiderProvider struct{}

func (p *AiderProvider) Name() string {
	return "Aider"
}

func (p *AiderProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "model: gpt-4\n"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".aider.conf.yml")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".aider.conf.yml")
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

func (p *AiderProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global Config
	if home != "" {
		path := filepath.Join(home, ".aider.conf.yml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: ".aider.conf.yml",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
	}

	// 2. Project Config
	if projectPath != "" {
		path := filepath.Join(projectPath, ".aider.conf.yml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: ".aider.conf.yml",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
	}
	return items, nil
}

func (p *AiderProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Aider status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
