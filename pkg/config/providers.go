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

	// 3. System Settings (Linux, macOS, Windows)
	sysSettingsLinux := "/etc/claude-code/managed-settings.json"
	if FileExists(sysSettingsLinux) {
		items = append(items, Item{
			Provider: p.Name(),
			Name:     "Managed Settings (Linux)",
			FileName: "managed-settings.json",
			Path:     sysSettingsLinux,
			Scope:    ScopeSystem,
			Format:   FormatJSON,
			Exists:   true,
		})
	}

	// macOS System Settings
	sysSettingsMac := "/Library/Application Support/ClaudeCode/managed-settings.json"
	if FileExists(sysSettingsMac) {
		items = append(items, Item{
			Provider: p.Name(),
			Name:     "Managed Settings (macOS)",
			FileName: "managed-settings.json",
			Path:     sysSettingsMac,
			Scope:    ScopeSystem,
			Format:   FormatJSON,
			Exists:   true,
		})
	}

	// Windows System Settings
	// Note: Go's filepath.Join handles OS separators, but for absolute Windows paths we usually need environment variables.
	// Hardcoding common path for now, ideal solution would check runtime.GOOS and use os.Getenv("ProgramData")
	sysSettingsWin := "C:\\ProgramData\\ClaudeCode\\managed-settings.json"
	if FileExists(sysSettingsWin) {
		items = append(items, Item{
			Provider: p.Name(),
			Name:     "Managed Settings (Windows)",
			FileName: "managed-settings.json",
			Path:     sysSettingsWin,
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

// --- OpenCode Provider ---

type OpenCodeProvider struct{}

func (p *OpenCodeProvider) Name() string {
	return "OpenCode"
}

func (p *OpenCodeProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

	// 1. Global Config
	// Linux/macOS: ~/.config/opencode/opencode.json
	if home != "" {
		path := filepath.Join(home, ".config", "opencode", "opencode.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "opencode.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Config
	if projectPath != "" {
		// opencode.json
		pathProj := filepath.Join(projectPath, "opencode.json")
		if FileExists(pathProj) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "opencode.json",
				Path:     pathProj,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// opencode.local.json (Secrets)
		pathLocal := filepath.Join(projectPath, "opencode.local.json")
		if FileExists(pathLocal) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Local Secrets",
				FileName: "opencode.local.json",
				Path:     pathLocal,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}
	return items, nil
}

// --- Crush Provider ---

type CrushProvider struct{}

func (p *CrushProvider) Name() string {
	return "Crush CLI"
}

func (p *CrushProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := GetUserHome()

	// 1. Global Config
	if home != "" {
		// Linux/macOS Main Config
		pathMain := filepath.Join(home, ".config", "crush", "crush.json")
		if FileExists(pathMain) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "crush.json",
				Path:     pathMain,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// Linux/macOS Providers Config
		pathProviders := filepath.Join(home, ".config", "crush", "providers.json")
		if FileExists(pathProviders) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Providers",
				FileName: "providers.json",
				Path:     pathProviders,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// Windows Config (Approximation without runtime.GOOS check logic here)
		// %LOCALAPPDATA%\crush\crush.json -> usually C:\Users\<User>\AppData\Local\crush\crush.json
		pathWin := filepath.Join(home, "AppData", "Local", "crush", "crush.json")
		if FileExists(pathWin) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config (Windows)",
				FileName: "crush.json",
				Path:     pathWin,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// 2. Project Config
	if projectPath != "" {
		// .crush.json (Hidden)
		pathHidden := filepath.Join(projectPath, ".crush.json")
		if FileExists(pathHidden) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config (Hidden)",
				FileName: ".crush.json",
				Path:     pathHidden,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// crush.json (Visible)
		pathVisible := filepath.Join(projectPath, "crush.json")
		if FileExists(pathVisible) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: "crush.json",
				Path:     pathVisible,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// .crushignore
		pathIgnore := filepath.Join(projectPath, ".crushignore")
		if FileExists(pathIgnore) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Ignore File",
				FileName: ".crushignore",
				Path:     pathIgnore,
				Scope:    ScopeProject,
				Format:   FormatTXT,
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
