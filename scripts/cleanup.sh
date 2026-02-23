#!/bin/bash
# Test cleanup script for lazycontainer
# Removes test containers, images, volumes, and networks

set -e

echo "=== Lazycontainer Test Cleanup ==="
echo ""

# Check if container CLI is available
if ! command -v container &> /dev/null; then
    echo "Error: 'container' CLI not found. Make sure Apple Containers is installed."
    exit 1
fi

echo "1. Stopping and removing test containers..."
for container in test-alpine-stopped test-alpine-running test-web test-logger; do
    if container list --all --format json | grep -q "\"${container}\""; then
        echo "   - Stopping $container..."
        container stop "$container" 2>/dev/null || true
        echo "   - Removing $container..."
        container delete "$container" 2>/dev/null || true
    fi
done

echo ""
echo "2. Removing test images..."
for image in alpine:latest busybox:latest; do
    if container image list --format json | grep -q "\"reference\":\"${image}\""; then
        echo "   - Removing $image..."
        container image delete "$image" 2>/dev/null || true
    fi
done

echo ""
echo "3. Removing test volumes..."
if container volume list --format json | grep -q '"test-volume"'; then
    echo "   - Removing test-volume..."
    container volume delete test-volume 2>/dev/null || true
fi

echo ""
echo "4. Removing test networks..."
if container network list --format json | grep -q '"test-network"'; then
    echo "   - Removing test-network..."
    container network delete test-network 2>/dev/null || true
fi

echo ""
echo "=== Cleanup Complete ==="
echo ""
echo "Remaining containers:"
container list --all
echo ""
echo "Remaining images:"
container image list
echo ""
echo "Remaining volumes:"
container volume list
echo ""
echo "Remaining networks:"
container network list
