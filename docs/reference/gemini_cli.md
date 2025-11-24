# Gemini CLI Configuration Reference

**Source:** https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/configuration.md

## Settings File Locations
*   **User Settings:** `~/.gemini/settings.json`
*   **Workspace Settings:** `./.gemini/settings.json` (Overrides user settings)
*   **System Defaults (Linux):** `/etc/gemini-cli/system-defaults.json`
*   **System Overrides (Linux):** `/etc/gemini-cli/settings.json`

## Format: JSON

## Configuration Options

### 1. UI & Appearance
```json
{
  "ui": {
    "theme": "GitHub",
    "hideBanner": true,
    "hideTips": false,
    "showStatusInTitle": true
  }
}
```

### 2. Model & Session
```json
{
  "model": {
    "name": "gemini-1.5-pro-latest",
    "maxSessionTurns": 10,
    "summarizeToolOutput": {
      "run_shell_command": {
        "tokenBudget": 100
      }
    }
  }
}
```

### 3. Tool Management & Sandbox
```json
{
  "tools": {
    "sandbox": "docker",
    "discoveryCommand": "bin/get_tools",
    "callCommand": "bin/call_tool",
    "exclude": ["write_file"]
  }
}
```

### 4. MCP Servers
Configure Model Context Protocol (MCP) servers.
```json
{
  "mcpServers": {
    "mainServer": {
      "command": "bin/mcp_server.py",
      "args": ["--verbose"],
      "env": {
        "API_KEY": "$MY_API_TOKEN"
      },
      "cwd": "./server-dir",
      "timeout": 30000
    }
  },
  "mcp": {
    "allowed": ["mainServer"],
    "excluded": ["experimental"]
  }
}
```

### 5. Context & Discovery
```json
{
  "context": {
    "fileName": ["CONTEXT.md", "GEMINI.md"],
    "includeDirectories": ["./docs", "./src/types"],
    "loadFromIncludeDirectories": true,
    "fileFiltering": {
      "respectGitIgnore": true,
      "enableRecursiveFileSearch": true
    }
  },
  "advanced": {
    "excludedEnvVars": ["DEBUG", "NODE_ENV"]
  }
}
```

### 6. Telemetry
```json
{
  "telemetry": {
    "enabled": true,
    "target": "local",
    "otlpEndpoint": "http://localhost:4317",
    "logPrompts": false
  }
}
```
