# Task: Advanced Config Validation

## 🎯 Objective
Prevent invalid configurations from being saved by verifying them against schemas.

## 📝 Description
Currently, we validate JSON syntax. We should validate against the actual JSON Schemas for each agent.

## ✅ Requirements
1.  **Schema Registry**: Maintain a mapping of `filename -> schema_url`.
2.  **Validation Engine**:
    -   Use `github.com/santhosh-tekuri/jsonschema/v5`.
    -   Validate before Save.
3.  **UI Feedback**:
    -   Show red squiggles in Monaco Editor.
    -   Show detailed error list below the editor.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/schema/validator.go`.
-   **Frontend**: Monaco Editor markers integration.
