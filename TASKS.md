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
- [ ] **Task**: Implement `GitProvider` to read `.gitconfig`.
- [ ] **Req**: Use `go-ini` or standard git config parsing.

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
- [ ] **MCP Injector**: Logic to parse `mcpServers` block and inject a new server entry. (In Progress - Issue #34)
- [ ] **File Watcher**: Auto-reload configs when changed on disk (fsnotify). (In Progress - Issue #35)
- [ ] **Auto-Schema Fetcher**: Scrape official docs/repos to update local references.
- [ ] **Add Config Wizard**: Create new config files from templates via UI.
- [ ] **Toast Notifications**: Replace alerts with non-blocking notifications.
- [x] **Cross-Platform Paths**: Refactor path logic into `pkg/util/paths`. (Completed)

## 🧩 Phase 6: Specialized Configs
- [ ] **Claude Extras**: Discover Subagents (`agents/*.md`), Hooks, and Custom Commands.
- [ ] **Gemini Extensions**: Discover standalone extension configurations (if applicable).