# Project Tasks & Roadmap

## 🧠 Context for Agents
This file serves as the Master Backlog. Agents (Jules, Claude, etc.) should look here to pick up tasks.
**Rule:** When picking a task, verify if a corresponding GitHub Issue exists. If not, request one.

## 🚨 Priority: Infrastructure & Safety
- [x] **CI/CD Pipeline**: Update `.github/workflows/ci.yml` to include `go test` and `golangci-lint`.
- [ ] **Linter Setup**: Add `golangci-lint` configuration to project root.
- [ ] **Agent Guidelines**: Create `CONTRIBUTING_AGENTS.md` with instructions on how to run tests and format code.

## 🔌 Phase 2: Backend Providers (Go)
*See `docs/reference/` for specs.*

### 2.1 Gemini CLI Provider
- [ ] **Ref:** `docs/reference/gemini_cli.md`
- [ ] **Task**: Implement `GeminiProvider` in `pkg/config/providers.go`.
- [ ] **Req**: Parse JSON. Handle `~/.gemini/settings.json` and `./.gemini/settings.json`.

### 2.2 Codex CLI Provider
- [ ] **Ref:** `docs/reference/codex_cli.md`
- [ ] **Task**: Implement `CodexProvider` in `pkg/config/providers.go`.
- [ ] **Req**: Add TOML support (`github.com/pelletier/go-toml/v2`). Parse `config.toml`.

### 2.3 Git Provider
- [ ] **Task**: Implement `GitProvider` to read `.gitconfig`.
- [ ] **Req**: Use `go-ini` or standard git config parsing.

## 💾 Phase 3: Core Logic (Go)
- [ ] **Config I/O Service**: Implement `ReadConfig(path)` and `SaveConfig(path, content)` in `pkg/config/service.go`.
- [ ] **Format Detection**: Auto-detect JSON/YAML/TOML based on extension.
- [ ] **Validation**: Ensure saved JSON is valid before writing to disk.

## 🎨 Phase 4: Frontend Foundation (React)
- [ ] **Component Library**: Install `lucide-react` for icons and setup basic Tailwind utility classes.
- [ ] **Sidebar Navigation**: Create a Sidebar that lists "Discovered Agents" dynamically.
- [ ] **Editor Layout**: Create a split view: `[Sidebar | Editor Area]`.
- [ ] **Config Context**: Create a React Context to store the list of discovered configs.

## 🚀 Phase 5: Advanced Features
- [ ] **MCP Injector**: Logic to parse `mcpServers` block and inject a new server entry.
- [ ] **File Watcher**: Auto-reload configs when changed on disk (fsnotify).