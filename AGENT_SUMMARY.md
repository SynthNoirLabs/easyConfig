# AI Agent & CLI Integration Summary

This document consolidates all learned and implemented aspects regarding AI agents, their respective Command-Line Interfaces (CLIs), and their integration into the `easyConfig` project's GitHub Actions workflows. It serves as a comprehensive reference for understanding the current setup, capabilities, and areas for potential review or enhancement.

---

## 1. AI Agents Overview & Configured Roles

`easyConfig` currently supports the following AI agents, each with a defined role in the development workflow, primarily orchestrated via label-based routing:

| Agent | Primary Role (via `agent-triage.yml`) | Configuration File Paths (as discovered by `easyConfig`) |
| :---- | :------------------------------------ | :------------------------------------------------------- |
| **Claude Code** | **Frontend & UI Specialist.** (React/TypeScript, UI components, design, UX) | `~/.claude/settings.json`, `~/.claude/claude_desktop_config.json`, `/etc/claude-code/managed-settings.json`, `/etc/claude-code/managed-mcp.json`, `./.claude/settings.json`, `./.claude/settings.local.json`, `./CLAUDE.md` |
| **Gemini CLI** | **Architect & Docs Specialist.** (Issue triage, documentation, architectural validation, high-level code review) | `~/.gemini/settings.json`, `/etc/gemini-cli/system-defaults.json`, `/etc/gemini-cli/settings.json`, `./.gemini/settings.json`, `./GEMINI.md` |
| **Codex CLI** | **Automation Engineer.** (Shell scripts, CI pipeline fixes, auto-fixing, specialized review) | `~/.codex/config.toml`, `./.codex/config.toml` |
| **Google Jules** | **Core Backend Developer.** (Complex Go logic, refactoring, multi-file architectural changes) | `~/.jules-mcp/data.json`, `./AGENTS.md` |
| **GitHub Copilot** | *Integrated via IDE / GitHub.com*. (Code completion, in-editor assistance, PR summarization) | `./.github/copilot-instructions.md`, `~/.copilot/mcp-config.json` |
| **OpenAI (Generic)** | *Generic LLM capabilities used by other agents.* (For custom prompts/tools) | `~/.config/openai/config.yaml` |

---

## 2. GitHub Actions Workflows (CI/CD)

The repository leverages GitHub Actions for Continuous Integration, Continuous Delivery, and advanced Agent Orchestration.

### 2.1 Core CI Pipeline (`.github/workflows/ci.yml`)

This workflow defines the quality gates for all code merged into `main`.

*   **Triggers:** `push` to `main`, `pull_request` targeting `main`.
*   **Key Jobs/Checks:**
    *   **Dependency Review:** `actions/dependency-review-action@v4` checks for vulnerable dependencies.
    *   **Go Backend Quality (`test-backend`):**
        *   Runs `go test -v -race -coverprofile=coverage.out ./pkg/...`
        *   **Enforces >80% Code Coverage.**
        *   Runs `golangci-lint` for static analysis and linting.
        *   Runs `aquasecurity/trivy-action@master` for vulnerability scanning.
    *   **Frontend Quality (`test-frontend`):**
        *   Builds the React app (`npm run build`).
        *   Runs `Biome` for linting and formatting.
    *   **Integration Tests (`test-integration`):**
        *   Executes `./scripts/run-integration.sh`.
        *   Builds a Docker image (`tests/integration/Dockerfile`) with necessary CLIs installed.
        *   Runs Go integration tests (`tests/integration/integration_test.go`) inside the container to validate `easyConfig`'s interaction with real CLI file paths.
    *   **Application Build (`build-app`):** Compiles the Wails application for Linux, Windows, and macOS, requiring all preceding jobs to pass.

### 2.2 Agent Orchestration Workflows

These workflows define how AI agents are triggered and interact with the repository.

*   **`.github/workflows/agent-triage.yml`**
    *   **Purpose:** Automated issue labeling and initial task routing.
    *   **Trigger:** `issues` of type `opened`.
    *   **Agent:** Claude Code (`anthropics/claude-code-action@v1`).
    *   **Logic:** Claude analyzes the issue title/body and applies relevant category labels (e.g., `backend`, `feature`) AND assigns an AI agent using labels like `agent:jules`, `agent:claude`, `agent:gemini`, `agent:codex`.

*   **`.github/workflows/gemini-agent.yml`**
    *   **Purpose:** Gemini CLI for PR reviews, issue triage, and interactive chat.
    *   **Triggers:**
        *   `pull_request` (`opened`, `synchronize`, `reopened`, `labeled`).
        *   `issues` (`opened`, `labeled`).
        *   `issue_comment` (contains `@gemini-cli`).
    *   **Agent:** Gemini CLI (`google-github-actions/run-gemini-cli@v1`).
    *   **Logic:**
        *   `gemini-review` job runs on PRs or when `gemini-review` label is added.
        *   `gemini-triage` job runs on issue open or when `gemini-triage` or `agent:gemini` label is added.
        *   `gemini-chat` job runs when `@gemini-cli` is mentioned in comments.
    *   **Authentication:** `GEMINI_API_KEY` GitHub Secret.

