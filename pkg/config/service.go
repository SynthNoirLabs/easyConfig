package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"easyConfig/pkg/settings"
	toml "github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// DiscoveryService manages the discovery of configurations across multiple providers
type DiscoveryService struct {
	providers       []Provider
	logger          *slog.Logger
	settingsService *settings.Service
	mu              sync.RWMutex
}

// NewDiscoveryService creates a new service with default providers
func NewDiscoveryService(logger *slog.Logger, settingsService *settings.Service) *DiscoveryService {
	if logger == nil {
		logger = slog.Default()
	}
	ds := &DiscoveryService{
		logger:          logger,
		settingsService: settingsService,
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
			&AiderProvider{},
			&GooseProvider{},
			&AmazonQProvider{},
			&CursorProvider{},
			&WindsurfProvider{},
			&ContinueProvider{},
			&MentatProvider{},
			&SweepProvider{},
			&PlandexProvider{},
			&OpenHandsProvider{},
		},
	}
	ds.loadDynamicProviders()
	return ds
}

// loadDynamicProviders scans the directories specified in the settings and loads any dynamic providers found.
func (s *DiscoveryService) loadDynamicProviders() {
	if s.settingsService == nil {
		return
	}
	cfg := s.settingsService.Get()
	for _, dir := range cfg.ProviderScanDirs {
		s.logger.Info("Scanning for dynamic providers", "directory", dir)
		files, err := os.ReadDir(dir)
		if err != nil {
			s.logger.Error("Failed to read provider scan directory", "directory", dir, "error", err)
			continue
		}

		for _, file := range files {
			if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
				defPath := filepath.Join(dir, file.Name())
				s.logger.Info("Found potential provider definition", "path", defPath)
				provider, err := NewDynamicProvider(defPath)
				if err != nil {
					s.logger.Error("Failed to load dynamic provider", "path", defPath, "error", err)
					continue
				}
				s.logger.Info("Registering dynamic provider", "name", provider.Name())
				s.RegisterProvider(provider)
			}
		}
	}
}

// RegisterProvider adds a new provider to the service
func (s *DiscoveryService) RegisterProvider(p Provider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers = append(s.providers, p)
}

// DiscoverAll iterates through all registered providers and collects configs in parallel.
// It respects the provided context for cancellation.
// projectPath: The root directory of the current project (optional)
func (s *DiscoveryService) DiscoverAll(ctx context.Context, projectPath string) ([]Item, error) {
	s.mu.RLock()
	// Make a copy of the slice to avoid holding the lock during the entire operation
	providers := make([]Provider, len(s.providers))
	copy(providers, s.providers)
	s.mu.RUnlock()

	var wg sync.WaitGroup
	resultsChan := make(chan []Item, len(providers))

	for _, p := range providers {
		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()

			// Early exit if context is cancelled
			select {
			case <-ctx.Done():
				return
			default:
			}

			items, err := p.Discover(projectPath)
			if err != nil {
				s.logger.Error("Error discovering for provider", "provider", p.Name(), "error", err)
				resultsChan <- nil // Send nil to signal completion with error
				return
			}

			// Check context again before sending result to avoid blocking
			select {
			case resultsChan <- items:
			case <-ctx.Done():
			}
		}(p)
	}

	// Goroutine to wait for all providers and then close the channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var allConfigs []Item
	for items := range resultsChan {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if items != nil {
			allConfigs = append(allConfigs, items...)
		}
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return allConfigs, nil
}

// CreateConfig finds the provider and creates a new config file for the given scope
func (s *DiscoveryService) CreateConfig(providerName string, scope Scope, projectPath string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
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
	//nolint:gosec // G304: Path is user-provided as intended feature
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
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		var js json.RawMessage
		if err := json.Unmarshal([]byte(content), &js); err != nil {
			return fmt.Errorf("invalid JSON content: %w", err)
		}
	case ".yaml", ".yml":
		var y any
		if err := yaml.Unmarshal([]byte(content), &y); err != nil {
			return fmt.Errorf("invalid YAML content: %w", err)
		}
	case ".toml":
		var t any
		if err := toml.Unmarshal([]byte(content), &t); err != nil {
			return fmt.Errorf("invalid TOML content: %w", err)
		}
	}

	// Write the content to the file with restrictive permissions (0600).
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// GetProviderStatuses iterates through all registered providers and returns their health status.
func (s *DiscoveryService) GetProviderStatuses() []ProviderStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var statuses []ProviderStatus

	for _, p := range s.providers {
		status := p.CheckStatus()
		statuses = append(statuses, status)
	}

	return statuses
}
