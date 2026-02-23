# lazycontainer

[![Test and Build](https://github.com/TuckerWarlock/lazycontainer/actions/workflows/test-build.yml/badge.svg)](https://github.com/TuckerWarlock/lazycontainer/actions/workflows/test-build.yml)

A terminal UI for managing Apple containers, inspired by [lazydocker](https://github.com/jesseduffield/lazydocker).


## Requirements

- macOS 26+ (Tahoe)
- Apple Containerization framework (`brew install --cask container`)
- Go 1.22+ (for building from source)

## Installation

### From Source

```bash
go install github.com/warl0ck/lazycontainer@latest
```

Or clone and build:

```bash
git clone https://github.com/warl0ck/lazycontainer.git
cd lazycontainer
go build -o lazycontainer .
```

## Usage

Make sure the container service is running:

```bash
container system start
```

Then launch the TUI:

```bash
./lazycontainer
```

### Keybindings

| Key | Action |
|-----|--------|
| `↑`/`k` | Move up |
| `↓`/`j` | Move down |
| `Enter` | Start container |
| `s` | Stop container |
| `d` | Delete container |
| `l` | View logs |
| `a` | Toggle show all containers |
| `r` | Refresh |
| `q` | Quit |

## Features

- View all containers with status indicators
- Start, stop, and delete containers
- View container details and logs
- Auto-refresh every 2 seconds
- Color-coded status (green = running, red = stopped)

## Development

```bash
# Run with debug logging
./lazycontainer -d

# Logs are written to ~/.config/lazycontainer/development.log
```

## License

MIT
