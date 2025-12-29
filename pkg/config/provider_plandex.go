package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// --- Plandex Provider ---

type PlandexProvider struct{}

func (p *PlandexProvider) Name() string {
	return "Plandex"
}

func (p *PlandexProvider) Create(scope Scope, projectPath string) (string, error) {
	if scope != ScopeProject {
		return "", fmt.Errorf("Plandex only supports project-level config (.plandex/)")
	}
	if projectPath == "" {
		return "", fmt.Errorf("project path is required")
	}
	dir := filepath.Join(projectPath, ".plandex")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create dir: %w", err)
	}
	path := filepath.Join(dir, "config.json")
	if FileExists(path) {
		return "", fmt.Errorf("file exists: %s", path)
	}
	if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
		return "", fmt.Errorf("failed to write: %w", err)
	}
	return path, nil
}

func (p *PlandexProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	if projectPath != "" {
		plandexDir := filepath.Join(projectPath, ".plandex")
		if info, err := os.Stat(plandexDir); err == nil && info.IsDir() {
			files, _ := os.ReadDir(plandexDir)
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				format := FormatTXT
				if filepath.Ext(f.Name()) == ".json" {
					format = FormatJSON
				}
				items = append(items, Item{
					Provider: p.Name(),
					Name:     "Project File: " + f.Name(),
					FileName: f.Name(),
					Path:     filepath.Join(plandexDir, f.Name()),
					Scope:    ScopeProject,
					Format:   format,
					Exists:   true,
				})
			}
		}
	}
	return items, nil
}

func (p *PlandexProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Plandex status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}

func (p *PlandexProvider) BinaryName() string {
	return "plandex"
}

func (p *PlandexProvider) VersionArgs() []string {
	return []string{"--version"}
}
