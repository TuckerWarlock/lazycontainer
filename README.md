# lazycontainer

[![Test and Build](https://github.com/TuckerWarlock/lazycontainer/actions/workflows/test-build.yml/badge.svg)](https://github.com/TuckerWarlock/lazycontainer/actions/workflows/test-build.yml)

A terminal UI for managing Apple containers on macOS 26+, inspired by [lazydocker](https://github.com/jesseduffield/lazydocker).

![lazycontainer screenshot placeholder]

## Requirements

- macOS 26+ (Tahoe)
- Apple Containerization framework: `brew install --cask container`
- Container service running: `container system start`

## Installation

```bash
git clone https://github.com/warl0ck/lazycontainer.git
cd lazycontainer
go build -o lazycontainer .
./lazycontainer
```

Or install directly:

```bash
go install github.com/warl0ck/lazycontainer@latest
```

## Usage

```bash
# Start the container service first
container system start

# Launch the TUI
./lazycontainer

# Debug mode (logs to ~/.config/lazycontainer/development.log)
./lazycontainer -d
```

## Keybindings

### Navigation

| Key | Action |
|-----|--------|
| `1` / `2` / `3` / `4` | Switch to containers / images / volumes / networks panel |
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `[` / `]` | Previous / next tab (logs, config, stats) |
| `Tab` | Cycle panels |
| `/` | Filter list |
| `q` | Quit |

### Container Actions

| Key | Action |
|-----|--------|
| `Enter` | Start container |
| `s` | Stop container (with confirmation) |
| `d` | Delete container (with confirmation) |
| `Space` | Toggle multi-select |
| `f` | Follow logs in real-time |
| `e` | Exec interactive shell |
| `v` | Stream live stats |

## Features

- **Four panels** — containers, images, volumes, networks
- **Real-time filtering** — `/` to filter any panel by name
- **Follow logs** — stream container logs live (`f`)
- **Interactive shell** — exec into running containers (`e`)
- **Streaming stats** — CPU, memory, network I/O (`v`)
- **Bulk operations** — Space to multi-select, then `s` or `d`
- **Confirmation dialogs** — destructive actions require confirmation
- **Color-coded status** — green = running, red = stopped
- **Mouse support** — click to select, scroll to navigate
- **Auto-refresh** — container list updates every 2 seconds

## Documentation

- [docs/IMPLEMENTATION.md](docs/IMPLEMENTATION.md) — feature status, lazydocker parity table, architecture
- [docs/TESTING.md](docs/TESTING.md) — test suite, manual test scenarios, CI setup

## Development

```bash
# Run tests
go test ./...

# Set up test containers for manual testing
./scripts/test_setup.sh
```

See [docs/TESTING.md](docs/TESTING.md) for the full testing guide.

## License

MIT
