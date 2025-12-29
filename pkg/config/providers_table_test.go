package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestProviders_Discover_TableDriven(t *testing.T) {
	if runtime.GOOS == "darwin" && os.Getenv("CI") != "" {
		t.Skip("Skipping provider discovery table test on darwin CI (permission/path differences)")
	}

	tests := []struct {
		name          string
		provider      Provider
		setupGlobal   func(t *testing.T, homeDir string)
		setupProject  func(t *testing.T, projectDir string)
		expectedCount int
		checkItems    func(t *testing.T, items []Item)
	}{
		{
			name:     "Claude Provider",
			provider: &ClaudeProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".claude")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create claude dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write settings.json: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				dir := filepath.Join(projectDir, ".claude")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create project claude dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write project settings.json: %v", err)
				}
				if err := os.WriteFile(filepath.Join(projectDir, "CLAUDE.md"), []byte("# Context"), 0o600); err != nil {
					t.Fatalf("Failed to write CLAUDE.md: %v", err)
				}

				// Subagents
				agentsDir := filepath.Join(projectDir, "agents")
				if err := os.MkdirAll(agentsDir, 0o750); err != nil {
					t.Fatalf("Failed to create agents dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(agentsDir, "coder.md"), []byte("# Coder"), 0o600); err != nil {
					t.Fatalf("Failed to write coder.md: %v", err)
				}

				// Custom Commands
				cmdDir := filepath.Join(projectDir, ".claude", "commands")
				if err := os.MkdirAll(cmdDir, 0o750); err != nil {
					t.Fatalf("Failed to create commands dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(cmdDir, "test.md"), []byte("# Command"), 0o600); err != nil {
					t.Fatalf("Failed to write test.md: %v", err)
				}
			},
			expectedCount: 5, // Global CLI, Project Settings, CLAUDE.md, Subagent, Command
			checkItems: func(t *testing.T, items []Item) {
				foundSubagent := false
				foundCommand := false
				for _, item := range items {
					if item.Name == "Subagent: coder.md" {
						foundSubagent = true
					}
					if item.Name == "Command: test.md" {
						foundCommand = true
					}
				}
				if !foundSubagent {
					t.Error("Expected to find subagent")
				}
				if !foundCommand {
					t.Error("Expected to find custom command")
				}
			},
		},
		{
			name:     "Gemini Provider",
			provider: &GeminiProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".gemini")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create gemini dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write gemini settings.json: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				dir := filepath.Join(projectDir, ".gemini")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create project gemini dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write project gemini settings.json: %v", err)
				}

				// Extensions
				extDir := filepath.Join(dir, "extensions", "my-ext")
				if err := os.MkdirAll(extDir, 0o750); err != nil {
					t.Fatalf("Failed to create extensions dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(extDir, "gemini-extension.json"), []byte(`{"displayName": "My Ext"}`), 0o600); err != nil {
					t.Fatalf("Failed to write gemini-extension.json: %v", err)
				}
			},
			expectedCount: 3, // Global, Project, Extension
			checkItems: func(t *testing.T, items []Item) {
				foundExt := false
				for _, item := range items {
					if item.Name == "Extension: my-ext" && item.Format == FormatJSON {
						foundExt = true
					}
				}
				if !foundExt {
					t.Error("Expected to find extension with JSON format")
				}
			},
		},
		{
			name:     "OpenCode Provider",
			provider: &OpenCodeProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".config", "opencode")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create opencode dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "opencode.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write opencode.json: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				if err := os.WriteFile(filepath.Join(projectDir, "opencode.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write project opencode.json: %v", err)
				}
			},
			expectedCount: 2,
		},
		{
			name:     "Crush Provider",
			provider: &CrushProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".config", "crush")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create crush dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "crush.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write crush.json: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				if err := os.WriteFile(filepath.Join(projectDir, "crush.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write project crush.json: %v", err)
				}
			},
			expectedCount: 2,
		},
		{
			name:     "Git Provider",
			provider: &GitProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				if err := os.WriteFile(filepath.Join(homeDir, ".gitconfig"), []byte("[user]"), 0o600); err != nil {
					t.Fatalf("Failed to write .gitconfig: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				dir := filepath.Join(projectDir, ".git")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create .git dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "config"), []byte("[core]"), 0o600); err != nil {
					t.Fatalf("Failed to write git config: %v", err)
				}
			},
			checkItems: func(t *testing.T, items []Item) {
				// Git provider may discover 2 or 3 items depending on whether /etc/gitconfig exists
				// We should have at least: global (.gitconfig) and project (.git/config)
				// Optionally: system (/etc/gitconfig)
				if len(items) < 2 {
					t.Errorf("Expected at least 2 items (global + project), got %d", len(items))
				}
				if len(items) > 3 {
					t.Errorf("Expected at most 3 items (global + project + system), got %d", len(items))
				}

				// Verify that we have the required scopes
				hasGlobal := false
				hasProject := false
				for _, item := range items {
					if item.Scope == ScopeGlobal {
						hasGlobal = true
					}
					if item.Scope == ScopeProject {
						hasProject = true
					}
				}
				if !hasGlobal {
					t.Error("Missing global scope config")
				}
				if !hasProject {
					t.Error("Missing project scope config")
				}
			},
		},
		{
			name:     "Aider Provider",
			provider: &AiderProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				if err := os.WriteFile(filepath.Join(homeDir, ".aider.conf.yml"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write global .aider.conf.yml: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				if err := os.WriteFile(filepath.Join(projectDir, ".aider.conf.yml"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write project .aider.conf.yml: %v", err)
				}
			},
			expectedCount: 2,
		},
		{
			name:     "Goose Provider",
			provider: &GooseProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".config", "goose")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create goose config dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write global config.yaml: %v", err)
				}
			},
			setupProject:  func(t *testing.T, projectDir string) {}, // No project config
			expectedCount: 1,
		},
		{
			name:     "Copilot Provider",
			provider: &CopilotProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".copilot")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create .copilot dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "mcp-config.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write global mcp-config.json: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				dir := filepath.Join(projectDir, ".github")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create .github dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "copilot-instructions.md"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write copilot-instructions.md: %v", err)
				}
			},
			expectedCount: 2,
		},
		{
			name:     "OpenAI Provider",
			provider: &OpenAIProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".config", "openai")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create openai config dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write global config.yaml: %v", err)
				}
			},
			setupProject:  func(t *testing.T, projectDir string) {},
			expectedCount: 1,
		},
		{
			name:     "Jules Provider",
			provider: &JulesProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".jules-mcp")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create .jules-mcp dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "data.json"), []byte("{}"), 0o600); err != nil {
					t.Fatalf("Failed to write global data.json: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				if err := os.WriteFile(filepath.Join(projectDir, "AGENTS.md"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write AGENTS.md: %v", err)
				}
			},
			expectedCount: 2,
		},
		{
			name:     "Codex Provider",
			provider: &CodexProvider{},
			setupGlobal: func(t *testing.T, homeDir string) {
				dir := filepath.Join(homeDir, ".codex")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create .codex dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write global config.toml: %v", err)
				}
			},
			setupProject: func(t *testing.T, projectDir string) {
				dir := filepath.Join(projectDir, ".codex")
				if err := os.MkdirAll(dir, 0o750); err != nil {
					t.Fatalf("Failed to create project .codex dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(""), 0o600); err != nil {
					t.Fatalf("Failed to write project config.toml: %v", err)
				}
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHome := t.TempDir()
			tempProject := t.TempDir()

			originalHome := os.Getenv("HOME")
			t.Setenv("HOME", tempHome)
			defer func() {
				if err := os.Setenv("HOME", originalHome); err != nil {
					t.Errorf("Failed to restore HOME: %v", err)
				}
			}()

			if tt.setupGlobal != nil {
				tt.setupGlobal(t, tempHome)
			}
			if tt.setupProject != nil {
				tt.setupProject(t, tempProject)
			}

			items, err := tt.provider.Discover(tempProject)
			if err != nil {
				t.Fatalf("Discover failed: %v", err)
			}

			// Only check expectedCount if it's non-zero (allows checkItems to be used alone)
			if tt.expectedCount > 0 && len(items) != tt.expectedCount {
				t.Errorf("Expected %d items, got %d", tt.expectedCount, len(items))
				for _, item := range items {
					t.Logf("Found item: %+v", item)
				}
			}

			if tt.checkItems != nil {
				tt.checkItems(t, items)
			}
		})
	}
}
