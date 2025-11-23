# Codex CLI Configuration Reference

**Source:** https://developers.openai.com/codex/local-config/
**See Also:** https://vladimirsiedykh.com/blog/codex-mcp-config-toml-shared-configuration-cli-vscode-setup-2025

## Settings File Locations
*   **Shared Config:** `~/.codex/config.toml` (Shared between CLI and VSCode extension)
*   **Project Config:** `./.codex/config.toml` (Overrides global)

## Format: TOML

## Key Settings

### 1. Model Selection
```toml
model = "gpt-4o"
model_provider = "openai"
# Optional: Azure settings
[model_providers.azure]
name = "Azure OpenAI"
base_url = "..."
```

### 2. Sandboxing
```toml
[sandbox]
enabled = true
```

### 3. MCP Servers
Codex uses a shared TOML format for MCP servers.
```toml
[mcp.servers.filesystem]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-filesystem", "/Users/me/desktop"]

[mcp.servers.github]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-github"]
env = { "GITHUB_TOKEN" = "..." }
```

## Context
*   **Prompts:** `~/.codex/prompts/*.md`
*   **Project Context:** `AGENTS.md` (in project root)

## Example `config.toml`
```toml
model = "gpt-5-codex"
model_provider = "openai"

[sandbox]
enabled = true

[mcp.servers.memory]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-memory"]
```