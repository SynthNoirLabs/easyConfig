package config

import (
	"fmt"
	"os"
	"path/filepath"

	"easyConfig/pkg/util/paths"
)

// --- Codex Provider ---

type CodexProvider struct{}

func (p *CodexProvider) Name() string {
	return "Codex CLI"
}

func (p *CodexProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "# Codex Config\n"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".codex", "config.toml")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".codex", "config.toml")
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

func (p *CodexProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global Config
	if home != "" {
		path := filepath.Join(home, ".codex", "config.toml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "config.toml",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatTOML,
				Exists:   true,
			})
		}
	}

	// 2. Project Config
	if projectPath != "" {
		path := filepath.Join(projectPath, ".codex", "config.toml")
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
