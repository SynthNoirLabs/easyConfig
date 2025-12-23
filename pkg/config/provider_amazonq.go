package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"easyConfig/pkg/util/paths"
)

// --- Amazon Q Provider ---

type AmazonQProvider struct{}

func (p *AmazonQProvider) Name() string {
	return "Amazon Q"
}

func (p *AmazonQProvider) Create(scope Scope, _ string) (string, error) {
	// Amazon Q typically uses a global mcp.json
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		if runtime.GOOS == "windows" {
			path = filepath.Join(home, ".aws", "amazonq", "mcp.json")
		} else {
			path = filepath.Join(home, ".aws", "amazonq", "mcp.json")
		}
	default:
		return "", fmt.Errorf("unsupported scope for Amazon Q (only global supported)")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
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

func (p *AmazonQProvider) Discover(_ string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	if home != "" {
		// Global mcp.json
		// Works for both Windows (User Profile) and Unix (Home) as they both map to GetHomeDir in our utils usually,
		// but specifically checking standard AWS paths.
		path := filepath.Join(home, ".aws", "amazonq", "mcp.json")
		if FileExists(path) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "MCP Config",
				FileName: "mcp.json",
				Path:     path,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	return items, nil
}

func (p *AmazonQProvider) CheckStatus() ProviderStatus {
	files, _ := p.Discover("")
	return ProviderStatus{
		ProviderName:    p.Name(),
		Health:          StatusHealthy,
		StatusMessage:   "Amazon Q status checking not implemented yet.",
		DiscoveredFiles: files,
		LastChecked:     time.Now().Format(time.RFC3339),
	}
}
