package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDynamicProvider(t *testing.T) {
	// 1. Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "dynamic-provider-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 2. Create a dummy provider definition file
	defContent := `
name: My Awesome Tool
binaryName: awesome
versionArgs: ["--version"]
files:
  - name: Global Config
    fileName: .awesome/config.json
    scope: global
    format: json
  - name: Project Config
    fileName: .awesome.yaml
    scope: project
    format: yaml
`
	defPath := filepath.Join(tempDir, "provider.yaml")
	if err := os.WriteFile(defPath, []byte(defContent), 0600); err != nil {
		t.Fatalf("Failed to write provider def: %v", err)
	}

	// 3. Create a new dynamic provider from the definition
	p, err := NewDynamicProvider(defPath)
	if err != nil {
		t.Fatalf("Failed to create dynamic provider: %v", err)
	}

	// 4. Verify the provider's name
	if p.Name() != "My Awesome Tool" {
		t.Errorf("Expected provider name to be 'My Awesome Tool', but got '%s'", p.Name())
	}

	// 5. Create dummy config files
	homeDir := tempDir // Use the temp dir as the home dir for this test
	projectDir := filepath.Join(tempDir, "my-project")
	os.MkdirAll(filepath.Join(homeDir, ".awesome"), 0750)
	os.WriteFile(filepath.Join(homeDir, ".awesome", "config.json"), []byte("{}"), 0600)
	os.MkdirAll(projectDir, 0750)
	os.WriteFile(filepath.Join(projectDir, ".awesome.yaml"), []byte("key: value"), 0600)

	// Override the home dir for this test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	// 6. Test discovery
	items, err := p.Discover(projectDir)
	if err != nil {
		t.Fatalf("Discovery failed: %v", err)
	}

	expectedItems := []Item{
		{
			Provider: "My Awesome Tool",
			Name:     "Global Config",
			FileName: ".awesome/config.json",
			Path:     filepath.Join(homeDir, ".awesome/config.json"),
			Scope:    ScopeGlobal,
			Format:   FormatJSON,
			Exists:   true,
		},
		{
			Provider: "My Awesome Tool",
			Name:     "Project Config",
			FileName: ".awesome.yaml",
			Path:     filepath.Join(projectDir, ".awesome.yaml"),
			Scope:    ScopeProject,
			Format:   FormatYAML,
			Exists:   true,
		},
	}

	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Discovered items do not match expected items.\nExpected: %v\nGot: %v", expectedItems, items)
	}
}
