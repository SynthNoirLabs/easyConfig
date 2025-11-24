package main

import (
	"context"
	"fmt"

	"easyConfig/pkg/config"
)

// App struct
type App struct {
	ctx              context.Context
	discoveryService *config.DiscoveryService
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
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// DiscoverConfigs returns all the discovered configurations
func (a *App) DiscoverConfigs(projectPath string) ([]config.ConfigItem, error) {
	return a.discoveryService.DiscoverAll(projectPath)
}

// ReadConfig reads the content of a configuration file
func (a *App) ReadConfig(path string) (string, error) {
	return a.discoveryService.ReadConfig(path)
}

// SaveConfig saves content to a configuration file
func (a *App) SaveConfig(path string, content string) error {
	return a.discoveryService.SaveConfig(path, content)
}
