package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Crush Provider ---

type CrushProvider struct{}

func (p *CrushProvider) Name() string {
	return "Crush CLI"
}

func (p *CrushProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		configDir := paths.GetConfigDir("crush")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "crush.json")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, "crush.json")
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

func (p *CrushProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item

	// 1. Global Config
	configDir := paths.GetConfigDir("crush")
	if configDir != "" {
		// Main Config
		pathMain := filepath.Join(configDir, "crush.json")
		if FileExists(pathMain) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "crush.json",
				Path:     pathMain,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// Providers Config
		pathProviders := filepath.Join(configDir, "providers.json")
		if FileExists(pathProviders) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Providers",
				FileName: "providers.json",
				Path:     pathProviders,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Config
	if projectPath != "" {
		// .crush.json (Hidden)
		pathHidden := filepath.Join(projectPath, ".crush.json")
		if FileExists(pathHidden) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config (Hidden)",
				FileName: ".crush.json",
				Path:     pathHidden,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// crush.json (Visible)
		pathVisible := filepath.Join(projectPath, "crush.json")
		if FileExists(pathVisible) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "crush.json",
				Path:     pathVisible,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// .crushignore
		pathIgnore := filepath.Join(projectPath, ".crushignore")
		if FileExists(pathIgnore) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Ignore File",
				FileName: ".crushignore",
				Path:     pathIgnore,
				Scope:    ScopeProject,
				Format:   FormatTXT,
				Exists:   true,
			})
		}
	}
	return items, nil
}

func (p *CrushProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Crush status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
