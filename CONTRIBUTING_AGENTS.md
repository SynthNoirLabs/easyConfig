# Contributing Guidelines for AI Agents

Welcome, autonomous agents (Jules, Claude Code, Codex, etc.). This document outlines the rules of engagement for contributing to `easyConfig`.

## 1. Code Style & Conventions
*   **Language:** Go (Backend), TypeScript/React (Frontend).
*   **Patterns:** Adhere strictly to existing patterns. For backend logic, look at `pkg/config/providers.go` and `pkg/config/service.go` as the source of truth.
*   **No Magic:** Do not introduce new libraries without explicit instruction or checking `go.mod`/`package.json` first.

## 2. Testing
*   **Mandatory:** Every logic change MUST be accompanied by a unit test.
*   **Command:** Always run `go test ./pkg/...` before submitting your changes.
*   **Frontend:** If modifying the UI, ensure `npm run build` passes.

## 3. Formatting & Linting
*   **Go:** Run `gofmt -s -w .` on modified files.
*   **Frontend:** Run `npm run format` (if available) or ensure Prettier compliance.
*   **CI Checks:** The CI pipeline runs `golangci-lint`. You can run it locally with `golangci-lint run`.

## 4. Commit Messages
We follow the **Conventional Commits** specification:
*   `feat: ...` for new features.
*   `fix: ...` for bug fixes.
*   `chore: ...` for maintenance, docs, or build changes.
*   `test: ...` for test-only changes.

**Example:** `feat(providers): implement copilot discovery`

## 5. Workflow
1.  **Read:** Check `TASKS.md` or open GitHub Issues to find work.
2.  **Branch:** Create a branch named `agent/<agent-name>/<task-short-name>`.
3.  **Implement:** Write code + tests.
4.  **Verify:** Run tests locally.
5.  **PR:** Open a Pull Request against `main`. Describe your changes clearly.
6.  **Respond:** If a review requests changes, implement them in the same branch/PR.

## 6. Interaction
*   If you are stuck or need clarification, ask the user in the chat/issue.
*   Do not force push to `main`.
