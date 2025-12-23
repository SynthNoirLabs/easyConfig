package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Mentat Provider ---

type MentatProvider struct{}

func (p *MentatProvider) Name() string {
	return "Mentat"
}

func (p *MentatProvider) Create(scope Scope, projectPath string) (string, error) {
	var path string
	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".mentat", ".env")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".mentat", "README.md")
	default:
		return "", fmt.Errorf("unsupported scope")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("failed to create dir: %w", err)
	}
	if FileExists(path) {
		return "", fmt.Errorf("file exists: %s", path)
	}
	content := "# Mentat Configuration\n"
	if filepath.Base(path) == ".env" {
		content = "OPENAI_API_KEY=\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return "", fmt.Errorf("failed to write: %w", err)
	}
	return path, nil
}

func (p *MentatProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global Config (~/.mentat/.env)
	if home != "" {
		path := filepath.Join(home, ".mentat", ".env")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Env",
				FileName: ".env",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatTXT,
				Exists:   true,
			})
		}
	}

	// 2. Project Config (.mentat/*)
	if projectPath != "" {
		mentatDir := filepath.Join(projectPath, ".mentat")
		if info, err := os.Stat(mentatDir); err == nil && info.IsDir() {
			files, _ := os.ReadDir(mentatDir)
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				format := FormatTXT
				if filepath.Ext(f.Name()) == ".md" {
					format = FormatMD
				}
				items = append(items, Item{
					Provider: p.Name(),
					Name:     "Project File: " + f.Name(),
					FileName: f.Name(),
					Path:     filepath.Join(mentatDir, f.Name()),
					Scope:    ScopeProject,
					Format:   format,
					Exists:   true,
				})
			}
		}
	}
	return items, nil
}

func (p *MentatProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Mentat status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
