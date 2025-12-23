# Task: Visual MCP Graph

## 🎯 Objective
Visualize the relationships between Agents, Configuration Files, and MCP Servers.

## 📝 Description
As the number of agents and MCP servers grows, it becomes hard to see which agent has access to what. A node-link diagram will make this clear.

## ✅ Requirements
1.  **Nodes**:
    -   **Agent Nodes** (Claude, Gemini, etc.)
    -   **Config File Nodes** (Global settings, Project settings)
    -   **MCP Server Nodes** (Postgres, GitHub, Filesystem)
2.  **Edges**:
    -   "Defines" (Config File -> MCP Server)
    -   "Reads" (Agent -> Config File)
3.  **Interaction**:
    -   Clicking a node opens the relevant config editor.

## 🛠️ Technical Implementation
-   **Frontend**: Use `reactflow` or `visx`.
-   **Data Source**: The existing `DiscoveryService` already has the graph data implicitly; we just need to structure it.
