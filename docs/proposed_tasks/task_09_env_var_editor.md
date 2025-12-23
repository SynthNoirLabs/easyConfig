# Task: Environment Variable Editor

## 🎯 Objective
A specialized GUI for managing `.env` files, which are the backbone of many agent configurations.

## 📝 Description
Many agents (Mentat, OpenHands) rely on `.env` files. A text editor is fine, but a specialized editor is better.

## ✅ Requirements
1.  **Parsing**: Parse `KEY=VALUE` lines, preserving comments.
2.  **Validation**: Warn about duplicates or invalid syntax.
3.  **Security**:
    -   Mask values by default (*****).
    -   Toggle visibility.
4.  **UX**:
    -   "Add Key" button.
    -   "Copy Value" button.

## 🛠️ Technical Implementation
-   **Backend**: Enhance `pkg/config/service.go` to parse `.env` as a structured map.
-   **Frontend**: `EnvEditor` component.
