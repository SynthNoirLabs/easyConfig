# Agent-Driven Development Architecture

This document details the architectural patterns and workflows established in `easyConfig`. It serves as a blueprint for other AI Agents (and their human counterparts) to understand how to build a repository that supports autonomous multi-agent collaboration.

---

## 1. The "Agent Routing" System (CI/CD)

We moved beyond simple "CI" to **Agent Orchestration**. Instead of a generic "AI fix," we use GitHub Labels to route tasks to the specific model best suited for the domain.

### Strategy
We implemented a label-based routing system in GitHub Actions (`.github/workflows/`).

*   **Trigger:** Apply a specific label to an Issue or Pull Request.
*   **The Specialists:**
    *   **`agent:claude`** $\rightarrow$ **Frontend & UI Specialist.**
        *   *Workflow:* `agent-interactive.yml`
        *   *Role:* React components, CSS styling, UX responsiveness, component logic.
    *   **`agent:gemini`** $\rightarrow$ **Architect & Doc Specialist.**
        *   *Workflow:* `gemini-agent.yml`
        *   *Role:* Issue triage, documentation updates, PR reviews, architectural validation.
    *   **`agent:codex`** $\rightarrow$ **Automation Engineer.**
        *   *Workflow:* `codex-agent.yml`
        *   *Role:* Shell scripts, CI pipeline fixes, auto-fixing failed unit tests.
    *   **`agent:jules`** $\rightarrow$ **Core Backend Developer.**
        *   *Platform:* Jules GitHub App.
        *   *Role:* Complex Go logic, refactoring, system-wide changes.

### Implementation Pattern
Each workflow uses a conditional check on the label name:
```yaml
if: |
  (github.event.action == 'labeled' && github.event.label.name == 'agent:claude') ||
  (contains(github.event.comment.body, '@claude'))
```

---

## 2. Integration Testing via Docker

To reliably test interactions with third-party CLI tools (Claude Code, Gemini CLI, Codex) without polluting the host machine or requiring user login, we built a containerized test environment.

### The Problem
AI Agents often hallucinate file paths or assume tools are installed.

### The Solution
1.  **`tests/integration/Dockerfile`**: A clean Ubuntu/Go environment that installs the *real* CLI binaries via `npm install -g ...`.
2.  **`scripts/run-integration.sh`**: A script that builds the container and mounts the repository into it.
3.  **`tests/integration/integration_test.go`**: Go tests that run *inside* the container. They mock the filesystem structure (`~/.claude/settings.json`) and verify that our application (`easyConfig`) correctly discovers and modifies these files in a real Linux environment.

**Takeaway:** Don't just unit test logic. Integration test the *environment* to ensure your agent's assumptions about the world are correct.

---

## 3. Backend Architecture: The Provider Pattern

The Go backend is designed for extensibility. Agents can add support for new tools without understanding the entire codebase.

*   **Interface (`pkg/config/types.go`):** Defines `Provider` with `Name()` and `Discover()`.
*   **Registry (`pkg/config/service.go`):** A central service that iterates over all registered providers.
*   **Safety First:**
    *   **Permissions:** `SaveConfig` enforces `0600` (User R/W only) permissions to protect API keys.
    *   **Validation:** JSON content is parsed and validated *before* writing to disk to prevent corruption.

---

## 4. Frontend Architecture: The Wails Bridge

*   **State Management:** `ConfigContext` acts as the bridge. It wraps the Wails backend calls (`window.go`) and provides a clean API (`readConfig`, `saveConfig`) to React components.
*   **Resilience:** The frontend is built to handle the absence of the backend (e.g., during browser-only development) by catching errors and logging them, rather than crashing the UI.
*   **Optimistic UI:** The `ConfigEditor` performs client-side validation (e.g., checking valid JSON) to provide instant feedback before engaging the backend.

---

## 5. Shared Context & Documentation

Agents cannot read your mind. They need explicit context.

*   **`TASKS.md`**: The Single Source of Truth for the backlog. All agents read this to know what to do next.
*   **`CONTRIBUTING_AGENTS.md`**: The Rulebook. It explicitly defines:
    *   Code Style (Go vs TS).
    *   Testing Requirements (Must run `go test`).
    *   Commit Message Format (`feat:`, `fix:`).

---

## Summary for Other Agents

If you are an AI working on this repo:
1.  **Read `TASKS.md`** first.
2.  **Check your Label:** Are you `agent:claude` working on frontend, or `agent:jules` on backend? Stay in your lane for maximum efficiency.
3.  **Use the Docker Test:** If you touch CLI integration logic, run `./scripts/run-integration.sh`.
4.  **Respect the Patterns:** Use the `Provider` interface for new tools. Use `ConfigContext` for UI state.
