package workflows

// TemplateID represents the ID of a workflow template
type TemplateID string

const (
	ClaudeComment TemplateID = "claude-comment"
	JulesLabel    TemplateID = "jules-label"
	CodexPR       TemplateID = "codex-pr"
	CopilotManual TemplateID = "copilot-manual"
)

// WorkflowTemplate holds the content and metadata for a workflow
type WorkflowTemplate struct {
	ID                TemplateID
	Name              string
	Description       string
	Agent             string
	Trigger           string
	Tags              []string
	DefaultFilename   string
	Content           string
	RequiredSecrets   []string
	SetupInstructions string
}

var templates = map[TemplateID]WorkflowTemplate{
	ClaudeComment: {
		ID:              ClaudeComment,
		Name:            "Claude comment triage",
		Description:     "Runs Claude when a comment includes /claude to triage issues/PRs.",
		Agent:           "Claude",
		Trigger:         "Comment",
		Tags:            []string{"claude", "comment", "triage"},
		DefaultFilename: "claude-comment.yml",
		Content: `name: Claude Agent
on:
  issue_comment:
    types: [created]

jobs:
  claude:
    if: contains(github.event.comment.body, '/claude')
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
      issues: write
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Run Claude
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # This is a placeholder command. Adjust based on actual CLI usage.
          npx @anthropic-ai/claude-code --prompt "${{ github.event.comment.body }}"
`,
		RequiredSecrets:   []string{"ANTHROPIC_API_KEY"},
		SetupInstructions: "Tip: Run `/install-github-app` in your Claude Code terminal for a guided setup.",
	},
	JulesLabel: {
		ID:              JulesLabel,
		Name:            "Jules label responder",
		Description:     "Triggers Jules when an issue gets the 'jules' label.",
		Agent:           "Jules",
		Trigger:         "Label",
		Tags:            []string{"jules", "label", "issues"},
		DefaultFilename: "jules-label.yml",
		Content: `name: Jules Agent
on:
  issues:
    types: [labeled]

jobs:
  jules:
    if: github.event.label.name == 'jules'
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
      issues: write
    steps:
      - uses: actions/checkout@v3
      - name: Run Jules
        env:
          JULES_API_KEY: ${{ secrets.JULES_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          echo "Triggering Jules Agent..."
          # Insert actual Jules CLI or API call here
`,
		RequiredSecrets:   []string{"JULES_API_KEY"},
		SetupInstructions: "Ensure you have the 'jules' label created in your repository.",
	},
	CodexPR: {
		ID:              CodexPR,
		Name:            "Codex PR reviewer",
		Description:     "Uses Codex/OpenAI to review pull requests on open/synchronize.",
		Agent:           "Codex",
		Trigger:         "PR",
		Tags:            []string{"codex", "openai", "pull_request"},
		DefaultFilename: "codex-pr.yml",
		Content: `name: Codex Review
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
    steps:
      - uses: actions/checkout@v3
      - name: Codex Code Review
        uses: openai/codex-action@v1
        with:
          openai_api_key: ${{ secrets.OPENAI_API_KEY }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
          model: 'gpt-4' # or specific codex model
`,
		RequiredSecrets:   []string{"OPENAI_API_KEY"},
		SetupInstructions: "Get your API key from platform.openai.com.",
	},
	CopilotManual: {
		ID:              CopilotManual,
		Name:            "Copilot manual dispatch",
		Description:     "Manually trigger Copilot CLI tasks via workflow_dispatch input.",
		Agent:           "Copilot",
		Trigger:         "Manual",
		Tags:            []string{"copilot", "manual"},
		DefaultFilename: "copilot-task.yml",
		Content: `name: Copilot Task
on:
  workflow_dispatch:
    inputs:
      task:
        description: 'Task description'
        required: true

jobs:
  copilot:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Copilot CLI
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh copilot suggest "${{ inputs.task }}"
`,
		RequiredSecrets:   []string{}, // GITHUB_TOKEN is automatic
		SetupInstructions: "This workflow uses the standard GITHUB_TOKEN.",
	},
}
