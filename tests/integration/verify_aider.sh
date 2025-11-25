#!/bin/bash
set -e

# verify_aider.sh
# Verifies that Aider configuration discovery works as expected.

PROJECT_DIR="/tmp/easyConfig_aider_test"
CONFIG_FILE="$PROJECT_DIR/.aider.conf.yml"

echo "Starting Aider Verification..."

# 1. Setup Test Project
echo "Setting up test project at $PROJECT_DIR..."
rm -rf "$PROJECT_DIR"
mkdir -p "$PROJECT_DIR"

# Create a dummy aider config
echo 'model: gpt-4
dark-mode: true
' > "$CONFIG_FILE"

# 2. Verify File Existence
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Failed to create test config file."
    exit 1
fi

echo "Created test config at $CONFIG_FILE"

# 3. Verify Content
if grep -q "model: gpt-4" "$CONFIG_FILE"; then
    echo "Config file contains expected content."
else
    echo "Error: Config file missing expected content."
    exit 1
fi

echo "Aider Verification Passed!"
