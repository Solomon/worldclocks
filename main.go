package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed defaults/worldclocks
var embeddedFiles embed.FS

type tickMsg time.Time

type timezone struct {
	tz      string
	country string
}

type model struct {
	zones      []timezone
	width      int
	height     int
	configPath string
}

// ensureConfig copies the embedded file if ~/.config/worldclocks doesn’t exist
func ensureConfig(path string) error {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = strings.Replace(path, "~", home, 1)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		data, err := fs.ReadFile(embeddedFiles, "defaults/worldclocks")
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("Created default config at %s\n", path)
	}
	return nil
}

// readConfig reads only timezones marked "Yes"
func readConfig(filename string) ([]timezone, error) {
	if strings.HasPrefix(filename, "~") {
		home, _ := os.UserHomeDir()
		filename = strings.Replace(filename, "~", home, 1)
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var zones []timezone
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(parts[2]), "yes") {
			zones = append(zones, timezone{
				tz:      strings.TrimSpace(parts[0]),
				country: strings.TrimSpace(parts[1]),
			})
		}
	}
	return zones, scanner.Err()
}

// openInEditor launches the system editor
func openInEditor(path string) error {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = strings.Replace(path, "~", home, 1)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}


func (m model) View() string {
	var boxes []string

	title := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("205")).
		Render("World Clocks")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 3).
		Margin(1, 1).
		Align(lipgloss.Center).
		Width(40)

	timeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		Align(lipgloss.Center).
		Render

	for _, z := range m.zones {
		loc, err := time.LoadLocation(z.tz)
		if err != nil {
			boxes = append(boxes, boxStyle.Render(fmt.Sprintf("%s\n[Invalid Zone]", z.tz)))
			continue
		}

		// Extract readable city name
		city := strings.ReplaceAll(z.tz[strings.Index(z.tz, "/")+1:], "_", " ")

		now := time.Now().In(loc)
		header := fmt.Sprintf("[%s]  %s", z.country, city)
		content := fmt.Sprintf("%s\n%s", header, now.Format("15:04:05  Mon, Jan 2"))
		boxes = append(boxes, boxStyle.Render(timeStyle(content)))
	}


	// Determine how many boxes fit per row based on terminal width
	termWidth := m.width
	if termWidth == 0 {
		termWidth = 120 // fallback if not yet resized
	}
	boxWidth := boxStyle.GetHorizontalFrameSize() + 40
	perRow := termWidth / boxWidth
	if perRow < 1 {
		perRow = 1
	}

	// Build rows
	var rows []string
	for i := 0; i < len(boxes); i += perRow {
		end := i + perRow
		if end > len(boxes) {
			end = len(boxes)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, boxes[i:end]...)
		rows = append(rows, row)
	}

	layout := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return title + "\n\n" + layout + "\n\nPress q to quit.\n"
}

func main() {
	configPath := "~/.config/worldclocks"

	if err := ensureConfig(configPath); err != nil {
		fmt.Println("Error creating config:", err)
		os.Exit(1)
	}

	zones, err := readConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	m := model{zones: zones, configPath: configPath}
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
