# Task: Team Config Export

## 🎯 Objective
Easily share a "starter pack" of agent configurations with a new team member.

## 📝 Description
Onboarding is hard. "Here, take my configs" is helpful, but sharing secrets is bad.

## ✅ Requirements
1.  **Selection**: Checkbox list of configs to export.
2.  **Sanitization**: **CRITICAL**. Regex scan to remove API Keys / Secrets before export. Replace with `<INSERT_KEY>`.
3.  **Packaging**: Create a `.zip` file.
4.  **Import**: "Import Team Pack" wizard.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/export/packager.go`.
