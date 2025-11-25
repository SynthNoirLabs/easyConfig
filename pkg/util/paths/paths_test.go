package paths

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	appName := "testApp"
	dir := GetConfigDir(appName)

	if dir == "" {
		t.Error("Expected non-empty config dir")
	}

	// Basic check to see if it contains the app name
	// Note: This might be OS specific, but generally true for XDG, AppData, etc.
	if filepath.Base(dir) != appName && filepath.Base(dir) != "testApp" {
		t.Logf("Warning: Config dir %s does not end with app name %s", dir, appName)
	}

	// Test XDG_CONFIG_HOME on Linux (or fallback)
	if runtime.GOOS == "linux" {
		tempConfig := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", tempConfig)

		dir = GetConfigDir(appName)
		expected := filepath.Join(tempConfig, appName)
		if dir != expected {
			t.Errorf("Expected %s, got %s", expected, dir)
		}
	}
}
