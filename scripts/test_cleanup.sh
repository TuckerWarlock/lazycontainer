#!/bin/bash
# Test cleanup script for lazycontainer
# Removes test containers, volumes, and networks

set -e

echo "=== Lazycontainer Test Cleanup ==="
echo ""

echo "1. Stopping test containers..."
container stop test-alpine-stopped test-alpine-running test-web 2>/dev/null || true

echo ""
echo "2. Removing test containers..."
container delete test-alpine-stopped test-alpine-running test-web 2>/dev/null || true

echo ""
echo "3. Removing test volume..."
container volume delete test-volume 2>/dev/null || true

echo ""
echo "4. Removing test network..."
container network delete test-network 2>/dev/null || true

echo ""
echo "5. Optionally removing images (uncomment if desired)..."
# container image delete alpine:latest busybox:latest 2>/dev/null || true

echo ""
echo "=== Cleanup Complete ==="
container list --all
