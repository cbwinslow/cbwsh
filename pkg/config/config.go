// Package config provides configuration management for cbwsh.
// It handles loading, saving, and validation of configuration files
// with support for YAML format and sensible defaults.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for cbwsh.
type Config struct {
	mu sync.RWMutex

	// Shell configuration
	Shell ShellConfig `yaml:"shell"`

	// UI configuration
	UI UIConfig `yaml:"ui"`

	// Plugins configuration
	Plugins PluginsConfig `yaml:"plugins"`

	// AI configuration
	AI AIConfig `yaml:"ai"`

	// SSH configuration
	SSH SSHConfig `yaml:"ssh"`

	// Secrets configuration
	Secrets SecretsConfig `yaml:"secrets"`

	// Keybindings
	Keybindings KeybindingsConfig `yaml:"keybindings"`
}

// ShellConfig holds shell-specific configuration.
type ShellConfig struct {
	// DefaultShell is the default shell type (bash or zsh).
	DefaultShell core.ShellType `yaml:"default_shell"`
	// HistorySize is the number of commands to keep in history.
	HistorySize int `yaml:"history_size"`
	// HistoryPath is the path to the history file.
	HistoryPath string `yaml:"history_path"`
	// StartupCommands are commands to run on shell startup.
	StartupCommands []string `yaml:"startup_commands"`
	// Environment holds additional environment variables.
	Environment map[string]string `yaml:"environment"`
	// Aliases holds command aliases.
	Aliases map[string]string `yaml:"aliases"`
}

// UIConfig holds UI-specific configuration.
type UIConfig struct {
	// Theme is the color theme name.
	Theme string `yaml:"theme"`
	// Layout is the default pane layout.
	Layout core.PaneLayout `yaml:"layout"`
	// ShowStatusBar toggles the status bar.
	ShowStatusBar bool `yaml:"show_status_bar"`
	// ShowLineNumbers toggles line numbers.
	ShowLineNumbers bool `yaml:"show_line_numbers"`
	// EnableAnimations toggles animations.
	EnableAnimations bool `yaml:"enable_animations"`
	// AnimationFPS is the frames per second for animations.
	AnimationFPS int `yaml:"animation_fps"`
	// SyntaxHighlighting toggles syntax highlighting.
	SyntaxHighlighting bool `yaml:"syntax_highlighting"`
	// HighlightTheme is the syntax highlighting theme.
	HighlightTheme string `yaml:"highlight_theme"`
	// MarkdownTheme is the markdown rendering theme.
	MarkdownTheme string `yaml:"markdown_theme"`
	// PromptStyle is the style of the prompt.
	PromptStyle string `yaml:"prompt_style"`
}

// PluginsConfig holds plugin-specific configuration.
type PluginsConfig struct {
	// Enabled lists enabled plugin names.
	Enabled []string `yaml:"enabled"`
	// Disabled lists explicitly disabled plugin names.
	Disabled []string `yaml:"disabled"`
	// Directory is the path to look for plugins.
	Directory string `yaml:"directory"`
	// AutoLoad toggles automatic plugin loading.
	AutoLoad bool `yaml:"auto_load"`
}

// AIConfig holds AI-specific configuration.
type AIConfig struct {
	// Provider is the AI provider to use.
	Provider core.AIProvider `yaml:"provider"`
	// APIKey is the API key for the AI service.
	APIKey string `yaml:"api_key"`
	// Model is the AI model to use.
	Model string `yaml:"model"`
	// MaxTokens is the maximum tokens for AI responses.
	MaxTokens int `yaml:"max_tokens"`
	// Temperature controls AI response randomness.
	Temperature float64 `yaml:"temperature"`
	// EnableSuggestions toggles AI command suggestions.
	EnableSuggestions bool `yaml:"enable_suggestions"`
	// LocalModelPath is the path to a local LLM model.
	LocalModelPath string `yaml:"local_model_path"`
	// OllamaURL is the URL for Ollama API.
	OllamaURL string `yaml:"ollama_url"`
	// OllamaModel is the model to use with Ollama.
	OllamaModel string `yaml:"ollama_model"`
	// EnableMonitoring toggles AI shell activity monitoring.
	EnableMonitoring bool `yaml:"enable_monitoring"`
	// MonitoringInterval is the interval for generating recommendations (in seconds).
	MonitoringInterval int `yaml:"monitoring_interval"`
}

// SSHConfig holds SSH-specific configuration.
type SSHConfig struct {
	// DefaultKeyPath is the default SSH key path.
	DefaultKeyPath string `yaml:"default_key_path"`
	// KnownHostsPath is the path to known_hosts file.
	KnownHostsPath string `yaml:"known_hosts_path"`
	// ConnectTimeout is the connection timeout in seconds.
	ConnectTimeout int `yaml:"connect_timeout"`
	// KeepAliveInterval is the keep-alive interval in seconds.
	KeepAliveInterval int `yaml:"keep_alive_interval"`
	// SavedHosts holds saved SSH host configurations.
	SavedHosts []core.SSHHost `yaml:"saved_hosts"`
}

