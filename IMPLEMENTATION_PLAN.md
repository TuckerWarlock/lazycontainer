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
   - `pkg/commands/container_types.go` - Container, Image, Volume, Network structs
   - `pkg/gui/gui.go` - Basic TUI with container panel
   - `pkg/i18n/english.go` - Translations (just created)
   - `pkg/i18n/i18n.go` - Translation loader (just created)
   - `pkg/tasks/tasks.go` - Task management (just created)
   - `pkg/gui/types/types.go` - MenuItem type (just created)

3. **Keybindings working**: j/k navigation, Enter (start), s (stop), d (delete), l (logs), a (toggle all), r (refresh), q (quit)

### In Progress - Copy Lazydocker Structure

The goal is to replicate lazydocker's full panel system for: Containers, Images, Volumes, Networks, and Config view.

## Files to Copy/Adapt from Lazydocker

### 1. Panel Infrastructure (`pkg/gui/panels/`)

**Source:** `/Users/warl0ck/Code/lazydocker/pkg/gui/panels/`

Files to create in `/Users/warl0ck/Code/lazycontainer/pkg/gui/panels/`:

```go
// filtered_list.go - Generic filtered list with sorting
// Already in lazydocker - copy and update imports

// list_panel.go - Base list panel with selection
// Needs: lcUtils.Clamp - either copy from lazycore or implement simple version

// context_state.go - Tab management for main panel
// Copy and update imports

// side_list_panel.go - Main reusable panel component
// Copy and update imports
```

**Key change needed:** Replace `github.com/jesseduffield/lazycore/pkg/utils` with local implementation:

```go
// Add to pkg/utils/utils.go:
func Clamp(value, min, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
```

### 2. Utility Functions (`pkg/utils/`)

**Source:** `/Users/warl0ck/Code/lazydocker/pkg/utils/utils.go`

Copy the following functions:
- `SplitLines`, `WithPadding`, `ColoredString`, `Decolorise`
- `RenderTable`, `getPadWidths`, `getPaddedDisplayStrings`
- `FormatBinaryBytes`, `FormatDecimalBytes`
- `FormatMap`, `FormatMapItem`
- `ApplyTemplate`, `GetGocuiAttribute`, `GetColorAttribute`
- `ColoredYamlString`, `MarshalIntoYaml`
- `Clamp` (add this)

**Dependencies needed in go.mod:**
```
github.com/go-errors/errors v1.0.2
github.com/goccy/go-yaml v1.11.0
github.com/mattn/go-runewidth v0.0.15
```

### 3. Presentation Layer (`pkg/gui/presentation/`)

Create display formatters for each resource type:

```go
// containers.go
func GetContainerDisplayStrings(container *commands.Container) []string {
    return []string{
        container.GetStatusSymbol(),
        container.GetName(),
        container.GetImage(),
        container.GetStatus(),
        container.GetPorts(),
    }
}

// images.go
func GetImageDisplayStrings(image *commands.Image) []string {
    return []string{
        image.Reference,
        utils.SafeTruncate(image.Digest, 12),
        utils.FormatBinaryBytes(int(image.Size)),
    }
}

// volumes.go
func GetVolumeDisplayStrings(volume *commands.Volume) []string {
    return []string{
        volume.Name,
        volume.Format,
        utils.FormatBinaryBytes(int(volume.SizeInBytes)),
    }
}

// networks.go
func GetNetworkDisplayStrings(network *commands.Network) []string {
    return []string{
        network.ID,
        network.Config.Type,
        network.Status.State,
    }
}
```

### 4. Views Structure (`pkg/gui/views.go`)

```go
type Views struct {
    Containers *gocui.View
    Images     *gocui.View
    Volumes    *gocui.View
    Networks   *gocui.View
    Main       *gocui.View
    Options    *gocui.View
    Information *gocui.View
    Confirmation *gocui.View
    Menu       *gocui.View
}
```

### 5. Panel Definitions

Each panel needs:
- View setup
- GetTableCells function
- Main tabs (Logs, Config, Stats, etc.)
- Filter/Sort functions

**containers_panel.go:**
```go
func (gui *Gui) getContainersPanel() *panels.SideListPanel[*commands.Container] {
    return &panels.SideListPanel[*commands.Container]{
        ContextState: &panels.ContextState[*commands.Container]{
            GetMainTabs: func() []panels.MainTab[*commands.Container] {
                return []panels.MainTab[*commands.Container]{
                    {Key: "logs", Title: gui.Tr.LogsTitle, Render: gui.renderContainerLogs},
                    {Key: "config", Title: gui.Tr.ConfigTitle, Render: gui.renderContainerConfig},
                    {Key: "stats", Title: gui.Tr.StatsTitle, Render: gui.renderContainerStats},
                }
            },
            GetItemContextCacheKey: func(c *commands.Container) string {
                return "containers-" + c.ID + "-" + c.Status
            },
        },
        ListPanel: panels.ListPanel[*commands.Container]{
            List: panels.NewFilteredList[*commands.Container](),
            View: gui.Views.Containers,
        },
        NoItemsMessage: gui.Tr.NoContainers,
        Gui:            gui,
        GetTableCells:  presentation.GetContainerDisplayStrings,
    }
}
```

