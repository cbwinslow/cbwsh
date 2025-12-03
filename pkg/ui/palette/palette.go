// Package palette provides a command palette (Ctrl+P) for cbwsh.
package palette

import (
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Command represents a command in the palette.
type Command struct {
	// ID is the unique identifier.
	ID string
	// Name is the display name.
	Name string
	// Description describes the command.
	Description string
	// Shortcut is the keyboard shortcut.
	Shortcut string
	// Category is the command category.
	Category string
	// Action is the function to execute.
	Action func() tea.Cmd
	// Keywords are search keywords.
	Keywords []string
}

// Palette is a command palette component.
type Palette struct {
	mu       sync.RWMutex
	commands []Command
	filtered []Command
	input    textinput.Model
	selected int
	visible  bool
	width    int
	height   int
	maxItems int
	styles   Styles
}

// Styles defines the palette styles.
type Styles struct {
	Background   lipgloss.Style
	Border       lipgloss.Style
	Input        lipgloss.Style
	Item         lipgloss.Style
	SelectedItem lipgloss.Style
	Description  lipgloss.Style
	Shortcut     lipgloss.Style
	Category     lipgloss.Style
	NoResults    lipgloss.Style
}

// DefaultStyles returns default palette styles.
func DefaultStyles() Styles {
	return Styles{
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Padding(1, 2),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")),
		Input: lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")),
		Item: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),
		SelectedItem: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")),
		Shortcut: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true),
		Category: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Bold(true),
		NoResults: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true).
			Padding(0, 1),
	}
}

// KeyMap defines key bindings for the palette.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Escape key.Binding
	Tab    key.Binding
}

// DefaultKeyMap returns default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "ctrl+p"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "ctrl+n"),
			key.WithHelp("↓", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc", "escape"),
			key.WithHelp("esc", "close"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
	}
}

// New creates a new command palette.
func New() *Palette {
	ti := textinput.New()
	ti.Placeholder = "Type a command..."
	ti.Prompt = "> "
	ti.CharLimit = 100
	ti.Width = 50

	return &Palette{
		commands: make([]Command, 0),
		filtered: make([]Command, 0),
		input:    ti,
		maxItems: 10,
		styles:   DefaultStyles(),
	}
}

// AddCommand adds a command to the palette.
func (p *Palette) AddCommand(cmd Command) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commands = append(p.commands, cmd)
	p.filter()
}

// AddCommands adds multiple commands to the palette.
func (p *Palette) AddCommands(cmds []Command) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commands = append(p.commands, cmds...)
	p.filter()
}

// RemoveCommand removes a command by ID.
func (p *Palette) RemoveCommand(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, cmd := range p.commands {
		if cmd.ID == id {
			p.commands = append(p.commands[:i], p.commands[i+1:]...)
			break
		}
	}
	p.filter()
}

// ClearCommands removes all commands.
func (p *Palette) ClearCommands() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commands = make([]Command, 0)
	p.filtered = make([]Command, 0)
}

// Open shows the palette.
func (p *Palette) Open() tea.Cmd {
	p.visible = true
	p.selected = 0
	p.input.Reset()
	p.input.Focus()
	p.filter()
	return textinput.Blink
}

// Close hides the palette.
func (p *Palette) Close() {
	p.visible = false
	p.input.Blur()
}

// IsVisible returns whether the palette is visible.
func (p *Palette) IsVisible() bool {
	return p.visible
}

// SetSize sets the palette size.
func (p *Palette) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.input.Width = width - 10
}

// SetStyles sets the palette styles.
func (p *Palette) SetStyles(styles Styles) {
	p.styles = styles
}

// Update handles input messages.
func (p *Palette) Update(msg tea.Msg) (bool, tea.Cmd) {
	if !p.visible {
		return false, nil
	}

	keys := DefaultKeyMap()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			p.Close()
			return true, nil

		case key.Matches(msg, keys.Up):
			p.selected--
			if p.selected < 0 {
				p.selected = len(p.filtered) - 1
				if p.selected < 0 {
					p.selected = 0
				}
			}
			return true, nil

		case key.Matches(msg, keys.Down), key.Matches(msg, keys.Tab):
			p.selected++
			if p.selected >= len(p.filtered) {
				p.selected = 0
			}
			return true, nil

		case key.Matches(msg, keys.Enter):
			p.mu.RLock()
			if p.selected >= 0 && p.selected < len(p.filtered) {
				cmd := p.filtered[p.selected]
				p.mu.RUnlock()
				p.Close()
				if cmd.Action != nil {
					return true, cmd.Action()
				}
				return true, nil
			}
			p.mu.RUnlock()
			return true, nil
		}
	}

	// Update text input
	var cmd tea.Cmd
	prevValue := p.input.Value()
	p.input, cmd = p.input.Update(msg)

	// Filter on input change
	if p.input.Value() != prevValue {
		p.filter()
		p.selected = 0
	}

	return true, cmd
}

