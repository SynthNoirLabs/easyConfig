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

func TestReadConfig_Success(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	expectedContent := "test content\nline 2"

	err := os.WriteFile(testFile, []byte(expectedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ds := NewDiscoveryService()
	content, err := ds.ReadConfig(testFile)
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	if content != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, content)
	}
}

func TestReadConfig_FileNotFound(t *testing.T) {
	ds := NewDiscoveryService()
	_, err := ds.ReadConfig("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestSaveConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "new content"

	ds := NewDiscoveryService()
	err := ds.SaveConfig(testFile, content)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify the file was written
	savedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(savedContent) != content {
		t.Errorf("Expected saved content %q, got %q", content, string(savedContent))
	}
}

func TestSaveConfig_ValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")
	validJSON := `{"key": "value", "number": 42}`

	ds := NewDiscoveryService()
	err := ds.SaveConfig(testFile, validJSON)
	if err != nil {
		t.Fatalf("SaveConfig failed for valid JSON: %v", err)
	}

	// Verify the file was written
	savedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(savedContent) != validJSON {
		t.Errorf("Expected saved content %q, got %q", validJSON, string(savedContent))
	}
}

func TestSaveConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")
	invalidJSON := `{"key": "value", "invalid"`

	ds := NewDiscoveryService()
	err := ds.SaveConfig(testFile, invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}

	// Verify the file was not created
	_, err = os.Stat(testFile)
	if !os.IsNotExist(err) {
		t.Error("Expected file not to exist after invalid JSON save")
	}
}

func TestSaveConfig_NonJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.yaml")
	content := "key: value\ninvalid: {json"

	ds := NewDiscoveryService()
	// Should not validate YAML files as JSON
	err := ds.SaveConfig(testFile, content)
	if err != nil {
		t.Fatalf("SaveConfig failed for non-JSON file: %v", err)
	}

	// Verify the file was written
	savedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(savedContent) != content {
		t.Errorf("Expected saved content %q, got %q", content, string(savedContent))
	}
}
