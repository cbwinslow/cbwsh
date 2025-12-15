// Package app provides the main application for cbwsh.
// It contains the Bubble Tea model and all core application logic.
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
	"github.com/cbwinslow/cbwsh/pkg/ai/monitor"
	"github.com/cbwinslow/cbwsh/pkg/config"
	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/logging"
	"github.com/cbwinslow/cbwsh/pkg/panes"
	"github.com/cbwinslow/cbwsh/pkg/plugins"
	"github.com/cbwinslow/cbwsh/pkg/privileges"
	"github.com/cbwinslow/cbwsh/pkg/process"
	"github.com/cbwinslow/cbwsh/pkg/secrets"
	"github.com/cbwinslow/cbwsh/pkg/shell"
	"github.com/cbwinslow/cbwsh/pkg/ssh"
	"github.com/cbwinslow/cbwsh/pkg/ui/aimonitor"
	"github.com/cbwinslow/cbwsh/pkg/ui/autocomplete"
	"github.com/cbwinslow/cbwsh/pkg/ui/highlight"
	"github.com/cbwinslow/cbwsh/pkg/ui/markdown"
	"github.com/cbwinslow/cbwsh/pkg/ui/menu"
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
	config           *config.Config
	paneManager      *panes.Manager
	pluginManager    *plugins.Manager
	secretsManager   *secrets.Manager
	sshManager       *ssh.Manager
	aiManager        *ai.Manager
	activityMonitor  *monitor.Monitor
	history          *shell.History
	jobManager       *process.JobManager
	privilegeManager *privileges.Manager
	logger           *logging.Logger

	// UI components
	input       textinput.Model
	spinner     spinner.Model
	styles      *styles.Styles
	highlighter *highlight.ShellHighlighter
	completer   *autocomplete.Completer
	mdRenderer  *markdown.Renderer
	menuBar     *menu.MenuBar
	monitorPane *aimonitor.MonitorPane

	// State
	mode          Mode
	width         int
	height        int
	ready         bool
	executing     bool
	showSidebar   bool
	showStatusBar bool
	showMenuBar   bool
	showMonitor   bool
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
	ToggleMenuBar   key.Binding
	ToggleMonitor   key.Binding
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
		ToggleMonitor: key.NewBinding(
			key.WithKeys("ctrl+m"),
			key.WithHelp("ctrl+m", "toggle AI monitor"),
		),
		CommandPalette: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "command palette"),
		),
		ToggleMenuBar: key.NewBinding(
			key.WithKeys("f10", "alt+m"),
			key.WithHelp("F10", "toggle menu"),
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

// New creates a new application model with the given config file path.
// If configPath is empty, it uses the default configuration.
func New(configPath string) (Model, error) {
	var cfg *config.Config
	var err error

	// Load configuration
	if configPath != "" {
		cfg, err = config.Load(configPath)
		if err != nil {
			return Model{}, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	} else {
		cfg, err = config.LoadFromDefaultPath()
		if err != nil {
			return Model{}, fmt.Errorf("failed to load config from default path: %w", err)
		}
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		return Model{}, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create text input component
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 80

	// Create spinner component
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Create markdown renderer with error handling
	mdRenderer, err := markdown.NewRenderer()
	if err != nil {
		return Model{}, fmt.Errorf("failed to create markdown renderer: %w", err)
	}

	// Create logger
	logger := logging.New(logging.WithLevel(logging.LevelInfo))
	logger.Info("Initializing cbwsh...")

	// Create menu bar with default menus
	menuBar := menu.NewMenuBar()
	for _, m := range menu.CreateDefaultMenus() {
		menuBar.AddMenu(m)
	}

	// Create activity monitor with error handling
	monitorCfg := &monitor.Config{
		OllamaURL:          cfg.AI.OllamaURL,
		OllamaModel:        cfg.AI.OllamaModel,
		MaxActivities:      100,
		MaxRecommendations: 50,
		AutoRecommend:      cfg.AI.EnableMonitoring,
		RecommendInterval:  time.Duration(cfg.AI.MonitoringInterval) * time.Second,
		MinActivityGap:     1 * time.Second,
	}
	activityMonitor := monitor.NewMonitor(monitorCfg)

	// Create monitor pane
	monitorPane := aimonitor.NewMonitorPane(activityMonitor)

	return Model{
		config:           cfg,
		paneManager:      panes.NewManager(cfg.Shell.DefaultShell),
		pluginManager:    plugins.NewManager(),
		secretsManager:   secrets.NewManager(cfg.Secrets.StorePath),
		sshManager:       ssh.NewManager("", time.Duration(cfg.SSH.ConnectTimeout)*time.Second),
		aiManager:        ai.NewManager(),
		activityMonitor:  activityMonitor,
		history:          shell.NewHistory(cfg.Shell.HistorySize, cfg.Shell.HistoryPath),
		jobManager:       process.NewJobManager(100),
		privilegeManager: privileges.NewManager(),
		logger:           logger,
		input:            ti,
		spinner:          s,
		styles:           styles.DefaultStyles(),
		highlighter:      highlight.NewShellHighlighter(),
		completer:        autocomplete.NewCompleter(),
		mdRenderer:       mdRenderer,
		menuBar:          menuBar,
		monitorPane:      monitorPane,
		mode:             ModeNormal,
		showStatusBar:    cfg.UI.ShowStatusBar,
		showMonitor:      cfg.AI.EnableMonitoring,
		commandOutput:    make([]outputLine, 0),
	}, nil
}

// validateConfig validates the configuration values
func validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	// Validate shell configuration
	if cfg.Shell.HistorySize < 0 {
		return fmt.Errorf("shell.history_size must be non-negative")
	}

	// Validate AI configuration
	if cfg.AI.MonitoringInterval < 0 {
		return fmt.Errorf("ai.monitoring_interval must be non-negative")
	}

	// Validate SSH configuration
	if cfg.SSH.ConnectTimeout < 0 {
		return fmt.Errorf("ssh.connect_timeout must be non-negative")
	}

	return nil
}

// Init initializes the application.
// It sets up the initial pane, loads command history, and starts background services.
func (m Model) Init() tea.Cmd {
	// Create initial pane with error handling
	if _, err := m.paneManager.Create(); err != nil {
		m.logger.Errorf("Failed to create initial pane: %v", err)
	}

	// Load history with error handling
	if err := m.history.Load(); err != nil {
		m.logger.Warnf("Failed to load command history: %v", err)
	}

	// Start activity monitor if enabled
	if m.showMonitor && m.activityMonitor != nil {
		m.activityMonitor.Start()
		m.logger.Info("AI activity monitor started")
	}

	m.logger.Info("Application initialized successfully")

	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
		aimonitor.Tick(),
	)
}

// Update handles messages and updates the application state.
// This is the main event handler for the Bubble Tea application.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update dimensions
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 10

		// Calculate available width for panes and monitor
		availableWidth := msg.Width
		monitorWidth := 0
		if m.showMonitor && m.monitorPane != nil {
			monitorWidth = msg.Width / 3 // Monitor takes 1/3 of width
			availableWidth = msg.Width - monitorWidth
			m.monitorPane.SetSize(monitorWidth, msg.Height-4)
		}

		m.paneManager.UpdateAllSizes(availableWidth, msg.Height-4)
		m.menuBar.SetWidth(msg.Width)
		m.ready = true
		m.logger.Debugf("Window resized to %dx%d", msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Handle menu bar input first when it's open
		if m.menuBar.IsOpen() {
			handled, cmd := m.menuBar.Update(msg)
			if handled {
				return m, cmd
			}
		}

		switch {
		case key.Matches(msg, keys.Quit):
			m.logger.Info("Application shutting down")
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

		case key.Matches(msg, keys.ToggleMonitor):
			m.showMonitor = !m.showMonitor
			if m.monitorPane != nil {
				m.monitorPane.Toggle()
			}
			if m.showMonitor && m.activityMonitor != nil {
				m.activityMonitor.Start()
			} else if m.activityMonitor != nil {
				m.activityMonitor.Stop()
			}
			// Force resize to accommodate the monitor pane
			if m.ready {
				return m.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			}
			return m, nil

		case key.Matches(msg, keys.ToggleMenuBar):
			m.menuBar.Toggle()
			m.showMenuBar = m.menuBar.IsOpen()
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
		m.logger.Debugf("Command completed: %s (exit code: %d)", result.Command, result.ExitCode)

		// Record activity to monitor
		if m.activityMonitor != nil && m.activityMonitor.IsEnabled() {
			pane := m.paneManager.ActivePane()
			workDir := "~"
			if pane != nil {
				workDir = pane.GetShellExecutor().GetWorkingDirectory()
			}
			m.activityMonitor.RecordCommand(&result, workDir)
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

	// Update monitor pane
	if m.monitorPane != nil {
		var cmd tea.Cmd
		var updatedPane *aimonitor.MonitorPane
		updatedPane, cmd = m.monitorPane.Update(msg)
		m.monitorPane = updatedPane
		cmds = append(cmds, cmd)
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

// runCommand executes a shell command in a background goroutine.
// It returns a tea.Cmd that sends the result back when complete.
func (m *Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		// Get active pane
		pane := m.paneManager.ActivePane()
		if pane == nil {
			m.logger.Error("No active pane available for command execution")
			return commandResultMsg{
				Command:  command,
				Error:    "No active pane available. Please create a pane first.",
				ExitCode: -1,
			}
		}

		// Execute command with context
		ctx := context.Background()
		m.logger.Debugf("Executing command: %s", command)
		
		result, err := pane.GetShellExecutor().Execute(ctx, command)
		if err != nil {
			m.logger.Errorf("Command execution failed: %v", err)
			return commandResultMsg{
				Command:  command,
				Error:    fmt.Sprintf("Execution error: %v", err),
				ExitCode: -1,
			}
		}

		return commandResultMsg(*result)
	}
}

// handleBuiltin handles built-in shell commands.
// Returns true if the command was handled as a builtin.
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
	case "jobs":
		// List background jobs
		jobs := m.jobManager.ListJobs()
		if len(jobs) == 0 {
			m.addOutput("No background jobs", false, 0)
		} else {
			for _, job := range jobs {
				m.addOutput(job.String(), false, 0)
			}
		}
		return true, m
	case "whoami":
		info := m.privilegeManager.GetUserInfo()
		if info != nil {
			m.addOutput(info.Username, false, 0)
		}
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

	// Menu bar (if visible)
	if m.showMenuBar || m.menuBar.IsOpen() {
		menuView := m.menuBar.View()
		sections = append(sections, menuView)
	}

	// Header
	header := m.renderHeader()
	sections = append(sections, header)

	// Main content area (with monitor pane if visible)
	var mainArea string
	if m.showMonitor && m.monitorPane != nil {
		// Split layout: content on left, monitor on right
		content := m.renderContent()
		monitorView := m.monitorPane.View()
		mainArea = lipgloss.JoinHorizontal(lipgloss.Top, content, monitorView)
	} else {
		mainArea = m.renderContent()
	}
	sections = append(sections, mainArea)

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
| Ctrl+M | Toggle AI monitor |
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

// Run starts the application with the given config file path.
// If configPath is empty, it uses the default configuration.
func Run(configPath string) error {
	// Create the application model with error handling
	model, err := New(configPath)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	// Create and run the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program and handle errors
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("application error: %w", err)
	}

	return nil
}
