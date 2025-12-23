package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

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

func (p *CodexProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()
	seen := map[string]bool{}
	add := func(it Item) {
		if !seen[it.Path] {
			items = append(items, it)
			seen[it.Path] = true
		}
	}

	// 1. Global Config
	if home != "" {
		path := filepath.Join(home, ".codex", "config.toml")
		if FileExists(path) {
			add(Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "config.toml",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatTOML,
				Exists:   true,
			})
		}

		// Some versions emit config.json
		jsonPath := filepath.Join(home, ".codex", "config.json")
		if FileExists(jsonPath) {
			add(Item{
				Provider: p.Name(),
				Name:     "Global Config (JSON)",
				FileName: "config.json",
				Path:     jsonPath,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		jsoncPath := filepath.Join(home, ".codex", "config.jsonc")
		if FileExists(jsoncPath) {
			add(Item{
				Provider: p.Name(),
				Name:     "Global Config (JSONC)",
				FileName: "config.jsonc",
				Path:     jsoncPath,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// Managed config (per docs)
		managedHome := filepath.Join(home, ".codex", "managed_config.toml")
		if FileExists(managedHome) {
			add(Item{
				Provider: p.Name(),
				Name:     "Managed Config",
				FileName: "managed_config.toml",
				Path:     managedHome,
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
			add(Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "config.toml",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatTOML,
				Exists:   true,
			})
		}

		jsonPath := filepath.Join(projectPath, ".codex", "config.json")
		if FileExists(jsonPath) {
			add(Item{
				Provider: p.Name(),
				Name:     "Project Config (JSON)",
				FileName: "config.json",
				Path:     jsonPath,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		jsoncPath := filepath.Join(projectPath, ".codex", "config.jsonc")
		if FileExists(jsoncPath) {
			add(Item{
				Provider: p.Name(),
				Name:     "Project Config (JSONC)",
				FileName: "config.jsonc",
				Path:     jsoncPath,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 3. System Managed Config
	sysManaged := "/etc/codex/managed_config.toml"
	const goosWindows = "windows"
	if runtime.GOOS == goosWindows {
		if home != "" {
			sysManaged = filepath.Join(home, ".codex", "managed_config.toml")
		}
	}
	if FileExists(sysManaged) {
		add(Item{
			Provider: p.Name(),
			Name:     "System Managed Config",
			FileName: "managed_config.toml",
			Path:     sysManaged,
			Scope:    ScopeSystem,
			Format:   FormatTOML,
			Exists:   true,
		})
	}

	return items, nil
}

func (p *CodexProvider) CheckStatus() ProviderStatus {
	const (
		msgConfigMissing = "Global config not found. Create one to get started."
		msgConfigOK      = "Configuration file found. (Authentication not yet verified)."
	)

	status := ProviderStatus{
		ProviderName: p.Name(),
		LastChecked:  time.Now().Format(time.RFC3339),
	}

	home := paths.GetHomeDir()
	if home == "" {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Home directory not found."
		return status
	}

	configPath := filepath.Join(home, ".codex", "config.toml")
	files, _ := p.Discover("")
	status.DiscoveredFiles = files

	if !FileExists(configPath) {
		status.Health = StatusUnhealthy
		status.StatusMessage = msgConfigMissing
	} else {
		status.Health = StatusHealthy
		status.StatusMessage = msgConfigOK
	}

	return status
}
