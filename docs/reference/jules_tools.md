# Jules Tools Configuration Reference

**Source:** https://jules.google/docs/cli/reference/

## Overview
Jules is an autonomous coding agent integrated with Google Cloud and GitHub. It relies primarily on cloud-side sessions but maintains some local state.

## Configuration Locations

### 1. Agent Context (Project Root)
*   **File:** `AGENTS.md` (Recommended)
*   **Purpose:** Defines the persona, architecture overview, and high-level instructions for Jules.
*   **Format:** Markdown.

### 2. Local Data Store
*   **Path:** `~/.jules-mcp/data.json`
*   **Purpose:** Persists local state, session IDs, and authentication tokens.
*   **Env Var:** `JULES_DATA_PATH` can override this location.

## API Interaction (Reference)
Jules is often driven via API or the GitHub App. Here is an example of a session creation request:

```bash
curl 'https://jules.googleapis.com/v1alpha/sessions' \
    -X POST \
    -H "Content-Type: application/json" \
    -H 'X-Goog-Api-Key: YOUR_API_KEY' \
    -d '{
      "prompt": "Refactor the login module",
      "sourceContext": {
        "source": "sources/github/owner/repo",
        "githubRepoContext": {
          "startingBranch": "main"
        }
      },
      "title": "Refactor Login"
    }'
```

## CLI Commands
*   `jules login`: Authenticate with Google.
*   `jules remote new`: Start a new task.
*   `jules remote list`: List active sessions.
*   `jules remote pull`: Apply changes from a session to local disk.
