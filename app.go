package main

import (
	"context"
	"fmt"
	"path/filepath"

	"easyConfig/pkg/config"
	"easyConfig/pkg/install"
	"easyConfig/pkg/schema"
	"easyConfig/pkg/util/paths"
	"easyConfig/pkg/watcher"
)

// App struct
type App struct {
	ctx              context.Context
	discoveryService *config.DiscoveryService
	watcherService   *watcher.Service
	installer        *install.Installer
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.discoveryService = config.NewDiscoveryService()
	a.watcherService = watcher.NewService()
	a.installer = install.NewInstaller()
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

// InstallMCPPackage installs an MCP server package
func (a *App) InstallMCPPackage(packageName string) error {
	// Get user's home directory to store the config
	homeDir := paths.GetHomeDir()
	if homeDir == "" {
		return fmt.Errorf("could not determine home directory")
	}

	// Store MCP configs in ~/.config/easyConfig/mcp/
	configDir := filepath.Join(homeDir, ".config", "easyConfig", "mcp")

	// Install the package (with verification)
	return a.installer.InstallPackage(a.ctx, packageName, configDir)
}
