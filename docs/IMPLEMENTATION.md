# Lazycontainer Implementation Status

A TUI for Apple's Containerization framework (macOS 26+), adapted from [lazydocker](https://github.com/jesseduffield/lazydocker).

---

## Current Status

### Phases Complete

**Phase 1 — Core TUI** ✅
**Phase 2 — Interactive Features** ✅
**Testing Infrastructure** ✅

---

## Lazydocker Feature Parity

| Feature | lazydocker | lazycontainer | Notes |
|---------|-----------|---------------|-------|
| Container list | ✅ | ✅ | |
| Start / Stop / Delete | ✅ | ✅ | |
| Container logs (static) | ✅ | ✅ | Uses `-n` not `--tail` |
| Container inspect/config | ✅ | ✅ | |
| Kill with signal | ✅ | ✅ | |
| Follow logs (`-f`) | ✅ | ✅ | 'f' keybinding |
| Exec / interactive shell | ✅ | ✅ | 'e' keybinding; `container exec <id> -- /bin/sh` |
| Streaming stats | ✅ | ✅ | 'v' keybinding; different JSON format |
| Bulk operations | ✅ | ✅ | Space to select; 's'/'d' for bulk stop/delete |
| Filtering (`/`) | ✅ | ✅ | Real-time filtering |
| Confirmation dialogs | ✅ | ✅ | y/Enter to confirm, n/Esc to cancel |
| Status bar | ✅ | ✅ | Basic — no animated spinner |
| Image list / delete | ✅ | ✅ | |
| Image pull | ✅ | ✅ | |
| Volume list / create / delete | ✅ | ✅ | |
| Network list / create / delete | ✅ | ✅ | |
| Panel navigation (1/2/3/4) | ✅ | ✅ | |
| Tab cycling ([/]) | ✅ | ✅ | logs / config / stats tabs |
| Mouse support | ✅ | ✅ | Click to select, scroll |
| Custom commands | ✅ | ❌ | Phase 3 |
| Menu system | ✅ | ❌ | Phase 3 |
| Options view (keybinding help) | ✅ | ❌ | Phase 3 |
| Config file (`~/.config/`) | ✅ | ❌ | Phase 3 |
| Copy to clipboard | ✅ | ❌ | Phase 4 |
| Container top | ✅ | ❌ | Needs investigation — Apple CLI may not support |
| Docker Compose / Services | ✅ | N/A | Apple containers have no compose equivalent |
| SSH tunneling | ✅ | N/A | Local containers only |

---

## What's Different from lazydocker

### Apple Container CLI vs Docker CLI

| Operation | Docker | Apple Container |
|-----------|--------|-----------------|
| Log tail | `--tail N` | `-n N` |
| Timestamp format | Unix epoch | CFAbsoluteTime (seconds since Jan 1, 2001) |
| Exec | `docker exec -it` | `container exec <id> -- cmd` |
| No compose | `docker-compose` | N/A |

### Architecture Differences

1. **No Docker Compose** — No Services/Projects panels
2. **No SSH tunneling** — Local containers only; no remote support
3. **Different stats JSON** — Apple's format uses `cpuUsageUsec`, `memoryUsageBytes`, etc.
4. **CFAbsoluteTime** — Timestamps are seconds since Jan 1, 2001 (not Unix epoch)

---

## Implementation Roadmap

### Phase 1: Core TUI ✅
- [x] Container, image, volume, network panels
- [x] Start / stop / delete operations
- [x] Logs and config tabs
- [x] Filtering (`/` to filter)
- [x] Confirmation dialogs
- [x] Status bar
- [x] Mouse support

### Phase 2: Interactive Features ✅
- [x] Streaming logs — `f` key, `container logs <id> -f`
- [x] Container exec/attach — `e` key, `container exec <id> -- /bin/sh`
- [x] Bulk operations — Space to multi-select, `s`/`d` for bulk actions
- [x] Streaming stats — `v` key, `container stats <id> --format json`

