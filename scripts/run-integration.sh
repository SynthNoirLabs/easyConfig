#!/bin/bash
set -e

# Directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Root of the repo
ROOT_DIR="$(dirname "$DIR")"

IMAGE_NAME="easyconfig-integration-test"

echo "building docker image..."
docker build -t "$IMAGE_NAME" -f "$ROOT_DIR/tests/integration/Dockerfile" "$ROOT_DIR/tests/integration"

echo "running integration tests inside docker..."
# Mount the project root to /app
# Use --rm to clean up container after run
# Run as current user (if possible) or root? The Dockerfile sets USER tester.
# We need to make sure /app is accessible.
docker run --rm \
    -v "$ROOT_DIR:/app" \
    -w /app \
    -e CGO_ENABLED=1 \
    "$IMAGE_NAME" \
    go test -v ./tests/integration/...
