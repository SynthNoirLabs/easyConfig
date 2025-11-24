# Agent Architecture Analysis Prompt

**Role:** You are a Principal AI Architect and Research Analyst.
**Objective:** Produce the definitive **`docs/AGENT_IMPLEMENTATION_PLAN.md`**.

## ⚠️ Critical Directive: Trust But Verify
The information provided in previous context is **preliminary** and potentially incomplete. Do **not** take it as the source of truth. You possess powerful tools (Web Search, Documentation Fetching, Reasoning). **Use them.**

---

## 🎯 Part 1: The Dual Mission

We are building `easyConfig`. This project has two distinct needs regarding AI agents.

### Mission A: Internal DevOps (The "Self-Hosting" Aspect)
*   **Goal:** This repo (`easyConfig`) is a playground for agent collaboration.
*   **Requirement:** Design a workflow system where agents actively participate.
*   **Deliverable:** Define the GitHub Actions (`.github/workflows/*.yml`) and permission sets.
*   **Challenge:** **Do NOT assume a single workflow.** Investigate:
    *   **Triage:** `agent-triage.yml` (Auto-labeling based on issue content).
    *   **Code Review:** `agent-review.yml` (Can Claude/Codex review PRs automatically?).
    *   **Fixing:** Triggering an agent to fix a specific lint error or test failure.

### Mission B: Product Features (The "Meta-Config" Aspect)
*   **Goal:** The `easyConfig` *application* allows users to manage their *own* local AI agent configurations.
*   **Requirement:** The app needs to discover, read, validate, and save configuration files.
*   **Deliverable:**
    *   Map the **exact file paths** (Global & Project scope) for each agent on **Linux, macOS, and Windows**.
    *   Map the **data schemas** (JSON/TOML/YAML keys).

---

## 🔎 Part 2: Deep-Dive Agent Research (Be Exhaustive)

You must investigate the following specifics for **Claude Code, Copilot, Jules, and Codex**:

### A. Claude Code (CLI & Web)
*   **Config Files:**
    *   `~/.claude.json` vs `~/.claude/settings.json` vs `.claude/settings.json`. Which is the current standard?
    *   **MCP Config:** Is it in `.mcp.json`, `mcp_servers.json`, or embedded in `settings.json`?
*   **Web Environment:**
    *   Does it support a `.devcontainer.json` style config?
    *   Can we provide a `setup.sh` hook for the web environment?

### B. Google Jules
*   **Repo Configuration:**
    *   Does it read a `.jules.yaml` or `jules.config.js`?
    *   **Environment:** How do we script the setup (npm install, go mod download) for Jules's cloud VM? (Look for "Configure repo environment" in docs).
*   **Context:** Confirm usage of `AGENTS.md`.

### C. GitHub Copilot (CLI & Agent)
*   **Config Files:**
    *   `~/.config/github-copilot/hosts.json` (Auth).
    *   `~/.config/github-copilot/mcp-config.json` (MCP Servers).
    *   `.github/copilot-instructions.md` (Project instructions).
*   **Manual Triggers:**
    *   How to enable "Assign to Copilot" in issues?

### D. OpenAI Codex (CLI)
*   **Config Files:**
    *   `~/.codex/config.toml` (Global).
    *   `.codex/config.toml` (Project).
*   **Deprecation Check:** Is the `codex` CLI still supported in late 2025?

---

## 📝 Output Specification: `docs/AGENT_IMPLEMENTATION_PLAN.md`

Your final artifact must be a comprehensive Markdown document containing:

1.  **Foundation Strategy:**
    *   Recommended CI/CD improvements (Caching, Biome, GolangCI-Lint).
    *   Local dev setup guide (Mise/Taskfile).
2.  **Agent Reference Matrix:** A detailed lookup table for each agent (Config Path, Format, Env Vars, Trigger Methods).
3.  **App Architecture:**
    *   Go Struct definitions for the "Meta-Config" feature.
    *   FileSystem Discovery logic.
4.  **DevOps Architecture:**
    *   YAML snippets for `.github/workflows/`.
    *   Instructions for Labels and Secrets.
5.  **Manual Config Guide:** A section dedicated to "Manual Interventions" required by the user (e.g., "Go to Jules Web UI > Settings > Environment").

**Go forth and research. Do not guess.**