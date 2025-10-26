# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`worldclocks` is a terminal-based world clock application built with Go using the Bubble Tea TUI framework. It displays multiple timezone clocks in a responsive grid layout.

## Build and Run Commands

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
The `readConfig()` function returns `[]string` containing timezone names, but the `View()` method expects a struct with `.tz` and `.country` fields. When modifying timezone handling, ensure consistency between the model's `zones` field type and how it's accessed in `View()`.

**Responsive Layout:**
The grid layout calculates boxes per row based on terminal width (`tea.WindowSizeMsg`). Each clock box has a fixed width of 40 chars plus frame size. The layout recalculates on terminal resize.

**Embedded Resources:**
The `defaults/worldclocks` file is embedded at compile time and extracted to `~/.config/worldclocks` if missing. Changes to the default config require rebuilding.

## Development Notes

- No tests currently exist in the codebase
- The application uses `time.LoadLocation()` for timezone handling - ensure valid IANA timezone names
- Editor integration: Pressing 'e' would typically open the config in `$EDITOR` (defaults to nano), but this feature appears to be planned rather than implemented
