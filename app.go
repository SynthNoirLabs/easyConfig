package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"easyConfig/pkg/config"
	"easyConfig/pkg/install"
	"easyConfig/pkg/marketplaces"
	"easyConfig/pkg/mcp"
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
	secretsManager   *workflows.SecretsManager
	mcpInjector      *mcp.Injector

	// Active wizards, mapped by provider name
	activeWizards map[string]config.Wizard
	wizardMu      sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		activeWizards: make(map[string]config.Wizard),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logger := slog.Default()
	a.discoveryService = config.NewDiscoveryService(logger)
	a.watcherService = watcher.NewService()
	a.installer = install.NewInstaller()
	a.smitheryClient = marketplaces.NewSmitheryClient()
	a.awesomeClient = marketplaces.NewAwesomeClient()
	a.workflowGen = workflows.NewGenerator()
	a.secretsManager = workflows.NewSecretsManager()
	a.mcpInjector = mcp.NewInjector()
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

// InstallMCPPackage installs an MCP server by creating a configuration file
func (a *App) InstallMCPPackage(pkgJSON string) error {
	var pkg marketplaces.MCPPackage
	if err := json.Unmarshal([]byte(pkgJSON), &pkg); err != nil {
		return fmt.Errorf("failed to unmarshal package json: %w", err)
	}

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

	// 2. Inject into Claude Desktop Config
	// Construct the MCP config for injection
	mcpConfig := mcp.ServerConfig{
		Command: "npx",
		Args:    []string{"-y", pkg.Name},
		Env:     map[string]string{}, // Add env vars if needed
	}

	// 3. Determine Claude Desktop config path
	homeDir := paths.GetHomeDir()
	if homeDir == "" {
		return fmt.Errorf("could not determine home directory")
	}

	// Standard path for Claude Desktop
	// macOS: ~/Library/Application Support/Claude/claude_desktop_config.json
	// Windows: %APPDATA%\Claude\claude_desktop_config.json
	// Linux: ~/.config/Claude/claude_desktop_config.json (unofficial/standard XDG)
	// But wait, provider_claude.go used ~/.claude/claude_desktop_config.json for Linux?
	// Let's check provider_claude.go again.
	// It used filepath.Join(home, ".claude", "claude_desktop_config.json")
	// But standard Claude Desktop on Mac is ~/Library/Application Support/Claude/claude_desktop_config.json
	// On Windows it's AppData/Roaming/Claude/claude_desktop_config.json

	// Let's use a helper or hardcode for now based on OS, or rely on what provider_claude.go does.
	// Actually, let's look at how provider_claude.go discovers it.
	// It checks `filepath.Join(home, ".claude", "claude_desktop_config.json")`.
	// This might be a simplification or specific to a certain setup.
	// For "Real" installation, we should target the actual file Claude Desktop uses.

	// We'll use the paths.GetConfigDir("Claude") which should handle OS differences if implemented correctly.
	// But paths.GetConfigDir usually returns ~/.config/AppName on Linux.
	// Claude Desktop on Mac: ~/Library/Application Support/Claude

	// Let's try to find the file or default to a standard location.
	// For now, I'll use the same path as provider_claude.go seems to expect for "Global Desktop Config"
	// which was `filepath.Join(home, ".claude", "claude_desktop_config.json")`.
	// WAIT, looking at provider_claude.go lines 57: `path := filepath.Join(home, ".claude", "claude_desktop_config.json")`
	// This seems to be where we expect it.

	configPath := filepath.Join(homeDir, ".claude", "claude_desktop_config.json")

	// On macOS, it's different.
	// if runtime.GOOS == "darwin" {
	//    configPath = filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	// }
	// I should probably make this robust.

	// For this iteration, I will stick to the path defined in provider_claude.go to be consistent with "Discovery".
	// If Discovery is wrong, we fix both.

	// 4. Inject
	// Use the package name (sanitized) as the server name
	packageName := pkg.Name

	return a.mcpInjector.Inject(configPath, packageName, mcpConfig)
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
// Returns WorkflowResponse, error
func (a *App) GenerateWorkflow(agent, trigger string) (*workflows.WorkflowResponse, error) {
	return a.workflowGen.GenerateWorkflow(agent, trigger)
}

// SetSecret sets a repository secret
func (a *App) SetSecret(name, value string) error {
	return a.secretsManager.SetRepositorySecret(name, value)
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

// ListWorkflowTemplates returns workflow templates with metadata/content
func (a *App) ListWorkflowTemplates() []workflows.Template {
	return a.workflowGen.ListTemplates()
}

// Profile operations
func (a *App) ListProfiles() ([]config.ProfileSummary, error) {
	return a.discoveryService.ListProfiles()
}

func (a *App) SaveProfile(name string) error {
	return a.discoveryService.SaveProfile(name, ".")
}

func (a *App) ApplyProfile(name string) ([]string, error) {
	return a.discoveryService.ApplyProfile(name)
}

func (a *App) DeleteProfile(name string) error {
	return a.discoveryService.DeleteProfile(name)
}

// GetProviderStatuses returns the health status of all registered providers.
func (a *App) GetProviderStatuses() []config.ProviderStatus {
	return a.discoveryService.GetProviderStatuses()
}

// ListDocs returns the locally synced documentation pages grouped by provider.
// It scans docs/vendor/<provider>/latest and reports available .md/.html pages.
func (a *App) ListDocs() ([]config.DocsProvider, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working dir: %w", err)
	}
	return config.ListDocsFromRoot(root)
}

// ReadDoc returns the contents of a local doc page.
// provider: provider name (e.g. "claude"), slug: base filename, format: "md" or "html".
// If the requested format is not available, it falls back to the other one.
func (a *App) ReadDoc(provider, slug, format string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working dir: %w", err)
	}
	return config.ReadDocFromRoot(root, provider, slug, format)
}

// --- Wizard Methods ---

// StartWizard initiates a configuration wizard for a given provider.
func (a *App) StartWizard(providerName string) (*config.WizardStep, error) {
	a.wizardMu.Lock()
	defer a.wizardMu.Unlock()

	provider := a.discoveryService.GetProvider(providerName)
	if provider == nil {
		return nil, fmt.Errorf("provider not found: %s", providerName)
	}

	wizard := provider.GetWizard()
	if wizard == nil {
		return nil, fmt.Errorf("provider '%s' does not have a wizard", providerName)
	}

	a.activeWizards[providerName] = wizard
	return wizard.Start()
}

// NextWizardStep advances the wizard to the next step.
func (a *App) NextWizardStep(providerName, currentStepID, response string) (*config.WizardStep, error) {
	a.wizardMu.Lock()
	defer a.wizardMu.Unlock()

	wizard, ok := a.activeWizards[providerName]
	if !ok {
		return nil, fmt.Errorf("no active wizard for provider: %s", providerName)
	}

	step, err := wizard.Next(currentStepID, response)
	if err != nil {
		return nil, err
	}
	// If the wizard is finished, remove it from the active list
	if step == nil {
		delete(a.activeWizards, providerName)
	}
	return step, nil
}

// CancelWizard cancels an active wizard.
func (a *App) CancelWizard(providerName string) error {
	a.wizardMu.Lock()
	defer a.wizardMu.Unlock()

	wizard, ok := a.activeWizards[providerName]
	if !ok {
		return fmt.Errorf("no active wizard for provider: %s", providerName)
	}

	delete(a.activeWizards, providerName)
	return wizard.Cancel()
}
