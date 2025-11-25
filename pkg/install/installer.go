package install

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// PackageType represents the type of MCP package
type PackageType string

const (
	PackageTypeNodeJS  PackageType = "nodejs"
	PackageTypePython  PackageType = "python"
	PackageTypeUnknown PackageType = "unknown"
)

// Installer handles MCP package installation and verification
type Installer struct {
	timeout time.Duration
}

// NewInstaller creates a new Installer with default timeout
func NewInstaller() *Installer {
	return &Installer{
		timeout: 10 * time.Second,
	}
}

// DetectPackageType attempts to determine if a package is Node.js or Python
// It tries npx first (most common), then uvx
func (i *Installer) DetectPackageType(ctx context.Context, packageName string) (PackageType, error) {
	// Try Node.js first (most MCP servers are Node.js)
	if err := i.checkToolExists("npx"); err == nil {
		if i.VerifyNodePackage(ctx, packageName) == nil {
			return PackageTypeNodeJS, nil
		}
	}

	// Try Python
	if err := i.checkToolExists("uvx"); err == nil {
		if i.VerifyPythonPackage(ctx, packageName) == nil {
			return PackageTypePython, nil
		}
	}

	// If uvx doesn't exist, try uv
	if err := i.checkToolExists("uv"); err == nil {
		if i.VerifyPythonPackage(ctx, packageName) == nil {
			return PackageTypePython, nil
		}
	}

	return PackageTypeUnknown, fmt.Errorf("package %s not found in npm or PyPI registry", packageName)
}

// VerifyNodePackage checks if a Node.js package exists by running npx --help
func (i *Installer) VerifyNodePackage(ctx context.Context, packageName string) error {
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	// npx -y <package>@latest --help
	cmd := exec.CommandContext(ctx, "npx", "-y", packageName+"@latest", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("package verification failed: %w (output: %s)", err, string(output))
	}
	return nil
}

// VerifyPythonPackage checks if a Python package exists by running uvx --help
func (i *Installer) VerifyPythonPackage(ctx context.Context, packageName string) error {
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	// Try uvx first
	cmd := exec.CommandContext(ctx, "uvx", packageName, "--help")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	// Fallback to uv tool run
	cmd = exec.CommandContext(ctx, "uv", "tool", "run", packageName, "--help")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("package verification failed: %w (output: %s)", err, string(output))
	}
	return nil
}

// checkToolExists verifies if a command-line tool is available
func (i *Installer) checkToolExists(tool string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		return fmt.Errorf("%s not found: %w", tool, err)
	}
	return nil
}

// ServerConfig represents an MCP server configuration
type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// GetServerConfig returns the configuration for an MCP server
func (i *Installer) GetServerConfig(packageName string, packageType PackageType) (*ServerConfig, error) {
	var config ServerConfig

	switch packageType {
	case PackageTypeNodeJS:
		config = ServerConfig{
			Command: "npx",
			Args:    []string{"-y", packageName + "@latest"},
		}
	case PackageTypePython:
		// Prefer uvx if available
		if i.checkToolExists("uvx") == nil {
			config = ServerConfig{
				Command: "uvx",
				Args:    []string{packageName},
			}
		} else {
			config = ServerConfig{
				Command: "uv",
				Args:    []string{"tool", "run", packageName},
			}
		}
	default:
		return nil, fmt.Errorf("unsupported package type: %s", packageType)
	}

	return &config, nil
}

// sanitizePackageName converts package names to valid identifiers
// e.g., @modelcontextprotocol/server-filesystem -> server-filesystem
func sanitizePackageName(name string) string {
	// Remove scope (e.g., @modelcontextprotocol/)
	if strings.Contains(name, "/") {
		parts := strings.Split(name, "/")
		name = parts[len(parts)-1]
	}
	// Replace invalid characters
	name = strings.ReplaceAll(name, "@", "")
	name = strings.ReplaceAll(name, " ", "-")
	return name
}

// InstallPackage verifies the package and returns its configuration
func (i *Installer) InstallPackage(ctx context.Context, packageName string) (*ServerConfig, error) {
	// Detect package type
	packageType, err := i.DetectPackageType(ctx, packageName)
	if err != nil {
		return nil, i.generateHelpfulError(err, packageName)
	}

	// Get server config
	config, err := i.GetServerConfig(packageName, packageType)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return config, nil
}

// generateHelpfulError creates user-friendly error messages
func (i *Installer) generateHelpfulError(err error, packageName string) error {
	errStr := err.Error()

	// Check for missing tools
	if i.checkToolExists("npx") != nil && i.checkToolExists("uvx") != nil {
		return fmt.Errorf("neither Node.js (npx) nor Python uv (uvx) found. Install one of:\n"+
			"  - Node.js from https://nodejs.org\n"+
			"  - Python uv: pip install uv\n"+
			"Original error: %w", err)
	}

	// Package not found
	if strings.Contains(errStr, "not found") {
		return fmt.Errorf("package '%s' not found in npm or PyPI registry. "+
			"Verify the package name at:\n"+
			"  - https://www.npmjs.com/search?q=%s\n"+
			"  - https://pypi.org/search/?q=%s\n"+
			"Original error: %w",
			packageName, packageName, packageName, err)
	}

	return fmt.Errorf("installation failed: %w", err)
}
