# Task: Secure Secrets Management

## 🎯 Objective
Move away from storing API keys in plain text files (like `.env` or `config.json`) and integrate with the system keychain.

## 📝 Description
Security is paramount. EasyConfig should encourage best practices by offering a built-in way to store and retrieve API keys using the OS's native secure storage (Keychain on macOS, Credential Manager on Windows, Secret Service API on Linux).

## ✅ Requirements
1.  **Library Integration**: Use `github.com/zalando/go-keyring` or similar.
2.  **UI Workflow**:
    -   When a user adds an MCP server requiring an API Key, offer a checkbox: "Store securely in Keychain".
    -   Instead of writing the key to JSON, write a placeholder (e.g., `{{KEYCHAIN:service_name}}`) or rely on the Agent's ability to read from env vars, and inject the env var at runtime (if EasyConfig launches the agent) - *Correction*: Since EasyConfig configures *other* tools, we might need to write a helper script or use the Agent's native secret support if available.
    -   *Alternative*: Create a local `.env.local` that is explicitly gitignored and warn the user.
3.  **Audit**: Scan existing configs for potential leaked keys and warn the user.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/security/keychain.go`.
-   **Frontend**: "Secrets" management drawer.
