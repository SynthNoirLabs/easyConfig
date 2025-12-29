package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Continue Provider ---

type ContinueProvider struct{}

func (p *ContinueProvider) Name() string {
	return "Continue"
}

func (p *ContinueProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "models: []\n"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".continue", "config.yaml")
	default:
		return "", fmt.Errorf("unsupported scope for Continue")
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

func (p *ContinueProvider) Discover(_ string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	if home != "" {
		// 1. Recommended config.yaml
		pathYaml := filepath.Join(home, ".continue", "config.yaml")
		if FileExists(pathYaml) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config (YAML)",
				FileName: "config.yaml",
				Path:     pathYaml,
				Scope:    ScopeGlobal,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
		// 2. Legacy config.json
		pathJson := filepath.Join(home, ".continue", "config.json")
		if FileExists(pathJson) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config (Legacy)",
				FileName: "config.json",
				Path:     pathJson,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// 3. Typescript config.ts (readonly/text)
		pathTs := filepath.Join(home, ".continue", "config.ts")
		if FileExists(pathTs) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config (TS)",
				FileName: "config.ts",
				Path:     pathTs,
				Scope:    ScopeGlobal,
				Format:   FormatTXT,
				Exists:   true,
			})
		}
	}

	return items, nil
}

func (p *ContinueProvider) BinaryName() string {
	return "continue"
}

func (p *ContinueProvider) VersionArgs() []string {
	return []string{"--version"}
}

func (p *ContinueProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Continue status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
