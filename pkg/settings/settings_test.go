package settings

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSettingsService(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "settings-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override the config dir to use the temp dir
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalConfigDir)

	s, err := NewService()
	if err != nil {
		t.Fatalf("Failed to create new settings service: %v", err)
	}

	// 1. Test initial state (should be empty)
	if len(s.Get().ProviderScanDirs) != 0 {
		t.Errorf("Expected initial ProviderScanDirs to be empty, but got %v", s.Get().ProviderScanDirs)
	}

	// 2. Test saving and getting settings
	newSettings := &Settings{
		ProviderScanDirs: []string{"/path/to/providers1", "/path/to/providers2"},
	}
	if err := s.Save(newSettings); err != nil {
		t.Fatalf("Failed to save settings: %v", err)
	}

	currentSettings := s.Get()
	if !reflect.DeepEqual(currentSettings, newSettings) {
		t.Errorf("Expected settings to be %v, but got %v", newSettings, currentSettings)
	}

	// 3. Test that the file was created
	expectedPath := filepath.Join(tempDir, "EasyConfig", "easyconfig.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Settings file was not created at %s", expectedPath)
	}

	// 4. Test that a new service loads the existing settings
	s2, err := NewService()
	if err != nil {
		t.Fatalf("Failed to create second settings service: %v", err)
	}
	if !reflect.DeepEqual(s2.Get(), newSettings) {
		t.Errorf("Expected second service to load settings %v, but got %v", newSettings, s2.Get())
	}
}
