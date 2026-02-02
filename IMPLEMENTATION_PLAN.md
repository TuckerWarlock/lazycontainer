# Lazycontainer Implementation Plan

A TUI for Apple's Containerization framework (macOS 26+), adapted from lazydocker.

## Current Status

### Completed

1. **Basic TUI working** - Container list, details panel, keybindings functional
2. **Core structure created**:
   - `main.go` - Entry point with CLI flags
   - `pkg/app/app.go` - Application wiring
   - `pkg/config/app_config.go` - Configuration
   - `pkg/commands/os.go` - OS command execution
   - `pkg/commands/container.go` - Apple container CLI wrapper
   - `pkg/commands/container_types.go` - Container, Image, Volume, Network structs with helper methods
   - `pkg/gui/gui.go` - Full TUI with panels system and IGui interface
   - `pkg/gui/views.go` - Views struct with all gocui views
   - `pkg/gui/keybindings.go` - Keybindings with mouse support
   - `pkg/i18n/english.go` - Translations
   - `pkg/i18n/i18n.go` - Translation loader
   - `pkg/tasks/tasks.go` - Task management for async operations
   - `pkg/gui/types/types.go` - MenuItem type

3. **Panel Infrastructure** (`pkg/gui/panels/`):
   - `filtered_list.go` - Generic filtered list with sorting
   - `list_panel.go` - Base list panel with selection
   - `context_state.go` - Tab management for main panel
   - `side_list_panel.go` - Main reusable panel component

4. **Presentation Layer** (`pkg/gui/presentation/`):
   - `containers.go` - Container display with colored status indicators
   - `images.go` - Image display formatter
   - `volumes.go` - Volume display formatter
   - `networks.go` - Network display formatter

5. **Utility Functions** (`pkg/utils/utils.go`):
   - `Clamp`, `RenderTable`, `ColoredString`, `WithPadding`
   - `SplitLines`, `FormatBinaryBytes`, `FormatDecimalBytes`

6. **Container Operations**:
   - List containers (all/running)
   - Start/Stop/Delete containers
   - View container logs (with `-n` flag for Apple CLI)
   - View container config/inspect
   - Kill container with signal

7. **Image Operations**:
   - List images
   - Pull images
   - Delete images

8. **Volume Operations**:
   - List volumes
   - Create volumes
   - Delete volumes

9. **Network Operations**:
   - List networks
   - Create networks
   - Delete networks

10. **Navigation**:
    - j/k or arrow keys for list navigation
    - 1/2/3/4 keys for panel switching
    - Tab cycling through panels
    - Mouse wheel scrolling
    - [ ] for main tab switching
    - Enter, s, d, x keybindings for actions

11. **Testing Infrastructure**:
    - `pkg/commands/container_test.go` - Unit tests for JSON parsing
    - `scripts/test_setup.sh` - Script to create test containers

### Partially Implemented

1. **Mouse Selection** - Basic mouse click binding exists but cursor highlight behavior needs refinement
2. **Stats View** - Framework in place but `container stats --format json` parsing not complete

---

## Comparison with Lazydocker: Missing Features

The following features exist in lazydocker but are NOT yet implemented in lazycontainer:

### High Priority - Core Functionality

#### 1. Confirmation Dialogs (`pkg/gui/confirmation_panel.go`)
- Modal confirmation before destructive actions (delete, stop, prune)
- lazydocker uses `CreateConfirmationPanel()` and `CreateMenu()`
- **Status:** Not implemented - currently actions happen immediately

#### 2. Filtering System (`pkg/gui/filtering.go`)
- `/` to enter filter mode
- Real-time filtering of panel lists
- Filter indicator in view title
- **Status:** Panel infrastructure supports it (`DisableFilter` field, `FilterString()`) but UI not connected

#### 3. Custom Commands (`pkg/config/user_config.go`)
- User-defined commands per resource type
- Template variables (`{{.Container.ID}}`, etc.)
- Custom keybindings for commands
- **Status:** Not implemented

#### 4. Status Manager (`pkg/gui/status_manager.go`)
- Bottom status line with messages
- Error display with colors
- Loading indicators
- **Status:** Not implemented - errors currently silent or logged

### Medium Priority - User Experience

#### 5. Container Attach/Exec (`pkg/commands/container.go`)
- Interactive shell into container
- lazydocker: `docker exec -it container /bin/sh`
- Apple CLI equivalent: `container exec <id> -- /bin/sh`
- **Status:** Not implemented

#### 6. Subprocess/PTY Support (`pkg/gui/sub_process.go`)
- Running interactive commands (attach, exec, logs -f)
- Terminal passthrough with proper PTY handling
- **Status:** Not implemented - required for attach/exec

#### 7. Bulk Operations (`pkg/gui/bulk_actions.go`)
- Select multiple items with space
- Perform action on all selected
- "Prune" operations (remove all stopped, dangling, etc.)
- **Status:** Not implemented

