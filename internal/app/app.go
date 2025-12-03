// Package app provides the main application for cbwsh.
package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cbwinslow/cbwsh/pkg/ai"
	"github.com/cbwinslow/cbwsh/pkg/config"
	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/panes"
	"github.com/cbwinslow/cbwsh/pkg/plugins"
	"github.com/cbwinslow/cbwsh/pkg/secrets"
	"github.com/cbwinslow/cbwsh/pkg/shell"
	"github.com/cbwinslow/cbwsh/pkg/ssh"
	"github.com/cbwinslow/cbwsh/pkg/ui/autocomplete"
	"github.com/cbwinslow/cbwsh/pkg/ui/highlight"
	"github.com/cbwinslow/cbwsh/pkg/ui/markdown"
	"github.com/cbwinslow/cbwsh/pkg/ui/styles"
)

// Mode represents the current application mode.
type Mode int

const (
	// ModeNormal is the normal command input mode.
	ModeNormal Mode = iota
	// ModeAI is the AI assistant mode.
	ModeAI
	// ModeSSH is the SSH connection mode.
	ModeSSH
	// ModeSecrets is the secrets management mode.
	ModeSecrets
	// ModeHelp is the help display mode.
	ModeHelp
	// ModeCommandPalette is the command palette mode.
	ModeCommandPalette
)

// Model is the main application model.
type Model struct {
	// Core components
	config         *config.Config
	paneManager    *panes.Manager
	pluginManager  *plugins.Manager
	secretsManager *secrets.Manager
	sshManager     *ssh.Manager
	aiManager      *ai.Manager
	history        *shell.History

	// UI components
	input       textinput.Model
	spinner     spinner.Model
	styles      *styles.Styles
	highlighter *highlight.ShellHighlighter
	completer   *autocomplete.Completer
	mdRenderer  *markdown.Renderer

	// State
	mode          Mode
	width         int
	height        int
	ready         bool
	executing     bool
	showSidebar   bool
	showStatusBar bool
	suggestions   []core.Suggestion
	selectedSugg  int
	commandOutput []outputLine
	lastError     string
}

type outputLine struct {
	content   string
	isCommand bool
	exitCode  int
}

// KeyMap defines the key bindings.
type KeyMap struct {
	Quit            key.Binding
	Help            key.Binding
	NewPane         key.Binding
	ClosePane       key.Binding
	NextPane        key.Binding
	PrevPane        key.Binding
	SplitVertical   key.Binding
	SplitHorizontal key.Binding
	ToggleSidebar   key.Binding
	CommandPalette  key.Binding
	AIAssist        key.Binding
	Execute         key.Binding
	Cancel          key.Binding
	Up              key.Binding
	Down            key.Binding
	Tab             key.Binding
	Clear           key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+q"),
			key.WithHelp("ctrl+q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("ctrl+?", "f1"),
			key.WithHelp("ctrl+?", "help"),
		),
		NewPane: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "new pane"),
		),
		ClosePane: key.NewBinding(
			key.WithKeys("ctrl+w"),
			key.WithHelp("ctrl+w", "close pane"),
		),
		NextPane: key.NewBinding(
			key.WithKeys("ctrl+tab", "ctrl+]"),
			key.WithHelp("ctrl+]", "next pane"),
		),
		PrevPane: key.NewBinding(
			key.WithKeys("ctrl+shift+tab", "ctrl+["),
			key.WithHelp("ctrl+[", "prev pane"),
		),
		SplitVertical: key.NewBinding(
			key.WithKeys("ctrl+\\"),
			key.WithHelp("ctrl+\\", "split vertical"),
		),
		SplitHorizontal: key.NewBinding(
			key.WithKeys("ctrl+-"),
			key.WithHelp("ctrl+-", "split horizontal"),
		),
		ToggleSidebar: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "toggle sidebar"),
		),
		CommandPalette: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "command palette"),
		),
		AIAssist: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "AI assist"),
		),
		Execute: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "execute"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "cancel"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "history up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "history down"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "autocomplete"),
		),
		Clear: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "clear"),
		),
	}
}

