package config

import (
	"path/filepath"
)

// --- Claude Provider ---

type ClaudeProvider struct{}

func (p *ClaudeProvider) Name() string {
	return "Claude Code"
}

func (p *ClaudeProvider) Discover(projectPath string) ([]ConfigItem, error) {
	var items []ConfigItem
	home := GetUserHome()

	// 1. Global Desktop Config
	if home != "" {
		path := filepath.Join(home, ".claude", "claude_desktop_config.json")
		if FileExists(path) {
			items = append(items, ConfigItem{
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
			items = append(items, ConfigItem{
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

	// 3. Project Settings
	if projectPath != "" {
		pathProj := filepath.Join(projectPath, ".claude", "settings.json")
		if FileExists(pathProj) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Project Settings",
				FileName: "settings.json",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	return items, nil
}

// --- Jules Provider ---

type JulesProvider struct{}

func (p *JulesProvider) Name() string {
	return "Jules"
}

func (p *JulesProvider) Discover(projectPath string) ([]ConfigItem, error) {
	var items []ConfigItem
	home := GetUserHome()

	// 1. Global Data
	if home != "" {
		path := filepath.Join(home, ".jules-mcp", "data.json")
		if FileExists(path) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "data.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Context (AGENTS.md)
	if projectPath != "" {
		pathProj := filepath.Join(projectPath, "AGENTS.md")
		if FileExists(pathProj) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Agents Context",
				FileName: "AGENTS.md",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}
	}

	return items, nil
}
