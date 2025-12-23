# Task: Agent Playground (Sanity Check)

## 🎯 Objective
Allow users to verify their configuration is working by running a minimal agent task.

## 📝 Description
After editing a config, the user wonders: "Did I break it?" The playground allows a quick test.

## ✅ Requirements
1.  **Execution Engine**: Ability to run shell commands (e.g., `claude -p "hello"`).
2.  **Output Capture**: Capture stdout/stderr.
3.  **Timeout**: Kill the process if it hangs (5s timeout).
4.  **Success Criteria**: Check exit code.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/runner/sanity.go`.
-   **Frontend**: "Test Config" button in the editor toolbar.