var keys = DefaultKeyMap()

// New creates a new application model.
func New() Model {
	cfg := config.Default()

	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	mdRenderer, _ := markdown.NewRenderer()

	return Model{
		config:         cfg,
		paneManager:    panes.NewManager(cfg.Shell.DefaultShell),
		pluginManager:  plugins.NewManager(),
		secretsManager: secrets.NewManager(cfg.Secrets.StorePath),
		sshManager:     ssh.NewManager("", time.Duration(cfg.SSH.ConnectTimeout)*time.Second),
		aiManager:      ai.NewManager(),
		history:        shell.NewHistory(cfg.Shell.HistorySize, cfg.Shell.HistoryPath),
		input:          ti,
		spinner:        s,
		styles:         styles.DefaultStyles(),
		highlighter:    highlight.NewShellHighlighter(),
		completer:      autocomplete.NewCompleter(),
		mdRenderer:     mdRenderer,
		mode:           ModeNormal,
		showStatusBar:  cfg.UI.ShowStatusBar,
		commandOutput:  make([]outputLine, 0),
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	// Create initial pane
	_, _ = m.paneManager.Create()

	// Load history
	_ = m.history.Load()

	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
	)
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 10
		m.paneManager.UpdateAllSizes(msg.Width, msg.Height-4)
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			_ = m.history.Save()
			return m, tea.Quit

		case key.Matches(msg, keys.Cancel):
			if m.executing {
				pane := m.paneManager.ActivePane()
				if pane != nil {
					_ = pane.GetShellExecutor().Interrupt()
				}
				m.executing = false
			}
			m.input.Reset()
			m.suggestions = nil
			return m, nil

		case key.Matches(msg, keys.Execute):
			if m.input.Value() != "" {
				return m.executeCommand()
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if len(m.suggestions) > 0 {
				m.selectedSugg--
				if m.selectedSugg < 0 {
					m.selectedSugg = len(m.suggestions) - 1
				}
				return m, nil
			}
			if cmd, ok := m.history.Previous(); ok {
				m.input.SetValue(cmd)
				m.input.SetCursor(len(cmd))
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if len(m.suggestions) > 0 {
				m.selectedSugg = (m.selectedSugg + 1) % len(m.suggestions)
				return m, nil
			}
			if cmd, ok := m.history.Next(); ok {
				m.input.SetValue(cmd)
				m.input.SetCursor(len(cmd))
			} else {
				m.input.Reset()
			}
			return m, nil

		case key.Matches(msg, keys.Tab):
			if len(m.suggestions) > 0 && m.selectedSugg >= 0 && m.selectedSugg < len(m.suggestions) {
				m.input.SetValue(m.suggestions[m.selectedSugg].Text)
				m.input.SetCursor(len(m.input.Value()))
				m.suggestions = nil
				return m, nil
			}
			// Get completions
			suggestions, _ := m.completer.Complete(m.input.Value(), m.input.Position())
			m.suggestions = suggestions
			m.selectedSugg = 0
			return m, nil

		case key.Matches(msg, keys.Clear):
			m.commandOutput = m.commandOutput[:0]
			return m, nil

		case key.Matches(msg, keys.Help):
			if m.mode == ModeHelp {
				m.mode = ModeNormal
			} else {
				m.mode = ModeHelp
			}
			return m, nil

		case key.Matches(msg, keys.NewPane):
			_, _ = m.paneManager.Create()
			return m, nil

		case key.Matches(msg, keys.ClosePane):
			active := m.paneManager.ActivePane()
			if active != nil && m.paneManager.Count() > 1 {
				_ = m.paneManager.Close(active.ID())
			}
			return m, nil

		case key.Matches(msg, keys.NextPane):
			_ = m.paneManager.NextPane()
			return m, nil

		case key.Matches(msg, keys.PrevPane):
			_ = m.paneManager.PrevPane()
			return m, nil

		case key.Matches(msg, keys.SplitVertical):
			_, _ = m.paneManager.Split(core.LayoutVerticalSplit)
			return m, nil

		case key.Matches(msg, keys.SplitHorizontal):
			_, _ = m.paneManager.Split(core.LayoutHorizontalSplit)
			return m, nil

		case key.Matches(msg, keys.ToggleSidebar):
			m.showSidebar = !m.showSidebar
			return m, nil

		case key.Matches(msg, keys.AIAssist):
			if m.mode == ModeAI {
				m.mode = ModeNormal
			} else {
				m.mode = ModeAI
			}
			return m, nil
		}

	case commandResultMsg:
		m.executing = false
		result := core.CommandResult(msg)
		m.addOutput(result.Output, false, result.ExitCode)
		if result.Error != "" {
			m.addOutput(result.Error, false, result.ExitCode)
		}
		return m, nil

	case spinner.TickMsg:
		if m.executing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update text input
	if m.mode == ModeNormal || m.mode == ModeAI {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)

		// Clear suggestions on input change
		if _, ok := msg.(tea.KeyMsg); ok {
			m.suggestions = nil
		}
	}

	return m, tea.Batch(cmds...)
}