#### 8. Container Top (`pkg/commands/container.go`)
- Show running processes in container
- lazydocker: `docker top`
- Apple CLI: Need to investigate if supported
- **Status:** Not implemented

### Lower Priority - Polish

#### 9. Information Panel (`pkg/gui/information_panel.go`)
- Bottom bar showing version, help hints
- Context-sensitive help
- **Status:** View exists but not populated

#### 10. Menu System (`pkg/gui/menu_panel.go`)
- Popup menus for complex choices
- E.g., signal selection for kill, restart policies
- **Status:** Types exist but UI not implemented

#### 11. Options View (`pkg/gui/options.go`)
- Show keybindings for current context
- Dynamic based on selected item type
- **Status:** View exists but not populated

#### 12. Config File Support (`pkg/config/`)
- `~/.config/lazycontainer/config.yml`
- Custom themes, keybindings, commands
- **Status:** Basic config exists but not user-configurable

#### 13. Mouse Mode Toggle
- Option to disable mouse (for tmux compatibility)
- **Status:** Not implemented

#### 14. Scrolling in Main Panel
- PgUp/PgDown, mouse wheel in main content
- **Status:** Basic scrolling exists but may need refinement

#### 15. Copy to Clipboard
- Copy container ID, logs, etc.
- **Status:** Not implemented

---

## Apple Container CLI Commands Reference

```bash
# List containers
container list --format json
container list --all --format json

# Container operations
container start <id>
container stop <id>
container delete <id>
container logs <id> -n <lines>        # Note: -n not --tail
container inspect --format json <id>
container kill <id> --signal <signal>
container exec <id> -- <command>      # Interactive execution

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

## Key Architectural Differences from Lazydocker

1. **No Docker Compose** - Apple containers don't have compose, so no Services/Project panels
2. **No SSH tunneling** - Local containers only
3. **Simpler stats** - Container stats may have different format
4. **No container attach** - May need to implement differently via `container exec`
5. **Different timestamp format** - Apple uses CFAbsoluteTime (seconds since Jan 1, 2001)
6. **Different CLI flags** - `-n` instead of `--tail` for logs

## Directory Structure

```
/Users/warl0ck/Code/lazycontainer/
├── main.go
├── go.mod
├── go.sum
├── IMPLEMENTATION_PLAN.md
├── scripts/
│   └── test_setup.sh           # Test infrastructure setup
├── pkg/
│   ├── app/
│   │   └── app.go
│   ├── commands/
│   │   ├── os.go
│   │   ├── container.go
│   │   ├── container_types.go
│   │   └── container_test.go
│   ├── config/
│   │   └── app_config.go
│   ├── gui/
│   │   ├── gui.go
│   │   ├── views.go
│   │   ├── keybindings.go
│   │   ├── panels/
│   │   │   ├── filtered_list.go
│   │   │   ├── list_panel.go
│   │   │   ├── context_state.go
│   │   │   └── side_list_panel.go
│   │   ├── presentation/
│   │   │   ├── containers.go
│   │   │   ├── images.go
│   │   │   ├── volumes.go
│   │   │   └── networks.go
│   │   └── types/
│   │       └── types.go
│   ├── i18n/
│   │   ├── i18n.go
│   │   └── english.go
│   ├── log/
│   │   └── log.go
│   ├── tasks/
│   │   └── tasks.go
│   └── utils/
│       └── utils.go
```

## Implementation Roadmap

### Phase 1: Core Polish (Current)
- [ ] Fix mouse selection highlight behavior
- [x] Add confirmation dialogs for destructive actions
- [ ] Implement status bar with error/info messages
- [ ] Connect filtering UI to panel infrastructure

### Phase 2: Interactive Features
- [ ] Implement container exec/attach with PTY support
- [ ] Add subprocess handling for interactive commands
- [ ] Implement streaming logs (`logs -f`)
- [ ] Add bulk selection and operations

### Phase 3: User Experience
- [ ] Config file support (`~/.config/lazycontainer/`)
- [ ] Custom commands per resource type
- [ ] Menu system for complex choices
- [ ] Options panel with keybinding help

### Phase 4: Polish
- [ ] Information panel with context help
- [ ] Copy to clipboard functionality
- [ ] Mouse mode toggle
- [ ] Theme customization

## Testing Commands

```bash
# Build and run
cd /Users/warl0ck/Code/lazycontainer
go build -o lazycontainer .
./lazycontainer

# Run tests
go test ./pkg/commands/...

# Set up test containers
./scripts/test_setup.sh
```

## Resume Instructions

1. Read this file: `/Users/warl0ck/Code/lazycontainer/IMPLEMENTATION_PLAN.md`
2. Read the lazydocker source for reference: `/Users/warl0ck/Code/lazydocker/`
3. Pick a feature from the Implementation Roadmap
4. Copy patterns from lazydocker, update imports from `github.com/jesseduffield/lazydocker` to `github.com/warl0ck/lazycontainer`
5. Replace Docker-specific code with Apple container CLI calls
