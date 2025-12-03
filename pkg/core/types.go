// Package core provides fundamental types, interfaces, and enums for the cbwsh shell.
package core

import (
	"context"
	"errors"
)

// Common errors.
var (
	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when a resource already exists.
	ErrAlreadyExists = errors.New("already exists")
	// ErrInvalidInput is returned when input is invalid.
	ErrInvalidInput = errors.New("invalid input")
)

// ShellType represents the type of shell to execute commands in.
type ShellType int

const (
	// ShellTypeBash represents the Bash shell.
	ShellTypeBash ShellType = iota
	// ShellTypeZsh represents the Zsh shell.
	ShellTypeZsh
)

// String returns the string representation of the shell type.
func (s ShellType) String() string {
	switch s {
	case ShellTypeBash:
		return "bash"
	case ShellTypeZsh:
		return "zsh"
	default:
		return "unknown"
	}
}

// PaneLayout defines how panes are arranged in the terminal.
type PaneLayout int

const (
	// LayoutSingle is a single pane layout.
	LayoutSingle PaneLayout = iota
	// LayoutHorizontalSplit splits panes horizontally.
	LayoutHorizontalSplit
	// LayoutVerticalSplit splits panes vertically.
	LayoutVerticalSplit
	// LayoutGrid arranges panes in a grid.
	LayoutGrid
	// LayoutCustom allows custom pane arrangements.
	LayoutCustom
)