type commandResultMsg core.CommandResult

func (m Model) executeCommand() (tea.Model, tea.Cmd) {
	command := strings.TrimSpace(m.input.Value())
	if command == "" {
		return m, nil
	}

	// Add to history
	m.history.Add(command)
	m.history.Reset()

	// Show command in output
	m.addOutput(m.getPrompt()+command, true, 0)

	// Handle built-in commands
	if handled, model := m.handleBuiltin(command); handled {
		m.input.Reset()
		return model, nil
	}

	// Execute command
	m.executing = true
	m.input.Reset()

	return m, tea.Batch(
		m.spinner.Tick,
		m.runCommand(command),
	)
}

func (m *Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		pane := m.paneManager.ActivePane()
		if pane == nil {
			return commandResultMsg{
				Command:  command,
				Error:    "No active pane",
				ExitCode: -1,
			}
		}

		result, err := pane.GetShellExecutor().Execute(context.Background(), command)
		if err != nil {
			return commandResultMsg{
				Command:  command,
				Error:    err.Error(),
				ExitCode: -1,
			}
		}

		return commandResultMsg(*result)
	}
}

func (m *Model) handleBuiltin(command string) (bool, tea.Model) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false, m
	}

	switch parts[0] {
	case "exit", "quit":
		return true, m
	case "clear":
		m.commandOutput = m.commandOutput[:0]
		return true, m
	case "cd":
		if len(parts) > 1 {
			pane := m.paneManager.ActivePane()
			if pane != nil {
				if err := pane.GetShellExecutor().SetWorkingDirectory(parts[1]); err != nil {
					m.lastError = err.Error()
				}
			}
		}
		return true, m
	case "help":
		m.mode = ModeHelp
		return true, m
	}

	return false, m
}

func (m *Model) addOutput(content string, isCommand bool, exitCode int) {
	m.commandOutput = append(m.commandOutput, outputLine{
		content:   content,
		isCommand: isCommand,
		exitCode:  exitCode,
	})

	// Limit output history
	if len(m.commandOutput) > 1000 {
		m.commandOutput = m.commandOutput[len(m.commandOutput)-1000:]
	}
}

func (m Model) getPrompt() string {
	pane := m.paneManager.ActivePane()
	if pane == nil {
		return "$ "
	}
	cwd := pane.GetShellExecutor().GetWorkingDirectory()
	return m.styles.Prompt.Render(cwd) + " " + m.styles.PromptSymbol.Render("❯") + " "
}

