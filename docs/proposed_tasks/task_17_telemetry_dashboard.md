# Task: Local Telemetry & Usage Dashboard

## 🎯 Objective
Show the user which agents they use the most (based on config modification activity).

## 📝 Description
"You've tweaked Claude's config 50 times this week, but haven't touched Codex."

## ✅ Requirements
1.  **Event Tracking**: Record "Save" events with timestamps locally (SQLite or JSON file).
2.  **Visualization**:
    -   Bar chart: "Activity by Agent".
    -   Timeline: "Recent Changes".
3.  **Privacy**: 100% local. No data leaves the machine.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/telemetry/local_store.go`.
-   **Frontend**: "Stats" tab using `recharts`.
