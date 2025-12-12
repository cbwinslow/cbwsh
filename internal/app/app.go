// Package app provides the main application logic for cbwsh.
//
// This package implements the Bubble Tea-based terminal user interface,
// managing the application state, handling user input, and coordinating
// between various subsystems like shell execution, AI integration,
// pane management, and more.
//
// The main entry point is the Run() function, which creates and starts
// the Bubble Tea program with proper initialization and error handling.
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

// New creates a new application model with all subsystems initialized.
//
// This function sets up:
//   - Configuration from default settings
//   - UI components (text input, spinner, markdown renderer)
//   - Shell subsystems (pane manager, history, job manager)
//   - AI components (manager, activity monitor, monitor pane)
//   - Security components (secrets manager, SSH manager, privileges)
//   - Logging infrastructure
//
// The returned Model is ready to be used with Bubble Tea's NewProgram.
// Any initialization errors are logged but do not prevent the application
// from starting with degraded functionality.
func New() Model {
	// Load configuration with defaults
	cfg := config.Default()

	// Initialize text input for command entry
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 80

	// Initialize spinner for command execution feedback
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize markdown renderer for help and AI responses
	// Error is ignored as the renderer has a fallback mode
	mdRenderer, _ := markdown.NewRenderer()

	// Create logger for application diagnostics
	logger := logging.New(logging.WithLevel(logging.LevelInfo))

	// Create menu bar with default menus (File, Edit, View, Help)
	menuBar := menu.NewMenuBar()
	for _, m := range menu.CreateDefaultMenus() {
		menuBar.AddMenu(m)
	}

	// Create activity monitor for AI-powered shell recommendations
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

	// Create monitor pane UI component
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
	}
}

// Init initializes the application after the model is created.
//
// This is called by Bubble Tea before the first Update. It performs
// startup tasks like:
//   - Creating the initial shell pane
//   - Loading command history from disk
//   - Starting the AI activity monitor if enabled
//   - Initializing UI component animations
//
// Returns a batch of commands to start the application's event loop.
func (m Model) Init() tea.Cmd {
	// Create initial pane for shell interaction
	if _, err := m.paneManager.Create(); err != nil {
		// Log error but continue - the user can create panes manually
		m.logger.Errorf("Failed to create initial pane: %v", err)
	}

	// Load command history from persistent storage
	if err := m.history.Load(); err != nil {
		// Log error but continue - history will be empty
		m.logger.Warnf("Failed to load command history: %v", err)
	}

	// Start activity monitor if enabled in configuration
	if m.showMonitor && m.activityMonitor != nil {
		m.activityMonitor.Start()
	}

	// Return initial commands to start UI animations
	return tea.Batch(
		textinput.Blink,   // Start cursor blinking
		m.spinner.Tick,    // Start spinner animation
		aimonitor.Tick(),  // Start AI monitor updates
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

// executeCommand processes and executes a user command.
//
// This function:
//   - Validates the command is not empty
//   - Adds the command to history for recall
//   - Displays the command in the output
//   - Checks if it's a built-in command (cd, exit, etc.)
//   - Executes external commands via the shell executor
//
// Returns the updated model and any commands to run.
func (m Model) executeCommand() (tea.Model, tea.Cmd) {
	command := strings.TrimSpace(m.input.Value())
	if command == "" {
		return m, nil
	}

	// Add to command history for up/down arrow recall
	m.history.Add(command)
	m.history.Reset()

	// Show command in output with current prompt
	m.addOutput(m.getPrompt()+command, true, 0)

	// Handle built-in commands (cd, exit, help, etc.)
	if handled, model := m.handleBuiltin(command); handled {
		m.input.Reset()
		return model, nil
	}

	// Execute external command via shell
	m.executing = true
	m.input.Reset()

	return m, tea.Batch(
		m.spinner.Tick,      // Show spinner while executing
		m.runCommand(command), // Run the command asynchronously
	)
}

// runCommand executes a shell command and returns the result as a message.
//
// This function runs asynchronously in the background and sends a
// commandResultMsg when complete. It handles:
//   - Checking for an active pane
//   - Executing the command via the shell executor
//   - Capturing output, errors, and exit codes
//   - Proper error handling and reporting
//
// Returns a tea.Cmd that will send a commandResultMsg when the command completes.
func (m *Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		// Ensure we have an active pane to run the command in
		pane := m.paneManager.ActivePane()
		if pane == nil {
			return commandResultMsg{
				Command:  command,
				Error:    "No active pane available. Press Ctrl+N to create a new pane.",
				ExitCode: -1,
			}
		}

		// Execute the command with a background context
		// Context allows for future cancellation support
		result, err := pane.GetShellExecutor().Execute(context.Background(), command)
		if err != nil {
			return commandResultMsg{
				Command:  command,
				Error:    fmt.Sprintf("Execution failed: %v", err),
				ExitCode: -1,
			}
		}

		return commandResultMsg(*result)
	}
}

// handleBuiltin processes built-in shell commands.
//
// Built-in commands are handled directly by cbwsh rather than being
// passed to the underlying shell. This includes:
//   - exit, quit: Exit the shell
//   - clear: Clear the output buffer
//   - cd: Change working directory
//   - help: Show help information
//   - jobs: List background jobs
//   - whoami: Show current user
//
// Returns:
//   - bool: true if the command was handled as a built-in
//   - tea.Model: the updated model (may be modified by the command)
func (m *Model) handleBuiltin(command string) (bool, tea.Model) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false, m
	}

	switch parts[0] {
	case "exit", "quit":
		// Exit handled by returning without action - caller will quit
		return true, m
		
	case "clear":
		// Clear output buffer
		m.commandOutput = m.commandOutput[:0]
		return true, m
		
	case "cd":
		// Change working directory
		if len(parts) > 1 {
			pane := m.paneManager.ActivePane()
			if pane != nil {
				if err := pane.GetShellExecutor().SetWorkingDirectory(parts[1]); err != nil {
					m.lastError = fmt.Sprintf("cd: %v", err)
					m.addOutput(m.lastError, false, 1)
				}
			} else {
				m.addOutput("cd: No active pane", false, 1)
			}
		} else {
			m.addOutput("cd: missing directory argument", false, 1)
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

// Run starts the cbwsh application.
//
// This is the main entry point for running the shell. It:
//   - Creates a new application model with all subsystems initialized
//   - Configures the Bubble Tea program with alt screen and mouse support
//   - Starts the main event loop
//   - Returns any errors that occur during initialization or execution
//
// The application runs in alt screen mode, preserving the terminal state
// before cbwsh started. Mouse support is enabled for UI interactions.
//
// Returns:
//   - nil on successful completion (user quit normally)
//   - error if the application failed to start or encountered a fatal error
func Run() error {
	// Create and configure the Bubble Tea program
	p := tea.NewProgram(
		New(),                      // Initialize application model
		tea.WithAltScreen(),        // Use alternate screen buffer
		tea.WithMouseCellMotion(),  // Enable mouse support for UI
	)
	
	// Run the program and return any errors
	// The first return value (final model) is ignored as we don't need it
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("application error: %w", err)
	}
	
	return nil
}
