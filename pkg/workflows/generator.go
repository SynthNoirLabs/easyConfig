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

// WorkflowResponse holds the generated workflow and metadata
type WorkflowResponse struct {
	Content           string   `json:"content"`
	RequiredSecrets   []string `json:"requiredSecrets"`
	SetupInstructions string   `json:"setupInstructions"`
}

// Template represents a workflow template returned to callers (UI/API).
type Template struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Agent             string   `json:"agent"`
	Trigger           string   `json:"trigger"`
	Tags              []string `json:"tags"`
	DefaultFilename   string   `json:"defaultFilename"`
	Content           string   `json:"content"`
	RequiredSecrets   []string `json:"requiredSecrets"`
	SetupInstructions string   `json:"setupInstructions"`
}

// GenerateWorkflow generates a workflow content for a specific agent and trigger
// Returns WorkflowResponse, error
func (g *Generator) GenerateWorkflow(agent, trigger string) (*WorkflowResponse, error) {
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
		return nil, fmt.Errorf("unsupported combination: agent=%s, trigger=%s", agent, trigger)
	}

	tmpl, ok := templates[templateID]
	if !ok {
		return nil, fmt.Errorf("template not found for ID: %s", templateID)
	}

	return &WorkflowResponse{
		Content:           tmpl.Content,
		RequiredSecrets:   tmpl.RequiredSecrets,
		SetupInstructions: tmpl.SetupInstructions,
	}, nil
}

// GetSupportedWorkflows returns a list of supported agent-trigger combinations
func (g *Generator) GetSupportedWorkflows() []string {
	workflows := make([]string, 0, len(templates))
	for _, tmpl := range templates {
		workflows = append(workflows, fmt.Sprintf("%s (%s)", tmpl.Agent, tmpl.Trigger))
	}
	return workflows
}

// ListTemplates returns all available templates with metadata and content
func (g *Generator) ListTemplates() []Template {
	result := make([]Template, 0, len(templates))
	for _, tmpl := range templates {
		result = append(result, Template{
			ID:                string(tmpl.ID),
			Name:              tmpl.Name,
			Description:       tmpl.Description,
			Agent:             tmpl.Agent,
			Trigger:           tmpl.Trigger,
			Tags:              tmpl.Tags,
			DefaultFilename:   tmpl.DefaultFilename,
			Content:           tmpl.Content,
			RequiredSecrets:   tmpl.RequiredSecrets,
			SetupInstructions: tmpl.SetupInstructions,
		})
	}
	return result
}
