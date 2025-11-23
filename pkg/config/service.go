package config

import (
	"fmt"
	"os"
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
	ds.RegisterProvider(&ClaudeProvider{})
	ds.RegisterProvider(&JulesProvider{})
	ds.RegisterProvider(&GeminiProvider{})
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
