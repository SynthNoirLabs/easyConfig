package config

import (
	"os"
	"path/filepath"
)

type ConfigType string

const (
	JSON     ConfigType = "json"
	TOML     ConfigType = "toml"
	MARKDOWN ConfigType = "markdown" // For AGENTS.md
)

type AgentTool struct {
	Name       string     `json:"name"`
	ConfigPath string     `json:"configPath"`
	Type       ConfigType `json:"type"`
	Exists     bool       `json:"exists"`
}

type DiscoveryService struct{}

func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{}
}

// Discover scans known locations (Skeleton Impl)
func (s *DiscoveryService) Discover() []AgentTool {
	home, _ := os.UserHomeDir()
	// TODO: Implement full scanning logic from PROMPT.md
	// Note: Add logic to check for env vars like JULES_DATA_PATH
	return []AgentTool{
		{Name: "Claude Code", ConfigPath: filepath.Join(home, ".claude", "settings.json"), Type: JSON},
		{Name: "Codex CLI", ConfigPath: filepath.Join(home, ".codex", "config.toml"), Type: TOML},
		{Name: "Jules", ConfigPath: "AGENTS.md", Type: MARKDOWN},
	}
}
