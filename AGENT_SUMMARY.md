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

```yaml
name: Wails CI/CD

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

env:
  NODE_OPTIONS: "--max-old-space-size=4096"

jobs:
  # 1. Security & Dependencies
  dependency-review:
    name: "Dependency Review"
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Dependency Review
        uses: actions/dependency-review-action@v4

  # 2. Backend Quality
  test-backend:
    name: "Go Backend (Test & Lint)"
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
          check-latest: true
          cache: true

      - name: Install Dependencies
        run: go mod tidy

      - name: Create Dummy Frontend Assets
        run: |
          mkdir -p frontend/dist
          touch frontend/dist/index.html

      - name: Run Unit Tests with Coverage
        run: go test -v -race -coverprofile=coverage.out ./pkg/...

      - name: Check Code Coverage
        run: |
          # Extract total coverage percentage (e.g., "90.4%")
          COVERAGE_Output=$(go tool cover -func=coverage.out | grep total:)
          echo "$COVERAGE_Output"
          
          # Parse percentage number
          PERCENT=$(echo "$COVERAGE_Output" | awk '{print $3}' | sed 's/%//')
          
          echo "Coverage: ${PERCENT}%"
          
          # Compare using awk for float comparison
          if awk "BEGIN {exit !($PERCENT < 80)}"; then
            echo "Error: Code coverage ($PERCENT%) is below the 80% threshold."
            exit 1
          fi

      - name: Run Linter
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

      - name: Run Trivy Vulnerability Scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          ignore-unfixed: true
          format: 'table'
          exit-code: '0' # Warning only for now to prevent build block
          severity: 'CRITICAL,HIGH'

  # 3. Frontend Quality
  test-frontend:
    name: "Frontend (Build & Typecheck)"
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./frontend
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: '18.x'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json

      - name: Install Dependencies
        run: npm ci || npm install

      - name: Setup Biome
        uses: biomejs/setup-biome@v2

      - name: Run Biome (Lint)
        run: biome ci .

      - name: Type Check (TSC)
        run: npm run build

  # 4. Integration Tests
  test-integration:
    name: "Integration Tests (Docker)"
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Run Integration Tests
        run: |
          chmod +x scripts/run-integration.sh
          ./scripts/run-integration.sh

  # 5. Full Application Build
  build-app:
    name: "Build Wails Binary"
    needs: [test-backend, test-frontend, test-integration]
    strategy:
      fail-fast: false
      matrix:
        build:
          - name: 'Linux'
            platform: 'linux/amd64'
            os: 'ubuntu-latest'
          - name: 'Windows'
            platform: 'windows/amd64'
            os: 'windows-latest'
          - name: 'macOS'
            platform: 'darwin/universal'
            os: 'macos-latest'

    runs-on: ${{ matrix.build.os }}

    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive

      - name: Build Wails App
        uses: dAppServer/wails-build-action@v3
        id: build
        with:
          build-name: 'easyConfig'
          build-platform: ${{ matrix.build.platform }}
          wails-version: 'v2.9.2'
          go-version: '1.23'
          node-version: '18.x'
          package: true

      - name: Upload Artifacts
        uses: actions/upload-artifact@v5
        with:
          name: ${{ matrix.build.name }}-build
          path: build/bin/*
```

### 2.2 Agent Orchestration Workflows

These workflows define how AI agents are triggered and interact with the repository.

#### 2.2.1 Agent Issue Triage (`.github/workflows/agent-triage.yml`)
*   **Purpose:** Automated issue labeling and initial task routing.
*   **Trigger:** `issues` of type `opened`.
*   **Agent:** Claude Code (`anthropics/claude-code-action@v1`).
*   **Logic:** Claude analyzes the issue title/body and applies relevant category labels (e.g., `backend`, `feature`) AND assigns an AI agent using labels like `agent:jules`, `agent:claude`, `agent:gemini`, `agent:codex`.

