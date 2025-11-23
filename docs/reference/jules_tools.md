# Jules Tools Configuration Reference

**Source:** https://jules.google/docs/cli/reference/

## Overview
Jules is designed to be autonomous and stateless on the client side, relying heavily on cloud sessions and project-based context.

## Configuration Locations

### 1. Agent Context (Project Root)
*   **File:** `AGENTS.md` or `JULES.md` (Project Root)
*   **Purpose:** Defines the persona, tools, and instructions for the agent.
*   **Format:** Markdown with optional frontmatter? (Search results imply standard Markdown instructions).

### 2. Local Data Store
*   **Path:** `~/.jules-mcp/data.json`
*   **Purpose:** Persists local state, session IDs, and possibly cached credentials.
*   **Env Var:** `JULES_DATA_PATH` can override this location.

### 3. Authentication
*   Managed via `jules login`.
*   Tokens likely stored in system keychain or `~/.config/google-jules/`.

## CLI Reference
*   `jules remote new`: Start a task.
*   `jules remote list`: List sessions.
*   `jules remote pull`: Apply changes.
*   `--theme [light|dark]`: UI preference.

## Example `AGENTS.md`
```markdown
# Project Context for Jules

## Architecture
This project uses a Provider Pattern in Go.

## Tools
- Use `go test ./...` to verify changes.
- Configs are in `pkg/config/`.
```