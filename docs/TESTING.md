# Lazycontainer Testing Guide

## Overview

This document describes the testing strategy for lazycontainer. We use a multi-layer approach:

1. **Unit Tests** - Test command execution and data parsing
2. **Integration Tests** - Test with real containers  
3. **Manual Test Scenarios** - Document expected user workflows

## Unit Tests

### Location
- `pkg/commands/container_test.go` - Command construction and parsing
- Run with: `go test ./pkg/commands -v`

### Current Coverage
- Container JSON parsing
- Image JSON parsing  
- Volume JSON parsing
- Network JSON parsing
- Streaming logs command construction
- Exec/attach command construction
- Stats streaming command construction

### Adding Tests
When adding new features:
```bash
go test ./pkg/commands -v
go test ./... -v  # All packages
```

## Integration Tests

### Setup
Run the integration test environment setup script:
```bash
./scripts/test_setup.sh
```

This creates:
- 3 running containers (test-alpine-running, test-web, test-logger)
- 1 stopped container (test-alpine-stopped)
- Test volume (test-volume, 250MB)
- Test network (test-network)

### Environment
Verify containers are running:
```bash
container list
```

Expected output:
```
ID                   IMAGE                            STATE    
test-alpine-running  docker.io/library/alpine:latest  running
test-web             docker.io/library/alpine:latest  running
test-logger          docker.io/library/alpine:latest  running
test-alpine-stopped  docker.io/library/alpine:latest  stopped
```

## Manual Test Scenarios

### Before Running
1. Build the app: `go build -o lazycontainer .`
2. Ensure containers are running: `./scripts/test_setup.sh`
3. Start the app: `./lazycontainer`

### Test Cases

#### Navigation & UI
- [ ] **Panel switching** (1/2/3/4 keys)
  - Press '1' → See containers panel
  - Press '2' → See images panel
  - Press '3' → See volumes panel
  - Press '4' → See networks panel

- [ ] **List navigation** (j/k or arrows)
  - Navigate up/down through list
  - Cursor highlights current item
  - Selection indicator shows on containers

- [ ] **Tab switching** ([/] keys for main tabs)
  - Press '[' → Previous tab
  - Press ']' → Next tab
  - See logs/config/stats tabs rotate

#### Filtering
- [ ] **Open filter** ('/')
  - Press '/'
  - Type filter text
  - List filters in real-time

- [ ] **Commit filter** (Enter)
  - Type filter text
  - Press Enter
  - Filter persists

- [ ] **Clear filter** (Esc)
  - Press Esc in filter mode
  - Filter is cleared
  - Original list shows

#### Container Operations (Single)
- [ ] **Start container** (Enter on stopped)
  - Select test-alpine-stopped
  - Press Enter
  - Status changes to running

- [ ] **Stop container** ('s' on running)
  - Select test-alpine-running
  - Press 's'
  - Confirmation dialog appears
  - Accept confirmation
  - Status changes to stopped
  - Undo: Start it again

- [ ] **Delete container** ('d' on stopped)
  - Select stopped container
  - Press 'd'
  - Confirmation dialog appears
  - Accept confirmation
  - Container disappears from list

#### Multi-Select & Bulk Operations
- [ ] **Toggle selection** (Space)
  - Select test-web
  - Press Space
  - Status shows "Selected: 1 container(s)"
  - Press Space again → deselected

- [ ] **Bulk stop** (Space multiple, press 's')
  - Select test-alpine-running (Space)
  - Select test-web (Space)
  - Status shows "Selected: 2 container(s)"
  - Press 's'
  - Confirmation: "Stop 2 container(s)?"
  - Accept confirmation
  - Both containers stopped

- [ ] **Bulk delete** (Space multiple, press 'd')
  - Select two stopped containers (Space each)
  - Status shows "Selected: 2 container(s)"
  - Press 'd'
  - Confirmation: "Delete 2 containers?"
  - Accept confirmation
  - Both deleted, selections cleared

#### Streaming Features
- [ ] **Follow logs** ('f')
  - Select test-logger
  - Press 'f'
  - Terminal shows: `+ container logs test-logger -f -n 100`
  - Logs stream in real-time
  - Press Ctrl+C → Returns to UI
  - Verify UI is responsive

- [ ] **Exec/shell** ('e')
  - Select test-alpine-running
  - Press 'e'
  - Terminal shows: `+ container exec test-alpine-running -- /bin/sh`
  - Get interactive shell prompt
  - Type: `echo "test"` → output appears
  - Type: `exit` → Returns to UI
  - Verify UI is responsive

- [ ] **Stats streaming** ('v')
  - Select test-alpine-running
  - Press 'v'
  - Terminal shows: `+ container stats test-alpine-running --format json`
  - Stats stream with CPU, memory, network I/O
  - Press Ctrl+C → Returns to UI
  - Verify UI is responsive

#### Error Handling
- [ ] **Operations on stopped container**
  - Select stopped container
  - Press 's' → Status: "Container is not running"
  - Press 'e' → Status: "Container is not running"
  - Press 'v' → Status: "Container is not running"

- [ ] **Operations with no selection**
  - If no container selected somehow
  - Press 'f', 'e', 'v' → Status: "No container selected"

#### State Management
- [ ] **Selections persist during navigation**
  - Select multiple containers (Space)
  - Navigate to images panel (2 key)
  - Navigate back to containers (1 key)
  - Selections should still be visible

- [ ] **Selections clear after operation**
  - Select containers (Space)
  - Status shows selection count
  - Perform bulk operation (delete/stop)
  - Selections clear
  - Status bar is empty

#### UI Responsiveness
- [ ] **UI resumes cleanly after subprocess**
  - Run streaming command (logs/exec/stats)
  - Exit with Ctrl+C
  - UI should respond immediately
  - No freezes or hangs

- [ ] **All keybindings work after subprocess**
  - After exiting subprocess
  - Test navigation (1/2/3/4, j/k, arrows)
  - Test filtering (/)
  - Test operations (s/d/f/e/v)

## Automated Testing Commands

### Build & Test
```bash
# Build
go build -o lazycontainer .

# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run specific test
go test -run TestGetContainerStatsStream -v
```

### Integration Test Flow
```bash
# 1. Set up test environment
./scripts/test_setup.sh

# 2. Verify containers are ready
container list

# 3. Build the app
go build -o lazycontainer .

# 4. Run manual test scenarios (see above)
./lazycontainer

# 5. Verify all scenarios pass
# Document any issues or missing features
```

## Regression Testing

### Before Merging Features
1. Run unit tests: `go test ./...`
2. Run manual test scenarios (see above)
3. Document any issues
4. Add new tests for new features

### Release Checklist
- [ ] All unit tests pass
- [ ] All manual test scenarios pass
- [ ] Build succeeds on Go 1.26
- [ ] GitHub Actions CI passes
- [ ] No new warnings or errors

## Known Issues & Limitations

- TUI testing is primarily manual (gocui doesn't have automated UI testing support)
- Integration tests require running Apple Container service (`container system start`)
- Some features (like copy to clipboard) require platform-specific implementations

## Future Improvements

- [ ] Automated UI testing framework (possibly with terminal recording/replay)
- [ ] Performance benchmarks
- [ ] Load testing (many containers)
- [ ] Multi-platform testing (macOS versions, arm64 vs x86)

## Contributing Tests

When adding features:
1. Add unit tests in appropriate `*_test.go` file
2. Add manual test scenarios to this document
3. Run full test suite before PR
4. Update this document as needed
