package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Windsurf Provider ---

type WindsurfProvider struct{}

func (p *WindsurfProvider) Name() string {
	return "Windsurf"
}

func (p *WindsurfProvider) Create(scope Scope, _ string) (string, error) {
	// Windsurf primarily relies on user settings in config dirs
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		// Attempt to resolve the correct User settings path
		baseDir := paths.GetConfigDir("Windsurf")
		if baseDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		// VS Code style paths: User/settings.json
		path = filepath.Join(baseDir, "User", "settings.json")

		// Adjust for Linux/macOS standard paths if GetConfigDir isn't sufficient
		// GetConfigDir often points to ~/.config/Windsurf (Linux) or ~/Library/Application Support/Windsurf (Mac)
		// which is correct for the base.
	default:
		return "", fmt.Errorf("unsupported scope (Windsurf Project settings are usually .vscode/settings.json or equivalent)")
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

func (p *WindsurfProvider) Discover(_ string) ([]Item, error) {
	var items []Item

	// 1. Global/User Settings
	// Linux: ~/.config/Windsurf/User/settings.json
	// macOS: ~/Library/Application Support/Windsurf/User/settings.json
	// Windows: %APPDATA%\Windsurf\User\settings.json

	// paths.GetConfigDir("Windsurf") typically returns the app data root
	configDir := paths.GetConfigDir("Windsurf")
	if configDir != "" {
		path := filepath.Join(configDir, "User", "settings.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "User Settings",
				FileName: "settings.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// Fallback check for Linux if XDG is weird or GetConfigDir differs
	if runtime.GOOS == "linux" {
		home := paths.GetHomeDir()
		if home != "" {
			path := filepath.Join(home, ".config", "Windsurf", "User", "settings.json")
			if FileExists(path) {
				// Avoid duplicate if GetConfigDir already found it
				found := false
				for _, i := range items {
					if i.Path == path {
						found = true
						break
					}
				}
				if !found {
					items = append(items, Item{
						Provider: p.Name(),
						Name:     "User Settings",
						FileName: "settings.json",
						Path:     path,
						Scope:    ScopeGlobal,
						Format:   FormatJSON,
						Exists:   true,
					})
				}
			}
		}
	}

	return items, nil
}

func (p *WindsurfProvider) BinaryName() string {
	return "windsurf"
}

func (p *WindsurfProvider) VersionArgs() []string {
	return []string{"--version"}
}

func (p *WindsurfProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Windsurf status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
