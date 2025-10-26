# World Clocks

A beautiful terminal-based world clock application built with Go and Bubble Tea. Display multiple timezone clocks in a responsive grid layout with a modern TUI interface.

![World Clocks Demo](https://via.placeholder.com/800x400?text=Add+Screenshot+Here)

## Features

- 🌍 Display multiple timezones simultaneously
- 🎨 Beautiful TUI with Lipgloss styling
- 📱 Responsive grid layout that adapts to terminal size
- ⚙️ Easy configuration via text file
- 🚀 Fast and lightweight
- ⌨️ Simple keyboard controls

## Installation

### macOS / Linux (Homebrew)

```bash
brew install yourusername/tap/worldclocks
```

### Arch Linux (AUR)

```bash
yay -S worldclocks-bin
# or
paru -S worldclocks-bin
```

### Debian / Ubuntu

Download the `.deb` package from the [releases page](https://github.com/yourusername/worldclocks/releases) and install:

```bash
sudo dpkg -i worldclocks_*_linux_amd64.deb
```

### Fedora / RHEL / CentOS

Download the `.rpm` package from the [releases page](https://github.com/yourusername/worldclocks/releases) and install:

```bash
sudo rpm -i worldclocks_*_linux_amd64.rpm
```

### Alpine Linux

Download the `.apk` package from the [releases page](https://github.com/yourusername/worldclocks/releases) and install:

```bash
sudo apk add --allow-untrusted worldclocks_*_linux_amd64.apk
```

### Manual Installation

Download the appropriate binary for your platform from the [releases page](https://github.com/yourusername/worldclocks/releases), extract it, and move to your PATH:

```bash
# Example for Linux
tar -xzf worldclocks_*_linux_amd64.tar.gz
sudo mv worldclocks /usr/local/bin/
```

### Build from Source

Requires Go 1.23 or later:

```bash
git clone https://github.com/yourusername/worldclocks.git
cd worldclocks

# Build as your user (Go must be in PATH)
make build

# Install as root (only copies the binary)
sudo make install
```

Or using Go directly:

```bash
go install github.com/yourusername/worldclocks@latest
```

**Note:** When building from source with Make, you must run `make build` and `sudo make install` as separate commands. This is because `sudo` doesn't have access to your user's Go installation.

## Usage

Simply run:

```bash
worldclocks
```

### Keyboard Controls

- **c** - Open configuration file in your editor
- **q** or **Ctrl+C** - Quit

## Configuration

On first run, a configuration file is created at `~/.config/worldclocks`.

The config file format is CSV with three columns:

```csv
# timezone,country_code,include(Yes/No)
America/New_York,US,Yes
Europe/London,GB,Yes
Asia/Tokyo,JP,Yes
Asia/Dubai,AE,No
```

- **timezone**: IANA timezone name (e.g., `America/New_York`)
- **country_code**: Two-letter country code for display
- **include**: Set to `Yes` to display, `No` to hide

### Editing Configuration

Press **c** while running the app to open the config file in your default editor (set via `$EDITOR` environment variable, defaults to `nano`).

After saving changes, restart the app to see updates.

## Development

### Building

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make install        # Install to /usr/local/bin
```

### Testing

```bash
make test           # Run tests
```

### Creating a Release

```bash
# Test release locally
make release-test

# Create actual release (requires git tag)
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions will automatically build and publish
```

## How It Works

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) for beautiful terminal UI
- Uses Go's `time` package for timezone handling
- Embedded default configuration for easy first-run experience

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Charm](https://charm.sh) for the amazing TUI libraries
- All contributors and users of this project
