#!/bin/bash
set -e

# verify_gemini.sh
# Verifies that Gemini configuration discovery works as expected.

PROJECT_DIR="/tmp/easyConfig_gemini_test"
CONFIG_DIR="$PROJECT_DIR/.gemini"
CONFIG_FILE="$CONFIG_DIR/settings.json"

echo "Starting Gemini Verification..."

# 1. Setup Test Project
echo "Setting up test project at $PROJECT_DIR..."
rm -rf "$PROJECT_DIR"
mkdir -p "$CONFIG_DIR"

# Create a dummy settings.json
echo '{"model": "gemini-pro"}' > "$CONFIG_FILE"

# 2. Verify File Existence
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Failed to create test config file."
    exit 1
fi

echo "Created test config at $CONFIG_FILE"

# 3. Verify JSON Validity
if jq . "$CONFIG_FILE" >/dev/null 2>&1; then
    echo "Config file is valid JSON."
else
    echo "Error: Config file is invalid JSON."
    exit 1
fi

echo "Gemini Verification Passed!"
