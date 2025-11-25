package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	}

	return items, nil
}

// --- OpenCode Provider ---

type OpenCodeProvider struct{}

func (p *OpenCodeProvider) Name() string {
	return "OpenCode"
}

func (p *OpenCodeProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		configDir := paths.GetConfigDir("opencode")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "opencode.json")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, "opencode.json")
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

func (p *OpenCodeProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item

	// 1. Global Config
	// Linux/macOS: ~/.config/opencode/opencode.json
	configDir := paths.GetConfigDir("opencode")
	if configDir != "" {
		path := filepath.Join(configDir, "opencode.json")
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

func (p *CrushProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		configDir := paths.GetConfigDir("crush")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "crush.json")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, "crush.json")
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

func (p *CrushProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item

	// 1. Global Config
	configDir := paths.GetConfigDir("crush")
	if configDir != "" {
		// Main Config
		pathMain := filepath.Join(configDir, "crush.json")
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

		// Providers Config
		pathProviders := filepath.Join(configDir, "providers.json")
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

	// 2. System Config (/etc/gitconfig)
	// Note: Windows path differs, usually in Program Files, skipping for simplicity/Linux focus
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

type CopilotProvider struct{}

func (p *CopilotProvider) Name() string {
	return "GitHub Copilot"
}

func (p *CopilotProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".copilot", "mcp-config.json")
	case ScopeProject:
		return "", fmt.Errorf("project creation not supported for Copilot yet")
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

func (p *CopilotProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

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

func (p *OpenAIProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "version: 1\n"
	var path string

	switch scope {
	case ScopeGlobal:
		configDir := paths.GetConfigDir("openai")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "config.yaml")
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

func (p *OpenAIProvider) Discover(_ string) ([]Item, error) {
	var items []Item

	configDir := paths.GetConfigDir("openai")
	if configDir != "" {
		path := filepath.Join(configDir, "config.yaml")
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

func (p *JulesProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".jules-mcp", "data.json")
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

func (p *JulesProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

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
				items = append(items, Item{
					Provider: p.Name(),
					Name:     "Extension: " + filepath.Base(match),
					FileName: filepath.Base(match),
					Path:     match,
					Scope:    ScopeProject,
					Format:   FormatTXT, // Default to text, could be JSON/JS
					Exists:   true,
				})
			}
		}
	}

	// Global Extensions (~/.gemini/extensions/*)
	if home != "" {
		globalExtPath := filepath.Join(home, ".gemini", "extensions", "*")
		if matches, err := filepath.Glob(globalExtPath); err == nil {
			for _, match := range matches {
				items = append(items, Item{
					Provider: p.Name(),
					Name:     "Extension: " + filepath.Base(match),
					FileName: filepath.Base(match),
					Path:     match,
					Scope:    ScopeGlobal,
					Format:   FormatTXT, // Default to text
					Exists:   true,
				})
			}
		}
	}

	return items, nil
}

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

func (p *CodexProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

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

// --- Aider Provider ---

type AiderProvider struct{}

func (p *AiderProvider) Name() string {
	return "Aider"
}

func (p *AiderProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "model: gpt-4\n"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".aider.conf.yml")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".aider.conf.yml")
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

func (p *AiderProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// 1. Global Config
	if home != "" {
		path := filepath.Join(home, ".aider.conf.yml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: ".aider.conf.yml",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
	}

	// 2. Project Config
	if projectPath != "" {
		path := filepath.Join(projectPath, ".aider.conf.yml")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: ".aider.conf.yml",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatYAML,
				Exists:   true,
			})
		}
	}
	return items, nil
}

// --- Goose Provider ---

type GooseProvider struct{}

func (p *GooseProvider) Name() string {
	return "Goose"
}

func (p *GooseProvider) Create(scope Scope, _ string) (string, error) {
	defaultContent := "version: 1\n"
	var path string

	switch scope {
	case ScopeGlobal:
		// Goose uses ~/.config/goose/config.yaml on Linux/macOS
		// On Windows it uses AppData/Block/goose/config/config.yaml
		// paths.GetConfigDir("goose") returns ~/.config/goose or AppData/goose
		// We might need special handling for Windows if "Block" is strictly required.
		// For now, we use standard XDG/AppData structure.
		configDir := paths.GetConfigDir("goose")
		if configDir == "" {
			return "", fmt.Errorf("config directory not found")
		}
		path = filepath.Join(configDir, "config.yaml")
	default:
		return "", fmt.Errorf("unsupported scope (Goose only supports global config)")
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

func (p *GooseProvider) Discover(_ string) ([]Item, error) {
	var items []Item

	// 1. Global Config
	configDir := paths.GetConfigDir("goose")
	if configDir != "" {
		path := filepath.Join(configDir, "config.yaml")
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
