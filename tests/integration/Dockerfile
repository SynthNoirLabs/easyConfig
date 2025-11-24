FROM golang:1.24

# Install dependencies (Node.js/NPM)
RUN apt-get update && apt-get install -y \
    curl \
    git \
    nodejs \
    npm \
    && rm -rf /var/lib/apt/lists/*

# Install CLIs
RUN npm install -g @anthropic-ai/claude-code
RUN npm install -g @google/gemini-cli
RUN npm install -g @openai/codex

# Create dummy user for testing (optional, good for simulating non-root)
# RUN useradd -m -s /bin/bash tester
# USER tester
# WORKDIR /home/tester

WORKDIR /app

# Default command
CMD ["bash"]