```yaml
name: Agent Issue Triage

on:
  issues:
    types: [opened]

permissions:
  issues: write
  contents: read

jobs:
  triage:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Run Claude Code for Triage
      uses: anthropics/claude-code-action@v1
      with:
        anthropic_api_key: ${{ secrets.ANTHROPIC_API_KEY }}
        prompt: |
          Analyze the issue title and body provided in the GitHub event context.
          
          1. Determine appropriate category labels (e.g., 'bug', 'feature', 'backend', 'frontend', 'CI', 'devops').
          
          2. Assign an AI Agent by applying ONE of the following labels based on the task type:
             - 'agent:jules': For complex backend logic, refactoring, or multi-file architectural changes.
             - 'agent:claude': For Frontend (React/TypeScript), UI components, and design work.
             - 'agent:gemini': For documentation, general triage, or simple logic checks.
             - 'agent:codex': For CI/CD scripts, automation, or one-off script fixes.
          
          Apply the selected labels using the available GitHub tools.
```

#### 2.2.2 Gemini Agent (`.github/workflows/gemini-agent.yml`)
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

```yaml
name: Gemini Agent

on:
  pull_request:
    types: [opened, synchronize, reopened, labeled]
  issues:
    types: [opened, labeled]
  issue_comment:
    types: [created]

permissions:
  contents: write
  pull-requests: write
  issues: write

jobs:
  # 1. Automated PR Review
  # Runs on PR Open, Sync, Reopen, OR when label 'gemini-review' is added
  gemini-review:
    if: |
      (github.event_name == 'pull_request') && 
      (
        github.event.action == 'opened' || 
        github.event.action == 'synchronize' || 
        github.event.action == 'reopened' || 
        (github.event.action == 'labeled' && github.event.label.name == 'gemini-review')
      )
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Gemini Code Review
        uses: google-github-actions/run-gemini-cli@v1
        with:
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
          command: 'review'

  # 2. Automated Issue Triage
  # Runs on Issue Open OR when label 'gemini-triage' or 'agent:gemini' is added
  gemini-triage:
    if: |
      (github.event_name == 'issues') && 
      (
        github.event.action == 'opened' || 
        (github.event.action == 'labeled' && (github.event.label.name == 'gemini-triage' || github.event.label.name == 'agent:gemini'))
      )
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Gemini Issue Triage
        uses: google-github-actions/run-gemini-cli@v1
        with:
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
          command: 'triage'

  # 3. Interactive Chat (@gemini-cli)
  gemini-chat:
    if: github.event_name == 'issue_comment' && contains(github.event.comment.body, '@gemini-cli')
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Gemini Chat
        uses: google-github-actions/run-gemini-cli@v1
        with:
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
          command: 'chat'
          issue_number: ${{ github.event.issue.number }}
          comment_body: ${{ github.event.comment.body }}
```

#### 2.2.3 Claude Interactive (`.github/workflows/agent-interactive.yml`)
*   **Purpose:** Interactive Claude Code assistance.
*   **Triggers:** `issue_comment` (contains `@claude`), `issues` (`labeled`), `pull_request` (`labeled`).
*   **Agent:** Claude Code (`anthropics/claude-code-action@v1`).
*   **Logic:** The `claude-interactive` job runs when `@claude` is mentioned or `agent:claude` label is applied.
*   **Authentication:** `ANTHROPIC_API_KEY` GitHub Secret.

```yaml
name: Agent Interactive (On-Demand)

on:
  issue_comment:
    types: [created]
  issues:
    types: [labeled]
  pull_request:
    types: [labeled]

permissions:
  contents: write
  issues: write
  pull-requests: write

jobs:
  claude-interactive:
    if: |
      (github.event_name == 'issue_comment' && contains(github.event.comment.body, '@claude')) ||
      (github.event.action == 'labeled' && github.event.label.name == 'agent:claude')
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Run Claude Code Action
      uses: anthropics/claude-code-action@v1
      with:
        anthropic_api_key: ${{ secrets.ANTHROPIC_API_KEY }}
```

#### 2.2.4 Codex Agent (`.github/workflows/codex-agent.yml`)
*   **Purpose:** Codex for PR reviews and auto-fixing CI failures.
*   **Triggers:**
    *   `pull_request` (`opened`, `synchronize`, `labeled`).
    *   `workflow_run` (on `Wails CI/CD` workflow `completed` with `failure`).
