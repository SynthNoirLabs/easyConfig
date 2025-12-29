# Provider Development Guide

This guide provides a comprehensive overview of how to add new AI CLI tool providers to EasyConfig.

## Table of Contents

- [Provider Interface Overview](#provider-interface-overview)
- [Creating a New Provider](#creating-a-new-provider)
- [Discovery Implementation](#discovery-implementation)
- [Create Implementation](#create-implementation)
- [Status Check Implementation](#status-check-implementation)
- [Testing Your Provider](#testing-your-provider)
- [Common Patterns & Best Practices](#common-patterns--best-practices)
- [Provider Checklist](#provider-checklist)
- [Existing Providers](#existing-providers)

## Provider Interface Overview

The `Provider` interface is the core of EasyConfig's discovery mechanism. It defines a set of methods that each provider must implement to be registered with the service.

```go
// pkg/config/types.go
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
}
```

## Creating a New Provider

To create a new provider, you'll need to create a new file in the `pkg/config` directory. The file should be named `provider_<provider_name>.go`. For example, a provider for a tool called "MyTool" would be in a file named `provider_mytool.go`.

Here's a basic template for a new provider:

```go
package config

import (
	"fmt"
	"path/filepath"
	"time"

	"easyConfig/pkg/util/paths"
)

type MyToolProvider struct{}

func (p *MyToolProvider) Name() string {
	return "MyTool"
}

func (p *MyToolProvider) Discover(projectPath string) ([]Item, error) {
	// Implementation details go here
	return nil, nil
}

func (p *MyToolProvider) Create(scope Scope, projectPath string) (string, error) {
	// Implementation details go here
	return "", nil
}

func (p *MyToolProvider) CheckStatus() ProviderStatus {
	// Implementation details go here
	return ProviderStatus{
		ProviderName: p.Name(),
		Health:       StatusUnknown,
		LastChecked:  time.Now().Format(time.RFC3339),
	}
}
```

Once you've created your provider, you'll need to register it in the `NewDiscoveryService` function in `pkg/config/service.go`.

```go
// pkg/config/service.go
func NewDiscoveryService(logger *slog.Logger) *DiscoveryService {
	// ...
	ds := &DiscoveryService{
		logger: logger,
		providers: []Provider{
			// ...
			&MyToolProvider{},
		},
	}
	return ds
}
```

## Discovery Implementation

The `Discover` method is responsible for finding all the configuration files for a given provider. It should return a slice of `Item` structs, where each struct represents a discovered configuration file.

### Discovery Patterns

Here's an example of how to implement the `Discover` method:

```go
func (p *MyToolProvider) Discover(projectPath string) ([]Item, error) {
	var items []Item
	home := paths.GetHomeDir()

	// Global scope
	if home != "" {
		globalPath := filepath.Join(home, ".mytool", "config.json")
		if FileExists(globalPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Global Config",
				FileName: "config.json",
				Path:     globalPath,
				Scope:    ScopeGlobal,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	// Project scope
	if projectPath != "" {
		projectPath := filepath.Join(projectPath, ".mytool.json")
		if FileExists(projectPath) {
			items = append(items, Item{
				Provider: p.Name(),
				Name:     "Project Config",
				FileName: ".mytool.json",
				Path:     projectPath,
				Scope:    ScopeProject,
				Format:   FormatJSON,
				Exists:   true,
			})
		}
	}

	return items, nil
}
```

### Multi-Platform Paths

When discovering configuration files, it's important to use cross-platform paths. The `pkg/util/paths` package provides a set of helper functions for this purpose.

- `paths.GetHomeDir()`: Returns the user's home directory.
- `paths.GetConfigDir(appName)`: Returns the OS-specific default configuration directory for an application.

The `FileExists` helper function is also available in the `config` package to check if a file exists.

## Create Implementation

The `Create` method is responsible for creating a new configuration file for a given provider. It should return the path of the created file, or an error if the file could not be created.

```go
func (p *MyToolProvider) Create(scope Scope, projectPath string) (string, error) {
	defaultContent := "{}"
	var path string

	switch scope {
	case ScopeGlobal:
		home := paths.GetHomeDir()
		if home == "" {
			return "", fmt.Errorf("home directory not found")
		}
		path = filepath.Join(home, ".mytool", "config.json")
	case ScopeProject:
		if projectPath == "" {
			return "", fmt.Errorf("project path is required")
		}
		path = filepath.Join(projectPath, ".mytool.json")
	default:
		return "", fmt.Errorf("unsupported scope")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return "", fmt.Errorf("failed to create dir: %w", err)
	}
	if FileExists(path) {
		return "", fmt.Errorf("file exists: %s", path)
	}
	if err := os.WriteFile(path, []byte(defaultContent), 0o600); err != nil {
		return "", fmt.Errorf("failed to write: %w", err)
	}
	return path, nil
}
```

## Status Check Implementation

The `CheckStatus` method is responsible for performing a health check on the provider's configuration. It should return a `ProviderStatus` struct that includes the provider's health status, a status message, and a list of discovered files.

```go
func (p *MyToolProvider) CheckStatus() ProviderStatus {
	status := ProviderStatus{
		ProviderName: p.Name(),
		LastChecked:  time.Now().Format(time.RFC3339),
	}

	home := paths.GetHomeDir()
	if home == "" {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Home directory not found."
		return status
	}

	globalPath := filepath.Join(home, ".mytool", "config.json")
	files, _ := p.Discover("")
	status.DiscoveredFiles = files

	if !FileExists(globalPath) {
		status.Health = StatusUnhealthy
		status.StatusMessage = "Global config not found. Use 'Create' to set one up."
	} else {
		status.Health = StatusHealthy
		status.StatusMessage = "Configuration files found."
	}

	return status
}
```

## Testing Your Provider

It's important to add unit tests for your provider. The tests should cover the `Discover`, `Create`, and `CheckStatus` methods.

Here's an example of how to test the `Discover` method:

```go
func TestMyProviderDiscover(t *testing.T) {
	// Setup: Create a temporary directory and a dummy config file
	tmpDir, err := os.MkdirTemp("", "test-discover")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configDir := filepath.Join(tmpDir, ".mytool")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	configFile := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Override the home directory to point to our temp directory
	// This is a simplified example. In a real test, you might use a library
	// to mock the home directory.
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)


	// Execution: Call the Discover method
	provider := &MyToolProvider{}
	items, err := provider.Discover("")

	// Assertion: Check that the config file was discovered
	if err != nil {
		t.Errorf("Discover() returned an error: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, but got %d", len(items))
	}

	if items[0].Path != configFile {
		t.Errorf("Expected path %s, but got %s", configFile, items[0].Path)
	}
}
```

## Common Patterns & Best Practices

- **Graceful Handling of Missing Files**: Your provider should not return an error if a configuration file is not found. Instead, it should return an empty slice of `Item` structs.
- **Cross-Platform Paths**: Always use the `pkg/util/paths` package to construct file paths.
- **Unit Tests**: Add unit tests for your provider to ensure that it's working correctly.
- **Registration**: Don't forget to register your provider in the `NewDiscoveryService` function in `pkg/config/service.go`.
- **Documentation**: Add a link to your provider's documentation in the `README.md` file.

## Provider Checklist

- [ ] Implements all interface methods
- [ ] Handles missing files gracefully
- [ ] Uses cross-platform paths
- [ ] Has unit tests
- [ ] Added to `NewDiscoveryService()` in `pkg/config/service.go`
- [ ] Documented in `README.md`

## Existing Providers

- [Aider](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_aider.go)
- [AmazonQ](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_amazonq.go)
- [Claude](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_claude.go)
- [Codex](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_codex.go)
- [Continue](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_continue.go)
- [Copilot](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_copilot.go)
- [Crush](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_crush.go)
- [Cursor](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_cursor.go)
- [Gemini](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_gemini.go)
- [Git](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_git.go)
- [Goose](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_goose.go)
- [Jules](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_jules.go)
- [Mentat](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_mentat.go)
- [OpenAI](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_openai.go)
- [Opencode](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_opencode.go)
- [OpenHands](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_openhands.go)
- [Plandex](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_plandex.go)
- [Sweep](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_sweep.go)
- [Windsurf](https://github.com/komod0/easyConfig/blob/main/pkg/config/provider_windsurf.go)
