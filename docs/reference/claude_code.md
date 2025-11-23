# Claude Code Configuration Reference

**Source:** https://code.claude.com/docs/en/settings

## Settings File Locations
*   **User settings:** `~/.claude/settings.json`
*   **Project settings:** `.claude/settings.json`
*   **Local settings:** `.claude/settings.local.json` (Not committed)

## Format: JSON

## Key Settings Categories

### 1. Permissions
Controls access to tools and files.
```json
"permissions": {
  "allow": ["Bash(npm run test:*)"],
  "deny": [
    "Bash(curl:*)",
    "Read(./secrets/**)",
    "Read(./.env)"
  ]
}
```

### 2. Environment Variables
Injected into the agent's session.
```json
"env": {
  "CLAUDE_CODE_ENABLE_TELEMETRY": "1",
  "OTEL_METRICS_EXPORTER": "otlp"
}
```

### 3. Sandbox
Isolates execution (Docker/Container logic).
```json
"sandbox": {
  "enabled": true,
  "autoAllowBashIfSandboxed": true,
  "excludedCommands": ["docker"],
  "network": {
    "allowLocalBinding": true
  }
}
```

### 4. Plugins & Marketplaces
Extend functionality.
```json
"enabledPlugins": {
  "formatter@company-tools": true
},
"extraKnownMarketplaces": {
  "company-tools": {
    "source": "github",
    "repo": "company/claude-plugins"
  }
}
```

## Full Example
```json
{
  "permissions": {
    "deny": ["Read(./secrets/**)"]
  },
  "env": {
    "CLAUDE_CODE_ENABLE_TELEMETRY": "0"
  },
  "sandbox": {
    "enabled": true
  },
  "statusLine": {
    "type": "command",
    "command": "~/.claude/statusline.sh"
  }
}
```