# AI Agents Context & Architecture Guide

**Goal:** Build a "Mission Control" dashboard for local AI CLI tools (Claude Code, Gemini, Codex, Jules).

## 🧠 Status Update (Nov 25, 2025)
The provider ecosystem has been significantly expanded. We now support 10+ AI tools.

### ✅ Active Providers
The following providers are fully implemented in `pkg/config/providers.go`:
1.  **Claude Code:** `settings.json`, `CLAUDE.md`, Subagents (`agents/*.md`).
2.  **Gemini CLI:** `settings.json`, `GEMINI.md`, Extensions (`extensions/`).
3.  **GitHub Copilot:** `config.json`, `mcp-config.json`, Instructions.
4.  **Codex CLI:** `config.toml`.
5.  **OpenAI:** `config.yaml`.
6.  **Jules:** `data.json`, `AGENTS.md`.
7.  **OpenCode:** `opencode.json`.
8.  **Crush:** `crush.json`, `.crushignore`.
9.  **Aider:** `.aider.conf.yml`.
10. **Goose:** `config.yaml`.
11. **Git:** `.gitconfig`.

## 📂 Documentation & Reference
*   [Claude Code Reference](docs/reference/claude_code.md)
*   [Gemini CLI Reference](docs/reference/gemini_cli.md)
*   [Codex CLI Reference](docs/reference/codex_cli.md)
*   [Jules Reference](docs/reference/jules_tools.md)

## 🏗️ Core Architecture: The Provider Pattern
The backend (`pkg/config`) uses a **Provider Pattern** to discover configurations.

### 1. The `Provider` Interface
Each CLI tool has a dedicated struct implementing:
```go
type Provider interface {
    Name() string
    Discover(projectPath string) ([]Item, error)
    Create(scope Scope, projectPath string) (string, error)
}
```

### 2. Scopes
*   **`ScopeGlobal`**: User's home directory.
*   **`ScopeProject`**: The repository root.
*   **`ScopeSystem`**: System-wide paths (e.g., `/etc`).

## 🚀 How to Contribute
1.  **Read `TASKS.md`** for the backlog.
2.  **Implement** new logic in `pkg/config` or `pkg/mcp`.
3.  **Test** using `go test ./...`.
4.  **Lint** using `golangci-lint` and `biome`.

## ⚠️ Important Rules
*   **Use `pkg/util/paths`**: Always use `paths.GetHomeDir()` or `paths.GetConfigDir()`. **Do NOT hardcode OS paths.**
*   **Handle missing files gracefully**: Discovery should not fail if a file is missing.
*   **Format:** Follow existing conventions.
