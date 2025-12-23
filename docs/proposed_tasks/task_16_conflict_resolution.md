# Task: Conflict Resolution Wizard

## 🎯 Objective
Detect and resolve port conflicts or file ownership issues between agents.

## 📝 Description
If Claude and Gemini both try to run an MCP server on port 3000, things break.

## ✅ Requirements
1.  **Port Scanning**: Parse configs to find "localhost:XXXX" or "port": XXXX.
2.  **Conflict Detection**: Map `Port -> [Agent1, Agent2]`.
3.  **UI**:
    -   Dashboard alert: "Port Conflict Detected".
    -   Resolution: "Change Port for Agent X".

## 🛠️ Technical Implementation
-   **Backend**: `pkg/analyzer/ports.go`.
