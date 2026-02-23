# Lazycontainer Testing Guide

## Overview

This document describes the testing strategy for lazycontainer. We use a multi-layer approach:

1. **Unit Tests** — Pure-logic tests that run anywhere, no container service required
2. **Integration Tests** — Test with real containers running locally
3. **Manual Test Scenarios** — Document expected user workflows for TUI behavior

---

## Unit Tests

### Running

```bash
go test ./...          # all packages
go test ./... -v       # verbose
go test -run TestClamp # specific test by name
```

### Test Suite (32 tests)

| Package | Tests | What's Covered |
|---------|-------|----------------|
| `pkg/commands` | 7 | JSON parsing (container, image, volume, network), command construction (logs stream, exec, stats stream) |
| `pkg/gui/panels` | 14 | `FilteredList` (set/get/filter/sort/index), `ListPanel` navigation with clamping, `SideListPanel` multi-select (toggle, count, clear) |
| `pkg/gui/presentation` | 3 | `GetContainerDisplayStrings` — column count, running vs stopped display, image prefix stripping |
| `pkg/utils` | 8 | `Clamp`, `Max`, `SplitLines`, `FormatBinaryBytes`, `FormatDecimalBytes`, `Decolorise`, `WithPadding`, `RenderTable` |

### What Is Not Unit Tested

- **TUI rendering** (`pkg/gui/gui.go`, views, keybindings) — gocui requires a real terminal; there is no headless test mode
- **Subprocess execution** (`pkg/gui/subprocess.go`) — PTY passthrough requires a real terminal
- **CLI integration** (`pkg/commands/container.go` at runtime) — executing `container list` etc. requires the Apple Container service

These are covered by manual test scenarios below.

---

## CI / Automated Checks

Every pull request to `main` runs `.github/workflows/test-build.yml`:

```
go vet ./...          → static analysis (catches type errors, unreachable code, etc.)
go test -v ./...      → all unit tests
go build -o lazycontainer .  → confirms the binary compiles
```

CI runs on `ubuntu-latest`.

### What CI Doesn't Automate (and Why)

**Coverage gates (e.g. minimum 40% threshold)**

We don't enforce a coverage minimum because:
- TUI code (the majority of the codebase) is untestable without a real terminal — a coverage gate would penalize us for code we cannot unit test
- The packages that *can* be unit tested are already at high coverage
- Worth revisiting once the testable surface grows

**Integration tests in CI (macOS runner)**

We don't run `container list`, `container start`, etc. in CI because:
- Apple's `container` CLI requires **macOS 26 (Tahoe)** — a beta OS not yet available on GitHub-hosted runners
- A self-hosted macOS 26 runner is an option for the future but out of scope now
- Integration tests are instead covered by the manual scenarios below

**Automated TUI testing**

gocui does not support headless/test mode. Terminal recording/replay frameworks (e.g., `vhs`) could be added in the future but are not practical now.

---

## Integration Tests

### Setup

Run the integration test environment setup script:

```bash
./scripts/test_setup.sh
```

This creates:
- 3 running containers (`test-alpine-running`, `test-web`, `test-logger`)
- 1 stopped container (`test-alpine-stopped`)
- Test volume (`test-volume`, 250MB)
- Test network (`test-network`)

### Verify

```bash
container list --all
```

Expected:
```
ID                   IMAGE                            STATE
test-alpine-running  docker.io/library/alpine:latest  running
test-web             docker.io/library/alpine:latest  running
test-logger          docker.io/library/alpine:latest  running
test-alpine-stopped  docker.io/library/alpine:latest  stopped
```

---

## Manual Test Scenarios

### Before Running

```bash
go build -o lazycontainer .
./scripts/test_setup.sh
./lazycontainer
```

### Navigation & UI

- [ ] **Panel switching** (1/2/3/4 keys)
  - Press '1' → See containers panel
  - Press '2' → See images panel
  - Press '3' → See volumes panel
  - Press '4' → See networks panel

- [ ] **List navigation** (j/k or arrows)
  - Navigate up/down through list
  - Cursor highlights current item

- [ ] **Tab switching** ([/] keys)
  - Press '[' → Previous tab (logs/config/stats)
  - Press ']' → Next tab

### Filtering

- [ ] **Open filter** ('/')
  - Press '/' → filter input appears
  - Type text → list filters in real-time

- [ ] **Commit filter** (Enter)
  - Type filter text, press Enter → filter persists

- [ ] **Clear filter** (Esc)
  - Press Esc in filter mode → filter cleared, original list restored

### Container Operations (Single)

- [ ] **Start container** (Enter on stopped)
  - Select `test-alpine-stopped`, press Enter → status changes to running

- [ ] **Stop container** ('s' on running)
  - Select `test-alpine-running`, press 's'
  - Confirmation dialog appears → accept → status changes to stopped

- [ ] **Delete container** ('d' on stopped)
  - Select stopped container, press 'd'
  - Confirmation dialog → accept → container disappears

### Multi-Select & Bulk Operations

- [ ] **Toggle selection** (Space)
  - Press Space → "Selected: 1 container(s)" in status bar
  - Press Space again → deselected

- [ ] **Bulk stop** (Space multiple, then 's')
  - Select `test-alpine-running` and `test-web` with Space
  - Press 's' → "Stop 2 container(s)?" → accept → both stopped

- [ ] **Bulk delete** (Space multiple, then 'd')
  - Select two stopped containers → press 'd' → "Delete 2 containers?" → both deleted
  - Selections clear after operation

### Streaming Features

- [ ] **Follow logs** ('f')
  - Select `test-logger`, press 'f'
  - Terminal shows: `+ container logs test-logger -f -n 100`
  - Logs stream in real-time
  - Ctrl+C → returns to UI, UI responsive

- [ ] **Exec/shell** ('e')
  - Select `test-alpine-running`, press 'e'
  - Terminal shows: `+ container exec test-alpine-running -- /bin/sh`
  - Interactive shell works; `exit` returns to UI

- [ ] **Stats streaming** ('v')
  - Select `test-alpine-running`, press 'v'
  - Terminal shows: `+ container stats test-alpine-running --format json`
  - Stats stream (CPU, memory, network); Ctrl+C returns to UI

### Error Handling

- [ ] **Operations on stopped container**
  - 's', 'e', 'v' on stopped → status bar shows "Container is not running"

- [ ] **Operations with no container selected**
  - 'f', 'e', 'v' → status bar shows "No container selected"

### State Management

- [ ] **Selections persist during panel navigation**
  - Space-select containers → switch to images (2) → back to containers (1) → selections still visible

- [ ] **Selections clear after bulk operation**
  - Bulk delete/stop → selections cleared, status bar empty

### UI Responsiveness

- [ ] **UI resumes cleanly after subprocess**
  - Run logs/exec/stats → exit with Ctrl+C → UI responds immediately, no hangs
  - All keybindings functional (navigation, filtering, operations)

---

## Regression Testing

### Before Merging

1. `go vet ./...` — no issues
2. `go test ./...` — all tests pass
3. Run manual test scenarios above
4. Add unit tests for any new pure-logic code

### Release Checklist

- [ ] All unit tests pass
- [ ] All manual test scenarios pass
- [ ] Build succeeds: `go build -o lazycontainer .`
- [ ] GitHub Actions CI passes
- [ ] No new `go vet` warnings

---

## Contributing Tests

When adding features:
1. Add unit tests for any pure-logic code (parsing, formatting, data structures)
2. Add manual test scenarios to this document for TUI behavior
3. Run `go test ./...` before opening a PR
