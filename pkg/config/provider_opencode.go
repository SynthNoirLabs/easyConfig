package config

import (
	"fmt"
	"os"
	"path/filepath"

	"easyConfig/pkg/util/paths"
)

// --- OpenCode Provider ---

type OpenCodeProvider struct{}

func (p *OpenCodeProvider) Name() string {
	return "OpenCode"
}

func (p *OpenCodeProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		configDir := paths.GetConfigDir("opencode")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "opencode.json")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, "opencode.json")
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

func (p *OpenCodeProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item

	// 1. Global Config
	// Linux/macOS: ~/.config/opencode/opencode.json
	configDir := paths.GetConfigDir("opencode")
	if configDir != "" {
		path := filepath.Join(configDir, "opencode.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "opencode.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Config
	if projectPath != "" {
		// opencode.json
		pathProj := filepath.Join(projectPath, "opencode.json")
		if FileExists(pathProj) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "opencode.json",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// opencode.local.json (Secrets)
		pathLocal := filepath.Join(projectPath, "opencode.local.json")
		if FileExists(pathLocal) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Local Secrets",
				FileName: "opencode.local.json",
				Path:     pathLocal,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}
	return items, nil
}
