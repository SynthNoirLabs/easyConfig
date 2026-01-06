package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExportProfiles(t *testing.T) {
	// Setup a temporary directory for profiles
	tmpDir := t.TempDir()
	profilesDir := filepath.Join(tmpDir, "profiles")
	os.MkdirAll(profilesDir, 0755)

	// Create a dummy profile file
	profile := Profile{
		Name:      "test-profile",
		Items:     []ProfileItem{{Provider: "claude", Scope: ScopeGlobal, Content: "content"}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	data, _ := json.Marshal(profile)
	os.WriteFile(filepath.Join(profilesDir, "test-profile.json"), data, 0644)

	// Create a DiscoveryService that points to the temp directory
	service := NewDiscoveryService(nil, nil)
	// Override the profilesRoot function to use our temp dir
	originalProfilesRoot := profilesRoot
	profilesRoot = func() (string, error) { return profilesDir, nil }
	defer func() { profilesRoot = originalProfilesRoot }()

	// Test exporting the profile
	exportedData, err := service.ExportProfiles([]string{"test-profile"})
	if err != nil {
		t.Fatalf("ExportProfiles failed: %v", err)
	}

	var exportContainer ExportData
	if err := json.Unmarshal(exportedData, &exportContainer); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	if len(exportContainer.Profiles) != 1 {
		t.Fatalf("Expected 1 profile, got %d", len(exportContainer.Profiles))
	}

	if exportContainer.Profiles[0].Name != "test-profile" {
		t.Errorf("Expected profile name 'test-profile', got '%s'", exportContainer.Profiles[0].Name)
	}
}
