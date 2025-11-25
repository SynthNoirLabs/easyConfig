# EasyConfig Project Context

## Project Overview
**EasyConfig** is a centralized Wails application designed to manage configuration files for various AI agents (Claude Code, Gemini CLI, Codex, Jules). It solves the problem of fragmented configuration ecosystems by providing a unified dashboard to find, edit, and extend agent capabilities, specifically focusing on Model Context Protocol (MCP) server injection.

## 🧠 Current Status (Nov 25, 2025)
**Milestone Reached:** Full Provider Support & Enhanced UX.

### ✅ Implemented Features
*   **Providers:** Added support for **OpenCode**, **Crush**, **Aider**, **Goose**, **Git**, and expanded **Claude** (Subagents) & **Gemini** (Extensions).
*   **UX:** Integrated **Monaco Editor**, **Toast Notifications**, and **Add/Delete Config** workflows.
*   **Backend:**
    *   **File Watcher:** Real-time updates via `fsnotify`.
    *   **MCP Injector:** Safe JSON modification logic.
    *   **Schema Fetcher:** Auto-downloads schemas for validation.
    *   **Path Helpers:** Centralized cross-platform logic in `pkg/util/paths`.

## Technology Stack
- **Backend:** Go (v1.24+)
- **Application Framework:** Wails (v2.11.0)
- **Frontend:** React (v18), TypeScript, Vite
- **Configuration Formats:** JSON, TOML, YAML, INI

## Directory Structure
- **`main.go`**: Application entry point.
- **`app.go`**: Main application logic (Wails bindings).
- **`pkg/`**:
  - **`config/`**: Providers and core Service logic.
  - **`watcher/`**: File system watcher service.
  - **`mcp/`**: MCP server injection logic.
  - **`schema/`**: JSON schema fetching logic.
  - **`util/paths/`**: Cross-platform path helpers.
- **`frontend/`**: React application.

## Development & Build Commands

### Prerequisites
- Go 1.24+
- Node.js 18+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0`)

### Running
```bash
wails dev
```

### Building
```bash
wails build
```