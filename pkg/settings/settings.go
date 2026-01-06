package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"easyConfig/pkg/util/paths"
)

const (
	settingsFile = "easyconfig.json"
)

// Settings defines the application's configurable settings
type Settings struct {
	ProviderScanDirs []string `json:"providerScanDirs"`
}

// Service manages the application settings
type Service struct {
	settings *Settings
	mu       sync.RWMutex
	path     string
}

// NewService creates a new settings service
func NewService() (*Service, error) {
	configDir := paths.GetConfigDir("EasyConfig")
	if configDir == "" {
		return nil, os.ErrNotExist
	}
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, settingsFile)
	s := &Service{
		path:     path,
		settings: &Settings{},
	}

	if err := s.load(); err != nil {
		// If the file doesn't exist, we can ignore the error
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return s, nil
}

// Get returns the current settings
func (s *Service) Get() *Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return a copy to prevent modification of the internal state
	settingsCopy := *s.settings
	return &settingsCopy
}

// Save persists the settings to disk
func (s *Service) Save(settings *Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.path, data, 0600); err != nil {
		return err
	}

	s.settings = settings
	return nil
}

// load reads the settings from disk
func (s *Service) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s.settings)
}