*   **Agent:** Codex CLI (`openai/codex-action@v1`).
*   **Logic:**
    *   `codex-review` job runs on PRs or when `codex-review` or `agent:codex` label is added.
    *   `codex-autofix` job runs if the main CI workflow fails (currently a mock for actual fix logic, but scaffolded for Codex CLI).
*   **Authentication:** `OPENAI_API_KEY` GitHub Secret.

```yaml
name: Codex Agent

on:
  pull_request:
    types: [opened, synchronize, labeled]
  workflow_run:
    workflows: ["Wails CI/CD"]
    types: [completed]

permissions:
  contents: write
  pull-requests: write
  issues: write

jobs:
  # 1. PR Code Review
  codex-review:
    if: |
      (github.event_name == 'pull_request') &&
      (
        github.event.action == 'opened' || 
        github.event.action == 'synchronize' || 
        (github.event.action == 'labeled' && (github.event.label.name == 'codex-review' || github.event.label.name == 'agent:codex'))
      )
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Run Codex Review
        uses: openai/codex-action@v1
        continue-on-error: true
        with:
          openai-api-key: ${{ secrets.OPENAI_API_KEY }}
          prompt: "Review this pull request. Identify bugs, security issues, and style violations."
          safety-strategy: read-only

  # 2. Auto-Fix CI Failures
  codex-autofix:
    if: github.event.workflow_run.conclusion == 'failure'
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4
        with:
          ref: ${{ github.event.workflow_run.head_sha }}

      - name: Install Codex CLI
        run: npm install -g @openai/codex

      - name: Analyze and Fix
        env:
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
        run: |
          # In a real scenario, we'd fetch the logs from the failed run
          echo "Analyzing failure for run ${{ github.event.workflow_run.id }}..."
          
          # Run codex fix (hypothetical CLI command based on standard patterns)
          # codex fix --auto-commit --push --branch "codex-fix/${{ github.run_id }}"
          
          echo "Auto-fix logic placeholder. Ensure @openai/codex is configured."
```

### 2.3 Release & Maintenance Workflows

#### 2.3.1 Release (`.github/workflows/release.yml`)
*   **Purpose:** Automates application releases.
*   **Trigger:** `push` to `tags` matching `v*` (e.g., `v1.0.0`).
*   **Logic:** Builds the Wails application for multiple platforms and creates a GitHub Release with the compiled binaries.

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    name: "Build & Release"
    strategy:
      matrix:
        platform: [linux/amd64, windows/amd64, darwin/universal]
        include:
          - platform: linux/amd64
            os: ubuntu-latest
          - platform: windows/amd64
            os: windows-latest
          - platform: darwin/universal
            os: macos-latest

    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
        with:
          submodules: recursive

      - name: Build Wails App
        uses: dAppServer/wails-build-action@v3
        id: build
        with:
          build-name: 'easyConfig'
          build-platform: ${{ matrix.platform }}
          wails-version: 'v2.9.2'
          go-version: '1.23'
          node-version: '18.x'
          package: true

      - name: Release
        uses: softprops/action-gh-release@v1
        if: startsWith(github.ref, 'refs/tags/')
        with:
          files: build/bin/*
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### 2.3.2 Stale (`.github/workflows/stale.yml`)
*   **Purpose:** Automatically manages stale issues and pull requests.
*   **Trigger:** `schedule` (daily).
*   **Logic:** Marks issues/PRs as stale after a period of inactivity and closes them if they remain stale.

```yaml
name: "Close stale issues and PRs"
on:
  schedule:
  - cron: "0 0 * * *"

jobs:
  stale:
    runs-on: ubuntu-latest
    permissions:
      issues: write
      pull-requests: write
    steps:
    - uses: actions/stale@v9
      with:
        repo-token: ${{ secrets.GITHUB_TOKEN }}
        stale-issue-message: 'This issue is stale because it has been open 60 days with no activity. Remove stale label or comment or this will be closed in 7 days.'
        stale-pr-message: 'This PR is stale because it has been open 45 days with no activity. Remove stale label or comment or this will be closed in 7 days.'
        stale-issue-label: 'no-issue-activity'
        stale-pr-label: 'no-pr-activity'
        days-before-stale: 60
        days-before-close: 7
```

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