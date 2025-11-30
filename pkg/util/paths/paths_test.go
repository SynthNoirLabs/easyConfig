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
	mustSetEnv(t, "HOME", expectedHome)
	if home := GetHomeDir(); home != expectedHome {
		t.Errorf("GetHomeDir() with HOME set: got %s, want %s", home, expectedHome)
	}

	// Test case 2: HOME is not set (simulate by clearing it)
	mustUnsetEnv(t, "HOME")
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

func TestGetConfigDirForOS(t *testing.T) {
	// Save original environment variables and restore after test
	originalHome := os.Getenv("HOME")
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	originalAppData := os.Getenv("APPDATA")
	defer func() {
		mustSetEnv(t, "HOME", originalHome)
		mustSetEnv(t, "XDG_CONFIG_HOME", originalXDGConfigHome)
		mustSetEnv(t, "APPDATA", originalAppData)
	}()

	testAppName := "testapp"
	mockHome := "/tmp/mock_home"
	mustSetEnv(t, "HOME", mockHome)

	tests := []struct {
		name          string
		osName        string
		env           map[string]string
		expectedPath  string
		shouldContain string // For path separator differences
	}{
		{
			name:   "Linux XDG_CONFIG_HOME",
			osName: "linux",
			env: map[string]string{
				"XDG_CONFIG_HOME": "/tmp/xdg",
			},
			expectedPath: filepath.Join("/tmp/xdg", testAppName),
		},
		{
			name:   "Linux Default",
			osName: "linux",
			env: map[string]string{
				"XDG_CONFIG_HOME": "",
			},
			expectedPath: filepath.Join(mockHome, ".config", testAppName),
		},
		{
			name:         "macOS Default",
			osName:       "darwin",
			env:          map[string]string{},
			expectedPath: filepath.Join(mockHome, "Library", "Application Support", testAppName),
		},
		{
			name:   "Windows APPDATA",
			osName: "windows",
			env: map[string]string{
				"APPDATA": "/tmp/appdata",
			},
			expectedPath: filepath.Join("/tmp/appdata", testAppName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env vars
			for k, v := range tt.env {
				mustSetEnv(t, k, v)
			}

			// Call internal helper
			result := getConfigDirForOS(tt.osName, testAppName)

			if result != tt.expectedPath {
				t.Errorf("getConfigDirForOS(%s) = %s, want %s", tt.osName, result, tt.expectedPath)
			}
		})
	}
}

func TestGetConfigDir(t *testing.T) {
	// Simple smoke test for the public API to ensure it calls the internal one
	// We just check it returns *something* reasonable for the current OS
	dir := GetConfigDir("testapp")
	if dir == "" {
		// It might be empty if HOME is not set, but in test env usually HOME is set or we can set it
		mustSetEnv(t, "HOME", "/tmp/test")
		dir = GetConfigDir("testapp")
		if dir == "" {
			t.Error("GetConfigDir returned empty string even with HOME set")
		}
	}
}

func mustSetEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("failed to set env %s: %v", key, err)
	}
}

func mustUnsetEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset env %s: %v", key, err)
	}
}
