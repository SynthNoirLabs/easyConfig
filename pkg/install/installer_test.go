package install

import (
	"context"
	"fmt"
	"testing"
)

func TestSanitizePackageName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"@modelcontextprotocol/server-filesystem", "server-filesystem"},
		{"simple-package", "simple-package"},
		{"@scope/package-name", "package-name"},
		{"package with spaces", "package-with-spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizePackageName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizePackageName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetServerConfig(t *testing.T) {
	installer := NewInstaller()

	tests := []struct {
		name        string
		packageName string
		packageType PackageType
		expectError bool
		wantCmd     string
	}{
		{
			name:        "Node.js package",
			packageName: "@test/server",
			packageType: PackageTypeNodeJS,
			expectError: false,
			wantCmd:     "npx",
		},
		{
			name:        "Python package",
			packageName: "test-server",
			packageType: PackageTypePython,
			expectError: false,
			wantCmd:     "uvx", // Assuming uvx is mocked or we just check logic
		},
		{
			name:        "Unknown package type",
			packageName: "test",
			packageType: PackageTypeUnknown,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily mock checkToolExists without refactoring Installer to use an interface or variable for lookPath.
			// For now, we'll assume the environment might not have tools, so we only check if it returns *a* config or error as expected.
			// Actually, GetServerConfig calls checkToolExists.
			// If we run this in an env without uvx, the Python test might default to "uv" or fail if neither exists?
			// Wait, the logic is: if uvx exists use it, else use uv.
			// It doesn't error if neither exists during GetServerConfig?
			// Ah, checkToolExists returns error.
			// In GetServerConfig: `if i.checkToolExists("uvx") == nil`. It just checks.
			// So it won't fail, it will just pick one path.

			config, err := installer.GetServerConfig(tt.packageName, tt.packageType)
			if (err != nil) != tt.expectError {
				t.Errorf("GetServerConfig() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && config != nil {
				if config.Command != tt.wantCmd && config.Command != "uv" {
					// Allow uv fallback for python
					t.Errorf("GetServerConfig() command = %v, want %v", config.Command, tt.wantCmd)
				}
			}
		})
	}
}

func TestCheckToolExists(t *testing.T) {
	installer := NewInstaller()

	// Test with a tool that should exist (sh/bash)
	err := installer.checkToolExists("sh")
	if err != nil {
		t.Errorf("checkToolExists(\"sh\") should not error, got: %v", err)
	}

	// Test with a tool that should not exist
	err = installer.checkToolExists("definitely-does-not-exist-tool-12345")
	if err == nil {
		t.Error("checkToolExists() should error for non-existent tool")
	}
}

func TestDetectPackageType_NoTools(t *testing.T) {
	// This test is environment-dependent
	// We can only really test the error path when tools don't exist
	installer := NewInstaller()
	ctx := context.Background()

	// Try with a non-existent package
	_, err := installer.DetectPackageType(ctx, "definitely-not-a-real-package-12345")
	if err == nil {
		t.Error("DetectPackageType() should error for non-existent package")
	}
}

func TestGenerateHelpfulError(t *testing.T) {
	installer := NewInstaller()

	tests := []struct {
		name        string
		err         error
		packageName string
		wantContain string
	}{
		{
			name:        "Package not found",
			err:         fmt.Errorf("package test not found in npm or PyPI registry"),
			packageName: "test-package",
			wantContain: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := installer.generateHelpfulError(tt.err, tt.packageName)
			if !contains(result.Error(), tt.wantContain) {
				t.Errorf("generateHelpfulError() = %v, want to contain %q", result, tt.wantContain)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