### 6. Main GUI Restructure

The main `gui.go` needs to be restructured to:

1. Hold `Views` and `Panels` structs
2. Implement the `IGui` interface for panels
3. Set up all panels on startup
4. Handle panel switching with Tab/arrow keys
5. Render to main panel based on selected item and tab

**Key IGui interface methods:**
```go
type IGui interface {
    HandleClick(v *gocui.View, itemCount int, selectedLine *int, handleSelect func() error) error
    NewSimpleRenderStringTask(getContent func() string) tasks.TaskFunc
    FocusY(selectedLine int, itemCount int, view *gocui.View)
    ShouldRefresh(contextKey string) bool
    GetMainView() *gocui.View
    IsCurrentView(*gocui.View) bool
    FilterString(view *gocui.View) string
    IgnoreStrings() []string
    Update(func() error)
    QueueTask(f func(ctx context.Context)) error
}
```

## Apple Container CLI Commands Reference

```bash
# List containers
container list --format json
container list --all --format json

# Container operations
container start <id>
container stop <id>
container delete <id>
container logs <id> --tail <n>
container inspect --format json <id>
container kill <id> --signal <signal>

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

## Go Dependencies to Add

```go
require (
    github.com/go-errors/errors v1.0.2
    github.com/goccy/go-yaml v1.11.0
    github.com/mattn/go-runewidth v0.0.15
)
```

## Implementation Order

1. **pkg/utils/utils.go** - Add all utility functions including Clamp
2. **pkg/gui/panels/** - Copy all 4 files, update imports
3. **pkg/gui/presentation/** - Create 4 display files
4. **pkg/gui/views.go** - Create views struct
5. **pkg/gui/*_panel.go** - Create 4 panel definition files
6. **pkg/gui/gui.go** - Restructure to use panels system
7. **pkg/gui/keybindings.go** - Separate keybindings file
8. **pkg/gui/layout.go** - Layout calculations
9. **pkg/gui/view_helpers.go** - View manipulation utilities

## Testing Commands

```bash
cd /Users/warl0ck/Code/lazycontainer
go build -o lazycontainer .
./lazycontainer
```

## Key Architectural Differences from Lazydocker

1. **No Docker Compose** - Apple containers don't have compose, so no Services/Project panels
2. **No SSH tunneling** - Local containers only
3. **Simpler stats** - Container stats may have different format
4. **No container attach** - May need to implement differently

## Files Already Created (Build Passing)

```
/Users/warl0ck/Code/lazycontainer/
├── main.go
├── go.mod
├── go.sum
├── IMPLEMENTATION_PLAN.md (this file)
├── pkg/
│   ├── app/
│   │   └── app.go
│   ├── commands/
│   │   ├── os.go
│   │   ├── container.go
│   │   └── container_types.go
│   ├── config/
│   │   └── app_config.go
│   ├── gui/
│   │   ├── gui.go (basic TUI working)
│   │   ├── panels/
│   │   │   ├── filtered_list.go ✓
│   │   │   ├── list_panel.go ✓
│   │   │   ├── context_state.go ✓
│   │   │   └── side_list_panel.go ✓
│   │   ├── presentation/ (empty, needs files)
│   │   └── types/
│   │       └── types.go ✓
│   ├── i18n/
│   │   ├── i18n.go ✓
│   │   └── english.go ✓
│   ├── log/
│   │   └── log.go
│   ├── tasks/
│   │   └── tasks.go ✓
│   └── utils/
│       └── utils.go ✓ (with Clamp, RenderTable, etc.)
```

## Next Steps (for new chat)

1. Create `pkg/gui/presentation/` files for display formatting
2. Create `pkg/gui/views.go` with Views struct
3. Restructure `pkg/gui/gui.go` to use the panels system
4. Create panel definition files (*_panel.go)

## Resume Instructions

In a new chat:
1. Read this file: `/Users/warl0ck/Code/lazycontainer/IMPLEMENTATION_PLAN.md`
2. Read the lazydocker source for reference: `/Users/warl0ck/Code/lazydocker/`
3. Continue from "Implementation Order" step 1
4. Copy files from lazydocker, update imports from `github.com/jesseduffield/lazydocker` to `github.com/warl0ck/lazycontainer`
5. Replace Docker-specific code with Apple container CLI calls
