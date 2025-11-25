package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveryService manages the discovery of configurations across multiple providers
type DiscoveryService struct {
	providers []Provider
}

// NewDiscoveryService creates a new service with default providers
func NewDiscoveryService() *DiscoveryService {
	ds := &DiscoveryService{
		providers: []Provider{
			&ClaudeProvider{},
			&JulesProvider{},
			&GeminiProvider{},
			&CopilotProvider{},
			&OpenAIProvider{},
			&CodexProvider{},
			&OpenCodeProvider{},
			&CrushProvider{},
			&GitProvider{},
		},
	}
	return ds
}

// RegisterProvider adds a new provider to the service
func (s *DiscoveryService) RegisterProvider(p Provider) {
	s.providers = append(s.providers, p)
}

// DiscoverAll iterates through all registered providers and collects configs
// projectPath: The root directory of the current project (optional)
func (s *DiscoveryService) DiscoverAll(projectPath string) ([]Item, error) {
	var allConfigs []Item

	for _, p := range s.providers {
		items, err := p.Discover(projectPath)
		if err != nil {
			// We log error but continue to next provider
			fmt.Printf("Error discovering for provider %s: %v\n", p.Name(), err)
			continue
		}
		allConfigs = append(allConfigs, items...)
	}

	return allConfigs, nil
}

// CreateConfig finds the provider and creates a new config file for the given scope
func (s *DiscoveryService) CreateConfig(providerName string, scope Scope, projectPath string) (string, error) {
	for _, p := range s.providers {
		if p.Name() == providerName {
			return p.Create(scope, projectPath)
		}
	}
	return "", fmt.Errorf("provider not found: %s", providerName)
}

// DeleteConfig removes a configuration file from disk
func (s *DiscoveryService) DeleteConfig(path string) error {
	// Verify it's a file and exists
	if !FileExists(path) {
		return fmt.Errorf("file not found or is a directory: %s", path)
	}

	// Remove the file
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Helper: GetUserHome returns the user's home directory safely
func GetUserHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// Helper: FileExists checks if a file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadConfig reads the content of a configuration file at the given path.
//
// Note: ReadConfig does not currently make use of any state on the
// DiscoveryService; it is defined as a method for potential future
// extensibility (for example, reading encrypted configs tied to a provider).
// It performs the read directly and returns a user-friendly error when the
// file cannot be found.
func (s *DiscoveryService) ReadConfig(path string) (string, error) {
	// Read the file content directly. Avoid a separate existence check to
	// prevent TOCTOU race conditions; rely on os.ReadFile and check the error.
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %s", path)
		}
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// SaveConfig writes content to a configuration file at the given path.
//
// If the file has a `.json` extension, it validates the JSON structure first.
// The method does not currently depend on any DiscoveryService state but is
// defined on the type for symmetry with ReadConfig and future extension.
//
// When writing files, SaveConfig uses a restrictive permission mode (0600)
// to avoid creating world-readable configuration files that may contain
// secrets.
func (s *DiscoveryService) SaveConfig(path, content string) error {
	// Validate JSON if the file is a .json file
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".json" {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(content), &js); err != nil {
			return fmt.Errorf("invalid JSON content: %w", err)
		}
	}

	// Write the content to the file with restrictive permissions (0600).
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
