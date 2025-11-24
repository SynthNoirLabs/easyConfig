package config

import (
	"path/filepath"
)

// --- Claude Provider ---

type ClaudeProvider struct{}

func (p *ClaudeProvider) Name() string {
	return "Claude Code"
}

func (p *ClaudeProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

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

	// 3. System Settings (Linux)
	sysSettings := "/etc/claude-code/managed-settings.json"
	if FileExists(sysSettings) {
		items = append(items, Item{
			Provider: p.Name(),
			Name:     "Managed Settings",
			FileName: "managed-settings.json",
			Path:     sysSettings,
			Scope:    ScopeSystem,
			Format:   FormatJSON,
			Exists:   true,
		})
	}
	sysMCP := "/etc/claude-code/managed-mcp.json"
	if FileExists(sysMCP) {
		items = append(items, Item{
			Provider: p.Name(),
			Name:     "Managed MCP",
			FileName: "managed-mcp.json",
			Path:     sysMCP,
			Scope:    ScopeSystem,
			Format:   FormatJSON,
			Exists:   true,
		})
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
	}

	return items, nil
}

// --- Copilot Provider ---

type CopilotProvider struct{}

func (p *CopilotProvider) Name() string {
	return "GitHub Copilot"
}

func (p *CopilotProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

	// 1. Global CLI Config
	if home != "" {
		path := filepath.Join(home, ".copilot", "mcp-config.json")
		if FileExists(path) {
			items = append(items, Item{
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
			items = append(items, Item{
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

func (p *OpenAIProvider) Discover(_ string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

	if home != "" {
		path := filepath.Join(home, ".config", "openai", "config.yaml")
		if FileExists(path) {
			items = append(items, Item{
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

func (p *JulesProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

	// 1. Global Data
	if home != "" {
		path := filepath.Join(home, ".jules-mcp", "data.json")
		if FileExists(path) {
			items = append(items, Item{
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
			items = append(items, Item{
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

func (p *GeminiProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

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
	sysDefaults := "/etc/gemini-cli/system-defaults.json"
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
	sysOverrides := "/etc/gemini-cli/settings.json"
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
	}

	return items, nil
}

// --- Codex Provider ---

type CodexProvider struct{}

func (p *CodexProvider) Name() string {
	return "Codex CLI"
}

func (p *CodexProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

	// 1. Global Config
	if home != "" {
		path := filepath.Join(home, ".codex", "config.toml")
		if FileExists(path) {
			items = append(items, Item{
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
			items = append(items, Item{
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
