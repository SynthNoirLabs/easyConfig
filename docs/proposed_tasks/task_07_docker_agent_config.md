# Task: Docker Agent Configurator

## 🎯 Objective
Support configuration management for agents that run exclusively or primarily in Docker containers (e.g., OpenHands, Devin).

## 📝 Description
Agents like OpenHands are complex to configure because their settings must be passed as Environment Variables to the Docker container. Editing a local file isn't enough; the user needs a correct `docker run` command.

## ✅ Requirements
1.  **Form Interface**:
    -   UI for settings like `SANDBOX_RUNTIME`, `LLM_API_KEY`, `WORKSPACE_MOUNT`.
2.  **Command Generator**:
    -   Real-time preview of the `docker run` command.
    -   Option to generate a `docker-compose.yml` file.
3.  **Launch Control**:
    -   Button to "Run Agent" which executes the docker command in a terminal.

## 🛠️ Technical Implementation
-   **Backend**: `pkg/docker/generator.go`.
-   **Frontend**: Dedicated "Docker Agents" section.
