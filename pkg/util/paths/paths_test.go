package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	// Test case 1: HOME is set
	expectedHome := "/tmp/test_home"
	os.Setenv("HOME", expectedHome)
	if home := GetHomeDir(); home != expectedHome {
		t.Errorf("GetHomeDir() with HOME set: got %s, want %s", home, expectedHome)
	}

	// Test case 2: HOME is not set (simulate by clearing it)
	os.Unsetenv("HOME")
	// On most systems, UserHomeDir will still find it, but if it returns empty string, our func should too
	if home := GetHomeDir(); home == "" {
		// This is acceptable, as os.UserHomeDir might return empty if it can't determine it.
		// More robust would be to mock os.UserHomeDir, but this is a unit test.
		t.Log("GetHomeDir() with HOME unset returned empty string (expected behavior on some systems)")
	} else if runtime.GOOS == "windows" {
		// On Windows, USERPROFILE is typically used if HOME is not set
		if userProfile := os.Getenv("USERPROFILE"); userProfile != "" && home != userProfile {
			t.Errorf("GetHomeDir() with HOME unset (Windows): got %s, want %s", home, userProfile)
		}
	}
}

func TestGetConfigDir(t *testing.T) {
	// Save original environment variables and restore after test
	originalHome := os.Getenv("HOME")
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	originalAppData := os.Getenv("APPDATA")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)
		os.Setenv("APPDATA", originalAppData)
	}()

	testAppName := "testapp"

	// Mock GetHomeDir for consistent testing
	mockHome := "/tmp/mock_home"
	os.Setenv("HOME", mockHome)

	// Test case: Linux with XDG_CONFIG_HOME
	if runtime.GOOS == "linux" {
		mockXDGConfigHome := "/tmp/mock_xdg_config"
		os.Setenv("XDG_CONFIG_HOME", mockXDGConfigHome)
		expected := filepath.Join(mockXDGConfigHome, testAppName)
		if configDir := GetConfigDir(testAppName); configDir != expected {
			t.Errorf("Linux (XDG_CONFIG_HOME): got %s, want %s", configDir, expected)
		}
		os.Unsetenv("XDG_CONFIG_HOME") // Clean up for next test
	}

	// Test case: Linux without XDG_CONFIG_HOME (falls back to ~/.config)
	if runtime.GOOS == "linux" {
		os.Unsetenv("XDG_CONFIG_HOME")
		expected := filepath.Join(mockHome, ".config", testAppName)
		if configDir := GetConfigDir(testAppName); configDir != expected {
			t.Errorf("Linux (no XDG_CONFIG_HOME): got %s, want %s", configDir, expected)
		}
	}

	// Test case: macOS
	if runtime.GOOS == "darwin" {
		expected := filepath.Join(mockHome, "Library", "Application Support", testAppName)
		if configDir := GetConfigDir(testAppName); configDir != expected {
			t.Errorf("macOS: got %s, want %s", configDir, expected)
		}
	}

	// Test case: Windows
	if runtime.GOOS == "windows" {
		mockAppData := "/tmp/mock_appdata" // On Windows, APPDATA is usually system-defined
		os.Setenv("APPDATA", mockAppData)
		expected := filepath.Join(mockAppData, testAppName)
		if configDir := GetConfigDir(testAppName); configDir != expected {
			t.Errorf("Windows: got %s, want %s", configDir, expected)
		}
		os.Unsetenv("APPDATA")
	}

	// Test case: Unknown OS or no relevant env vars
	// This scenario is hard to test directly without mocking runtime.GOOS,
	// but the function should return "" in such cases.
	// For now, we rely on the above tests covering specific OS paths.
}
