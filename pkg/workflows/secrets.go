package workflows

import (
	"fmt"
	"os/exec"
)

// SecretsManager handles setting repository secrets
type SecretsManager struct{}

// NewSecretsManager creates a new SecretsManager
func NewSecretsManager() *SecretsManager {
	return &SecretsManager{}
}

// SetRepositorySecret sets a GitHub repository secret using the gh CLI
func (s *SecretsManager) SetRepositorySecret(name, value string) error {
	// Check if gh is installed
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed or not in PATH")
	}

	// Run gh secret set
	cmd := exec.Command("gh", "secret", "set", name, "--body", value)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set secret: %s: %s", err, string(output))
	}

	return nil
}
