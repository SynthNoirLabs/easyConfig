import type { config, marketplaces, workflows } from "../../wailsjs/go/models";

const demoConfigItems: config.Item[] = [
  {
    provider: "Claude Code",
    name: "Desktop MCP Config",
    fileName: "claude_desktop_config.json",
    path: "/demo/.claude/claude_desktop_config.json",
    scope: "global",
    format: "json",
    exists: true,
  },
  {
    provider: "Gemini CLI",
    name: "Project Settings",
    fileName: "settings.json",
    path: "/demo/project/.gemini/settings.json",
    scope: "project",
    format: "json",
    exists: true,
  },
  {
    provider: "OpenCode",
    name: "Workspace Config",
    fileName: "opencode.json",
    path: "/demo/project/opencode.json",
    scope: "project",
    format: "json",
    exists: true,
  },
  {
    provider: "GitHub Copilot",
    name: "Instructions",
    fileName: "copilot-instructions.md",
    path: "/demo/project/.github/copilot-instructions.md",
    scope: "project",
    format: "markdown",
    exists: true,
  },
];

const demoConfigContents = new Map<string, string>([
  [
    "/demo/.claude/claude_desktop_config.json",
    JSON.stringify(
      {
        mcpServers: {
          filesystem: {
            command: "npx",
            args: ["-y", "@modelcontextprotocol/server-filesystem", "/demo"],
          },
          github: {
            command: "docker",
            args: ["run", "-i", "--rm", "ghcr.io/github/github-mcp-server"],
            env: {
              GITHUB_TOKEN: String.raw`\${GITHUB_TOKEN}`,
            },
          },
        },
      },
      null,
      2,
    ),
  ],
  [
    "/demo/project/.gemini/settings.json",
    JSON.stringify(
      {
        model: "gemini-2.5-pro",
        sandbox: {
          enabled: true,
          mode: "workspace-write",
        },
        telemetry: {
          enabled: false,
        },
      },
      null,
      2,
    ),
  ],
  [
    "/demo/project/opencode.json",
    JSON.stringify(
      {
        workspace: {
          trust: "trusted",
          autoApprove: false,
        },
        tools: ["terminal", "search", "edit"],
      },
      null,
      2,
    ),
  ],
  [
    "/demo/project/.github/copilot-instructions.md",
    `# Copilot Instructions

- Prefer small, reviewable changes.
- Keep provider logic in \`pkg/config\`.
- Use \`pkg/util/paths\` for cross-platform paths.
`,
  ],
]);

const demoProviderStatuses = [
  {
    providerName: "Claude Code",
    health: "healthy",
    statusMessage: "CLI detected and desktop config found",
    discoveredFiles: [demoConfigItems[0]],
    lastChecked: new Date("2026-03-07T08:45:00Z").toISOString(),
  },
  {
    providerName: "Gemini CLI",
    health: "healthy",
    statusMessage: "Project settings discovered",
    discoveredFiles: [demoConfigItems[1]],
    lastChecked: new Date("2026-03-07T08:45:00Z").toISOString(),
  },
  {
    providerName: "OpenCode",
    health: "unknown",
    statusMessage: "Config present, binary not verified in browser demo mode",
    discoveredFiles: [demoConfigItems[2]],
    lastChecked: new Date("2026-03-07T08:45:00Z").toISOString(),
  },
] as unknown as config.ProviderStatus[];

