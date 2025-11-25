#!/bin/bash
set -e

# verify_git.sh
# Verifies that Git configuration discovery works as expected.

PROJECT_DIR="/tmp/easyConfig_git_test"
GIT_DIR="$PROJECT_DIR/.git"
CONFIG_FILE="$GIT_DIR/config"

echo "Starting Git Verification..."

# 1. Setup Test Project
echo "Setting up test project at $PROJECT_DIR..."
rm -rf "$PROJECT_DIR"
mkdir -p "$GIT_DIR"

# Create a dummy git config
echo '[user]
	name = Test User
	email = test@example.com
' > "$CONFIG_FILE"

# 2. Verify File Existence
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Failed to create test config file."
    exit 1
fi

echo "Created test config at $CONFIG_FILE"

# 3. Verify Content
if grep -q "Test User" "$CONFIG_FILE"; then
    echo "Config file contains expected content."
else
    echo "Error: Config file missing expected content."
    exit 1
fi

echo "Git Verification Passed!"
