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
}

// ProviderStatusReport provides a detailed health check of a single provider.
type ProviderStatusReport struct {
	ProviderName string `json:"providerName"`
	Installed    bool   `json:"installed"`  // Tool binary found
	Configured   bool   `json:"configured"` // Config file exists
	Valid        bool   `json:"valid"`      // Config passes validation
	Message      string `json:"message"`    // Status description
	Version      string `json:"version"`    // Tool version if available
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
	// BinaryName returns the name of the tool's binary.
	BinaryName() string
	// VersionArgs returns the arguments to pass to the binary to get its version.
	VersionArgs() []string
}