*   **`.github/workflows/agent-interactive.yml`**
    *   **Purpose:** Interactive Claude Code assistance.
    *   **Triggers:** `issue_comment` (contains `@claude`), `issues` (`labeled`), `pull_request` (`labeled`).
    *   **Agent:** Claude Code (`anthropics/claude-code-action@v1`).
    *   **Logic:** The `claude-interactive` job runs when `@claude` is mentioned or `agent:claude` label is applied.
    *   **Authentication:** `ANTHROPIC_API_KEY` GitHub Secret.

*   **`.github/workflows/codex-agent.yml`**
    *   **Purpose:** Codex for PR reviews and auto-fixing CI failures.
    *   **Triggers:**
        *   `pull_request` (`opened`, `synchronize`, `labeled`).
        *   `workflow_run` (on `Wails CI/CD` workflow `completed` with `failure`).
    *   **Agent:** Codex CLI (`openai/codex-action@v1`).
    *   **Logic:**
        *   `codex-review` job runs on PRs or when `codex-review` or `agent:codex` label is added.
        *   `codex-autofix` job runs if the main CI workflow fails (currently a mock for actual fix logic, but scaffolded for Codex CLI).
    *   **Authentication:** `OPENAI_API_KEY` GitHub Secret.

### 2.3 Release & Maintenance Workflows

*   **`.github/workflows/release.yml`**
    *   **Purpose:** Automates application releases.
    *   **Trigger:** `push` to `tags` matching `v*` (e.g., `v1.0.0`).
    *   **Logic:** Builds the Wails application for multiple platforms and creates a GitHub Release with the compiled binaries.

*   **`.github/workflows/stale.yml`**
    *   **Purpose:** Automatically manages stale issues and pull requests.
    *   **Trigger:** `schedule` (daily).
    *   **Logic:** Marks issues/PRs as stale after a period of inactivity and closes them if they remain stale.

---

## 3. CLI Configuration & Discovery (`pkg/config/`)

The `easyConfig` backend is designed to discover and interact with the configuration files of these CLIs across different scopes.

### Core Principles
*   **Provider Pattern:** New CLIs are added by implementing the `Provider` interface.
*   **Scoped Discovery:** Configuration files are searched in:
    *   `Global`: User's home directory (e.g., `~/.claude/settings.json`).
    *   `Project`: Current project directory (e.g., `./.gemini/settings.json`).
    *   `System`: System-wide locations (e.g., `/etc/claude-code/managed-settings.json`).
*   **Secure I/O:**
    *   `SaveConfig` enforces `0o600` file permissions to protect sensitive data.
    *   `SaveConfig` performs validation (e.g., JSON syntax check) before writing.

### Local Reference Documentation (`docs/reference/`)

This directory contains detailed information, including example schemas and file paths, for each supported CLI. This localizes the knowledge needed for agents to understand and interact with the configuration files.

*   `claude_code.md`
*   `gemini_cli.md`
*   `codex_cli.md`
*   `jules_tools.md`

---

## 4. Authentication Strategy

*   **GitHub Secrets:** All API keys (`GEMINI_API_KEY`, `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`) are stored as GitHub Secrets for secure access by workflows.
*   **No Direct Login:** Workflows cannot perform interactive browser-based logins typical of personal AI subscriptions. Authentication relies on API keys or Workload Identity Federation for enterprise Google Cloud services.

---

## 5. Areas for Review / Further Development

*   **`ANTHROPIC_API_KEY`:** You **must manually set** this secret in the GitHub repository settings for Claude workflows to function (I could not retrieve it from your environment).
*   **Jules Integration:** While `easyConfig` can discover Jules' local data, explicit GitHub Actions for Jules beyond the Bot App are not configured. If a workflow-based trigger is desired, further investigation into Jules' API or action is needed.
*   **Codex `codex-autofix`:** The auto-fix logic for Codex in `codex-agent.yml` is currently a placeholder and needs to be fully implemented to fetch logs and apply fixes.
*   **`File Watcher` (Issue #35):** This will be a significant backend task for real-time config updates.
*   **`MCP Injector` (Issue #34):** Implementing the logic to manipulate MCP server configurations within JSON files.
*   **Frontend UI Enhancements (Issue #36):** Upgrading the editor to Monaco for a richer user experience.

---

This `AGENT_SUMMARY.md` provides a holistic view of the AI-native development environment you've built.
