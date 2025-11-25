package main

import (
	"context"
	"fmt"

	"easyConfig/pkg/config"
	"easyConfig/pkg/watcher"
)

// App struct
type App struct {
	ctx              context.Context
	discoveryService *config.DiscoveryService
	watcherService   *watcher.Service
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
