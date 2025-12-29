package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Goose Provider ---

type GooseProvider struct{}

func (p *GooseProvider) Name() string {
	return "Goose"
}

func (p *GooseProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "version: 1\n"
	var path string

	switch scope {
	case ScopeGlobal:
		// Goose default locations:
		// Linux/macOS: ~/.config/goose/config.yaml
		// Windows: %APPDATA%/Block/goose/config/config.yaml
		configDir := paths.GetConfigDir("goose")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "config.yaml")
		// If Windows Block path exists, prefer it
		if block := os.Getenv("APPDATA"); block != "" {
			blockPath := filepath.Join(block, "Block", "goose", "config", "config.yaml")
			if _, err := os.Stat(filepath.Dir(blockPath)); err == nil {
				path = blockPath
			}
		}
	default:
		return "", fmt.Errorf("unsupported scope (Goose only supports global config)")
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

func (p *GooseProvider) Discover(_ string) ([]Item, error) {
	var items []Item

	// 1. Global Config (XDG/AppData)
	configDir := paths.GetConfigDir("goose")
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
	// 2. Windows Block path
	if appData := os.Getenv("APPDATA"); appData != "" {
		blockPath := filepath.Join(appData, "Block", "goose", "config", "config.yaml")
		if FileExists(blockPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Block Config",
				FileName: "config.yaml",
				Path:     blockPath,
				Scope:    ScopeGlobal,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
	}
	return items, nil
}

func (p *GooseProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Goose status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}

func (p *GooseProvider) BinaryName() string {
	return "goose"
}

func (p *GooseProvider) VersionArgs() []string {
	return []string{"--version"}
}
