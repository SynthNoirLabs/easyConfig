package config

import (
	"os"
	"path/filepath"
	"testing"
)

type MockProvider struct {
	name  string
	items []ConfigItem
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Discover(projectPath string) ([]ConfigItem, error) {
	return m.items, nil
}

func TestDiscoveryService_DiscoverAll(t *testing.T) {
	ds := NewDiscoveryService()

	mock1 := &MockProvider{
		name: "Mock1",
		items: []ConfigItem{
			{Provider: "Mock1", Name: "Config1", Path: "/tmp/1", Scope: ScopeGlobal},
		},
	}
	ds.RegisterProvider(mock1)

	results, err := ds.DiscoverAll("/tmp")
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 config, got %d", len(results))
	}
	if results[0].Provider != "Mock1" {
		t.Errorf("Expected provider Mock1, got %s", results[0].Provider)
	}
}

func TestGeminiProvider_Discover(t *testing.T) {
	// Setup: Create dummy files
	homeDir := GetUserHome()
	projectDir, _ := os.Getwd()

	// Global
	os.MkdirAll(filepath.Join(homeDir, ".gemini"), 0755)
	os.WriteFile(filepath.Join(homeDir, ".gemini", "settings.json"), []byte("{}"), 0644)

	// Project
	os.MkdirAll(filepath.Join(projectDir, ".gemini"), 0755)
	os.WriteFile(filepath.Join(projectDir, ".gemini", "settings.json"), []byte("{}"), 0644)

	provider := &GeminiProvider{}
	items, err := provider.Discover(projectDir)
	if err != nil {
		t.Fatalf("GeminiProvider.Discover failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 config items, got %d", len(items))
	}

	globalFound := false
	projectFound := false
	for _, item := range items {
		if item.Scope == ScopeGlobal {
			globalFound = true
		}
		if item.Scope == ScopeProject {
			projectFound = true
		}
	}

	if !globalFound {
		t.Error("Expected to find a global config, but didn't")
	}
	if !projectFound {
		t.Error("Expected to find a project config, but didn't")
	}

	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(homeDir, ".gemini"))
		os.RemoveAll(filepath.Join(projectDir, ".gemini"))
	})
}
