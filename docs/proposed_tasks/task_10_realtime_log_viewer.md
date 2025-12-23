# Task: Real-time Log Viewer

## 🎯 Objective
Provide visibility into what the agents are actually doing by tailing their log files.

## 📝 Description
When an agent fails, the user has to hunt for log files. EasyConfig should bring these logs to the dashboard.

## ✅ Requirements
1.  **Log Discovery**: Know where each agent stores logs (e.g., `~/.codex/codex.log`).
2.  **Tailing**: Use Go's `fsnotify` or polling to read new lines.
3.  **Streaming**: Stream lines to the frontend via Wails Events.
4.  **Filtering**: Allow filtering logs by "Error", "Warning", "Info".

## 🛠️ Technical Implementation
-   **Backend**: `pkg/logs/tailer.go`.
-   **Frontend**: `LogViewer` component with auto-scroll.
