package config

import (
	"os"
	"testing"
)

func TestProviders_Create(t *testing.T) {
	tests := []struct {
		name        string
		provider    Provider
		scope       Scope
		expectError bool
	}{
		{
			name:     "Claude Provider - Create Project",
			provider: &ClaudeProvider{},
			scope:    ScopeProject,
		},
		{
			name:     "Claude Provider - Create Global",
			provider: &ClaudeProvider{},
			scope:    ScopeGlobal,
		},
		{
			name:     "Gemini Provider - Create Project",
			provider: &GeminiProvider{},
			scope:    ScopeProject,
		},
		{
			name:     "OpenCode Provider - Create Project",
			provider: &OpenCodeProvider{},
			scope:    ScopeProject,
		},
		{
			name:     "Crush Provider - Create Project",
			provider: &CrushProvider{},
			scope:    ScopeProject,
		},
		{
			name:     "Git Provider - Create Project",
			provider: &GitProvider{},
			scope:    ScopeProject,
		},
		{
			name:     "Aider Provider - Create Project",
			provider: &AiderProvider{},
			scope:    ScopeProject,
		},
		{
			name:     "Goose Provider - Create Global",
			provider: &GooseProvider{},
			scope:    ScopeGlobal,
		},
		{
			name:        "Copilot Provider - Create Project",
			provider:    &CopilotProvider{},
			scope:       ScopeProject,
			expectError: true,
		},
		{
			name:     "OpenAI Provider - Create Global",
			provider: &OpenAIProvider{},
			scope:    ScopeGlobal,
		},
		{
			name:        "Jules Provider - Create Project",
			provider:    &JulesProvider{},
			scope:       ScopeProject,
			expectError: true,
		},
		{
			name:     "Codex Provider - Create Project",
			provider: &CodexProvider{},
			scope:    ScopeProject,
		},
		{
			name:        "Claude Provider - Create Project (Empty Path)",
			provider:    &ClaudeProvider{},
			scope:       ScopeProject,
			expectError: true, // projectPath is empty (default in test loop if not set, but we pass tempProject. We need to pass empty)
		},
		{
			name:        "Claude Provider - Create Unsupported Scope",
			provider:    &ClaudeProvider{},
			scope:       ScopeSystem,
			expectError: true,
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

			// Special case for empty project path test
			projectPath := tempProject
			if tt.name == "Claude Provider - Create Project (Empty Path)" {
				projectPath = ""
			}

			path, err := tt.provider.Create(tt.scope, projectPath)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Create failed: %v", err)
				}
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("Created file does not exist: %s", path)
				}
			}
		})
	}
}
