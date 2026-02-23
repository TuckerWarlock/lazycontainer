# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Lazycontainer is a terminal UI for managing Apple containers on macOS 26+ (Tahoe), adapted from [lazydocker](https://github.com/jesseduffield/lazydocker). It wraps the Apple `container` CLI with an interactive TUI built using gocui.

## Build and Run Commands

```bash
# Build
go build -o lazycontainer .

# Run
./lazycontainer

# Run with debug logging (logs to ~/.config/lazycontainer/development.log)
./lazycontainer -d

# Run tests
go test ./pkg/commands/...

# Set up test containers/volumes/networks
./scripts/test_setup.sh
```

## Prerequisites

- macOS 26+ (Tahoe)
- Apple Containerization framework: `brew install --cask container`
- Container service must be running: `container system start`

## Architecture

### Package Structure

```
pkg/
├── app/           # Application bootstrap (wires everything together)
├── commands/      # Apple container CLI wrapper
│   ├── container.go       # CLI operations (list, start, stop, etc.)
│   └── container_types.go # Container, Image, Volume, Network structs
├── config/        # Application configuration
├── gui/           # Terminal UI (gocui-based)
│   ├── gui.go         # Main GUI struct, Run loop, refresh logic
│   ├── views.go       # gocui view definitions
│   ├── keybindings.go # Keyboard/mouse bindings
│   ├── panels/        # Reusable panel components
│   │   ├── side_list_panel.go  # Main panel abstraction
│   │   ├── filtered_list.go    # Generic filtered/sorted list
│   │   └── context_state.go    # Tab management (logs/config/stats)
│   └── presentation/  # Display formatters for each resource type
├── i18n/          # Translations (english.go)
├── tasks/         # Async task management for non-blocking UI
└── utils/         # Utilities (formatting, colors, clamping)
```

### Key Patterns (from lazydocker)

1. **SideListPanel[T]**: Generic panel for listing resources with selection, filtering, and tabs
2. **ContextState**: Manages tabs (logs/config/stats) per resource type with caching
3. **TaskManager**: Queues async operations to keep UI responsive
4. **IGui interface**: Abstraction allowing panels to interact with GUI without tight coupling

### Apple Container CLI Differences from Docker

| Operation | Docker | Apple Container |
|-----------|--------|-----------------|
| Log tail | `--tail N` | `-n N` |
| Timestamp format | Unix epoch | CFAbsoluteTime (seconds since Jan 1, 2001) |
| Exec | `docker exec -it` | `container exec <id> -- cmd` |
| No compose | `docker-compose` | N/A (not supported) |

## Reference Project

The lazydocker source at `/Users/warl0ck/Code/lazydocker/` serves as the reference implementation. When adding features:

1. Find the equivalent in lazydocker
2. Copy patterns, update imports from `github.com/jesseduffield/lazydocker` to `github.com/warl0ck/lazycontainer`
3. Replace Docker CLI calls with Apple container CLI equivalents
4. **For design questions, use lazydocker as the authoritative reference** - match its behavior exactly rather than asking the user

## Development Workflow

After implementing new functionality:
1. Run `go build -o lazycontainer .` to verify compilation
2. **Always run `go test ./...` to ensure tests pass**
3. Test manually with `./lazycontainer` if applicable
4. **Update `IMPLEMENTATION_PLAN.md`** to mark completed items and track progress

## Cleanup

When creating temporary design docs or plans (e.g., `docs/plans/*.md`), remind the user to delete them after implementation is complete. These are working documents, not permanent repo documentation.

## Implementation Status

See `IMPLEMENTATION_PLAN.md` for detailed status. Key missing features:
- Filtering UI (`/` to filter)
- Container attach/exec with PTY
- Bulk operations