// View renders the palette.
func (p *Palette) View() string {
	if !p.visible {
		return ""
	}

	var sb strings.Builder

	// Input
	sb.WriteString(p.styles.Input.Render(p.input.View()))
	sb.WriteString("\n")

	p.mu.RLock()
	filtered := p.filtered
	p.mu.RUnlock()

	// Results
	if len(filtered) == 0 {
		sb.WriteString(p.styles.NoResults.Render("No matching commands"))
	} else {
		for i, cmd := range filtered {
			if i >= p.maxItems {
				break
			}

			style := p.styles.Item
			if i == p.selected {
				style = p.styles.SelectedItem
			}

			line := cmd.Name
			if cmd.Description != "" {
				line += " " + p.styles.Description.Render(cmd.Description)
			}
			if cmd.Shortcut != "" {
				line += " " + p.styles.Shortcut.Render(cmd.Shortcut)
			}

			sb.WriteString(style.Render(line))
			sb.WriteString("\n")
		}
	}

	content := sb.String()
	content = p.styles.Border.Render(p.styles.Background.Render(content))

	return content
}

func (p *Palette) filter() {
	query := strings.ToLower(p.input.Value())

	if query == "" {
		p.filtered = make([]Command, len(p.commands))
		copy(p.filtered, p.commands)
		return
	}

	p.filtered = make([]Command, 0)
	for _, cmd := range p.commands {
		if p.matches(cmd, query) {
			p.filtered = append(p.filtered, cmd)
		}
	}
}

func (p *Palette) matches(cmd Command, query string) bool {
	// Check name
	if strings.Contains(strings.ToLower(cmd.Name), query) {
		return true
	}

	// Check description
	if strings.Contains(strings.ToLower(cmd.Description), query) {
		return true
	}

	// Check category
	if strings.Contains(strings.ToLower(cmd.Category), query) {
		return true
	}

	// Check keywords
	for _, kw := range cmd.Keywords {
		if strings.Contains(strings.ToLower(kw), query) {
			return true
		}
	}

	return false
}

// Selected returns the currently selected command.
func (p *Palette) Selected() *Command {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.selected >= 0 && p.selected < len(p.filtered) {
		return &p.filtered[p.selected]
	}
	return nil
}

// Query returns the current search query.
func (p *Palette) Query() string {
	return p.input.Value()
}

// SetQuery sets the search query.
func (p *Palette) SetQuery(query string) {
	p.input.SetValue(query)
	p.filter()
}

// DefaultCommands returns a list of default commands for cbwsh.
func DefaultCommands() []Command {
	return []Command{
		{
			ID:          "new_pane",
			Name:        "New Pane",
			Description: "Create a new terminal pane",
			Shortcut:    "Ctrl+N",
			Category:    "Panes",
			Keywords:    []string{"create", "add", "terminal"},
		},
		{
			ID:          "close_pane",
			Name:        "Close Pane",
			Description: "Close the active pane",
			Shortcut:    "Ctrl+W",
			Category:    "Panes",
			Keywords:    []string{"remove", "delete"},
		},
		{
			ID:          "split_vertical",
			Name:        "Split Vertical",
			Description: "Split pane vertically",
			Shortcut:    "Ctrl+\\",
			Category:    "Panes",
			Keywords:    []string{"divide", "layout"},
		},
		{
			ID:          "split_horizontal",
			Name:        "Split Horizontal",
			Description: "Split pane horizontally",
			Shortcut:    "Ctrl+-",
			Category:    "Panes",
			Keywords:    []string{"divide", "layout"},
		},
		{
			ID:          "toggle_sidebar",
			Name:        "Toggle Sidebar",
			Description: "Show/hide the sidebar",
			Shortcut:    "Ctrl+B",
			Category:    "View",
			Keywords:    []string{"hide", "show", "panel"},
		},
		{
			ID:          "ai_assist",
			Name:        "AI Assist",
			Description: "Open AI assistant",
			Shortcut:    "Ctrl+A",
			Category:    "AI",
			Keywords:    []string{"help", "suggest", "chat"},
		},
		{
			ID:          "clear_screen",
			Name:        "Clear Screen",
			Description: "Clear terminal output",
			Shortcut:    "Ctrl+L",
			Category:    "View",
			Keywords:    []string{"clean", "reset"},
		},
		{
			ID:          "show_help",
			Name:        "Show Help",
			Description: "Display help information",
			Shortcut:    "Ctrl+?",
			Category:    "Help",
			Keywords:    []string{"documentation", "guide"},
		},
		{
			ID:          "quit",
			Name:        "Quit",
			Description: "Exit the application",
			Shortcut:    "Ctrl+Q",
			Category:    "Application",
			Keywords:    []string{"exit", "close"},
		},
		{
			ID:          "settings",
			Name:        "Settings",
			Description: "Open settings",
			Category:    "Application",
			Keywords:    []string{"config", "preferences", "options"},
		},
		{
			ID:          "theme",
			Name:        "Change Theme",
			Description: "Switch color theme",
			Category:    "View",
			Keywords:    []string{"colors", "appearance"},
		},
		{
			ID:          "ssh_connect",
			Name:        "SSH Connect",
			Description: "Connect to SSH host",
			Category:    "SSH",
			Keywords:    []string{"remote", "server"},
		},
		{
			ID:          "git_status",
			Name:        "Git Status",
			Description: "Show git repository status",
			Category:    "Git",
			Keywords:    []string{"changes", "modified"},
		},
		{
			ID:          "secrets_manager",
			Name:        "Secrets Manager",
			Description: "Manage encrypted secrets",
			Category:    "Security",
			Keywords:    []string{"password", "keys", "credentials"},
		},
	}
}
