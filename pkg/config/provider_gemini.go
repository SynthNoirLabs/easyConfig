package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	seen := map[string]bool{}
	add := func(it Item) {
		if !seen[it.Path] {
			items = append(items, it)
			seen[it.Path] = true
		}
	}

	// 1. Global settings
	if home != "" {
		path := filepath.Join(home, ".gemini", "settings.json")
		if FileExists(path) {
			add(Item{
				Provider: p.Name(),
				Name:     "User Settings",
				FileName: "settings.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}

		// Optional legacy/global config.json (XDG)
		if cfgDir := paths.GetConfigDir("gemini-cli"); cfgDir != "" {
			cfgJSON := filepath.Join(cfgDir, "config.json")
			if FileExists(cfgJSON) {
				add(Item{
					Provider: p.Name(),
					Name:     "Global Config",
					FileName: "config.json",
					Path:     cfgJSON,
					Scope:    ScopeGlobal,
					Format:   FormatJSON,
					Exists:   true,
				})
			}
		}
	}

	// 2. System Settings
	sysCandidates := []string{}
	if dir := paths.GetConfigDir("gemini-cli"); dir != "" {
		sysCandidates = append(sysCandidates,
			filepath.Join(dir, "system-defaults.json"),
			filepath.Join(dir, "settings.json"))
	}
	if runtime.GOOS == "linux" {
		sysCandidates = append(sysCandidates,
			"/etc/gemini-cli/system-defaults.json",
			"/etc/gemini-cli/settings.json")
	}
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("PROGRAMDATA"); appData != "" {
			sysCandidates = append(sysCandidates,
				filepath.Join(appData, "gemini-cli", "system-defaults.json"),
				filepath.Join(appData, "gemini-cli", "settings.json"))
		}
	}
	for _, candidate := range sysCandidates {
		if FileExists(candidate) {
			name := "System Config"
			if strings.Contains(candidate, "system-defaults") {
				name = "System Defaults"
			} else if strings.Contains(candidate, "settings.json") {
				name = "System Overrides"
			}
			add(Item{
				Provider: p.Name(),
				Name:     name,
				FileName: filepath.Base(candidate),
				Path:     candidate,
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
			add(Item{
				Provider: p.Name(),
				Name:     "Workspace Settings",
				FileName: "settings.json",
				Path:     path,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
		// Optional project config.json
		pathCfg := filepath.Join(projectPath, ".gemini", "config.json")
		if FileExists(pathCfg) {
			add(Item{
				Provider: p.Name(),
				Name:     "Workspace Config",
				FileName: "config.json",
				Path:     pathCfg,
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

		// Project Extensions via fast scan
		exts, _ := fastWalk(projectPath, 6, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			inExtDir := strings.Contains(path, string(filepath.Separator)+".gemini"+string(filepath.Separator)+"extensions"+string(filepath.Separator))
			if !inExtDir {
				return false
			}
			name := strings.ToLower(d.Name())
			return name == "gemini-extension.json" ||
				name == "gemini.md" ||
				strings.HasSuffix(name, ".json") ||
				strings.HasSuffix(name, ".js") ||
				strings.HasSuffix(name, ".ts") ||
				strings.HasSuffix(name, ".toml")
		})
		for _, match := range exts {
			name := "Extension: " + filepath.Base(filepath.Dir(match))
			format := FormatTXT
			if strings.HasSuffix(match, ".json") {
				format = FormatJSON
			} else if strings.HasSuffix(match, ".md") {
				format = FormatMD
			} else if strings.HasSuffix(match, ".toml") {
				format = FormatTOML
			}
			add(Item{
				Provider: p.Name(),
				Name:     name,
				FileName: filepath.Base(match),
				Path:     match,
				Scope:    ScopeProject,
				Format:   format,
				Exists:   true,
			})
		}
	}

	// Global Extensions (~/.gemini/extensions/*)
	if home != "" {
		exts, _ := fastWalk(filepath.Join(home, ".gemini"), 4, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			inExtDir := strings.Contains(path, string(filepath.Separator)+"extensions"+string(filepath.Separator))
			if !inExtDir {
				return false
			}
			name := strings.ToLower(d.Name())
			return name == "gemini-extension.json" ||
				name == "gemini.md" ||
				strings.HasSuffix(name, ".json") ||
				strings.HasSuffix(name, ".js") ||
				strings.HasSuffix(name, ".ts") ||
				strings.HasSuffix(name, ".toml")
		})
		for _, match := range exts {
			name := "Extension: " + filepath.Base(filepath.Dir(match))
			format := FormatTXT
			if strings.HasSuffix(match, ".json") {
				format = FormatJSON
			} else if strings.HasSuffix(match, ".md") {
				format = FormatMD
			} else if strings.HasSuffix(match, ".toml") {
				format = FormatTOML
			}
			add(Item{
				Provider: p.Name(),
				Name:     name,
				FileName: filepath.Base(match),
				Path:     match,
				Scope:    ScopeGlobal,
				Format:   format,
				Exists:   true,
			})
		}
	}

	return items, nil
}
