# Agent Implementation Plan

**Author:** Principal AI Architect and Research Analyst
**Date:** November 23, 2025
**Objective:** To define the foundational strategy, CI/CD modernization, local development standardization, and the dual-mission architecture for AI agent integration within the `easyConfig` project (Wails/Go/React).

-----

## 1. Foundation Strategy: The Bedrock

A robust, standardized foundation is essential for a Wails/Go/React project, facilitating both human and AI collaboration.

### 1.1. CI/CD Modernization (GitHub Actions)

We will optimize the CI pipeline (`.github/workflows/ci.yml`) for speed, security, and reliability, adhering to 2025 best practices.

#### 1.1.1. Optimized Caching

We will leverage the integrated caching mechanisms provided by the setup actions, which is the recommended approach for optimal performance.

```yaml
# Snippet from .github/workflows/ci.yml

      - name: Setup Go
        uses: actions/setup-go@v6 # Use latest major version
        with:
          go-version-file: 'go.mod'
          cache: true # Caches Go modules (GOMODCACHE) and build cache (GOCACHE)

      - name: Setup Node.js
        uses: actions/setup-node@v6 # Use latest major version
        with:
          node-version-file: 'frontend/package.json'
          cache: 'npm'
          cache-dependency-path: 'frontend/package-lock.json'
```

#### 1.1.2. Linting and Formatting

**Backend (Go 1.24+): GolangCI-Lint**
Configuration optimized for modern Go, balancing thoroughness with performance.

  * **Configuration:** `.golangci.yml`

```yaml
run:
  timeout: 5m
  go: '1.24'

linters-settings:
  gci:
    # Ensure imports are grouped: standard, default, and local prefixes
    local-prefixes: github.com/easyConfig/easyConfig # Adjust module name if needed
  gofumpt:
    extra-rules: true # Enforce stricter formatting
  govet:
    check-shadowing: true
  gocyclo:
    min-complexity: 15
  misspell:
    locale: US

linters:
  # Use a 'disable-all' approach for explicit configuration
  disable-all: true
  enable:
    # Core/Correctness
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    # Style, Formatting, and Best Practices
    - bodyclose
    - exportloopref
    - gci
    - gocyclo
    - gofumpt
    - goimports
    - misspell
    - revive
    - unconvert
    - whitespace
    # Security
    - gosec
```

  * **CI Execution:**

```yaml
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.60.x # Use latest stable version
```

**Frontend (React/TS): Biome.js**
Biome is the standard, high-performance replacement for ESLint/Prettier.

  * **Configuration:** `biome.json`.
  * **CI Execution:**

```yaml
      - name: Setup Biome
        uses: biomejs/setup-biome@v2

      - name: Run Biome (Lint and Format Check)
        run: biome ci frontend/
```

#### 1.1.3. Security

We will implement a multi-layered "shift-left" security approach.

1.  **SCA:** Enable **Dependabot** via `.github/dependabot.yml`.

2.  **SAST:** Implement **GitHub CodeQL** (`.github/workflows/codeql.yml`).

3.  **Secret Scanning:** Ensure GitHub Secret Scanning and **Push Protection** are enabled in repository settings.

4.  **Dependency Review:** Block PRs introducing vulnerable dependencies.

    ```yaml
    # Add to ci.yml
    - name: 'Dependency Review'
      uses: actions/dependency-review-action@v5
    ```

5.  **Vulnerability & Misconfiguration Scanning:** Use **Trivy**.

```yaml
      - name: Run Trivy scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scanners: 'vuln,secret,misconfig'
          format: 'sarif'
          output: 'trivy-results.sarif'
          severity: 'CRITICAL,HIGH'

      - name: Upload Trivy results to GitHub Security
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'
```

### 1.2. Local Development Experience

We recommend **Mise** for tool version management and **Taskfile** for task execution.

#### 1.2.1. Tool Management (Mise)

