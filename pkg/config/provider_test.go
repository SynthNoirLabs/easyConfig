package config

import (
	"os"
	"path/filepath"
	"testing"
)

type MockProvider struct {
	name  string
	items []Item
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Discover(_ string) ([]Item, error) {
	return m.items, nil
}

func TestDiscoveryService_DiscoverAll(t *testing.T) {
	ds := NewDiscoveryService()
	ds.providers = []Provider{} // Clear default providers for testing

	mock1 := &MockProvider{
		name: "Mock1",
		items: []Item{
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

	err := os.WriteFile(testFile, []byte(expectedContent), 0o600)
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

func TestGeminiProvider_Discover(t *testing.T) {
	// Setup: Create temporary directories for testing
	tempHome := t.TempDir()
	tempProject := t.TempDir()

	// Override HOME environment variable for testing
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global
	if err := os.MkdirAll(filepath.Join(tempHome, ".gemini"), 0o755); err != nil {
		t.Fatalf("Failed to create global .gemini directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempHome, ".gemini", "settings.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write global settings.json: %v", err)
	}

	// Project
	if err := os.MkdirAll(filepath.Join(tempProject, ".gemini"), 0o755); err != nil {
		t.Fatalf("Failed to create project .gemini directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempProject, ".gemini", "settings.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write project settings.json: %v", err)
	}

	provider := &GeminiProvider{}
	items, err := provider.Discover(tempProject)
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
}

func TestCodexProvider_Discover(t *testing.T) {
	// Setup: Create temporary directories for testing
	tempHome := t.TempDir()
	tempProject := t.TempDir()

	// Override HOME environment variable for testing
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global
	if err := os.MkdirAll(filepath.Join(tempHome, ".codex"), 0o755); err != nil {
		t.Fatalf("Failed to create global .codex directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempHome, ".codex", "config.toml"), []byte("# config"), 0o600); err != nil {
		t.Fatalf("Failed to write global config.toml: %v", err)
	}

	// Project
	if err := os.MkdirAll(filepath.Join(tempProject, ".codex"), 0o755); err != nil {
		t.Fatalf("Failed to create project .codex directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempProject, ".codex", "config.toml"), []byte("# config"), 0o600); err != nil {
		t.Fatalf("Failed to write project config.toml: %v", err)
	}

	provider := &CodexProvider{}
	items, err := provider.Discover(tempProject)
	if err != nil {
		t.Fatalf("CodexProvider.Discover failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 config items, got %d", len(items))
	}
}

func TestClaudeProvider_Discover(t *testing.T) {
	tempHome := t.TempDir()
	tempProject := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global
	if err := os.MkdirAll(filepath.Join(tempHome, ".claude"), 0o755); err != nil {
		t.Fatalf("Failed to create global .claude directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempHome, ".claude", "claude_desktop_config.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write desktop config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempHome, ".claude", "settings.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write cli config: %v", err)
	}

	// Project
	if err := os.MkdirAll(filepath.Join(tempProject, ".claude"), 0o755); err != nil {
		t.Fatalf("Failed to create project .claude directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempProject, ".claude", "settings.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write project config: %v", err)
	}

	provider := &ClaudeProvider{}
	items, err := provider.Discover(tempProject)
	if err != nil {
		t.Fatalf("ClaudeProvider.Discover failed: %v", err)
	}

	// Expecting: Desktop Config (Global), CLI Settings (Global), Project Settings (Project)
	expectedCount := 3
	if len(items) < expectedCount {
		t.Errorf("Expected at least %d config items, got %d", expectedCount, len(items))
	}
}

func TestCopilotProvider_Discover(t *testing.T) {
	tempHome := t.TempDir()
	tempProject := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global
	if err := os.MkdirAll(filepath.Join(tempHome, ".copilot"), 0o755); err != nil {
		t.Fatalf("Failed to create global .copilot directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempHome, ".copilot", "mcp-config.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write global mcp-config: %v", err)
	}

	// Project
	if err := os.MkdirAll(filepath.Join(tempProject, ".github"), 0o755); err != nil {
		t.Fatalf("Failed to create project .github directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempProject, ".github", "copilot-instructions.md"), []byte("# Instructions"), 0o600); err != nil {
		t.Fatalf("Failed to write project instructions: %v", err)
	}

	provider := &CopilotProvider{}
	items, err := provider.Discover(tempProject)
	if err != nil {
		t.Fatalf("CopilotProvider.Discover failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 config items, got %d", len(items))
	}
}

func TestOpenAIProvider_Discover(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global
	if err := os.MkdirAll(filepath.Join(tempHome, ".config", "openai"), 0o755); err != nil {
		t.Fatalf("Failed to create global .config/openai directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempHome, ".config", "openai", "config.yaml"), []byte("key: value"), 0o600); err != nil {
		t.Fatalf("Failed to write global config.yaml: %v", err)
	}

	provider := &OpenAIProvider{}
	items, err := provider.Discover("") // No project path
	if err != nil {
		t.Fatalf("OpenAIProvider.Discover failed: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 config item, got %d", len(items))
	}
}

func TestJulesProvider_Discover(t *testing.T) {
	tempHome := t.TempDir()
	tempProject := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global
	if err := os.MkdirAll(filepath.Join(tempHome, ".jules-mcp"), 0o755); err != nil {
		t.Fatalf("Failed to create global .jules-mcp directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempHome, ".jules-mcp", "data.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write global data.json: %v", err)
	}

	// Project
	if err := os.WriteFile(filepath.Join(tempProject, "AGENTS.md"), []byte("# Agents"), 0o600); err != nil {
		t.Fatalf("Failed to write project AGENTS.md: %v", err)
	}

	provider := &JulesProvider{}
	items, err := provider.Discover(tempProject)
	if err != nil {
		t.Fatalf("JulesProvider.Discover failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 config items, got %d", len(items))
	}
}

func TestOpenCodeProvider_Discover(t *testing.T) {
	tempHome := t.TempDir()
	tempProject := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global: ~/.config/opencode/opencode.json
	globalDir := filepath.Join(tempHome, ".config", "opencode")
	if err := os.MkdirAll(globalDir, 0o755); err != nil {
		t.Fatalf("Failed to create global opencode dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, "opencode.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write global opencode.json: %v", err)
	}

	// Project: opencode.json and opencode.local.json
	if err := os.WriteFile(filepath.Join(tempProject, "opencode.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write project opencode.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempProject, "opencode.local.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write project opencode.local.json: %v", err)
	}

	provider := &OpenCodeProvider{}
	items, err := provider.Discover(tempProject)
	if err != nil {
		t.Fatalf("OpenCodeProvider.Discover failed: %v", err)
	}

	// Expecting: Global, Project, Local Secrets
	expectedCount := 3
	if len(items) != expectedCount {
		t.Errorf("Expected %d config items, got %d", expectedCount, len(items))
	}
}

func TestCrushProvider_Discover(t *testing.T) {
	tempHome := t.TempDir()
	tempProject := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Global: ~/.config/crush/crush.json & providers.json
	globalDir := filepath.Join(tempHome, ".config", "crush")
	if err := os.MkdirAll(globalDir, 0o755); err != nil {
		t.Fatalf("Failed to create global crush dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, "crush.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write global crush.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(globalDir, "providers.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write global providers.json: %v", err)
	}

	// Project: .crush.json, crush.json, .crushignore
	if err := os.WriteFile(filepath.Join(tempProject, ".crush.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write project .crush.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempProject, "crush.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("Failed to write project crush.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempProject, ".crushignore"), []byte("node_modules"), 0o600); err != nil {
		t.Fatalf("Failed to write project .crushignore: %v", err)
	}

	provider := &CrushProvider{}
	items, err := provider.Discover(tempProject)
	if err != nil {
		t.Fatalf("CrushProvider.Discover failed: %v", err)
	}

	// Expecting: Global Main, Global Providers, Project Hidden, Project Visible, Project Ignore
	expectedCount := 5
	if len(items) != expectedCount {
		t.Errorf("Expected %d config items, got %d", expectedCount, len(items))
	}
}
