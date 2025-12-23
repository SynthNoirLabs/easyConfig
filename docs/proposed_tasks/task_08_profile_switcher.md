# Task: Global Profile Switcher

## 🎯 Objective
Allow users to switch between "Work", "Personal", and "Hobby" contexts instantly.

## 📝 Description
A user might use a secure, enterprise-hosted LLM for work (Claude Enterprise) but a personal API key for hobby projects. Manually editing configs to switch endpoints is tedious.

## ✅ Requirements
1.  **Profile Definition**:
    -   A "Profile" is a named set of overrides.
2.  **Switching Mechanism**:
    -   **Symlink Swapping**: Point `~/.claude/config.json` to `~/.claude/config-work.json`.
    -   **Env Var Injection**: Set `CLAUDE_CONFIG_DIR` or similar if supported.
3.  **UI**:
    -   Dropdown in the top header: "Current Profile: Personal".

## 🛠️ Technical Implementation
-   **Backend**: `pkg/profiles/manager.go`.
