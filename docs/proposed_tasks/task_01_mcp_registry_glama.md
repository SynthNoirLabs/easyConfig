# Task: MCP Marketplace Integration (Glama)

## 🎯 Objective
Integrate Glama's public MCP registry into the EasyConfig "Marketplace" tab to allow users to browse and install MCP servers easily.

## 📝 Description
Glama offers a comprehensive index of over 6,500 MCP servers. Integrating this into EasyConfig will transform it from a configuration editor into a full-blown capability manager for AI agents.

## ✅ Requirements
1.  **API Integration**: Reverse engineer or find the public API endpoint for Glama's registry (or scrape if permitted/necessary, though API is preferred).
2.  **UI Components**:
    -   **Search Bar**: Filter by name, category, or description.
    -   **Card View**: Display Server Name, Description, Author, and "Install" button.
    -   **Detail View**: Show full readme/docs for the server.
3.  **Installation Logic**:
    -   When "Install" is clicked, prompt the user to select which Agent (Claude, Gemini, etc.) to inject the server into.
    -   Use the existing `McpManager` logic to generate the config snippet.
    -   Handle `npx` vs `docker` installation commands based on user preference.

## 🔗 References
-   [Glama Website](https://glama.ai)
-   [Model Context Protocol](https://modelcontextprotocol.io)

## 🛠️ Technical Implementation
-   **Backend**: Add `pkg/marketplaces/glama.go` to fetch and cache the registry list.
-   **Frontend**: Update `frontend/src/components/Marketplace.tsx` to support a "Glama" tab.
