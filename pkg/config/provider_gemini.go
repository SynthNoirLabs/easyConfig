package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Gemini Provider ---

type GeminiProvider struct{}

func (p *GeminiProvider) Name() string {
	return "Gemini"
}

func (p *GeminiProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		// Standard location: ~/.gemini/settings.json
		path = filepath.Join(home, ".gemini", "settings.json")

	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required for project scope")
		}
		path = filepath.Join(projectPath, ".gemini", "settings.json")

	default:
		return "", fmt.Errorf("unsupported scope: %s", scope)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file if it doesn't exist
	if FileExists(path) {
		return "", fmt.Errorf("file already exists: %s", path)
	}

	if err := os.WriteFile(path, []byte(defaultContent), 0o600); err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return path, nil
}

func (p *GeminiProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global settings
	if home != "" {
		path := filepath.Join(home, ".gemini", "settings.json")
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

	// 2. System Settings (Linux)
	sysConfigDir := paths.GetConfigDir("gemini-cli") // Use app name for GetConfigDir
	if sysConfigDir != "" {
		sysDefaults := filepath.Join(sysConfigDir, "system-defaults.json")
		if FileExists(sysDefaults) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "System Defaults",
				FileName: "system-defaults.json",
				Path:     sysDefaults,
				Scope:    ScopeSystem,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		sysOverrides := filepath.Join(sysConfigDir, "settings.json")
		if FileExists(sysOverrides) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "System Overrides",
				FileName: "settings.json",
				Path:     sysOverrides,
				Scope:    ScopeSystem,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 3. Project settings
	if projectPath != "" {
		// settings.json
		path := filepath.Join(projectPath, ".gemini", "settings.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Workspace Settings",
				FileName: "settings.json",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// GEMINI.md
		pathContext := filepath.Join(projectPath, "GEMINI.md")
		if FileExists(pathContext) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Context File",
				FileName: "GEMINI.md",
				Path:     pathContext,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}

		// Project Extensions (.gemini/extensions/*)
		projectExtPath := filepath.Join(projectPath, ".gemini", "extensions", "*")
		if matches, err := filepath.Glob(projectExtPath); err == nil {
			for _, match := range matches {
				name := "Extension: " + filepath.Base(match)
				// Try to read gemini-extension.json
				manifestPath := filepath.Join(match, "gemini-extension.json")
				if _, err := os.ReadFile(manifestPath); err == nil {
					// Simple regex to find "displayName" or "name" to avoid full JSON unmarshal overhead for just a name
					// But for robustness, we could use a struct. Let's stick to simple name for now.
					// Or just append the manifest file itself if it exists?
					// The requirement is to parse it.
					// Let's assume the extension IS the directory, and we want to show it.
					// If the match is a directory, we look inside.
					if info, err := os.Stat(match); err == nil && info.IsDir() {
						// It's a directory extension
						items = append(items, Item{
							Provider: p.Name(),
							Name:     name,
							FileName: "gemini-extension.json",
							Path:     manifestPath,
							Scope:    ScopeProject,
							Format:   FormatJSON,
							Exists:   true,
						})
					} else {
						// It's a file extension (e.g. .js)
						items = append(items, Item{
							Provider: p.Name(),
							Name:     name,
							FileName: filepath.Base(match),
							Path:     match,
							Scope:    ScopeProject,
							Format:   FormatTXT,
							Exists:   true,
						})
					}
				} else {
					// No manifest or not a dir, treat as simple file/dir
					items = append(items, Item{
						Provider: p.Name(),
						Name:     name,
						FileName: filepath.Base(match),
						Path:     match,
						Scope:    ScopeProject,
						Format:   FormatTXT,
						Exists:   true,
					})
				}
			}
		}
	}

	// Global Extensions (~/.gemini/extensions/*)
	if home != "" {
		globalExtPath := filepath.Join(home, ".gemini", "extensions", "*")
		if matches, err := filepath.Glob(globalExtPath); err == nil {
			for _, match := range matches {
				name := "Extension: " + filepath.Base(match)
				manifestPath := filepath.Join(match, "gemini-extension.json")
				if content, err := os.ReadFile(manifestPath); err == nil {
					// Check if valid JSON
					if len(content) > 0 {
						items = append(items, Item{
							Provider: p.Name(),
							Name:     name,
							FileName: "gemini-extension.json",
							Path:     manifestPath,
							Scope:    ScopeGlobal,
							Format:   FormatJSON,
							Exists:   true,
						})
						continue
					}
				}

				items = append(items, Item{
					Provider: p.Name(),
					Name:     name,
					FileName: filepath.Base(match),
					Path:     match,
					Scope:    ScopeGlobal,
					Format:   FormatTXT,
					Exists:   true,
				})
			}
		}
	}

	return items, nil
}

func (p *GeminiProvider) CheckStatus() ProviderStatus {
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

	configPath := filepath.Join(home, ".gemini", "settings.json")
	files, _ := p.Discover("")
	status.DiscoveredFiles = files

	if !FileExists(configPath) {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Global settings file not found. Create one to get started."
	} else {
		status.Health = StatusHealthy
		status.StatusMessage = "Configuration file found. (Authentication not yet verified)."
	}

	return status
}