const demoMarketplacePackages: marketplaces.MCPPackage[] = [
  {
    name: "filesystem-mcp",
    description: "Browse and edit local project files through MCP.",
    vendor: "modelcontextprotocol",
    source: JSON.stringify({
      name: "filesystem-mcp",
      version: "1.2.0",
      url: "https://github.com/modelcontextprotocol/servers",
    }),
    version: "1.2.0",
    stars: 1840,
    downloads: 92000,
    tags: ["filesystem", "core", "local"],
    verified: true,
  },
  {
    name: "github-mcp",
    description:
      "Pull requests, issues, and workflows from GitHub in one MCP endpoint.",
    vendor: "github",
    source: JSON.stringify({
      name: "github-mcp",
      version: "0.18.0",
      url: "https://github.com/github/github-mcp-server",
    }),
    version: "0.18.0",
    stars: 3320,
    downloads: 64000,
    tags: ["github", "code-review", "automation"],
    verified: true,
  },
  {
    name: "postgres-mcp",
    description:
      "Query and inspect PostgreSQL databases safely from your agent.",
    vendor: "supabase",
    source: JSON.stringify({
      name: "postgres-mcp",
      version: "0.9.4",
      url: "https://github.com/supabase-community/postgres-mcp",
    }),
    version: "0.9.4",
    stars: 870,
    downloads: 12800,
    tags: ["database", "postgres", "analytics"],
    verified: false,
  },
];

const demoWorkflowTemplates: workflows.Template[] = [
  {
    id: "claude-review",
    name: "Claude Review",
    description: "Run Claude on pull requests for architectural and UX review.",
    agent: "Claude",
    trigger: "pull_request",
    tags: ["review", "frontend", "agent"],
    defaultFilename: "claude-review.yml",
    content: `name: Claude Review

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: anthropics/claude-code-action@v1
`,
    requiredSecrets: ["ANTHROPIC_API_KEY"],
    setupInstructions:
      "Add the API key as a repository secret, then save this template into .github/workflows.",
  },
  {
    id: "gemini-triage",
    name: "Gemini Triage",
    description: "Label and summarize new issues with Gemini CLI.",
    agent: "Gemini",
    trigger: "issues",
    tags: ["triage", "issues", "automation"],
    defaultFilename: "gemini-triage.yml",
    content: `name: Gemini Triage

on:
  issues:
    types: [opened]

jobs:
  triage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: google-github-actions/run-gemini-cli@v0
`,
    requiredSecrets: ["GEMINI_API_KEY"],
    setupInstructions:
      "Enable the workflow on issue creation and verify that your Google/Gemini credentials are configured.",
  },
];

export function isBrowserDemoMode(): boolean {
  const maybeWindow = globalThis as typeof globalThis & {
    go?: { main?: { App?: unknown } };
  };

  return !maybeWindow.go?.main?.App;
}

export function isWailsUnavailableError(error: unknown): boolean {
  const text = String(error ?? "");
  return (
    isBrowserDemoMode() ||
    text.includes("window.go") ||
    text.includes("window.runtime") ||
    text.includes("is not a function") ||
    text.includes("undefined")
  );
}

export function getDemoConfigs(): config.Item[] {
  return demoConfigItems.map((item) => ({ ...item }));
}

export function readDemoConfig(path: string): string {
  return (
    demoConfigContents.get(path) ??
    `// Demo content unavailable for ${path}\n// Connect the Wails backend for live data.`
  );
}

export function saveDemoConfig(path: string, content: string): void {
  demoConfigContents.set(path, content);
}

export function deleteDemoConfig(path: string): void {
  const index = demoConfigItems.findIndex((item) => item.path === path);
  if (index >= 0) {
    demoConfigItems.splice(index, 1);
  }
  demoConfigContents.delete(path);
}

export function getDemoProviderStatuses(): config.ProviderStatus[] {
  return demoProviderStatuses.map((status) => ({
    ...status,
    discoveredFiles: status.discoveredFiles?.map((item) => ({ ...item })),
  })) as unknown as config.ProviderStatus[];
}

export function getDemoMarketplacePackages(): marketplaces.MCPPackage[] {
  return demoMarketplacePackages.map((pkg) => ({
    ...pkg,
    tags: pkg.tags ? [...pkg.tags] : [],
  }));
}

export function getDemoWorkflowTemplates(): workflows.Template[] {
  return demoWorkflowTemplates.map((template) => ({
    ...template,
    tags: template.tags ? [...template.tags] : [],
    requiredSecrets: template.requiredSecrets
      ? [...template.requiredSecrets]
      : [],
  }));
}
