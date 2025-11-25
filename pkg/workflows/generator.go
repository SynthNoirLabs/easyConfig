package workflows

import (
	"fmt"
	"strings"
)

// Generator handles the generation of GitHub Actions workflows
type Generator struct{}

// NewGenerator creates a new Generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateWorkflow generates a workflow content for a specific agent and trigger
func (g *Generator) GenerateWorkflow(agent, trigger string) (string, error) {
	// Map agent+trigger to TemplateID
	// This is a simple mapping for now, could be more sophisticated
	var templateID TemplateID

	key := fmt.Sprintf("%s-%s", strings.ToLower(agent), strings.ToLower(trigger))

	switch key {
	case "claude-comment":
		templateID = ClaudeComment
	case "jules-label":
		templateID = JulesLabel
	case "codex-pr":
		templateID = CodexPR
	case "copilot-manual":
		templateID = CopilotManual
	default:
		return "", fmt.Errorf("unsupported combination: agent=%s, trigger=%s", agent, trigger)
	}

	content, ok := templates[templateID]
	if !ok {
		return "", fmt.Errorf("template not found for ID: %s", templateID)
	}

	return content, nil
}

// GetSupportedWorkflows returns a list of supported agent-trigger combinations
func (g *Generator) GetSupportedWorkflows() []string {
	return []string{
		"Claude (Comment)",
		"Jules (Label)",
		"Codex (PR)",
		"Copilot (Manual)",
	}
}
