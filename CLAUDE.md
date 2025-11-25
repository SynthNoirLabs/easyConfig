# Claude Code Project Context

## Project: EasyConfig
**EasyConfig** is a unified dashboard for managing configuration files for AI CLI tools (Claude Code, Gemini, Codex, Jules, Aider, Goose, etc.). It allows discovery, creation, editing, and deletion of configs across Global and Project scopes.

## 🧠 Current Status (Nov 25, 2025)
The project has reached a major milestone with the completion of the initial feature set.

### ✅ Completed Features
- **Frontend UX:**
  - **Monaco Editor:** Replaced basic textarea with rich code editor (syntax highlighting).
  - **Notifications:** Replaced alerts with `sonner` toast notifications.
  - **Wizards:** "Add Config" modal to create new configurations from templates.
  - **Actions:** Delete, Reset, and Reload configuration files.
- **Backend Providers:**
  - **Supported Tools:** Claude Code (including Subagents), Gemini CLI (including Extensions), GitHub Copilot, Codex, OpenAI, Jules, OpenCode, Crush, Aider, Goose, Git.
  - **Cross-Platform:** robust path handling via `pkg/util/paths`.
- **Core Services:**
  - **File Watcher:** `fsnotify` service to auto-reload changes from disk.
  - **MCP Injector:** Logic to safely inject `mcpServers` into JSON configs.
  - **Schema Fetcher:** Automated retrieval of JSON schemas for validation.

## 🛠 Architecture
- **Stack:** Go (Wails) + React/TypeScript.
- **Pattern:** Provider-based discovery in `pkg/config/providers.go`.
- **Context:** `ConfigContext.tsx` manages state and Wails bridge.

## 🚀 Next Steps for Agents
*   Maintain and update providers as tool locations change.
*   Expand "Auto-Schema Fetcher" to support more tools.
*   Refine MCP injection logic for complex nested configs.
