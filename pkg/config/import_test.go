package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestImportProfiles(t *testing.T) {
	// Setup a temporary directory for profiles
	tmpDir := t.TempDir()
	profilesDir := filepath.Join(tmpDir, "profiles")
	os.MkdirAll(profilesDir, 0755)

	// Create a DiscoveryService that points to the temp directory
	service := NewDiscoveryService(nil)
	originalProfilesRoot := profilesRoot
	profilesRoot = func() (string, error) { return profilesDir, nil }
	defer func() { profilesRoot = originalProfilesRoot }()

	// Create some export data
	exportData := ExportData{
		Version: "1.0",
		Profiles: []ExportedProfile{
			{Name: "new-profile", Configs: []ExportedConfig{{Provider: "claude", Content: "test"}}},
		},
	}
	data, _ := json.Marshal(exportData)

	// Test importing the data
	results, err := service.ImportProfiles(data, ImportStrategySkip)
	if err != nil {
		t.Fatalf("ImportProfiles failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Status != ImportStatusSuccess {
		t.Errorf("Expected status 'success', got '%s'", results[0].Status)
	}

	// Verify the profile was created
	_, err = os.Stat(filepath.Join(profilesDir, "new-profile.json"))
	if os.IsNotExist(err) {
		t.Error("Expected profile file to be created, but it was not")
	}
}