### Testing Infrastructure ✅
- [x] 32 unit tests across `pkg/commands`, `pkg/gui/panels`, `pkg/gui/presentation`, `pkg/utils`
- [x] `go vet` in CI
- [x] `scripts/test_setup.sh` for integration test environment
- [x] Manual test scenarios documented in `docs/TESTING.md`

### Phase 3: User Experience
- [ ] Config file (`~/.config/lazycontainer/config.yml`) — custom themes, keybindings
- [ ] Custom commands per resource type with template variables
- [ ] Menu system for complex choices (signal selection, etc.)
- [ ] Options view showing available keybindings in context

### Phase 4: Polish
- [ ] Information panel with context-sensitive help
- [ ] Copy to clipboard (container ID, logs)
- [ ] Mouse mode toggle (for tmux compatibility)
- [ ] Theme customization

---

## Apple Container CLI Reference

```bash
# List containers
container list --format json
container list --all --format json

# Container operations
container start <id>
container stop <id>
container delete <id>
container logs <id> -n <lines>        # Note: -n not --tail
container logs <id> -f -n <lines>     # Follow
container inspect --format json <id>
container kill <id> --signal <signal>
container exec <id> -- <command>      # Interactive execution
container stats <id> --format json    # Streaming stats

# Images
container image list --format json
container image pull <ref>
container image delete <ref>

# Volumes
container volume list --format json
container volume create <name>
container volume delete <name>

# Networks
container network list --format json
container network create <name>
container network delete <name>

# System
container system status
container system start
container system stop
```

---

## Directory Structure

```
lazycontainer/
├── main.go
├── go.mod / go.sum
├── docs/
│   ├── IMPLEMENTATION.md     ← this file
│   └── TESTING.md
├── scripts/
│   └── test_setup.sh         # Creates test containers/volumes/networks
├── pkg/
│   ├── app/app.go            # Application wiring
│   ├── commands/
│   │   ├── container.go      # Apple container CLI wrapper
│   │   ├── container_types.go
│   │   ├── container_test.go # 7 unit tests
│   │   └── os.go
│   ├── config/app_config.go
│   ├── gui/
│   │   ├── gui.go            # Main GUI struct, run loop, refresh
│   │   ├── views.go          # gocui view definitions
│   │   ├── keybindings.go    # Keyboard/mouse bindings
│   │   ├── subprocess.go     # PTY subprocess runner (logs/exec/stats)
│   │   ├── panels/
│   │   │   ├── filtered_list.go       # Generic filtered/sorted list (6 tests)
│   │   │   ├── list_panel.go          # Base navigation (3 tests)
│   │   │   ├── context_state.go       # Tab management
│   │   │   └── side_list_panel.go     # Main panel + multi-select (5 tests)
│   │   ├── presentation/
│   │   │   ├── containers.go          # (3 tests)
│   │   │   ├── images.go
│   │   │   ├── volumes.go
│   │   │   └── networks.go
│   │   └── types/types.go
│   ├── i18n/
│   ├── tasks/tasks.go        # Async task management
│   └── utils/utils.go        # Formatting utilities (8 tests)
```

---

## Development Commands

```bash
# Build
go build -o lazycontainer .

# Run
./lazycontainer

# Debug mode (logs to ~/.config/lazycontainer/development.log)
./lazycontainer -d

# Tests
go vet ./...
go test ./...
go test ./... -v

# Set up test environment
./scripts/test_setup.sh
```

---

## Reference Implementation

The lazydocker source at `/Users/warl0ck/Code/lazydocker/` serves as the reference. When adding features:

1. Find the equivalent in lazydocker
2. Copy patterns, update imports: `github.com/jesseduffield/lazydocker` → `github.com/warl0ck/lazycontainer`
3. Replace Docker CLI calls with Apple `container` CLI equivalents
4. Update this file and `docs/TESTING.md`
