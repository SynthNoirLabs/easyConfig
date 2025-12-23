# Task: Sub-Agent Discovery & Visualization

## 🎯 Objective
Visualize the hierarchy of sub-agents (especially for Claude and systems using `AGENTS.md`).

## 📝 Description
Complex agentic systems use multiple specialized agents. `AGENTS.md` defines these. We should parse this and show a tree view.

## ✅ Requirements
1.  **Parser**: Markdown parser that extracts agent definitions from `AGENTS.md`.
2.  **Visualization**:
    -   Tree Diagram.
    -   Details panel showing the "Prompt/Persona" for each sub-agent.
3.  **Editing**: Allow modifying a sub-agent's instructions directly from the UI.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/parsers/agents_md.go`.
-   **Frontend**: `SubAgentTree` component.
