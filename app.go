package main

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"easyConfig/pkg/config"
	"easyConfig/pkg/install"
	"easyConfig/pkg/marketplaces"
	"easyConfig/pkg/schema"
	"easyConfig/pkg/util/paths"
	"easyConfig/pkg/watcher"
	"easyConfig/pkg/workflows"
)

// App struct
type App struct {
	ctx              context.Context
	discoveryService *config.DiscoveryService
	watcherService   *watcher.Service
	installer        *install.Installer
	smitheryClient   *marketplaces.SmitheryClient
	awesomeClient    *marketplaces.AwesomeClient
	workflowGen      *workflows.Generator
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
	a.smitheryClient = marketplaces.NewSmitheryClient()
	a.awesomeClient = marketplaces.NewAwesomeClient()
	a.workflowGen = workflows.NewGenerator()
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

// FetchPopularServers fetches popular MCP servers from Smithery and Awesome lists
func (a *App) FetchPopularServers() ([]marketplaces.MCPPackage, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allPackages []marketplaces.MCPPackage
	var errors []error

	// Fetch from Smithery
	wg.Add(1)
	go func() {
		defer wg.Done()
		pkgs, err := a.smitheryClient.FetchPopularServers()
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errors = append(errors, fmt.Errorf("smithery error: %w", err))
		} else {
			allPackages = append(allPackages, pkgs...)
		}
	}()

	// Fetch from Awesome Lists
	wg.Add(1)
	go func() {
		defer wg.Done()
		pkgs, err := a.awesomeClient.FetchServers()
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errors = append(errors, fmt.Errorf("awesome list error: %w", err))
		} else {
			allPackages = append(allPackages, pkgs...)
		}
	}()

	wg.Wait()

	// Deduplicate packages based on name
	seen := make(map[string]bool)
	uniquePackages := []marketplaces.MCPPackage{}
	for _, pkg := range allPackages {
		if !seen[pkg.Name] {
			seen[pkg.Name] = true
			uniquePackages = append(uniquePackages, pkg)
		}
	}

	// If we have at least some packages, return them even if one source failed
	if len(uniquePackages) > 0 {
		return uniquePackages, nil
	}

	// If everything failed, return combined error
	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to fetch servers: %v", errors)
	}

	return uniquePackages, nil
}

// GenerateWorkflow generates a GitHub Actions workflow content
func (a *App) GenerateWorkflow(agent, trigger string) (string, error) {
	return a.workflowGen.GenerateWorkflow(agent, trigger)
}

// SaveWorkflow saves the workflow content to .github/workflows/
func (a *App) SaveWorkflow(filename, content string) error {
	// Determine project root (assuming current working directory for now, or passed from frontend)
	// For this MVP, we'll use the current working directory or a specific project path if we had one in context.
	// Ideally, the frontend should pass the project path.
	// Let's assume the user wants to save it to the current directory where the app is running (or we could ask for a path).
	// However, `easyConfig` is often run *in* the project root.

	// Better approach: Use the path from DiscoveryService if available, or default to "."
	projectPath := "."

	workflowsDir := filepath.Join(projectPath, ".github", "workflows")
	if err := paths.EnsureDir(workflowsDir); err != nil {
		return fmt.Errorf("failed to create workflows directory: %w", err)
	}

	fullPath := filepath.Join(workflowsDir, filename)

	// Use DiscoveryService to save (it handles file writing)
	return a.discoveryService.SaveConfig(fullPath, content)
}

// GetSupportedWorkflows returns the list of supported workflows
func (a *App) GetSupportedWorkflows() []string {
	return a.workflowGen.GetSupportedWorkflows()
}
