#!/bin/bash
set -e

# Check if gh is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed."
    echo "Please install it: https://cli.github.com/"
    exit 1
fi

# Check if user is logged in
if ! gh auth status &> /dev/null; then
    echo "Error: You are not logged into GitHub CLI."
    echo "Please run: gh auth login"
    exit 1
fi

BRANCH="main"
if [ -n "$1" ]; then
    BRANCH="$1"
fi

echo "🚀 Triggering release workflow on branch: $BRANCH"
gh workflow run release.yml --ref "$BRANCH"

echo "✅ Workflow triggered successfully!"
echo "👀 Watch the progress with:"
echo "   gh run watch $(gh run list --workflow=release.yml --limit 1 --json databaseId --jq '.[0].databaseId')"
