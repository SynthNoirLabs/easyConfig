# Task: Custom Provider Plugins

## 🎯 Objective
Allow the community to add support for new agents without modifying the EasyConfig source code.

## 📝 Description
We can't support every agent natively. A plugin system allows extensibility.

## ✅ Requirements
1.  **Plugin Format**: Simple YAML or JSON file defining:
    -   Provider Name.
    -   Config Paths (Global/Project).
    -   Format (JSON/TOML/YAML).
2.  **Loader**: Logic to scan `~/.config/easyconfig/plugins/*.yaml`.
3.  **Dynamic Provider**: A generic `Provider` implementation that reads these definitions.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/config/provider_dynamic.go`.
