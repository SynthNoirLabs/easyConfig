# Task: MCP Marketplace Integration (Smithery)

## 🎯 Objective
Integrate Smithery's MCP registry as a secondary source for MCP servers, focusing on their "CLI" and "Hosted" features.

## 📝 Description
Smithery provides a distinct set of managed and community MCP servers. It also offers a CLI tool that simplifies installation. This task involves integrating Smithery's catalog and potentially using their CLI if available on the user's system.

## ✅ Requirements
1.  **Catalog Fetching**: Implement logic to fetch the list of available servers from Smithery.
2.  **Smithery CLI Support**:
    -   Check if `smithery` CLI is installed via `exec.LookPath`.
    -   If installed, offer a "Quick Install" option that runs `smithery install <server>`.
3.  **Manual Install Fallback**: If CLI is missing, fallback to standard `npx` or config injection methods.

## 🔗 References
-   [Smithery.ai](https://smithery.ai)

## 🛠️ Technical Implementation
-   **Backend**: Add `pkg/marketplaces/smithery.go`.
-   **Comparison**: Ideally, show if a server exists on both Glama and Smithery and let the user choose the source.
