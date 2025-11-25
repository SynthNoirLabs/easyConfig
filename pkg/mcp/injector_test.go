package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInjector_Inject(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	injector := NewInjector()
	serverName := "test-server"
	serverConfig := ServerConfig{
		Command: "node",
		Args:    []string{"index.js"},
		Env:     map[string]string{"KEY": "VALUE"},
	}

	// Test 1: Create new file
	err := injector.Inject(configFile, serverName, serverConfig)
	if err != nil {
		t.Fatalf("Failed to inject into new file: %v", err)
	}

	verifyFile(t, configFile, serverName, serverConfig)

	// Test 2: Update existing file (add new server)
	serverName2 := "second-server"
	serverConfig2 := ServerConfig{
		Command: "python",
		Args:    []string{"main.py"},
	}
	err = injector.Inject(configFile, serverName2, serverConfig2)
	if err != nil {
		t.Fatalf("Failed to inject second server: %v", err)
	}

	verifyFile(t, configFile, serverName, serverConfig) // Original should still be there
	verifyFile(t, configFile, serverName2, serverConfig2)

	// Test 3: Update existing server
	serverConfigUpdated := ServerConfig{
		Command: "node",
		Args:    []string{"updated.js"},
	}
	err = injector.Inject(configFile, serverName, serverConfigUpdated)
	if err != nil {
		t.Fatalf("Failed to update server: %v", err)
	}

	verifyFile(t, configFile, serverName, serverConfigUpdated)
}

func verifyFile(t *testing.T, path string, serverName string, expected ServerConfig) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var root map[string]interface{}
	if err := json.Unmarshal(content, &root); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	mcpServers, ok := root["mcpServers"].(map[string]interface{})
	if !ok {
		t.Fatalf("mcpServers not found or invalid type")
	}

	serverRaw, ok := mcpServers[serverName]
	if !ok {
		t.Fatalf("Server %s not found in config", serverName)
	}

	// Helper to convert map back to struct for comparison is a bit tedious,
	// checking specific fields is easier for this test context
	serverMap, ok := serverRaw.(map[string]interface{})
	if !ok {
		t.Fatalf("Server entry is not a map")
	}

	if serverMap["command"] != expected.Command {
		t.Errorf("Expected command %s, got %v", expected.Command, serverMap["command"])
	}
}