// View renders the application.
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var sections []string

	// Header
	header := m.renderHeader()
	sections = append(sections, header)

	// Main content
	content := m.renderContent()
	sections = append(sections, content)

	// Input area
	inputArea := m.renderInput()
	sections = append(sections, inputArea)

	// Status bar
	if m.showStatusBar {
		statusBar := m.renderStatusBar()
		sections = append(sections, statusBar)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderHeader() string {
	title := m.styles.Header.Render("🐚 cbwsh")
	mode := ""
	switch m.mode {
	case ModeAI:
		mode = m.styles.Info.Render(" [AI Mode]")
	case ModeSSH:
		mode = m.styles.Info.Render(" [SSH Mode]")
	case ModeHelp:
		mode = m.styles.Info.Render(" [Help]")
	}
	return title + mode
}

func (m Model) renderContent() string {
	if m.mode == ModeHelp {
		return m.renderHelp()
	}

	// Render command output
	var lines []string
	outputHeight := m.height - 6 // Leave room for header, input, status
	startIdx := 0
	if len(m.commandOutput) > outputHeight {
		startIdx = len(m.commandOutput) - outputHeight
	}

	for i := startIdx; i < len(m.commandOutput); i++ {
		line := m.commandOutput[i]
		if line.isCommand {
			lines = append(lines, m.highlighter.HighlightCommand(line.content))
		} else if line.exitCode != 0 {
			lines = append(lines, m.styles.Error.Render(line.content))
		} else {
			lines = append(lines, line.content)
		}
	}

	content := strings.Join(lines, "\n")
	return content
}

func (m Model) renderInput() string {
	prompt := m.getPrompt()

	if m.executing {
		return prompt + m.spinner.View() + " Running..."
	}

	inputView := m.input.View()

	// Render suggestions
	if len(m.suggestions) > 0 {
		suggView := m.renderSuggestions()
		return prompt + inputView + "\n" + suggView
	}

	return prompt + inputView
}

func (m Model) renderSuggestions() string {
	var lines []string

	for i, sugg := range m.suggestions {
		style := m.styles.Suggestion
		if i == m.selectedSugg {
			style = m.styles.SelectedSuggestion
		}

		line := style.Render(sugg.Text)
		if sugg.Description != "" {
			line += " " + m.styles.SuggestionDesc.Render(sugg.Description)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderStatusBar() string {
	pane := m.paneManager.ActivePane()
	shellType := "bash"
	cwd := "~"
	if pane != nil {
		shellType = pane.GetShellExecutor().GetShellType().String()
		cwd = pane.GetShellExecutor().GetWorkingDirectory()
	}

	left := fmt.Sprintf(" %s | %s", shellType, cwd)
	right := fmt.Sprintf("Panes: %d ", m.paneManager.Count())

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	return m.styles.StatusBar.Render(left + strings.Repeat(" ", gap) + right)
}

func (m Model) renderHelp() string {
	helpText := `# cbwsh Help

## Key Bindings

| Key | Action |
|-----|--------|
| Ctrl+Q | Quit |
| Ctrl+C | Cancel current command |
| Enter | Execute command |
| Tab | Autocomplete |
| ↑/↓ | Navigate history |
| Ctrl+L | Clear screen |
| Ctrl+N | New pane |
| Ctrl+W | Close pane |
| Ctrl+] | Next pane |
| Ctrl+[ | Previous pane |
| Ctrl+\ | Split vertical |
| Ctrl+- | Split horizontal |
| Ctrl+B | Toggle sidebar |
| Ctrl+A | AI assist mode |
| Ctrl+? | Help |

## Built-in Commands

- **cd** - Change directory
- **clear** - Clear screen
- **exit** - Exit shell
- **help** - Show this help

Press any key to return...
`

	rendered, err := m.mdRenderer.Render(helpText)
	if err != nil {
		return helpText
	}
	return rendered
}

// Run starts the application.
func Run() error {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
