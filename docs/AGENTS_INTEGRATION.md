# AI Agent Integration Guide

This repository is designed to be worked on by multiple AI agents. Below is the guide on how to invoke and collaborate with them.

## 1. Google Jules (Autonomous Coding Agent)
**Best for:** Complex feature implementation, refactoring, and multi-file tasks.

*   **Trigger (CLI):**
    ```bash
    # Requires 'zen-mcp-server' or similar integration
    /jules "Refactor the config provider to use a factory pattern"
    ```
*   **Trigger (Web):** Go to [jules.google.com](https://jules.google.com), select this repo, and prompt.
*   **Trigger (Issue):** (If configured) Assign the issue to the Jules bot or use a specific label (check repo settings).
*   **Status:** Check [Jules Dashboard](https://jules.google.com/dashboard).

## 2. Claude Code (Web/CLI)
**Best for:** Rapid iteration, PR reviews, and "one-shot" fixes.

*   **Setup:** Install the [Claude Code GitHub App](https://docs.anthropic.com/en/docs/claude-code) (if available) or use the CLI.
*   **Trigger (Issue):** Mention `@Claude` in a GitHub Issue to have it attempt a fix (requires App).
*   **Trigger (CLI):**
    ```bash
    claude "Fix the bug in providers.go regarding JSON parsing"
    ```
*   **CI/CD Integration:** We use the `anthropics/claude-code-action` in `.github/workflows/agent-dispatch.yml` (planned) for automated fixes.

## 3. GitHub Copilot
**Best for:** Code completion, quick questions, and "in-editor" assistance.

*   **Issues:** Assign **Copilot** to an issue to have it start working.
*   **PRs:** Use `@github` in PR comments to ask Copilot to explain, summarize, or suggest fixes.
*   **CLI:**
    ```bash
    gh copilot explain "why did the CI fail?"
    gh copilot suggest "create a struct for Gemini config"
    ```

## 4. OpenAI Codex
**Best for:** CI/CD Auto-fixes and scripted modifications.

*   **Auto-Fix:** We have configured a workflow that triggers Codex on CI failure to analyze logs and suggest a patch.
*   **CLI:**
    ```bash
    codex exec "find all TODOs and list them in TODO.md"
    ```

## Workflow for Agents
1.  **Pick an Issue:** Look at `TASKS.md` or GitHub Issues.
2.  **Create Branch:** `agent/<agent-name>/<task-name>`.
3.  **Implement & Test:** Run `go test ./pkg/...` locally.
4.  **Open PR:** Use specific PR templates if available.
5.  **Respond to Feedback:** Listen for comments from human reviewers or other agents.
