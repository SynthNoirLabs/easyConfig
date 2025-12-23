package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInjector_Inject(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create initial config
	initialConfig := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"existing": map[string]interface{}{
				"command": "echo",
			},
		},
	}
	data, _ := json.Marshal(initialConfig)
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	injector := NewInjector()

	// Test Inject
	newServer := ServerConfig{
		Command: "node",
		Args:    []string{"server.js"},
	}

	err := injector.Inject(configPath, "new-server", newServer)
	if err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	// Verify injection
	//nolint:gosec // G304: configPath is constructed from t.TempDir().
	content, _ := os.ReadFile(configPath)
	var config map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	servers := config["mcpServers"].(map[string]interface{})
	if _, ok := servers["new-server"]; !ok {
		t.Error("New server not found in config")
	}
	if _, ok := servers["existing"]; !ok {
		t.Error("Existing server removed from config")
	}
}

func TestInjector_Inject_NewFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "new-config.json")

	injector := NewInjector()

	newServer := ServerConfig{
		Command: "node",
	}

	err := injector.Inject(configPath, "new-server", newServer)
	if err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file not created")
	}
}
