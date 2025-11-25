package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoveryService(t *testing.T) {
	// Setup
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to restore HOME: %v", err)
		}
	}()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := NewDiscoveryService(logger)

	// Test RegisterProvider
	// (Already registered in NewDiscoveryService, but let's check if we can add a dummy one if needed,
	// or just rely on default ones)
	if len(service.providers) == 0 {
		t.Error("Expected default providers to be registered")
	}

	// Test DiscoverAll
	// Create a dummy config to be discovered
	claudeDir := filepath.Join(tempHome, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write settings.json: %v", err)
	}

	items, err := service.DiscoverAll(tempHome) // Pass project path (using tempHome for simplicity as it has .claude)
	if err != nil {
		t.Errorf("DiscoverAll failed: %v", err)
	}
	if len(items) == 0 {
		t.Error("Expected to discover at least one item")
	}

	// Test ReadConfig
	// Use the path from the discovered item
	content, err := service.ReadConfig(items[0].Path)
	if err != nil {
		t.Errorf("ReadConfig failed: %v", err)
	}
	if content != "{}" {
		t.Errorf("Expected content '{}', got '%s'", content)
	}

	// Test SaveConfig
	newContent := `{"test": true}`
	err = service.SaveConfig(items[0].Path, newContent)
	if err != nil {
		t.Errorf("SaveConfig failed: %v", err)
	}

	readBack, _ := os.ReadFile(items[0].Path)
	if string(readBack) != newContent {
		t.Errorf("Expected saved content '%s', got '%s'", newContent, string(readBack))
	}

	// Test CreateConfig
	// Create a new project config
	tempProject := t.TempDir()
	newPath, err := service.CreateConfig("Claude Code", ScopeProject, tempProject)
	if err != nil {
		t.Errorf("CreateConfig failed: %v", err)
	}
	if newPath == "" {
		t.Error("Expected new item path to be set")
	}
	t.Logf("Created config at: %s", newPath)
	if !FileExists(newPath) {
		t.Errorf("File does not exist immediately after creation: %s", newPath)
	}

	// Test DeleteConfig
	err = service.DeleteConfig(newPath)
	if err != nil {
		t.Errorf("DeleteConfig failed: %v", err)
	}
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Error("Expected file to be deleted")
	}

	// Test DeleteConfig - Non-existent
	err = service.DeleteConfig(newPath) // Already deleted
	if err == nil {
		t.Error("Expected error when deleting non-existent file")
	}

	// Test ReadConfig - Non-existent
	_, err = service.ReadConfig(newPath)
	if err == nil {
		t.Error("Expected error when reading non-existent file")
	}

	// Test RegisterProvider
	service.RegisterProvider(&ClaudeProvider{})
	if len(service.providers) != 12 { // 11 default + 1 added
		t.Errorf("Expected 12 providers, got %d", len(service.providers))
	}

	// Test DiscoverAll with failing provider
	failingProvider := &MockFailingProvider{}
	service.RegisterProvider(failingProvider)
	items, err = service.DiscoverAll(tempHome)
	if err != nil {
		t.Errorf("DiscoverAll should not fail even if one provider fails: %v", err)
	}
	// Should still have items from other providers
	if len(items) == 0 {
		t.Error("Expected items even with failing provider")
	}
}

type MockFailingProvider struct{}

func (m *MockFailingProvider) Name() string { return "Failing" }
func (m *MockFailingProvider) Discover(projectPath string) ([]Item, error) {
	return nil, os.ErrPermission // Simulate error
}
func (m *MockFailingProvider) Create(scope Scope, projectPath string) (string, error) {
	return "", os.ErrPermission
}

func TestGetUserHome(t *testing.T) {
	home := GetUserHome()
	if home == "" {
		t.Error("Expected non-empty home dir")
	}
}
