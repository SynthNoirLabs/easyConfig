package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ServerConfig represents the configuration for a single MCP server
type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// Injector handles injecting MCP server configurations
type Injector struct{}

// NewInjector creates a new Injector
func NewInjector() *Injector {
	return &Injector{}
}

// Inject adds or updates an MCP server configuration in the target file
func (i *Injector) Inject(configPath string, serverName string, config ServerConfig) error {
	// 1. Read the file
	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, create a new one with just this server
			return i.createConfigFile(configPath, serverName, config)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 2. Parse JSON
	// We use map[string]interface{} to preserve other fields we don't know about
	var root map[string]interface{}
	if err := json.Unmarshal(content, &root); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 3. Navigate/Create mcpServers
	mcpServersRaw, ok := root["mcpServers"]
	var mcpServers map[string]interface{}

	if !ok || mcpServersRaw == nil {
		mcpServers = make(map[string]interface{})
		root["mcpServers"] = mcpServers
	} else {
		mcpServers, ok = mcpServersRaw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("mcpServers field is not an object")
		}
	}

	// 4. Add/Update the server entry
	// Convert ServerConfig to map to ensure clean JSON structure matching the file's style
	// (though struct marshaling works, inserting struct into map[string]interface{} is fine)
	mcpServers[serverName] = config

	// Update the root map
	root["mcpServers"] = mcpServers

	// 5. Write back
	return i.writeConfigFile(configPath, root)
}

func (i *Injector) createConfigFile(path string, serverName string, config ServerConfig) error {
	root := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			serverName: config,
		},
	}
	return i.writeConfigFile(path, root)
}

func (i *Injector) writeConfigFile(path string, data map[string]interface{}) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal with indentation
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, bytes, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