Mise ensures consistent runtimes (Go, Node) across all environments.

  * **Configuration (`.mise.toml` at root):**

    ```toml
    [tools]
    go = '1.24.x'
    node = '22.x'
    task = 'latest'
    # Ensure consistent linter versions
    golangci-lint = '1.60.0'
    ```

#### 1.2.2. Task Runner (Taskfile)

Taskfile provides robust, cross-platform task execution.

  * **Configuration (`Taskfile.yml` at root):**

    ```yaml
    version: '3'

tasks:
  setup:
    desc: Install dependencies and setup environment
    cmds:
      - mise install
      - go mod download
      - npm install --prefix frontend
      - pre-commit install

  dev:
    desc: Run the application in development mode
    cmds:
      - wails dev

  lint:go:
    desc: Lint Go backend
    cmds:
      - golangci-lint run

  lint:web:
    desc: Lint Web frontend (Biome)
    cmds:
      - npx @biomejs/biome check frontend/

  lint:
    desc: Run all linters
    cmds:
      - task: lint:go
      - task: lint:web
    ```

#### 1.2.3. Pre-commit Hooks (`.pre-commit-config.yaml`)

```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.6.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: detect-private-key

  # Go Hooks
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.60.0
    hooks:
      - id: golangci-lint
        # Optimization: Run fast locally by only analyzing changed files
        args: [--fast]

  # Web Hooks (Biome)
  - repo: https://github.com/biomejs/biome
    rev: v1.8.3
    hooks:
      - id: biome
        # Check and apply safe fixes automatically
        args: ["check", "--write"]
        files: ^frontend/
```

-----

## 2. Agent Reference Matrix (Mission B)

This matrix details the configuration locations and integration vectors for the target AI agents, crucial for the `easyConfig` application's "Meta-Config" feature.

| Agent | Status (2025) | Integration Vectors | Global/User Config Path(s) | System/Managed Config Path(s) | Context/Instruction File(s) (Repo) | Auth Storage Location (Sensitive) |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **GitHub Copilot** | Active | IDE, CLI (`gh copilot`), Actions | VS Code User Settings [1]<br>`~/.copilot/mcp-config.json` (CLI) | N/A | `.github/copilot-instructions.md`<br>`.github/agents/` | Copilot Auth [2] (`hosts.json`) |
| **Claude Code** | Active | CLI (`claude`), GitHub App/Action, SDK | User settings (varies) | Claude Managed Settings [3] | `CLAUDE.md` | OS Keychain (CLI)<br>`ANTHROPIC_API_KEY` (Env Var) |
| **OpenAI (Codex/GPT)** | Active | CLI (`openai`), API | `~/.config/openai/config.yaml` (Inferred) | N/A | `AGENTS.md` | `OPENAI_API_KEY` (Env Var) |
| **Google Jules** | Active | Web UI, GitHub App, CLI (`jules`), API | N/A | N/A | `AGENTS.md` | Google OAuth (Web/CLI)<br>API Key (API) |

**Path References:**

  * **[1] VS Code User settings.json:**
      * Windows: `%APPDATA%\Code\User\settings.json`
      * macOS: `~/Library/Application Support/Code/User/settings.json`
      * Linux: `~/.config/Code/User/settings.json`
  * **[2] Copilot Auth (hosts.json):**
      * Linux/macOS: `~/.config/github-copilot/hosts.json` (or `$XDG_CONFIG_HOME`)
      * Windows: `%LOCALAPPDATA%\github-copilot\hosts.json`
  * **[3] Claude Managed Settings (managed-settings.json):**
      * Windows: `C:\ProgramData\ClaudeCode\`
      * macOS: `/Library/Application Support/ClaudeCode/`
      * Linux: `/etc/claude-code/`

-----

## 3. App Architecture (Mission B: Meta-Config)

The `easyConfig` application will discover, validate, and manage these configurations.

### 3.1. Go Struct Definitions (Interface-Based)

We will use an interface-based approach in the Go backend for extensibility.

```go
package agentconfig

