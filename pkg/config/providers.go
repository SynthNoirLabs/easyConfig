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

// --- Copilot Provider ---

type CopilotProvider struct{}

func (p *CopilotProvider) Name() string {
	return "GitHub Copilot"
}

func (p *CopilotProvider) Discover(projectPath string) ([]ConfigItem, error) {
	var items []ConfigItem
	home := GetUserHome()

	// 1. Global CLI Config
	if home != "" {
		path := filepath.Join(home, ".copilot", "mcp-config.json")
		if FileExists(path) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "CLI Config",
				FileName: "mcp-config.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Context
	if projectPath != "" {
		pathProj := filepath.Join(projectPath, ".github", "copilot-instructions.md")
		if FileExists(pathProj) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Instructions",
				FileName: "copilot-instructions.md",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}
	}
	return items, nil
}

// --- OpenAI Provider ---

type OpenAIProvider struct{}

func (p *OpenAIProvider) Name() string {
	return "OpenAI"
}

func (p *OpenAIProvider) Discover(projectPath string) ([]ConfigItem, error) {
	var items []ConfigItem
	home := GetUserHome()

	if home != "" {
		path := filepath.Join(home, ".config", "openai", "config.yaml")
		if FileExists(path) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "config.yaml",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatYAML,
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

// --- Gemini Provider ---

type GeminiProvider struct{}

func (p *GeminiProvider) Name() string {
	return "Gemini"
}

func (p *GeminiProvider) Discover(projectPath string) ([]ConfigItem, error) {
	var items []ConfigItem
	home := GetUserHome()

	// 1. Global settings
	if home != "" {
		path := filepath.Join(home, ".gemini", "settings.json")
		if FileExists(path) {
			items = append(items, ConfigItem{
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

	// 2. Project settings
	if projectPath != "" {
		path := filepath.Join(projectPath, ".gemini", "settings.json")
		if FileExists(path) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Workspace Settings",
				FileName: "settings.json",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	return items, nil
}

// --- Codex Provider ---

type CodexProvider struct{}

func (p *CodexProvider) Name() string {
	return "Codex CLI"
}

func (p *CodexProvider) Discover(projectPath string) ([]ConfigItem, error) {
	var items []ConfigItem
	home := GetUserHome()

	// 1. Global Config
	if home != "" {
		path := filepath.Join(home, ".codex", "config.toml")
		if FileExists(path) {
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "config.toml",
				Path:     path,
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
			items = append(items, ConfigItem{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "config.toml",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatTOML,
				Exists:   true,
			})
		}
	}

	return items, nil
}
