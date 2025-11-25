package config

import (
	"fmt"
	"os"
	"path/filepath"

	"easyConfig/pkg/util/paths"
)

// --- Copilot Provider ---

type CopilotProvider struct{}

func (p *CopilotProvider) Name() string {
	return "GitHub Copilot"
}

func (p *CopilotProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".copilot", "mcp-config.json")
	case ScopeProject:
		return "", fmt.Errorf("project creation not supported for Copilot yet")
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

func (p *CopilotProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global CLI Config
	if home != "" {
		path := filepath.Join(home, ".copilot", "mcp-config.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "CLI Config",
				FileName: "mcp-config.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Context
	if projectPath != "" {
		pathProj := filepath.Join(projectPath, ".github", "copilot-instructions.md")
		if FileExists(pathProj) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Instructions",
				FileName: "copilot-instructions.md",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}
	}
	return items, nil
}
