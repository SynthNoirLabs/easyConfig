FROM ubuntu:latest

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    golang-go \
    nodejs \
    npm \
    && rm -rf /var/lib/apt/lists/*

# Install CLIs
RUN npm install -g @anthropic-ai/claude-code
RUN npm install -g @google/gemini-cli
RUN npm install -g @openai/codex

# Create dummy user for testing
RUN useradd -m -s /bin/bash tester
USER tester
WORKDIR /home/tester

# Setup project structure
RUN mkdir -p /home/tester/project
WORKDIR /home/tester/project

# Copy integration test files (to be mounted or copied)
# COPY . .

# Default command
CMD ["bash"]
