package config

import (
	"fmt"
	"os"
	"path/filepath"

	"easyConfig/pkg/util/paths"
)

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
	seen := map[string]bool{}
	add := func(it Item) {
		if !seen[it.Path] {
			items = append(items, it)
			seen[it.Path] = true
		}
	}

	// 1. Global Config
	// Linux/macOS: ~/.config/opencode/opencode.json
	configDir := paths.GetConfigDir("opencode")
	if configDir != "" {
		for _, fname := range []string{"opencode.json", "opencode.jsonc"} {
			path := filepath.Join(configDir, fname)
			if FileExists(path) {
				add(Item{
					Provider: p.Name(),
					Name:     "Global Config",
					FileName: fname,
					Path:     path,
					Scope:    ScopeGlobal,
					Format:   FormatJSON,
					Exists:   true,
				})
			}
		}
	}

	// 1b. OPENCODE_CONFIG env override
	if cfgEnv := os.Getenv("OPENCODE_CONFIG"); cfgEnv != "" && FileExists(cfgEnv) {
		add(Item{
			Provider: p.Name(),
			Name:     "Env Config",
			FileName: filepath.Base(cfgEnv),
			Path:     cfgEnv,
			Scope:    ScopeGlobal,
			Format:   FormatJSON,
			Exists:   true,
		})
	}

	// 1c. OPENCODE_CONFIG_DIR custom dir
	if cfgDirEnv := os.Getenv("OPENCODE_CONFIG_DIR"); cfgDirEnv != "" {
		for _, fname := range []string{"opencode.json", "opencode.jsonc"} {
			path := filepath.Join(cfgDirEnv, fname)
			if FileExists(path) {
				add(Item{
					Provider: p.Name(),
					Name:     "Custom Config",
					FileName: fname,
					Path:     path,
					Scope:    ScopeGlobal,
					Format:   FormatJSON,
					Exists:   true,
				})
			}
		}
	}

	// 2. Project Config
	if projectPath != "" {
		paths, _ := fastWalk(projectPath, 4, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			base := filepath.Base(path)
			return base == "opencode.json" || base == "opencode.jsonc" || base == "opencode.local.json"
		})
		for _, pth := range paths {
			name := "Project Config"
			if filepath.Base(pth) == "opencode.local.json" {
				name = "Local Secrets"
			}
			add(Item{
				Provider: p.Name(),
				Name:     name,
				FileName: filepath.Base(pth),
				Path:     pth,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}
	return items, nil
}
