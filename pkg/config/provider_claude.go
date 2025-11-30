package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	seen := make(map[string]bool)
	add := func(it Item) {
		if !seen[it.Path] {
			items = append(items, it)
			seen[it.Path] = true
		}
	}

	// 1. Global Desktop Config
	if home != "" {
		path := filepath.Join(home, ".claude", "claude_desktop_config.json")
		if FileExists(path) {
			add(Item{
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
			add(Item{
				Provider: p.Name(),
				Name:     "CLI Settings",
				FileName: "settings.json",
				Path:     pathCLI,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// OS config dir (macOS ~/Library/Application Support/Claude, Win %APPDATA%/Claude, Linux ~/.config/Claude)
		if cfgDir := paths.GetConfigDir("Claude"); cfgDir != "" {
			pathDesktop := filepath.Join(cfgDir, "claude_desktop_config.json")
			if FileExists(pathDesktop) {
				add(Item{
					Provider: p.Name(),
					Name:     "Desktop Config",
					FileName: "claude_desktop_config.json",
					Path:     pathDesktop,
					Scope:    ScopeGlobal,
					Format:   FormatJSON,
					Exists:   true,
				})
			}
		}

		// Global memory file
		globalMemory := filepath.Join(home, ".claude", "CLAUDE.md")
		if FileExists(globalMemory) {
			add(Item{
				Provider: p.Name(),
				Name:     "Global Memory",
				FileName: "CLAUDE.md",
				Path:     globalMemory,
				Scope:    ScopeGlobal,
				Format:   FormatMD,
				Exists:   true,
			})
		}

		// Global commands/agents/hooks under ~/.claude
		if globalPaths, _ := fastWalk(filepath.Join(home, ".claude"), 4, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
				return false
			}
			return strings.Contains(path, string(filepath.Separator)+"commands"+string(filepath.Separator)) ||
				strings.Contains(path, string(filepath.Separator)+"agents"+string(filepath.Separator)) ||
				strings.Contains(path, string(filepath.Separator)+"hooks"+string(filepath.Separator))
		}); len(globalPaths) > 0 {
			for _, gp := range globalPaths {
				base := filepath.Base(gp)
				switch {
				case strings.Contains(gp, string(filepath.Separator)+"commands"+string(filepath.Separator)):
					add(Item{
						Provider: p.Name(),
						Name:     "Command: " + base,
						FileName: base,
						Path:     gp,
						Scope:    ScopeGlobal,
						Format:   FormatMD,
						Exists:   true,
					})
				case strings.Contains(gp, string(filepath.Separator)+"hooks"+string(filepath.Separator)):
					add(Item{
						Provider: p.Name(),
						Name:     "Hook: " + base,
						FileName: base,
						Path:     gp,
						Scope:    ScopeGlobal,
						Format:   FormatMD,
						Exists:   true,
					})
				default:
					add(Item{
						Provider: p.Name(),
						Name:     "Subagent: " + base,
						FileName: base,
						Path:     gp,
						Scope:    ScopeGlobal,
						Format:   FormatMD,
						Exists:   true,
					})
				}
			}
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
			add(Item{
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
			add(Item{
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
		// settings.json (direct lookups)
		pathProj := filepath.Join(projectPath, ".claude", "settings.json")
		if FileExists(pathProj) {
			add(Item{
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
			add(Item{
				Provider: p.Name(),
				Name:     "Local Settings",
				FileName: "settings.local.json",
				Path:     pathLocal,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// CLAUDE.md (root)
		pathMemory := filepath.Join(projectPath, "CLAUDE.md")
		if FileExists(pathMemory) {
			add(Item{
				Provider: p.Name(),
				Name:     "Memory File",
				FileName: "CLAUDE.md",
				Path:     pathMemory,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}
		// CLAUDE.local.md (project-local memory)
		pathMemoryLocal := filepath.Join(projectPath, "CLAUDE.local.md")
		if FileExists(pathMemoryLocal) {
			add(Item{
				Provider: p.Name(),
				Name:     "Local Memory",
				FileName: "CLAUDE.local.md",
				Path:     pathMemoryLocal,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}

		// Deep scan (fast, with exclusions) for Claude files, commands, agents, hooks.
		matches, _ := fastWalk(projectPath, 6, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			lower := strings.ToLower(d.Name())
			if lower == "claude.md" || lower == "claude.local.md" {
				return true
			}
			if strings.HasSuffix(lower, ".md") &&
				(strings.Contains(path, string(filepath.Separator)+".claude"+string(filepath.Separator)+"commands"+string(filepath.Separator)) ||
					strings.Contains(path, string(filepath.Separator)+".claude"+string(filepath.Separator)+"agents"+string(filepath.Separator)) ||
					strings.Contains(path, string(filepath.Separator)+".claude"+string(filepath.Separator)+"hooks"+string(filepath.Separator)) ||
					strings.Contains(path, string(filepath.Separator)+"agents"+string(filepath.Separator))) {
				return true
			}
			if (lower == "settings.json" || lower == "settings.local.json") &&
				strings.Contains(path, string(filepath.Separator)+".claude"+string(filepath.Separator)) {
				return true
			}
			return false
		})

		for _, match := range matches {
			name := filepath.Base(match)
			switch strings.ToLower(name) {
			case "claude.md":
				add(Item{
					Provider: p.Name(),
					Name:     "Memory File",
					FileName: "CLAUDE.md",
					Path:     match,
					Scope:    ScopeProject,
					Format:   FormatMD,
					Exists:   true,
				})
			case "claude.local.md":
				add(Item{
					Provider: p.Name(),
					Name:     "Local Memory",
					FileName: "CLAUDE.local.md",
					Path:     match,
					Scope:    ScopeProject,
					Format:   FormatMD,
					Exists:   true,
				})
			case "settings.json":
				add(Item{
					Provider: p.Name(),
					Name:     "Project Settings",
					FileName: "settings.json",
					Path:     match,
					Scope:    ScopeProject,
					Format:   FormatJSON,
					Exists:   true,
				})
			case "settings.local.json":
				add(Item{
					Provider: p.Name(),
					Name:     "Local Settings",
					FileName: "settings.local.json",
					Path:     match,
					Scope:    ScopeProject,
					Format:   FormatJSON,
					Exists:   true,
				})
			default:
				if strings.Contains(match, string(filepath.Separator)+"commands"+string(filepath.Separator)) {
					add(Item{
						Provider: p.Name(),
						Name:     "Command: " + filepath.Base(match),
						FileName: filepath.Base(match),
						Path:     match,
						Scope:    ScopeProject,
						Format:   FormatMD,
						Exists:   true,
					})
				} else if strings.Contains(match, string(filepath.Separator)+"hooks"+string(filepath.Separator)) {
					add(Item{
						Provider: p.Name(),
						Name:     "Hook: " + filepath.Base(match),
						FileName: filepath.Base(match),
						Path:     match,
						Scope:    ScopeProject,
						Format:   FormatMD,
						Exists:   true,
					})
				} else {
					add(Item{
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
		}
	}

	return items, nil
}
