package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"easyConfig/pkg/config"
	"easyConfig/pkg/marketplaces"
	"easyConfig/pkg/schema"
	"easyConfig/pkg/util/paths"
	"easyConfig/pkg/watcher"
)

// App struct
type App struct {
	ctx              context.Context
	discoveryService *config.DiscoveryService
	watcherService   *watcher.Service
	smitheryClient   *marketplaces.SmitheryClient
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Initialize logger (slog default is fine for now, or we can configure it)
	logger := slog.Default()
	a.discoveryService = config.NewDiscoveryService(logger)
	a.watcherService = watcher.NewService()
	a.smitheryClient = marketplaces.NewSmitheryClient()
	if a.watcherService != nil {
		a.watcherService.Start(ctx)
	}
}

// shutdown is called at application termination
func (a *App) shutdown(_ context.Context) {
	if a.watcherService != nil {
		a.watcherService.Close()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// DiscoverConfigs returns all the discovered configurations
func (a *App) DiscoverConfigs(projectPath string) ([]config.Item, error) {
	items, err := a.discoveryService.DiscoverAll(projectPath)
	if err != nil {
		return nil, err
	}

	// Watch all discovered files
	if a.watcherService != nil {
		for _, item := range items {
			if item.Exists {
				_ = a.watcherService.Add(item.Path)
			}
		}
	}

	return items, nil
}

// ReadConfig reads the content of a configuration file
func (a *App) ReadConfig(path string) (string, error) {
	return a.discoveryService.ReadConfig(path)
}

// SaveConfig saves content to a configuration file
func (a *App) SaveConfig(path, content string) error {
	return a.discoveryService.SaveConfig(path, content)
}

// DeleteConfig deletes a configuration file
func (a *App) DeleteConfig(path string) error {
	// Stop watching before delete to avoid error logs
	if a.watcherService != nil {
		_ = a.watcherService.Remove(path)
	}
	return a.discoveryService.DeleteConfig(path)
}

// CreateConfig creates a new configuration file
func (a *App) CreateConfig(providerName, scope, projectPath string) (string, error) {
	// Convert string scope to config.Scope
	var cfgScope config.Scope
	switch scope {
	case "global":
		cfgScope = config.ScopeGlobal
	case "project":
		cfgScope = config.ScopeProject
	default:
		return "", fmt.Errorf("invalid scope: %s", scope)
	}
	return a.discoveryService.CreateConfig(providerName, cfgScope, projectPath)
}

// FetchSchemas downloads the latest configuration schemas for supported tools
func (a *App) FetchSchemas() error {
	// Use easyConfig's own config directory to store schemas
	configDir := paths.GetConfigDir("easyConfig")
	if configDir == "" {
		// Fallback to local directory if standard path fails
		configDir = "."
	}
	schemaDir := filepath.Join(configDir, "schemas")

	fetcher := schema.NewFetcher()
	return fetcher.FetchAllSchemas(schemaDir)
}

// FetchPopularServers returns a list of popular MCP servers from Smithery
func (a *App) FetchPopularServers() ([]marketplaces.MCPPackage, error) {
	if a.smitheryClient == nil {
		return nil, fmt.Errorf("smithery client not initialized")
	}
	return a.smitheryClient.FetchPopularServers()
}

// InstallMCPPackage installs an MCP server by creating a configuration file
func (a *App) InstallMCPPackage(pkg marketplaces.MCPPackage) error {
	// For now, we'll create a JSON config file in the easyConfig directory
	// In a real scenario, this might involve `npm install` or `pip install`
	// Here we just create a config file that references the server.

	configDir := paths.GetConfigDir("easyConfig")
	if configDir == "" {
		return fmt.Errorf("failed to get config directory")
	}

	// Create mcp-servers directory if it doesn't exist
	mcpDir := filepath.Join(configDir, "mcp-servers")
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return fmt.Errorf("failed to create mcp-servers directory: %w", err)
	}

	filename := fmt.Sprintf("%s.json", pkg.Name)
	filePath := filepath.Join(mcpDir, filename)

	// Create a simple config structure for the MCP server
	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			pkg.Name: map[string]interface{}{
				"command": "npx", // Assumption for now, or use pkg metadata if available
				"args":    []string{"-y", pkg.Name},
				"url":     pkg.URL,
				"version": pkg.Version,
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
