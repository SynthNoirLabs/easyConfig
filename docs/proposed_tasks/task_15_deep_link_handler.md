# Task: Deep Link Protocol Handler

## 🎯 Objective
Enable opening EasyConfig from the terminal or web browser via `easyconfig://`.

## 📝 Description
Integration with the OS to handle custom URL schemes. Useful for "Open in EasyConfig" buttons on websites or terminal output.

## ✅ Requirements
1.  **Protocol Registration**: Register `easyconfig://` during installation (NSIS for Windows, .desktop for Linux, Info.plist for Mac).
2.  **Arguments**: Support `open?path=/path/to/config.json`.
3.  **App Logic**: Parse startup args and navigate the React router to the correct page.

## 🛠️ Technical Implementation
-   **Wails**: Use `wails.json` and platform-specific configs to register the scheme.
-   **Frontend**: React Router handling.
