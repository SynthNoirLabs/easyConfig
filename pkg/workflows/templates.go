package workflows

// TemplateID represents the ID of a workflow template
type TemplateID string

const (
	ClaudeComment TemplateID = "claude-comment"
	JulesLabel    TemplateID = "jules-label"
	CodexPR       TemplateID = "codex-pr"
	CopilotManual TemplateID = "copilot-manual"
)

var templates = map[TemplateID]string{
	ClaudeComment: `name: Claude Agent
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
	JulesLabel: `name: Jules Agent
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
	CodexPR: `name: Codex Review
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
	CopilotManual: `name: Copilot Task
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
}
