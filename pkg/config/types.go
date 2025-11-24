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

// Provider defines the interface for a tool configuration provider
type Provider interface {
	// Name returns the unique name of the provider (e.g. "Claude Code")
	Name() string
	// Discover searches for configs. projectPath can be empty if no project is open.
	Discover(projectPath string) ([]Item, error)
}
