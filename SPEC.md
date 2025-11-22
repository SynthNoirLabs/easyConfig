# EasyConfig: AI Agent Configuration Hub

## 1. Core Philosophy
**"One Dashboard, Every Agent, Every Config."**
A centralized Wails application to manage the fragmented configuration ecosystems of CLI AI tools (Claude Code, Gemini CLI, Codex, Jules, etc.).

## 2. Key Capabilities

### A. Multi-Level Configuration Management
The app must handle the hierarchy of config files:
1.  **Global Scope:** User-level configs (e.g., `~/.claude/config.json`, `~/.gemini/config.json`).
2.  **Project Scope:** Local directory configs (e.g., `./.claude.json`, `.codex/config.toml`) found in developer workspaces.
3.  **Environment Scope:** Read-only view of active ENV vars (e.g., `ANTHROPIC_API_KEY`) that override files.

### B. Universal Editor
- **Formats:** Support reading/writing JSON, YAML, TOML.
- **Raw Mode:** Monaco-based text editor for direct control.
- **GUI Mode:** Form-based fields for common settings (Temperature, Model, Max Tokens).

### C. The "Smart" Layer (MCP & Skills)
A dedicated interface to extend agent capabilities:
- **MCP Injection:** "Add PostgreSQL Server" -> Automatically appends the correct JSON/TOML block to the agent's `mcpServers` config.
- **Skill Registry:** A searchable browser for:
    - MCP Servers (from `glama` or official lists).
    - Agent Skills/Tools (specific to the agent's plugin architecture).
    - Sub-agents (if supported).

### D. Future Roadmap (Registry)
- Online browsing of community MCPs.
- One-click installation of toolchains.

## 3. Architecture

### Backend (Go)
- **`DiscoveryService`:**
    - Scans standard paths (`~`, XDG config folders).
    - Scans specific "Dev" directories for project-level configs.
- **`ParserService`:**
    - robust unmarshaling/marshaling of JSON/TOML/YAML.
    - *Critical:* Attempt to preserve comments in YAML/TOML if possible.
- **`McpManager`:**
    - Logic to generate the correct config snippet for a given MCP server and Agent combination.

### Frontend (React + Shadcn UI)
- **Sidebar:** Grouped by Agent (Claude, Gemini, etc.).
- **Main View:**
    - **Tabs:** Global | Project 1 | Project 2.
    - **Content:** Config Editor & "Add Capability" Search Bar.

## 4. Supported Agents (Initial)
| Agent | Global Path | Project Path | Format |
| :--- | :--- | :--- | :--- |
| **Claude Code** | `~/.claude/config.json` | `cwd/claude.json` | JSON |
| **Gemini CLI** | `~/.gemini/settings.json` | `cwd/.gemini.json` | JSON |
| **Codex** | `~/.codex/config.toml` | `.codex/config.toml` | TOML |
| **Jules** | `~/.jules/config.yaml` | `AGENTS.md` | YAML/MD |

