package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	seen := map[string]bool{}
	add := func(it Item) {
		if !seen[it.Path] {
			items = append(items, it)
			seen[it.Path] = true
		}
	}

	// Helper to check and add a simple file
	checkFile := func(path, name string, scope Scope, format Format) {
		if FileExists(path) {
			add(Item{
				Provider: p.Name(),
				Name:     name,
				FileName: filepath.Base(path),
				Path:     path,
				Scope:    scope,
				Format:   format,
				Exists:   true,
			})
		}
	}

	home := paths.GetHomeDir()

	// 1. Global Settings & Extensions
	if home != "" {
		// Legacy Global
		checkFile(filepath.Join(home, ".gemini", "settings.json"), "User Settings (Legacy)", ScopeGlobal, FormatJSON)

		// XDG Global
		if cfgDir := paths.GetConfigDir("gemini-cli"); cfgDir != "" {
			checkFile(filepath.Join(cfgDir, "config.json"), "Global Config", ScopeGlobal, FormatJSON)
		}

		// Global Extensions (Legacy .gemini/extensions)
		exts, _ := fastWalk(filepath.Join(home, ".gemini"), 4, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			// Look for extension manifests or context files
			return strings.Contains(path, filepath.Join(".gemini", "extensions")) &&
				(strings.EqualFold(d.Name(), "gemini-extension.json") ||
					strings.EqualFold(d.Name(), "gemini.md"))
		})
		for _, match := range exts {
			add(Item{
				Provider: p.Name(),
				Name:     "Global Extension: " + filepath.Base(filepath.Dir(match)),
				FileName: filepath.Base(match),
				Path:     match,
				Scope:    ScopeGlobal,
				Format:   detectFormat(match),
				Exists:   true,
			})
		}
	}

	// 2. System Settings
	systemPaths := []string{}
	if dir := paths.GetConfigDir("gemini-cli"); dir != "" {
		systemPaths = append(systemPaths, filepath.Join(dir, "system-defaults.json"))
	}
	if runtime.GOOS == "linux" {
		systemPaths = append(systemPaths, "/etc/gemini-cli/system-defaults.json")
	}
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("PROGRAMDATA"); appData != "" {
			systemPaths = append(systemPaths, filepath.Join(appData, "gemini-cli", "system-defaults.json"))
		}
	}
	for _, sp := range systemPaths {
		checkFile(sp, "System Defaults", ScopeSystem, FormatJSON)
	}

	// 3. Project Settings
	if projectPath != "" {
		// New Standard: .mcp.json
		checkFile(filepath.Join(projectPath, ".mcp.json"), "MCP Config", ScopeProject, FormatJSON)

		// New Standard: .agent/
		// Rules
		rules, _ := fastWalk(filepath.Join(projectPath, ".agent", "rules"), 1, func(path string, d os.DirEntry) bool {
			return !d.IsDir() && strings.HasSuffix(d.Name(), ".md")
		})
		for _, r := range rules {
			add(Item{
				Provider: p.Name(),
				Name:     "Rule: " + strings.TrimSuffix(filepath.Base(r), ".md"),
				FileName: filepath.Base(r),
				Path:     r,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}

		// Workflows
		workflows, _ := fastWalk(filepath.Join(projectPath, ".agent", "workflows"), 1, func(path string, d os.DirEntry) bool {
			return !d.IsDir() && strings.HasSuffix(d.Name(), ".md")
		})
		for _, w := range workflows {
			add(Item{
				Provider: p.Name(),
				Name:     "Workflow: " + strings.TrimSuffix(filepath.Base(w), ".md"),
				FileName: filepath.Base(w),
				Path:     w,
				Scope:    ScopeProject,
				Format:   FormatMD,
				Exists:   true,
			})
		}

		// Legacy Project Configs
		checkFile(filepath.Join(projectPath, ".gemini", "settings.json"), "Workspace Settings (Legacy)", ScopeProject, FormatJSON)
		checkFile(filepath.Join(projectPath, ".gemini", "config.json"), "Workspace Config (Legacy)", ScopeProject, FormatJSON)
		checkFile(filepath.Join(projectPath, "GEMINI.md"), "Context File", ScopeProject, FormatMD)

		// Legacy Project Extensions
		exts, _ := fastWalk(projectPath, 6, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			inExtDir := strings.Contains(path, filepath.Join(".gemini", "extensions"))
			if !inExtDir {
				return false
			}
			name := strings.ToLower(d.Name())
			return name == "gemini-extension.json" || name == "gemini.md"
		})
		for _, match := range exts {
			add(Item{
				Provider: p.Name(),
				Name:     "Extension: " + filepath.Base(filepath.Dir(match)),
				FileName: filepath.Base(match),
				Path:     match,
				Scope:    ScopeProject,
				Format:   detectFormat(match),
				Exists:   true,
			})
		}
	}

	return items, nil
}

func detectFormat(path string) Format {
	if strings.HasSuffix(path, ".json") {
		return FormatJSON
	} else if strings.HasSuffix(path, ".md") {
		return FormatMD
	} else if strings.HasSuffix(path, ".toml") {
		return FormatTOML
	}
	return FormatTXT
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