// ConfigMetadata holds information about where and how a config is stored.
type ConfigMetadata struct {
    Path   string `json:"path"`
    Scope  string `json:"scope"` // "user", "project", "system"
    Format string `json:"format"` // "json", "toml", "md", "yaml"
}

// AgentConfig is the interface implemented by all agent configuration structs.
type AgentConfig interface {
    GetAgentName() string
    GetMetadata() ConfigMetadata
    Validate() error
    // Save() error // For implementing write functionality
}

// CopilotVSCodeSettings represents relevant parts of VS Code settings.json.
// Note: This requires careful parsing as settings.json contains many other settings.
type CopilotVSCodeSettings struct {
    Metadata              ConfigMetadata `json:"-"`
    EnableAutoCompletions *bool          `json:"github.copilot.editor.enableAutoCompletions,omitempty"`
    LocaleOverride        *string        `json:"github.copilot.chat.localeOverride,omitempty"`
    // ... other relevant keys
}

func (c CopilotVSCodeSettings) GetAgentName() string { return "GitHub Copilot" }
func (c CopilotVSCodeSettings) GetMetadata() ConfigMetadata { return c.Metadata }
// ... Validate implementation

// MarkdownContextFile represents CLAUDE.md, AGENTS.md, etc.
type MarkdownContextFile struct {
    Metadata ConfigMetadata `json:"-"`
    Content  string         `json:"content"`
}
// ... Implementation
```

### 3.2. FileSystem Discovery Logic

A dedicated `discovery` package will handle cross-platform file location, respecting standards like the XDG Base Directory Specification.

```go
package discovery

import (
    "os"
    "path/filepath"
    "runtime"
    // "github.com/easyConfig/easyConfig/internal/agentconfig"
)

// GetPlatformPaths returns OS-specific configuration roots.
func GetPlatformPaths() (configHome, localAppData, programData string) {
    home, _ := os.UserHomeDir()

    switch runtime.GOOS {
    case "windows":
        configHome = os.Getenv("APPDATA") // Roaming
        localAppData = os.Getenv("LOCALAPPDATA")
        programData = os.Getenv("ProgramData")

    case "darwin":
        // macOS often uses Application Support for GUI, .config for CLI
        configHome = filepath.Join(home, ".config")
        localAppData = filepath.Join(home, "Library", "Application Support")
        programData = "/Library/Application Support"

    case "linux":
        // Adhere to XDG Base Directory Specification
        if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
            configHome = xdgConfig
        } else {
            configHome = filepath.Join(home, ".config")
        }
        programData = "/etc"
    }
    return
}

// DiscoverAll scans the system for configurations.
func DiscoverAll(projectRoot string) {
    configHome, localAppData, programData := GetPlatformPaths()

    // 1. User/Global Discovery (e.g., VS Code, Copilot CLI)
    // ...

    // 2. System Discovery (e.g., Claude Managed Settings)
    // ...

    // 3. Project Specific Discovery (e.g., CLAUDE.md)
    // ...
}
```

### 3.3. Security Note on Auth

`easyConfig` must identify authentication files (like Copilot's `hosts.json`) but **must never** read, display, or persist the actual tokens. The application should only verify the file's existence and report its status.

-----

## 4. DevOps Architecture (Mission A: Self-Hosting)

We will integrate agents into our GitHub workflow to automate maintenance tasks.

### 4.1. Agent Selection and Prerequisites

  * **Primary Agent:** Claude Code (Review, Triage, Interactive).
  * **Secondary Agent:** Jules (Asynchronous Fixing).
  * **Required Secrets:** `ANTHROPIC_API_KEY` must be added to the repository secrets.

### 4.2. Repository Context (`CLAUDE.md` / `AGENTS.md`)

A context file at the repository root guides the agents.

```markdown
# easyConfig Project Context

This repository contains the source code for `easyConfig`, a desktop application built with Wails.

## Tech Stack
- Backend: Go 1.24+
- Frontend: React, TypeScript
- Framework: Wails v3
- Linting: GolangCI-Lint (Backend, see `.golangci.yml`), Biome (Frontend, see `biome.json`)
- Tooling: Mise, Taskfile (see `Taskfile.yml`)

