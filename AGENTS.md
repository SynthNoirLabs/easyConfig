# AI Agents Context & Architecture Guide

**Goal:** Build a "Mission Control" dashboard for local AI CLI tools (Claude Code, Gemini, Codex, Jules).

## 📂 Documentation & Reference
We have gathered official configuration references for the supported CLIs. **Always check these files before guessing config paths or formats.**
*   [Claude Code Reference](docs/reference/claude_code.md)
*   [Gemini CLI Reference](docs/reference/gemini_cli.md)
*   [Codex CLI Reference](docs/reference/codex_cli.md)
*   [Jules Reference](docs/reference/jules_tools.md)

## 🏗️ Core Architecture: The Provider Pattern
The backend (`pkg/config`) uses a **Provider Pattern** to discover configurations.

### 1. The `Provider` Interface
Each CLI tool has a dedicated struct (e.g., `ClaudeProvider`, `GeminiProvider`) that implements:
```go
type Provider interface {
    Name() string
    Discover(projectPath string) ([]ConfigItem, error)
}
```

### 2. Scopes
Configurations are found in specific scopes. You must correctly identify them:
*   **`ScopeGlobal`**: User's home directory (e.g., `~/.claude/settings.json`).
*   **`ScopeProject`**: The repository root (e.g., `./.codex/config.toml`).
*   **`ScopeSystem`**: System-wide paths (rare, but possible).

## 🚀 How to Contribute
If you are an agent (Jules, Claude, etc.) assigned to a task:
1.  **Read `TASKS.md`** to find your assigned item.
2.  **Check `docs/reference/`** for the specific tool's config details.
3.  **Implement the Provider** in `pkg/config/providers.go` (or a new file if it gets large).
4.  **Register the Provider** in `app.go` (or wherever `NewDiscoveryService` is called).
5.  **Test** your provider using `go test ./pkg/config/...`.

## ⚠️ Important Rules
*   **Do NOT hardcode paths** if a helper exists (use `GetUserHome()`).
*   **Handle missing files gracefully**. It is NOT an error if a config file is missing; just return an empty list.
*   **Respect existing conventions**. Match the style of `pkg/config/providers.go`.