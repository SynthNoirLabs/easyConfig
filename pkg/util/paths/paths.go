package paths

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetHomeDir returns the user's home directory.
func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// GetConfigDir returns the OS-specific default configuration directory for an application.
// For Linux/macOS, it follows XDG Base Directory Specification if XDG_CONFIG_HOME is set,
// otherwise ~/.config/<appName>.
// For Windows, it returns %APPDATA%/<appName>.
// For macOS, it returns ~/Library/Application Support/<appName>.
func GetConfigDir(appName string) string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, appName)
		}
	case "darwin": // macOS
		if home := GetHomeDir(); home != "" {
			return filepath.Join(home, "Library", "Application Support", appName)
		}
	case "linux": // Linux and other Unix-like systems
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			return filepath.Join(xdgConfigHome, appName)
		}
		if home := GetHomeDir(); home != "" {
			return filepath.Join(home, ".config", appName)
		}
	}
	return "" // Fallback if no specific directory can be determined
}

// EnsureDir ensures that a directory exists, creating it if necessary.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
