#!/bin/bash
# Test setup script for lazycontainer
# Creates test containers, volumes, and networks for testing the TUI

set -e

echo "=== Lazycontainer Test Setup ==="
echo ""

# Check if container CLI is available
if ! command -v container &> /dev/null; then
    echo "Error: 'container' CLI not found. Make sure Apple Containers is installed."
    exit 1
fi

echo "1. Pulling alpine image..."
container image pull alpine:latest || {
    echo "Warning: Failed to pull alpine, trying busybox..."
    container image pull busybox:latest
}

echo ""
echo "2. Creating test volume..."
container volume create test-volume 2>/dev/null || echo "Volume 'test-volume' may already exist"

echo ""
echo "3. Creating test network..."
container network create test-network 2>/dev/null || echo "Network 'test-network' may already exist"

echo ""
echo "4. Creating test containers..."

# Container 1: Simple alpine container (will be stopped)
echo "   - Creating 'test-alpine-stopped'..."
container run --name test-alpine-stopped -d alpine:latest sleep 1 2>/dev/null || true
# It will stop after 1 second

# Container 2: Long-running alpine container
echo "   - Creating 'test-alpine-running'..."
container run --name test-alpine-running -d alpine:latest sleep 3600 2>/dev/null || {
    echo "   Container may already exist, starting it..."
    container start test-alpine-running 2>/dev/null || true
}

# Container 3: Container with port mapping
echo "   - Creating 'test-web'..."
container run --name test-web -d -p 8080:80 alpine:latest sh -c "while true; do echo -e 'HTTP/1.1 200 OK\n\nHello from lazycontainer test!' | nc -l -p 80; done" 2>/dev/null || {
    echo "   Container may already exist"
    container start test-web 2>/dev/null || true
}

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Containers:"
container list --all
echo ""
echo "Images:"
container image list
echo ""
echo "Volumes:"
container volume list
echo ""
echo "Networks:"
container network list
echo ""
echo "You can now run: ./lazycontainer"
