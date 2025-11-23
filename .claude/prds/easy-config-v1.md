---
name: easy-config-v1
description: Centralized Wails application to manage AI agent configurations (Claude, Gemini, Codex)
status: backlog
created: 2025-11-22T12:00:00Z
---

# PRD: easy-config-v1

## Executive Summary
EasyConfig is a "Mission Control" desktop application (Wails + React) that centralizes the management of AI Agent configurations. It allows developers to view, edit, and extend the settings of tools like Claude Code, Gemini CLI, and Codex from a single dashboard, solving the fragmentation of config files scattered across the OS.

## Problem Statement
AI CLI tools store configurations in various formats (JSON, TOML, YAML) and locations (Home dir, Project dir, XDG paths). Developers struggle to:
1.  Remember where each agent's config lives.
2.  Safely edit configs without breaking syntax.
3.  Manually copy-paste complex MCP server configurations (JSON blocks) into these files.

## User Stories
- **As a Developer**, I want to see all my installed AI agents in one list so I know what is configured.
- **As a Developer**, I want to edit my global Claude config and my project-specific Codex config in the same UI.
- **As a Power User**, I want to "Inject" an MCP server (like Postgres) into Claude Code by just clicking "Add", without manually editing JSON.
- **As a User**, I want a raw editor fallback in case the UI doesn't support a specific new setting.

## Requirements

### Functional Requirements
1.  **Agent Discovery**:
    - Automatically detect config files for Claude Code, Gemini CLI, Codex, and Jules.
    - Support Global (`~/.`) and Project-scoped configurations.
2.  **Universal Editor**:
    - Read/Write support for JSON, YAML, and TOML.
    - Monaco-based "Raw Mode" for direct text editing.
3.  **Smart MCP Injection**:
    - Database of common MCP servers (Filesystem, Git, Postgres).
    - Logic to append `mcpServers` blocks to existing configs without syntax errors.
4.  **Safe Saving**:
    - Validate JSON/TOML syntax before saving to disk.

### Non-Functional Requirements
- **Cross-Platform**: Run on Linux, Windows, and macOS (Wails).
- **Fast Startup**: <2 seconds to load and scan configs.
- **Secure**: Do not expose API keys in plain text in the UI unless requested.

## Success Criteria
- Successfully detect installed agents on the user's machine.
- Successfully read and write a change to `~/.claude/config.json`.
- Successfully inject a new MCP server block into a config file via the UI.

## Constraints & Assumptions
- User has read/write access to their home directory.
- Agents follow standard configuration paths (as defined in SPEC.md).
- Wails v2 runtime environment is available.

## Out of Scope
- Managing the *execution* or *logs* of the agents (Config only).
- "One-click install" of the agents themselves (we only manage existing configs).

## Dependencies
- **Backend**: Go 1.23+, Wails v2.
- **Frontend**: React 18, Vite, Shadcn UI, Monaco Editor.
