# Project Tasks & Roadmap

## 🧠 Context for Agents
This file serves as the Master Backlog. Agents (Jules, Claude, etc.) should look here to pick up tasks.
**Rule:** When picking a task, verify if a corresponding GitHub Issue exists. If not, request one.

## 🚨 Priority: Infrastructure & Safety
- [x] **CI/CD Pipeline**: Update `.github/workflows/ci.yml` to include `go test` and `golangci-lint`.
- [x] **Linter Setup**: Add `golangci-lint` configuration to project root.
- [x] **Agent Guidelines**: Create `CONTRIBUTING_AGENTS.md` with instructions on how to run tests and format code.

## 🔌 Phase 2: Backend Providers (Go)
*See `docs/reference/` for specs.*

### 2.1 Gemini CLI Provider
- [x] **Ref:** `docs/reference/gemini_cli.md`
- [x] **Task**: Implement `GeminiProvider` in `pkg/config/providers.go`.
- [x] **Req**: Parse JSON. Handle `~/.gemini/settings.json` and `./.gemini/settings.json`.

### 2.2 Codex CLI Provider
- [x] **Ref:** `docs/reference/codex_cli.md`
- [x] **Task**: Implement `CodexProvider` in `pkg/config/providers.go`.
- [x] **Req**: Add TOML support (`github.com/pelletier/go-toml/v2`). Parse `config.toml`.

### 2.3 Git Provider
- [x] **Task**: Implement `GitProvider` to read `.gitconfig`.
- [x] **Req**: Use `go-ini` or standard git config parsing.

### 2.4 New Providers (Expansion)
- [x] **Amazon Q**: Implement `pkg/config/provider_amazonq.go` (`~/.aws/amazonq/mcp.json`).
- [x] **Cursor**: Implement `pkg/config/provider_cursor.go` (`~/.cursor/cli-config.json`).
- [x] **Windsurf**: Implement `pkg/config/provider_windsurf.go` (`~/.config/Windsurf/User/settings.json`).
- [x] **Continue**: Implement `pkg/config/provider_continue.go` (`~/.continue/config.yaml`).
- [x] **Mentat**: Implement `pkg/config/provider_mentat.go` (`.mentat/` and `~/.mentat/.env`).
- [x] **Sweep**: Implement `pkg/config/provider_sweep.go` (`.sweep.yaml`).
- [x] **Plandex**: Implement `pkg/config/provider_plandex.go` (`.plandex/`).
- [x] **OpenHands**: Implement `pkg/config/provider_openhands.go` (`config.toml`).

## 💾 Phase 3: Core Logic (Go)
- [x] **Config I/O Service**: Implement `ReadConfig(path)` and `SaveConfig(path, content)` in `pkg/config/service.go`.
- [x] **Format Detection**: Auto-detect JSON/YAML/TOML based on extension.
- [x] **Validation**: Ensure saved JSON is valid before writing to disk.

## 🎨 Phase 4: Frontend Foundation (React)
- [x] **Component Library**: Install `lucide-react` for icons and setup basic Tailwind utility classes.
- [x] **Sidebar Navigation**: Create a Sidebar that lists "Discovered Agents" dynamically.
- [x] **Editor Layout**: Create a split view: `[Sidebar | Editor Area]`.
- [x] **Config Context**: Create a React Context to store the list of discovered configs.
- [x] **Config Editor**: Create a text editor component to view/edit config content.

## 🚀 Phase 5: Advanced Features
- [x] **MCP Injector**: Logic to parse `mcpServers` block and inject a new server entry.
- [x] **File Watcher**: Auto-reload configs when changed on disk (fsnotify).
- [ ] **Auto-Schema Fetcher**: Scrape official docs/repos to update local references. (Partially implemented)
- [x] **Add Config Wizard**: Create new config files from templates via UI.
- [x] **Toast Notifications**: Replace alerts with non-blocking notifications.
- [x] **Cross-Platform Paths**: Refactor path logic into `pkg/util/paths`.

## 🧩 Phase 6: Specialized Configs
- [x] **Claude Extras**: Discover Subagents (`agents/*.md`), Hooks, and Custom Commands.
- [x] **Gemini Extensions**: Discover standalone extension configurations.