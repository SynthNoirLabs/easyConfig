# Task: Unified Context File Generator

## 🎯 Objective
Create a "Context Wizard" that generates `CLAUDE.md`, `GEMINI.md`, and `.cursorrules` simultaneously from a single set of user inputs.

## 📝 Description
Developers often use multiple AI tools (e.g., Cursor for coding, Claude for architecture, Gemini for CLI). They currently have to maintain three separate context files. This feature unifies them.

## ✅ Requirements
1.  **Input Form**:
    -   **Project Name & Description**
    -   **Tech Stack** (Select: React, Go, Python, etc.)
    -   **Coding Style** (e.g., "Functional", "OOP", "Concise", "Verbose")
    -   **Architecture** (e.g., "Hexagonal", "MVC")
2.  **Template Engine**:
    -   Create Go templates for each file format.
    -   Map the single input form to the specific sections required by each agent.
3.  **Output Generation**:
    -   Button: "Generate All".
    -   Writes files to the project root.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/context/generator.go` with `text/template`.
-   **Frontend**: New "Context" tab in the Sidebar.
