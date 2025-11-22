# Project: AI Agent Mission Control (Wails v2)

## 1. Project Manifesto
**Goal:** Build a "Mission Control" dashboard for local AI CLI tools (Claude Code, Gemini CLI, Codex CLI, Jules).
**Philosophy:** "Configuration over Execution." The app reads/edits the *settings* (MCP servers, Permissions) of these agents.
**Workflow:** Managed via `ccpm`. All features must be implemented as Epics/Tasks.

## 2. Technology Stack
- **Runtime:** Wails v2 (Go 1.23+)
- **Frontend:** React 18 + TypeScript + Vite
- **UI:** Shadcn UI + Tailwind CSS (Pre-installed in template)
- **Linting:** Biome (`biome.json`), golangci-lint (`.golangci.yml`)

## 3. Domain Logic: "Smart Discovery"
The app must scan the OS for known config files. Do NOT hardcode user paths.

| Agent | Config Format | Locations (Priority Order) | Docs |
| :--- | :--- | :--- | :--- |
| **Claude Code** | JSON | `~/.claude/claude_desktop_config.json`<br>`~/.claude/settings.json`<br>`**/.claude/settings.json` (Recursive) | [Docs](https://docs.claude.com/code/reference/configuration) |
| **Claude Web** | JSON | Check `~/.claude/settings.json` for `env.CLAUDE_CODE_ON_THE_WEB` keys. | [Web Docs](https://code.claude.com/docs/en/claude-code-on-the-web) |
| **Gemini** | JSON | `~/.gemini/settings.json`<br>`~/.config/google/ai-studio/config.json` | [Docs](https://ai.google.dev/gemini-api/docs/cli) |
| **Codex** | TOML | `~/.codex/config.toml`<br>`.codex/config.toml` (Project scope) | [Docs](https://platform.openai.com/docs/guides/codex-cli) |
| **Jules** | JSON/Markdown | `AGENTS.md` (Project Root)<br>`~/.jules-mcp/data.json`<br>Env: `JULES_DATA_PATH` | [Jules Docs](https://jules.google/docs/) |

## 4. Architecture & Milestones

### Milestone 1: Core Engine (Go)
- **Structs:** Create strict Go structs for all agents.
- **Discovery:** Implement `DiscoverConfigs()` scanning the paths above.
- **Jules Support:** Parse `AGENTS.md` to extract "Context" and "Tools" sections.
- **Safety:** Use `map[string]interface{}` for unknown fields to prevent data loss on save.

### Milestone 2: Dashboard (React)
- **Sidebar:** List detected agents.
- **Status Cards:** Show "Found" vs "Not Configured".

### Milestone 3: Editors
- **Forms:** User-friendly toggles for common settings.
- **MCP Table:** CRUD interface for MCP Servers (translates to JSON/TOML).
- **Raw Editor:** Monaco/Textarea fallback.

## 5. Creative Freedom & Self-Improvement
**You are an intelligent engineer, not a script.**
- **Improve the Prompt:** If you find a better way to structure this `PROMPT.md` as you work, propose updates.
- **Refine the UI:** Don't just build a "list". If Shadcn has a better component (e.g., `Accordion` for nested configs), use it.
- **Suggest Features:** If you see a missing feature (e.g., "Backup Configs before Save"), add it to the backlog.