// SecretsConfig holds secrets-specific configuration.
type SecretsConfig struct {
	// StorePath is the path to the secrets store.
	StorePath string `yaml:"store_path"`
	// EncryptionAlgorithm is the encryption algorithm to use.
	EncryptionAlgorithm string `yaml:"encryption_algorithm"`
	// KeyDerivation is the key derivation function.
	KeyDerivation string `yaml:"key_derivation"`
}

// KeybindingsConfig holds keybinding configuration.
type KeybindingsConfig struct {
	// Global keybindings
	Quit            string `yaml:"quit"`
	Help            string `yaml:"help"`
	NewPane         string `yaml:"new_pane"`
	ClosePane       string `yaml:"close_pane"`
	NextPane        string `yaml:"next_pane"`
	PrevPane        string `yaml:"prev_pane"`
	SplitVertical   string `yaml:"split_vertical"`
	SplitHorizontal string `yaml:"split_horizontal"`
	ToggleSidebar   string `yaml:"toggle_sidebar"`
	CommandPalette  string `yaml:"command_palette"`
	AIAssist        string `yaml:"ai_assist"`
}

// Default returns a new Config with default values.
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".cbwsh")

	return &Config{
		Shell: ShellConfig{
			DefaultShell: core.ShellTypeBash,
			HistorySize:  10000,
			HistoryPath:  filepath.Join(configDir, "history"),
			Environment:  make(map[string]string),
			Aliases:      make(map[string]string),
		},
		UI: UIConfig{
			Theme:              "default",
			Layout:             core.LayoutSingle,
			ShowStatusBar:      true,
			ShowLineNumbers:    false,
			EnableAnimations:   true,
			AnimationFPS:       60,
			SyntaxHighlighting: true,
			HighlightTheme:     "monokai",
			MarkdownTheme:      "dark",
			PromptStyle:        "default",
		},
		Plugins: PluginsConfig{
			Enabled:   []string{},
			Disabled:  []string{},
			Directory: filepath.Join(configDir, "plugins"),
			AutoLoad:  true,
		},
		AI: AIConfig{
			Provider:           core.AIProviderNone,
			MaxTokens:          2048,
			Temperature:        0.7,
			EnableSuggestions:  true,
			OllamaURL:          "http://localhost:11434",
			OllamaModel:        "llama2",
			EnableMonitoring:   false,
			MonitoringInterval: 30,
		},
		SSH: SSHConfig{
			DefaultKeyPath:    filepath.Join(homeDir, ".ssh", "id_rsa"),
			KnownHostsPath:    filepath.Join(homeDir, ".ssh", "known_hosts"),
			ConnectTimeout:    30,
			KeepAliveInterval: 60,
			SavedHosts:        []core.SSHHost{},
		},
		Secrets: SecretsConfig{
			StorePath:           filepath.Join(configDir, "secrets.enc"),
			EncryptionAlgorithm: "AES-256-GCM",
			KeyDerivation:       "argon2id",
		},
		Keybindings: KeybindingsConfig{
			Quit:            "ctrl+q",
			Help:            "ctrl+?",
			NewPane:         "ctrl+n",
			ClosePane:       "ctrl+w",
			NextPane:        "ctrl+tab",
			PrevPane:        "ctrl+shift+tab",
			SplitVertical:   "ctrl+\\",
			SplitHorizontal: "ctrl+-",
			ToggleSidebar:   "ctrl+b",
			CommandPalette:  "ctrl+p",
			AIAssist:        "ctrl+a",
		},
	}
}

// Load loads configuration from a file.
// If the file doesn't exist, it returns the default configuration.
// Returns an error if the file exists but cannot be read or parsed.
func Load(path string) (*Config, error) {
	cfg := Default()

	// Read configuration file
	data, err := os.ReadFile(path)
	if err != nil {
		// Return default config if file doesn't exist
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML configuration
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// LoadFromDefaultPath loads configuration from the default path.
// Returns default configuration if the default path cannot be determined
// or the file doesn't exist.
func LoadFromDefaultPath() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Default(), nil
	}

	configPath := filepath.Join(homeDir, ".cbwsh", "config.yaml")
	return Load(configPath)
}

// Save saves configuration to a file.
// Creates the directory if it doesn't exist.
func (c *Config) Save(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Ensure directory exists with secure permissions (owner only)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal configuration to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveToDefaultPath saves configuration to the default path.
func (c *Config) SaveToDefaultPath() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".cbwsh", "config.yaml")
	return c.Save(configPath)
}

// GetShellConfig returns the shell configuration.
func (c *Config) GetShellConfig() ShellConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Shell
}

// SetShellConfig sets the shell configuration.
func (c *Config) SetShellConfig(shell ShellConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Shell = shell
}

// GetUIConfig returns the UI configuration.
func (c *Config) GetUIConfig() UIConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.UI
}

// SetUIConfig sets the UI configuration.
func (c *Config) SetUIConfig(ui UIConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.UI = ui
}

// GetAIConfig returns the AI configuration.
func (c *Config) GetAIConfig() AIConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AI
}

// SetAIConfig sets the AI configuration.
func (c *Config) SetAIConfig(ai AIConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AI = ai
}
