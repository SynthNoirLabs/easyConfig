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
		providers: []Provider{},
	}
	// Providers will be registered here or via a Register method
	return ds
}

// RegisterProvider adds a new provider to the service
func (s *DiscoveryService) RegisterProvider(p Provider) {
	s.providers = append(s.providers, p)
}

// DiscoverAll iterates through all registered providers and collects configs
// projectPath: The root directory of the current project (optional)
func (s *DiscoveryService) DiscoverAll(projectPath string) ([]ConfigItem, error) {
	var allConfigs []ConfigItem

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

// ReadConfig reads the content of a configuration file at the given path
func (s *DiscoveryService) ReadConfig(path string) (string, error) {
	// Check if file exists
	if !FileExists(path) {
		return "", fmt.Errorf("file not found: %s", path)
	}

	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// SaveConfig writes content to a configuration file at the given path
// If the file has a .json extension, it validates the JSON structure first
func (s *DiscoveryService) SaveConfig(path string, content string) error {
	// Validate JSON if the file is a .json file
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".json" {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(content), &js); err != nil {
			return fmt.Errorf("invalid JSON content: %w", err)
		}
	}

	// Write the content to the file
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
