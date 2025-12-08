package config

// Scope defines where a configuration file is located
type Scope string

const (
	ScopeGlobal  Scope = "global"
	ScopeProject Scope = "project"
	ScopeSystem  Scope = "system"
)

// Format defines the file format
type Format string

const (
	FormatJSON Format = "json"
	FormatTOML Format = "toml"
	FormatYAML Format = "yaml"
	FormatTXT  Format = "text"
	FormatMD   Format = "markdown"
	FormatINI  Format = "ini"
)

// Item represents a discovered configuration file
type Item struct {
	Provider string `json:"provider"` // e.g., "Claude", "Git"
	Name     string `json:"name"`     // Display name e.g. "Global Config"
	FileName string `json:"fileName"` // Actual filename e.g. "config.json"
	Path     string `json:"path"`     // Absolute path
	Scope    Scope  `json:"scope"`
	Format   Format `json:"format"`
	Exists   bool   `json:"exists"`
}

// HealthStatus defines the overall health of a provider's configuration
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
	StatusUnknown   HealthStatus = "unknown"
)

// ProviderStatus represents the health and configuration status of a provider
type ProviderStatus struct {
	ProviderName    string       `json:"providerName"`
	Health          HealthStatus `json:"health"`
	StatusMessage   string       `json:"statusMessage,omitempty"`
	DiscoveredFiles []Item       `json:"discoveredFiles,omitempty"`
	LastChecked     string       `json:"lastChecked"` // ISO 8601 format
	HasWizard       bool         `json:"hasWizard"`
}

// WizardStep represents a single step in a configuration wizard.
type WizardStep struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	// More fields like question type, options, etc. can be added here.
}

// Wizard defines the interface for a multi-step, interactive configuration wizard.
type Wizard interface {
	// Start begins the wizard and returns the first step.
	Start() (*WizardStep, error)
	// Next takes the current step ID and the user's response, and returns the next step.
	// If the wizard is finished, it returns a nil step.
	Next(currentStepID, response string) (*WizardStep, error)
	// Cancel aborts the wizard.
	Cancel() error
}

// Provider defines the interface for a tool configuration provider
type Provider interface {
	// Name returns the unique name of the provider (e.g. "Claude Code")
	Name() string
	// Discover searches for configs. projectPath can be empty if no project is open.
	Discover(projectPath string) ([]Item, error)
	// Create generates a new default configuration file for the given scope
	// scope: "global" or "project"
	// projectPath: required if scope is "project"
	// Returns the path of the created file or error
	Create(scope Scope, projectPath string) (string, error)
	// CheckStatus performs a health check on the provider's configuration
	CheckStatus() ProviderStatus
	// GetWizard returns the configuration wizard for this provider, if available.
	GetWizard() Wizard
}