## Coding Standards
- Adhere strictly to linting configurations.
- **CRITICAL:** Go code must prioritize cross-platform compatibility. Filesystem operations in `/internal/discovery/` must correctly handle Windows, macOS, and Linux paths (including XDG standards).
- TypeScript must utilize modern React features and strict typing.

## Key Directories
- `/internal/agentconfig/`: Go structs and configuration models.
- `/internal/discovery/`: Filesystem discovery logic.
- `/frontend/src/`: React/TS source code.
```

### 4.3. Workflow: Agent Code Review (`.github/workflows/agent-review.yml`)

Automated review of Pull Requests using the official Claude action.

```yaml
name: Agent PR Review (Claude)

on:
  pull_request:
    types: [opened, synchronize]

# Ensure secure permissions
permissions:
  contents: read
  pull-requests: write # Allows Claude to comment on the PR

jobs:
  claude-review:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v5

    - name: Run Claude Code Action for Review
      uses: anthropics/claude-code-action@.claude/prds/easy-config-v1.md
      with:
        # Utilize the built-in /review command which uses the PR context and CLAUDE.md
        prompt: "/review"
        anthropic_api_key: ${{ secrets.ANTHROPIC_API_KEY }}
```

### 4.4. Workflow: Agent Triage (`.github/workflows/agent-triage.yml`)

Auto-labeling issues based on content.

```yaml
name: Agent Issue Triage

on:
  issues:
    types: [opened]

jobs:
  triage:
    runs-on: ubuntu-latest
    permissions:
      issues: write
      contents: read
    steps:
    - name: Checkout repository
      uses: actions/checkout@v5

    - name: Run Claude Code for Triage
      # Using the base action for a custom prompt
      uses: anthropics/claude-code-base-action@.claude/prds/easy-config-v1.md
      with:
        anthropic_api_key: ${{ secrets.ANTHROPIC_API_KEY }}
        prompt: |
          Analyze the issue title and body provided in the GitHub event context.
          Determine the appropriate labels from the repository's existing labels (e.g., 'bug', 'feature', 'backend', 'frontend', 'CI').
          Apply the labels to the issue using the available GitHub tools. Do not add comments.
```

### 4.5. Workflow: Interactive Agent (`.github/workflows/agent-interactive.yml`)

Allowing developers to invoke agents on demand via comments (e.g., ` @claude /refactor this`).

```yaml
name: Agent Interactive (On-Demand)

on:
  issue_comment:
    types: [created]

permissions:
  contents: read
  issues: write
  pull-requests: write

jobs:
  claude-interactive:
    # Only run if the comment specifically mentions @claude (or other configured agent name)
    if: contains(github.event.comment.body, ' @claude')
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v5

    - name: Run Claude Code Action
      uses: anthropics/claude-code-action@.claude/prds/easy-config-v1.md
      with:
        # The prompt is automatically derived from the comment content
        anthropic_api_key: ${{ secrets.ANTHROPIC_API_KEY }}
```

### 4.6. Jules Integration

Jules integrates via a GitHub App ([https://jules.google/](https://jules.google/)). This requires no YAML workflow as it uses webhooks. We will configure Jules to trigger on a specific label (e.g., `jules-fix`) applied to an issue, allowing it to autonomously process the task and open a Pull Request.

-----

## 5. Future Proofing

The architecture is designed for extensibility:

1.  **Interface-Driven Backend (Mission B):** The use of the `AgentConfig` interface in the Go backend allows new agents to be integrated by simply implementing the interface and adding a corresponding discovery function, without altering core application logic.
2.  **Robust Discovery Logic:** The `discovery` package's adherence to cross-platform standards (XDG, OS-specific paths) ensures it can adapt to new configuration locations.
3.  **DevOps Adaptability (Mission A):** The reliance on standard GitHub triggers (Actions, comments, labels) ensures that new AI tools can be integrated into the DevOps workflow as they release their own integrations.

```