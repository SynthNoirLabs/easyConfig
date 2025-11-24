# Claude Code Configuration Reference

**Source:** https://docs.claude.com/en/docs/claude-code/settings

## Settings File Locations
*   **User settings:** `~/.claude/settings.json` (Global user config)
*   **Project settings:** `.claude/settings.json` (Shared project config)
*   **Local settings:** `.claude/settings.local.json` (Local overrides, not committed)
*   **Managed Settings (Linux):** `/etc/claude-code/managed-settings.json`
*   **Managed MCP (Linux):** `/etc/claude-code/managed-mcp.json`

## Format: JSON

## Configuration Options

### 1. Model Selection
To permanently set the default model:
```json
{
  "model": "claude-3-7-sonnet-latest"
}
```

### 2. Permissions & Sandbox
Controls access to tools, files, and network.
```json
{
  "permissions": {
    "allow": [
      "Bash(npm run lint)",
      "Bash(npm run test:*)"
    ],
    "deny": [
      "Bash(curl:*)",
      "Read(./.env)",
      "Read(./secrets/**)"
    ]
  },
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": true,
    "excludedCommands": ["docker"],
    "network": {
      "allowUnixSockets": ["/var/run/docker.sock"],
      "allowLocalBinding": true
    }
  }
}
```

### 3. Plugins & Marketplaces
Enable specific plugins from defined marketplaces.
```json
{
  "enabledPlugins": {
    "formatter@company-tools": true,
    "deployer@company-tools": true
  },
  "extraKnownMarketplaces": {
    "company-tools": {
      "source": "github",
      "repo": "company/claude-plugins"
    }
  }
}
```

### 4. Environment Variables & Telemetry
Inject env vars and configure OpenTelemetry.
```json
{
  "env": {
    "CLAUDE_CODE_ENABLE_TELEMETRY": "1",
    "OTEL_METRICS_EXPORTER": "otlp",
    "ANTHROPIC_AUTH_TOKEN": "sk-..." // Static API key if not using OAuth
  }
}
```

## Context Files
*   **`CLAUDE.md`**: Placed in the project root to provide instructions, architecture notes, and conventions to the agent.
