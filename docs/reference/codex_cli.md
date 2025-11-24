# Codex CLI Configuration Reference

**Source:** https://github.com/openai/codex/blob/main/docs/config.md

## Settings File Locations
*   **Shared Config:** `~/.codex/config.toml` (Global user config)
*   **Project Config:** `./.codex/config.toml` (Overrides global config)

## Format: TOML

## Configuration Options

### 1. Model & Provider
```toml
model = "o3"
model_provider = "openai"

[model_providers.openai]
name = "OpenAI"
base_url = "https://api.openai.com/v1"
env_key = "OPENAI_API_KEY"
request_max_retries = 4
stream_max_retries = 10
```

### 2. Profiles
Define preset configurations switchable via `--profile`.
```toml
[profiles.dev]
model = "gpt-4o"
approval_policy = "always"

[profiles.ci]
model = "o3"
approval_policy = "never"
```

### 3. MCP Servers
Configure tools using the Model Context Protocol.
```toml
[mcp.servers.filesystem]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-filesystem", "/Users/me/desktop"]

[mcp.servers.github]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-github"]
env = { "GITHUB_TOKEN" = "..." }
cwd = "/path/to/run"
```

### 4. Sandboxing & Environment
```toml
[sandbox]
enabled = true

[shell_environment_policy]
inherit = "none"
set = { PATH = "/usr/bin", TERM = "xterm-256color" }
```

### 5. Tool Management
```toml
startup_timeout_sec = 20
tool_timeout_sec = 30
enabled_tools = ["search", "summarize"]
disabled_tools = ["dangerous_tool"]
```
