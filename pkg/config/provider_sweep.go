package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// --- Sweep Provider ---

type SweepProvider struct{}

func (p *SweepProvider) Name() string {
	return "Sweep"
}

func (p *SweepProvider) Create(scope Scope, projectPath string) (string, error) {
	if scope != ScopeProject {
		return "", fmt.Errorf("Sweep only supports project-level config (.sweep.yaml)")
	}
	if projectPath == "" {
		return "", fmt.Errorf("project path is required")
	}
	path := filepath.Join(projectPath, ".sweep.yaml")
	defaultContent := "branch: main\n"

	if FileExists(path) {
		return "", fmt.Errorf("file exists: %s", path)
	}
	if err := os.WriteFile(path, []byte(defaultContent), 0o600); err != nil {
		return "", fmt.Errorf("failed to write: %w", err)
	}
	return path, nil
}

func (p *SweepProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	if projectPath != "" {
		path := filepath.Join(projectPath, ".sweep.yaml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: ".sweep.yaml",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
	}
	return items, nil
}

func (p *SweepProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Sweep status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}

func (p *SweepProvider) BinaryName() string {
	return "sweep"
}

func (p *SweepProvider) VersionArgs() []string {
	return []string{"--version"}
}
