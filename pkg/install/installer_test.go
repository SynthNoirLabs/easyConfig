package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

func TestCreateServerConfig(t *testing.T) {
	tempDir := t.TempDir()
	installer := NewInstaller()

	tests := []struct {
		name        string
		packageName string
		packageType PackageType
		expectError bool
	}{
		{
			name:        "Node.js package",
			packageName: "@test/server",
			packageType: PackageTypeNodeJS,
			expectError: false,
		},
		{
			name:        "Python package",
			packageName: "test-server",
			packageType: PackageTypePython,
			expectError: false,
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
			err := installer.CreateServerConfig(tt.packageName, tt.packageType, tempDir)
			if (err != nil) != tt.expectError {
				t.Errorf("CreateServerConfig() error = %v, expectError %v", err, tt.expectError)
			}

			if !tt.expectError {
				// Verify file was created
				expectedFile := filepath.Join(tempDir, "mcp-"+sanitizePackageName(tt.packageName)+".json")
				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Errorf("Expected config file %s was not created", expectedFile)
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
