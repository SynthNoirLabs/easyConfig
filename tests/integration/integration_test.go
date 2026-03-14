package main_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"easyConfig/pkg/config"
)

// TestIntegration_Discovery verifies that the backend logic correctly
// identifies configuration files created in a mock environment that
// mimics the Docker container structure.
//
// Note: This test runs on the host machine but simulates the paths.
// For true integration inside Docker, we would compile this binary
// and run it inside the container.
func TestIntegration_Discovery(t *testing.T) {
	// 1. Setup Mock Environment
	tempHome := t.TempDir()
	tempProject := t.TempDir()

	// Set HOME to tempHome to simulate user directory
	t.Setenv("HOME", tempHome)

	// 2. Install Mock Config Files (mimicking what CLIs would create)

	// Claude
	createFile(t, filepath.Join(tempHome, ".claude", "settings.json"), "{}")
	createFile(t, filepath.Join(tempProject, ".claude", "settings.json"), "{}")

	// Gemini
	createFile(t, filepath.Join(tempHome, ".gemini", "settings.json"), "{}")
	createFile(t, filepath.Join(tempProject, "GEMINI.md"), "# Context")

	// Codex
	createFile(t, filepath.Join(tempHome, ".codex", "config.toml"), "# config")

	// 3. Run Discovery
	ds := config.NewDiscoveryService(nil, nil)
	items, err := ds.DiscoverAll(context.Background(), tempProject)
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	// 4. Verify Results
	expectedProviders := map[string]int{
		"Claude Code": 2, // Global CLI + Project Settings
		"Gemini":      2, // Global Settings + Project Context
		"Codex CLI":   1, // Global Config
	}

	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Provider]++
	}

	for provider, expected := range expectedProviders {
		if counts[provider] < expected {
			t.Errorf("Provider %s: expected at least %d configs, got %d", provider, expected, counts[provider])
		}
	}
}

// TestIntegration_CLI_Interaction simulates the interaction between
// easyConfig modifying a file and a CLI tool reading it.
// Since we cannot easily run the actual proprietary CLIs in this test environment
// without authentication, we mock the CLI behavior by checking file integrity.
func TestIntegration_CLI_Interaction(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	// 1. Setup Codex Config
	configPath := filepath.Join(tempHome, ".codex", "config.toml")
	createFile(t, configPath, "theme = 'dark'\n")

	// 2. easyConfig modifies the file
	ds := config.NewDiscoveryService(nil, nil)

	// Verify read
	content, err := ds.ReadConfig(configPath)
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}
	if !strings.Contains(content, "theme = 'dark'") {
		t.Errorf("Expected content to contain 'theme = dark', got %s", content)
	}

	// Modify
	newContent := "theme = 'light'\nverbose = true\n"
	err = ds.SaveConfig(configPath, newContent)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// 3. Verify file on disk (simulating CLI reading it)
	//nolint:gosec // G304: configPath is constructed from t.TempDir().
	savedContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read file from disk: %v", err)
	}
	if string(savedContent) != newContent {
		t.Errorf("File on disk does not match expected content.\nExpected:\n%s\nGot:\n%s", newContent, string(savedContent))
	}
}

func createFile(t *testing.T, path, content string) {
	err := os.MkdirAll(filepath.Dir(path), 0o750)
	if err != nil {
		t.Fatalf("Failed to create directory for %s: %v", path, err)
	}
	err = os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
}
