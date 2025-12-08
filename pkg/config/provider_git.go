package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Git Provider ---

type GitProvider struct{}

func (p *GitProvider) Name() string {
	return "Git"
}

func (p *GitProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "[user]\n\tname = Your Name\n\temail = your.email@example.com\n"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".gitconfig")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".git", "config")
	default:
		return "", fmt.Errorf("unsupported scope")
	}

	// Ensure directory exists (especially for .git/config)
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

func (p *GitProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global Config (~/.gitconfig)
	if home != "" {
		path := filepath.Join(home, ".gitconfig")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: ".gitconfig",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatINI, // Git config is INI-like
				Exists:   true,
			})
		}
	}

	// 2. System Config
	// Windows path differs, usually in Program Files.
	// We only check /etc/gitconfig on non-Windows systems for now.
	if runtime.GOOS != "windows" {
		sysPath := "/etc/gitconfig"
		if FileExists(sysPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "System Config",
				FileName: "gitconfig",
				Path:     sysPath,
				Scope:    ScopeSystem,
				Format:   FormatINI,
				Exists:   true,
			})
		}
	}

	// 3. Project Config (.git/config)
	if projectPath != "" {
		path := filepath.Join(projectPath, ".git", "config")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "config",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatINI,
				Exists:   true,
			})
		}
	}

	return items, nil
}

func (p *GitProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Git status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}

func (p *GitProvider) GetWizard() Wizard {
	return nil
}
