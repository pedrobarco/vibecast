# Vibecast

_A fast, keyboard-driven TUI for managing and playing IPTV playlists._

[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/github/license/pedrobarco/vibecast)](./LICENSE)

---

## Table of Contents

- [Overview](#overview)
- [Demo](#demo)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Supported Players](#supported-players)
- [Supported Playlist Formats](#supported-playlist-formats)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

**Vibecast** is a terminal user interface (TUI) for managing, searching, and playing IPTV playlists.
It is player-agnostic: you can use it with any player that supports opening URLs from the command line (e.g., VLC, mpv, etc.).

**Features:**

- Add, remove, and edit IPTV playlists (local files or URLs)
- Search and filter channels with fuzzy matching
- Bookmark (favorite) channels per playlist
- Play channels in your preferred player (VLC, mpv, etc.)
- Keyboard-driven, Vim-like navigation
- Portable YAML configuration

---

## Demo

<!--
![Vibecast Demo](https://user-images.githubusercontent.com/yourusername/vibecast-demo.gif)
-->

---

## Installation

### Using Go

```bash
go install github.com/pedrobarco/vibecast/cmd/vibecast@latest
```

### Prebuilt Binaries

Prebuilt binaries may be provided on the [Releases page](https://github.com/pedrobarco/vibecast/releases).

---

## Usage

### Run Vibecast

```bash
vibecast
```

### Keybindings

| Key      | Action                          |
| -------- | ------------------------------- |
| `j`/`k`  | Move down/up                    |
| `enter`  | Select / Play channel           |
| `/`      | Search/filter channels          |
| `b`      | Toggle bookmarks view           |
| `m`      | Mark/unmark channel as bookmark |
| `esc`    | Go back / Cancel / Clear filter |
| `ctrl+c` | Quit                            |

### Example Workflow

1. Start Vibecast.
2. Add a playlist (local file or URL).
3. Select a playlist to view channels.
4. Use `/` to search, `b` to show only bookmarks, `m` to mark favorites.
5. Press `enter` to play a channel in your configured player.

---

## Configuration

- The config file is stored at:
  `$HOME/.config/vibecast/config.yaml`

- Example config:

```yaml
playlists:
  - name: My IPTV
    path: /path/to/playlist.m3u
  - name: Remote List
    path: https://example.com/iptv.m3u
favourites:
  "My IPTV":
    - Channel 1
    - Channel 2
  "Remote List":
    - News Channel
```

- You can copy this file to another machine and Vibecast will work out-of-the-box.

---

## Supported Players

Vibecast is player-agnostic. By default, it will try to open channels using the system's default handler for URLs:

- **macOS:** `open -a VLC <url>` (or your default player)
- **Linux:** `xdg-open <url>`
- **Windows:** `start <url>`

You can configure your system to use your preferred player for streaming URLs.

---

## Supported Playlist Formats

- **M3U** (local files or remote URLs)

(Other formats may be supported in the future.)

---

## Contributing

Contributions are welcome! Please open issues or pull requests.
See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## License

[MIT](./LICENSE) © 2025 Pedro Barco
