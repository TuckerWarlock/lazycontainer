# Copilot Instructions for lazycontainer

## Build, Test, and Lint

### Build
```bash
go build -o lazycontainer .
```

### Tests
```bash
go test ./...           # Run all tests
go test -v ./...        # Verbose output
go test -run TestName   # Run specific test by name (e.g., go test -run TestClamp)
go test ./pkg/...       # Run tests for pkg directory only
```

**Test Coverage**: 32 unit tests across `pkg/commands`, `pkg/gui/panels`, `pkg/gui/presentation`, and `pkg/utils`. Tests verify JSON parsing, list filtering/sorting, panel navigation, and utility functions. TUI rendering and subprocess execution require manual testing (see docs/TESTING.md).

### Linting and Static Analysis
```bash
go vet ./...    # Static analysis (catches type errors, unreachable code, etc.)
```

### Debug Mode
```bash
./lazycontainer -d    # Enables debug logging to ~/.config/lazycontainer/development.log
```

## High-Level Architecture

**lazycontainer** is a terminal UI for managing Apple containers on macOS 26+ (built with gocui, adapted from lazydocker).

### Package Structure

- **pkg/app** — Application bootstrap, wires together config, GUI, and task manager
- **pkg/commands** — Wrapper around Apple `container` CLI; parses output into structured types (Container, Image, Volume, Network)
- **pkg/config** — App configuration (paths, version info, debug mode)
- **pkg/gui** — Terminal UI engine
  - **gui.go** — Main run loop, refresh logic, task coordination
  - **views.go** — gocui view definitions (containers panel, logs view, etc.)
  - **keybindings.go** — Keyboard and mouse handlers
  - **panels/** — Reusable panel components:
    - `SideListPanel[T]` — Generic abstraction for resource lists with selection and filtering
    - `FilteredList[T]` — Generic filtered/sorted list with clamped indexing
    - `ContextState` — Tab management for logs/config/stats tabs
  - **presentation/** — Display formatters for each resource type
  - **subprocess.go** — PTY passthrough for exec/logs streaming
- **pkg/i18n** — Translations (english.go)
- **pkg/tasks** — TaskManager for async operations (keeps UI responsive)
- **pkg/utils** — Utilities (formatting, colors, clamping, table rendering)

### Key Design Patterns

1. **SideListPanel[T]**: Generic panel for listing resources with multi-select, filtering by name (`/`), and sorting. Tracks selection state and filtered indices independently.

2. **ContextState**: Manages multiple tabs (logs, config, stats) per resource. Each tab has its own cached data to avoid redundant API calls. Tab switching is done with `[` and `]`.

3. **TaskManager**: Queues async operations (e.g., `container list`) to prevent UI freezes. Main run loop checks for completed tasks and updates view.

4. **IGui Interface**: Panels interact with GUI through abstraction instead of direct coupling, enabling testability.

### Apple Container CLI vs Docker

The `container` CLI behaves differently from Docker in key ways:

| Operation | Docker | Apple Container |
|-----------|--------|-----------------|
| List containers | `docker ps` | `container list` |
| Log tail | `docker logs --tail N` | `container logs -n N` |
| Timestamp format | Unix epoch | CFAbsoluteTime (seconds since Jan 1, 2001) |
| Exec interactive | `docker exec -it id cmd` | `container exec id -- cmd` (PTY via stdin/stdout) |
| No compose equivalent | `docker-compose` | N/A |
| Image pull | `docker pull name` | `container image pull name` |
| Volume operations | `docker volume` | `container volume` |
| Network operations | `docker network` | `container network` |

See `pkg/commands/container_types.go` for the exact JSON structures returned by each command.

## Key Conventions

### Naming and File Organization

- Test files: `*_test.go` (unit tests only; no integration tests in CI)
- Package-level utility functions: `utils.go` in each package
- Presentation logic: `pkg/gui/presentation/{resource}_test.go` for display formatting
- Panel types: `{resource}_panel.go` in `pkg/gui/panels`

### Error Handling

- Use `go-errors/errors` library for stack traces: `errors.Wrap(err, "context")` or `errors.New("message")`
- Log errors with context: `logrus.WithError(err).Error("operation failed")`
- Return early for error cases; don't nest deeply

### Testing

- Use standard `testing` package (no test frameworks)
- Test table patterns for multiple scenarios (see `utils_test.go`)
- Mock the IGui interface for panel tests (see `panels/*_test.go`)
- Manual test scenarios in `docs/TESTING.md` for TUI behavior (filtering, keybindings, etc.)

### Async Operations

- All long-running operations (CLI calls, file I/O) must use TaskManager, not blocking calls
- TaskManager enqueues tasks and main run loop checks `tasks.GetCompletedTasks()` to update views
- Never block the main GUI run loop

### Color and Formatting

- Use `fatih/color` for terminal colors (e.g., `color.GreenString("text")`)
- Use `FormatBinaryBytes` for file sizes (binary: 1024), `FormatDecimalBytes` for network I/O (decimal: 1000)
- Use `WithPadding` for column alignment in tables
- Use `Decolorise` to strip ANSI codes when needed

### Keybindings

Global keybindings (see `keybindings.go`):
- `1/2/3/4` — Switch panels (containers/images/volumes/networks)
- `j/↓` and `k/↑` — Move up/down in lists
- `[` and `]` — Previous/next tab (logs/config/stats)
- `Tab` — Cycle panels
- `/` — Filter current panel by name
- `Space` — Multi-select (toggles selection)
- `Enter` — Start/apply action
- `s` — Stop container (with confirmation)
- `d` — Delete (with confirmation)
- `f` — Follow logs (streaming)
- `e` — Exec interactive shell
- `v` — Stream live stats
- `c` — Copy resource ID to clipboard
- `C` — Copy logs to clipboard
- `t` — Cycle themes (default/light/dark)
- `?` — Show context-sensitive help
- `m` — Toggle mouse support (useful for tmux)
- `q` — Quit

Resource-specific keybindings can be added in panel handlers.

### Confirmation Dialogs

Destructive actions (stop, delete, remove volume/network) trigger a confirmation panel. Pattern:
1. Show `ConfirmationPanel` with action description
2. User presses `Enter` or `y` to confirm, `n` or `Esc` to cancel
3. On confirm, enqueue task with TaskManager

### Dependencies and Imports

- **TUI**: `github.com/jesseduffield/gocui`
- **CLI parsing**: `github.com/integrii/flaggy`
- **Logging**: `github.com/sirupsen/logrus`
- **Color**: `github.com/fatih/color`
- **Error handling**: `github.com/go-errors/errors`
- **Utilities**: `github.com/samber/lo` (functional utilities)

## Reference Project

lazydocker (`/Users/warl0ck/Code/lazydocker/`) is the reference implementation. When adding new features:
1. Find the equivalent feature in lazydocker
2. Copy the pattern and update imports from `github.com/jesseduffield/lazydocker` to `github.com/warl0ck/lazycontainer`
3. Replace Docker CLI calls with Apple Container CLI equivalents
4. For design questions about behavior, use lazydocker as the authoritative reference—match it exactly rather than inventing alternatives

## Implementation Status

See `docs/IMPLEMENTATION.md` for detailed feature parity with lazydocker. 

**Phase 1–2 (Core & Interactive)**: Complete ✅
- Container/image/volume/network management
- Interactive shell access
- Streaming logs and stats
- Real-time filtering and bulk operations
- Confirmation dialogs

**Phase 3 (User Experience)**: Not implemented
- Custom commands, menu system, config file

**Phase 4 (Polish)**: Complete ✅
- Copy to clipboard (`c`, `C`)
- Theme customization (`t` cycles 3 themes)
- Context-sensitive help (`?`)
- Mouse mode toggle (`m` for tmux compatibility)
