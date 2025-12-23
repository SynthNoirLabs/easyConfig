# Task: System Dark Mode Sync & Theming

## 🎯 Objective
Polish the UI by respecting the OS theme preference automatically.

## 📝 Description
The app currently defaults to a theme or has a manual toggle. It should auto-switch.

## ✅ Requirements
1.  **Detection**: Use Wails `runtime.SystemDarkTheme()` (or JS `window.matchMedia`).
2.  **Listener**: React to OS changes in real-time.
3.  **Theme Engine**: Ensure all Tailwind classes support `dark:` variants correctly.
4.  **Monaco**: Dynamically switch Monaco theme (`vs-dark` vs `vs-light`).

## 🛠️ Technical Implementation
-   **Frontend**: `ThemeContext.tsx` updates.
