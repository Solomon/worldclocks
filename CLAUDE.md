# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`worldclocks` is a terminal-based world clock application built with Go using the Bubble Tea TUI framework. It displays multiple timezone clocks in a responsive grid layout.

## Build and Run Commands

### Using Makefile (Recommended)

```bash
# Build the application
make build

# Run the application
make run

# Install to /usr/local/bin (requires sudo)
sudo make install

# Uninstall from system
sudo make uninstall

# Clean build artifacts
make clean

# Update dependencies
make deps

# Test GoReleaser configuration locally
make release-test

# Show all available commands
make help
```

### Using Go directly

```bash
# Build the application
go build -o worldclocks .

# Run the application
go run main.go

# Install dependencies
go mod download

# Update dependencies
go mod tidy
```

## Architecture

### Single-File Application Structure

The entire application is contained in `main.go` with the following key components:

**Bubble Tea Model:**
- `model` struct: Holds application state (timezone list, terminal dimensions, config path)
- `Init()`: Sets up the 1-second tick for clock updates
- `Update()`: Handles tick messages, window resize events, and keyboard input (q/Ctrl+C to quit)
- `View()`: Renders the clock grid using Lipgloss styling

**Configuration System:**
- Default config embedded via `//go:embed defaults/worldclocks`
- Config file location: `~/.config/worldclocks`
- `ensureConfig()`: Copies embedded default on first run
- `readConfig()`: Parses CSV-format config (timezone,country_code,include)
  - Only processes timezones with "Yes" in the include column
  - Skips empty lines and comments (#)

**Config File Format:**
```
# timezone,country_code,include(Yes/No)
America/New_York,US,Yes
Europe/London,GB,No
```

### Key Implementation Details

**Timezone Data Structure:**
The application uses a `timezone` struct with `tz` (timezone name) and `country` (country code) fields. The `readConfig()` function parses the CSV config file and returns `[]timezone`. Ensure any config file modifications maintain the three-column CSV format: `timezone,country_code,include`.

**Responsive Layout:**
The grid layout calculates boxes per row based on terminal width (`tea.WindowSizeMsg`). Each clock box has a fixed width of 40 chars plus frame size. The layout recalculates on terminal resize.

**Embedded Resources:**
The `defaults/worldclocks` file is embedded at compile time and extracted to `~/.config/worldclocks` if missing. Changes to the default config require rebuilding.

## User Interface

**Keyboard Controls:**
- `c`: Opens the config file in the default editor (`$EDITOR`, defaults to nano)
- `q` or `Ctrl+C`: Quit the application

After editing the config with `c`, the application quits so you can restart it to see changes.

## Distribution and Release Process

This project uses **GoReleaser** for automated multi-platform releases.

### Release Workflow

1. **Testing locally:**
   ```bash
   make release-test
   # Creates packages in dist/ without publishing
   ```

2. **Creating a release:**
   ```bash
   # Commit all changes
   git add . && git commit -m "Release v1.0.0"

   # Create and push a version tag
   git tag v1.0.0
   git push origin main
   git push origin v1.0.0
   ```

3. **Automated builds:**
   - GitHub Actions (`.github/workflows/release.yml`) automatically triggers on tag push
   - Creates binaries for: Linux (amd64, arm64), macOS (amd64, arm64), Windows
   - Generates packages: `.deb`, `.rpm`, `.apk`
   - Creates GitHub release with binaries and checksums
   - Optional: Updates Homebrew tap if configured

### GoReleaser Configuration

**File:** `.goreleaser.yaml`

Key sections:
- `builds`: Cross-platform binary compilation
- `archives`: Creates tar.gz/zip archives
- `nfpms`: Generates Linux packages (.deb, .rpm, .apk)
- `brews`: Homebrew formula (requires separate tap repo)
- `aurs`: Arch Linux AUR package (optional)

**Important:** Before first release, update these fields in `.goreleaser.yaml`:
- `homepage`: Your GitHub repo URL
- `maintainer`: Your name and email
- `repository.owner`: Your GitHub username
- `repository.name`: Should be `homebrew-tap` if using Homebrew

### Embedded Files in Releases

The `defaults/worldclocks` file is embedded via `//go:embed` directive and included in all binaries. Changes to default config require a new release build.

## Development Notes

- No tests currently exist in the codebase
- The application uses `time.LoadLocation()` for timezone handling - ensure valid IANA timezone names
- All releases are deterministic and reproducible thanks to `-trimpath` flag
