package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Claude Provider ---

type ClaudeProvider struct{}

func (p *ClaudeProvider) Name() string {
	return "Claude Code"
}

func (p *ClaudeProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".claude", "settings.json")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".claude", "settings.json")
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

func (p *ClaudeProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global Desktop Config
	if home != "" {
		path := filepath.Join(home, ".claude", "claude_desktop_config.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Desktop Config",
				FileName: "claude_desktop_config.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// 2. Global CLI Settings
		pathCLI := filepath.Join(home, ".claude", "settings.json")
		if FileExists(pathCLI) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "CLI Settings",
				FileName: "settings.json",
				Path:     pathCLI,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 3. System Settings
	// Linux /etc/claude-code/managed-settings.json
	// macOS /Library/Application Support/ClaudeCode/managed-settings.json
	// Windows C:\ProgramData\ClaudeCode\managed-settings.json
	sysConfigDir := paths.GetConfigDir("ClaudeCode")
	if sysConfigDir != "" {
		sysSettingsPath := filepath.Join(sysConfigDir, "managed-settings.json")
		if FileExists(sysSettingsPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Managed Settings",
				FileName: "managed-settings.json",
				Path:     sysSettingsPath,
				Scope:    ScopeSystem,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		sysMCPPath := filepath.Join(sysConfigDir, "managed-mcp.json")
		if FileExists(sysMCPPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Managed MCP",
				FileName: "managed-mcp.json",
				Path:     sysMCPPath,
				Scope:    ScopeSystem,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 4. Project Settings
	if projectPath != "" {
		// settings.json
		pathProj := filepath.Join(projectPath, ".claude", "settings.json")
		if FileExists(pathProj) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Settings",
				FileName: "settings.json",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// settings.local.json
		pathLocal := filepath.Join(projectPath, ".claude", "settings.local.json")
		if FileExists(pathLocal) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Local Settings",
				FileName: "settings.local.json",
				Path:     pathLocal,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// CLAUDE.md
		pathMemory := filepath.Join(projectPath, "CLAUDE.md")
		if FileExists(pathMemory) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Memory File",
				FileName: "CLAUDE.md",
				Path:     pathMemory,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}

		// Subagents (agents/*.md)
		agentsPath := filepath.Join(projectPath, "agents", "*.md")
		if matches, err := filepath.Glob(agentsPath); err == nil {
			for _, match := range matches {
				items = append(items, Item{
					Provider: p.Name(),
					Name:     "Subagent: " + filepath.Base(match),
					FileName: filepath.Base(match),
					Path:     match,
					Scope:    ScopeProject,
					Format:   FormatMD,
					Exists:   true,
				})
			}
		}

		// Custom Commands (.claude/commands/*.md)
		commandsPath := filepath.Join(projectPath, ".claude", "commands", "*.md")
		if matches, err := filepath.Glob(commandsPath); err == nil {
			for _, match := range matches {
				items = append(items, Item{
					Provider: p.Name(),
					Name:     "Command: " + filepath.Base(match),
					FileName: filepath.Base(match),
					Path:     match,
					Scope:    ScopeProject,
					Format:   FormatMD,
					Exists:   true,
				})
			}
		}
	}

	return items, nil
}

func (p *ClaudeProvider) CheckStatus() ProviderStatus {
	status := ProviderStatus{
		ProviderName: p.Name(),
		LastChecked:  time.Now().Format(time.RFC3339),
	}

	// For Claude, we'll check for the global settings file.
	// A more robust check could validate the contents, check for auth keys, etc.
	home := paths.GetHomeDir()
	if home == "" {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Home directory not found."
		return status
	}

	cliSettingsPath := filepath.Join(home, ".claude", "settings.json")
	files, _ := p.Discover("") // project path is empty for global checks
	status.DiscoveredFiles = files

	if !FileExists(cliSettingsPath) {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Global CLI settings not found. Use 'Create' to set one up."
	} else {
		status.Health = StatusHealthy
		status.StatusMessage = "Configuration files found. (Authentication not yet verified)."
	}

	return status
}
