package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Copilot Provider ---

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
	seen := map[string]bool{}
	add := func(it Item) {
		if !seen[it.Path] {
			items = append(items, it)
			seen[it.Path] = true
		}
	}

	// 1. Global CLI Config
	if home != "" {
		path := filepath.Join(home, ".copilot", "mcp-config.json")
		if FileExists(path) {
			add(Item{
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
		paths, _ := fastWalk(projectPath, 5, func(path string, d os.DirEntry) bool {
			if d.IsDir() {
				return false
			}
			base := filepath.Base(path)
			if strings.HasPrefix(base, "copilot") && (strings.HasSuffix(base, ".md") || strings.HasSuffix(base, ".json")) {
				return true
			}
			// in .github directory
			return strings.Contains(path, string(filepath.Separator)+".github"+string(filepath.Separator)) &&
				(strings.Contains(base, "copilot") || base == "mcp-config.json")
		})
		for _, pth := range paths {
			format := FormatTXT
			switch filepath.Ext(pth) {
			case ".md":
				format = FormatMD
			case ".json":
				format = FormatJSON
			}
			add(Item{
				Provider: p.Name(),
				Name:     filepath.Base(pth),
				FileName: filepath.Base(pth),
				Path:     pth,
				Scope:    ScopeProject,
				Format:   format,
				Exists:   true,
			})
		}
	}
	return items, nil
}

func (p *CopilotProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusUnknown,
		StatusMessage:   "Copilot status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