// String returns the string representation of the layout.
func (l PaneLayout) String() string {
	switch l {
	case LayoutSingle:
		return "single"
	case LayoutHorizontalSplit:
		return "horizontal"
	case LayoutVerticalSplit:
		return "vertical"
	case LayoutGrid:
		return "grid"
	case LayoutCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// PluginType defines the type of plugin.
type PluginType int

const (
	// PluginTypeCommand adds new commands.
	PluginTypeCommand PluginType = iota
	// PluginTypeUI modifies the UI.
	PluginTypeUI
	// PluginTypeHook provides lifecycle hooks.
	PluginTypeHook
	// PluginTypeFormatter formats output.
	PluginTypeFormatter
)

// String returns the string representation of the plugin type.
func (p PluginType) String() string {
	switch p {
	case PluginTypeCommand:
		return "command"
	case PluginTypeUI:
		return "ui"
	case PluginTypeHook:
		return "hook"
	case PluginTypeFormatter:
		return "formatter"
	default:
		return "unknown"
	}
}

// CommandResult represents the result of executing a shell command.
type CommandResult struct {
	// Command is the original command that was executed.
	Command string
	// Output is the stdout from the command.
	Output string
	// Error is the stderr from the command.
	Error string
	// ExitCode is the exit code of the command.
	ExitCode int
	// Duration is how long the command took to execute.
	Duration int64
}

// SSHConnectionState represents the state of an SSH connection.
type SSHConnectionState int

const (
	// SSHDisconnected indicates no active connection.
	SSHDisconnected SSHConnectionState = iota
	// SSHConnecting indicates connection in progress.
	SSHConnecting
	// SSHConnected indicates an active connection.
	SSHConnected
	// SSHError indicates a connection error.
	SSHError
)

// String returns the string representation of the SSH connection state.
func (s SSHConnectionState) String() string {
	switch s {
	case SSHDisconnected:
		return "disconnected"
	case SSHConnecting:
		return "connecting"
	case SSHConnected:
		return "connected"
	case SSHError:
		return "error"
	default:
		return "unknown"
	}
}

// AIProvider represents different AI service providers.
type AIProvider int

const (
	// AIProviderNone indicates no AI provider.
	AIProviderNone AIProvider = iota
	// AIProviderOpenAI uses OpenAI API.
	AIProviderOpenAI
	// AIProviderAnthropic uses Anthropic API.
	AIProviderAnthropic
	// AIProviderGemini uses Google Gemini API.
	AIProviderGemini
	// AIProviderLocal uses a local LLM.
	AIProviderLocal
)

// String returns the string representation of the AI provider.
func (a AIProvider) String() string {
	switch a {
	case AIProviderNone:
		return "none"
	case AIProviderOpenAI:
		return "openai"
	case AIProviderAnthropic:
		return "anthropic"
	case AIProviderGemini:
		return "gemini"
	case AIProviderLocal:
		return "local"
	default:
		return "unknown"
	}
}

// Executor defines the interface for command execution.
type Executor interface {
	// Execute runs a command and returns the result.
	Execute(ctx context.Context, command string) (*CommandResult, error)
	// ExecuteAsync runs a command asynchronously.
	ExecuteAsync(ctx context.Context, command string) (<-chan *CommandResult, error)
	// Interrupt stops the currently running command.
	Interrupt() error
	// SetShellType sets the shell type (bash/zsh).
	SetShellType(shellType ShellType) error
	// GetShellType returns the current shell type.
	GetShellType() ShellType
	// SetWorkingDirectory sets the current working directory.
	SetWorkingDirectory(path string) error
	// GetWorkingDirectory returns the current working directory.
	GetWorkingDirectory() string
	// SetEnvironment sets environment variables.
	SetEnvironment(env map[string]string) error
	// GetEnvironment returns current environment variables.
	GetEnvironment() map[string]string
}

// SecretsManager defines the interface for secrets management.
type SecretsManager interface {
	// Store securely stores a secret.
	Store(key string, value []byte) error
	// Retrieve gets a stored secret.
	Retrieve(key string) ([]byte, error)
	// Delete removes a stored secret.
	Delete(key string) error
	// List returns all stored secret keys.
	List() ([]string, error)
	// Exists checks if a secret exists.
	Exists(key string) bool
}

// SSHManager defines the interface for SSH connection management.
type SSHManager interface {
	// Connect establishes an SSH connection.
	Connect(ctx context.Context, host string, port int, user string) error
	// Disconnect closes the current SSH connection.
	Disconnect() error
	// Execute runs a command on the remote host.
	Execute(ctx context.Context, command string) (*CommandResult, error)
	// State returns the current connection state.
	State() SSHConnectionState
	// ListSavedHosts returns a list of saved SSH hosts.
	ListSavedHosts() ([]SSHHost, error)
	// SaveHost saves an SSH host configuration.
	SaveHost(host SSHHost) error
	// RemoveHost removes a saved SSH host.
	RemoveHost(name string) error
}

// SSHHost represents a saved SSH host configuration.
type SSHHost struct {
	Name       string
	Host       string
	Port       int
	User       string
	KeyPath    string
	Passphrase string
}

// Plugin defines the interface that all plugins must implement.
type Plugin interface {
	// Name returns the plugin name.
	Name() string
	// Type returns the plugin type.
	Type() PluginType
	// Version returns the plugin version.
	Version() string
	// Initialize sets up the plugin.
	Initialize(ctx context.Context) error
	// Shutdown cleans up the plugin.
	Shutdown(ctx context.Context) error
	// Enabled returns whether the plugin is enabled.
	Enabled() bool
	// Enable enables the plugin.
	Enable() error
	// Disable disables the plugin.
	Disable() error
}

// PluginManager defines the interface for plugin management.
type PluginManager interface {
	// Register adds a plugin to the manager.
	Register(plugin Plugin) error
	// Unregister removes a plugin from the manager.
	Unregister(name string) error
	// Get returns a plugin by name.
	Get(name string) (Plugin, bool)
	// List returns all registered plugins.
	List() []Plugin
	// ListByType returns plugins of a specific type.
	ListByType(pluginType PluginType) []Plugin
	// Initialize initializes all registered plugins.
	Initialize(ctx context.Context) error
	// Shutdown shuts down all registered plugins.
	Shutdown(ctx context.Context) error
}

// AIAgent defines the interface for AI-powered agents.
type AIAgent interface {
	// Name returns the agent name.
	Name() string
	// Provider returns the AI provider.
	Provider() AIProvider
	// Query sends a query to the AI agent.
	Query(ctx context.Context, prompt string) (string, error)
	// StreamQuery sends a query and streams the response.
	StreamQuery(ctx context.Context, prompt string) (<-chan string, error)
	// SuggestCommand suggests a command based on natural language.
	SuggestCommand(ctx context.Context, description string) (string, error)
	// ExplainCommand explains what a command does.
	ExplainCommand(ctx context.Context, command string) (string, error)
	// FixError suggests a fix for an error.
	FixError(ctx context.Context, command string, err string) (string, error)
}

// Pane defines the interface for terminal panes.
type Pane interface {
	// ID returns the pane identifier.
	ID() string
	// Title returns the pane title.
	Title() string
	// SetTitle sets the pane title.
	SetTitle(title string)
	// IsActive returns whether the pane is active.
	IsActive() bool
	// Activate makes this pane active.
	Activate()
	// Deactivate makes this pane inactive.
	Deactivate()
	// Width returns the pane width.
	Width() int
	// Height returns the pane height.
	Height() int
	// SetSize sets the pane dimensions.
	SetSize(width, height int)
	// GetExecutor returns the pane's executor.
	GetExecutor() Executor
}

// PaneManager defines the interface for pane management.
type PaneManager interface {
	// Create creates a new pane.
	Create() (Pane, error)
	// Close closes a pane by ID.
	Close(id string) error
	// Get returns a pane by ID.
	Get(id string) (Pane, bool)
	// Active returns the active pane.
	Active() Pane
	// SetActive sets the active pane.
	SetActive(id string) error
	// List returns all panes.
	List() []Pane
	// Layout returns the current layout.
	Layout() PaneLayout
	// SetLayout sets the pane layout.
	SetLayout(layout PaneLayout) error
	// Split splits the active pane.
	Split(direction PaneLayout) (Pane, error)
}

// Highlighter defines the interface for syntax highlighting.
type Highlighter interface {
	// Highlight applies syntax highlighting to text.
	Highlight(text string, language string) (string, error)
	// HighlightCommand highlights a shell command.
	HighlightCommand(command string) (string, error)
}

// Autocompleter defines the interface for autocompletion.
type Autocompleter interface {
	// Complete returns completion suggestions for the given input.
	Complete(input string, cursorPos int) ([]Suggestion, error)
	// AddProvider adds a completion provider.
	AddProvider(provider CompletionProvider)
}

// Suggestion represents an autocompletion suggestion.
type Suggestion struct {
	Text        string
	Description string
	Category    string
}

// CompletionProvider provides completion suggestions.
type CompletionProvider interface {
	// Name returns the provider name.
	Name() string
	// Provide returns suggestions for the given context.
	Provide(input string, cursorPos int) ([]Suggestion, error)
}

// Animator defines the interface for animations.
type Animator interface {
	// Start starts the animation.
	Start()
	// Stop stops the animation.
	Stop()
	// Update updates the animation state.
	Update() string
	// SetValue sets the target value for spring animations.
	SetValue(value float64)
	// GetValue returns the current animated value.
	GetValue() float64
}

// ProgressReporter defines the interface for progress reporting.
type ProgressReporter interface {
	// Start starts progress tracking.
	Start(total int)
	// Increment increments progress by one.
	Increment()
	// IncrementBy increments progress by a specific amount.
	IncrementBy(amount int)
	// SetMessage sets the progress message.
	SetMessage(message string)
	// Finish completes the progress.
	Finish()
	// View returns the progress view string.
	View() string
}

// MarkdownRenderer defines the interface for markdown rendering.
type MarkdownRenderer interface {
	// Render renders markdown to styled terminal output.
	Render(markdown string) (string, error)
	// SetWidth sets the render width.
	SetWidth(width int)
	// SetTheme sets the rendering theme.
	SetTheme(theme string) error
}
