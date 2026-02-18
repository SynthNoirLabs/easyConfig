//go:build ignore
// +build ignore

// This file is an example template and is excluded from the build.
// Copy it to pkg/config/ and remove the build constraint to use it.

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// MyToolProvider is a template for creating a new provider.
//
// To use this template, follow these steps:
//
// 1. Rename the file to `provider_<your_tool_name>.go`.
// 2. Replace all instances of `MyTool` with the name of your tool.
// 3. Implement the `Discover`, `Create`, and `CheckStatus` methods.
// 4. Register your provider in `pkg/config/service.go`.
type MyToolProvider struct{}

// Name returns the name of the provider.
func (p *MyToolProvider) Name() string {
	return "MyTool"
}

// Discover finds all the configuration files for the provider.
func (p *MyToolProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// Global scope
	if home != "" {
		globalPath := filepath.Join(home, ".mytool", "config.json")
		// FileExists is a helper function in the config package
		if FileExists(globalPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "config.json",
				Path:     globalPath,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// Project scope
	if projectPath != "" {
		projectPath := filepath.Join(projectPath, ".mytool.json")
		if FileExists(projectPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: ".mytool.json",
				Path:     projectPath,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	return items, nil
}

// Create creates a new configuration file for the provider.
func (p *MyToolProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".mytool", "config.json")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".mytool.json")
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

// CheckStatus performs a health check on the provider's configuration.
func (p *MyToolProvider) CheckStatus() ProviderStatus {
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

	globalPath := filepath.Join(home, ".mytool", "config.json")
	files, _ := p.Discover("")
	status.DiscoveredFiles = files

	if !FileExists(globalPath) {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Global config not found. Use 'Create' to set one up."
	} else {
		status.Health = StatusHealthy
		status.StatusMessage = "Configuration files found."
	}

	return status
}
