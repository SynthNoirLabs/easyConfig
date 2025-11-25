package schema

// Type identifies the type of schema
type Type string

const (
	TypeJSON Type = "json"
	TypeDocs Type = "docs" // For future use
)

// Info contains metadata about where to find a schema
type Info struct {
	ToolName string `json:"toolName"`
	Type     Type   `json:"type"`
	URL      string `json:"url"` // URL to JSON Schema or documentation
}

// Registry holds the mapping of tools to their schema information
var Registry = []Info{
	{
		ToolName: "Gemini",
		Type:     TypeJSON,
		URL:      "https://raw.githubusercontent.com/google-gemini/gemini-cli/main/schemas/settings.schema.json",
	},
	{
		ToolName: "OpenCode",
		Type:     TypeJSON,
		URL:      "https://opencode.ai/config.json",
	},
	{
		ToolName: "Crush",
		Type:     TypeJSON,
		URL:      "https://charm.land/crush.json",
	},
	// Placeholder for documentation-based schemas (future implementation)
	{
		ToolName: "Claude Code",
		Type:     TypeDocs,
		URL:      "https://docs.anthropic.com/en/docs/claude-code/configuration",
	},
	{
		ToolName: "Aider",
		Type:     TypeDocs,
		URL:      "https://aider.chat/docs/config/aider_conf.html",
	},
}
