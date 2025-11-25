package config

import (
	"fmt"
	"os"
	"path/filepath"

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
		// Goose uses ~/.config/goose/config.yaml on Linux/macOS
		// On Windows it uses AppData/Block/goose/config/config.yaml
		// paths.GetConfigDir("goose") returns ~/.config/goose or AppData/goose
		// We might need special handling for Windows if "Block" is strictly required.
		// For now, we use standard XDG/AppData structure.
		configDir := paths.GetConfigDir("goose")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "config.yaml")
	default:
		return "", fmt.Errorf("unsupported scope (Goose only supports global config)")
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

func (p *GooseProvider) Discover(_ string) ([]Item, error) {
	var items []Item

	// 1. Global Config
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
	return items, nil
}
