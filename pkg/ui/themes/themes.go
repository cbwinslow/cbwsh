// Package themes provides theme management with hot-reloading for cbwsh.
package themes

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// Common errors.
var (
	ErrThemeNotFound = errors.New("theme not found")
	ErrInvalidTheme  = errors.New("invalid theme configuration")
)

// ColorScheme defines the colors for a theme.
type ColorScheme struct {
	// Background colors
	Background string `yaml:"background"`
	Foreground string `yaml:"foreground"`
	Selection  string `yaml:"selection"`
	Comment    string `yaml:"comment"`

	// Syntax colors
	Keyword  string `yaml:"keyword"`
	String   string `yaml:"string"`
	Number   string `yaml:"number"`
	Function string `yaml:"function"`
	Variable string `yaml:"variable"`
	Operator string `yaml:"operator"`
	Type     string `yaml:"type"`

	// UI colors
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary"`
	Success   string `yaml:"success"`
	Warning   string `yaml:"warning"`
	Error     string `yaml:"error"`
	Info      string `yaml:"info"`

	// Terminal colors (ANSI)
	Black         string `yaml:"black"`
	Red           string `yaml:"red"`
	Green         string `yaml:"green"`
	Yellow        string `yaml:"yellow"`
	Blue          string `yaml:"blue"`
	Magenta       string `yaml:"magenta"`
	Cyan          string `yaml:"cyan"`
	White         string `yaml:"white"`
	BrightBlack   string `yaml:"bright_black"`
	BrightRed     string `yaml:"bright_red"`
	BrightGreen   string `yaml:"bright_green"`
	BrightYellow  string `yaml:"bright_yellow"`
	BrightBlue    string `yaml:"bright_blue"`
	BrightMagenta string `yaml:"bright_magenta"`
	BrightCyan    string `yaml:"bright_cyan"`
	BrightWhite   string `yaml:"bright_white"`
}

// Theme represents a complete theme.
type Theme struct {
	// Name is the theme name.
	Name string `yaml:"name"`
	// Author is the theme author.
	Author string `yaml:"author"`
	// Description describes the theme.
	Description string `yaml:"description"`
	// Version is the theme version.
	Version string `yaml:"version"`
	// Colors is the color scheme.
	Colors ColorScheme `yaml:"colors"`
	// IsDark indicates if this is a dark theme.
	IsDark bool `yaml:"is_dark"`
}

// Manager manages themes with hot-reloading support.
type Manager struct {
	mu            sync.RWMutex
	themes        map[string]*Theme
	currentTheme  string
	themesDir     string
	watchInterval time.Duration
	watching      bool
	stopWatch     chan struct{}
	onChange      []func(*Theme)
}

// NewManager creates a new theme manager.
func NewManager(themesDir string) *Manager {
	m := &Manager{
		themes:        make(map[string]*Theme),
		themesDir:     themesDir,
		watchInterval: time.Second,
		stopWatch:     make(chan struct{}),
		onChange:      make([]func(*Theme), 0),
	}

	// Register built-in themes
	m.registerBuiltinThemes()

	return m
}

// Load loads a theme by name.
func (m *Manager) Load(name string) (*Theme, error) {
	m.mu.RLock()
	theme, ok := m.themes[name]
	m.mu.RUnlock()

	if ok {
		return theme, nil
	}

	// Try to load from file
	theme, err := m.loadFromFile(name)
	if err != nil {
		return nil, ErrThemeNotFound
	}

	m.mu.Lock()
	m.themes[name] = theme
	m.mu.Unlock()

	return theme, nil
}

// Current returns the current theme.
func (m *Manager) Current() *Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentTheme == "" {
		return m.themes["default"]
	}
	return m.themes[m.currentTheme]
}

// SetCurrent sets the current theme.
func (m *Manager) SetCurrent(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.themes[name]; !ok {
		return ErrThemeNotFound
	}

	m.currentTheme = name
	theme := m.themes[name]

	// Call change callbacks
	for _, cb := range m.onChange {
		cb(theme)
	}

	return nil
}

// List returns all available theme names.
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.themes))
	for name := range m.themes {
		names = append(names, name)
	}
	return names
}

// ListThemes returns all available themes.
func (m *Manager) ListThemes() []*Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()

	themes := make([]*Theme, 0, len(m.themes))
	for _, theme := range m.themes {
		themes = append(themes, theme)
	}
	return themes
}

// Register registers a theme.
func (m *Manager) Register(theme *Theme) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.themes[theme.Name] = theme
}

// Unregister removes a theme.
func (m *Manager) Unregister(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.themes, name)
}

// Reload reloads all themes from disk.
func (m *Manager) Reload() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Re-register built-in themes
	m.registerBuiltinThemes()

	// Load themes from directory
	if m.themesDir != "" {
		if err := m.loadThemesFromDir(); err != nil {
			return err
		}
	}

	// Notify of theme change
	if theme, ok := m.themes[m.currentTheme]; ok {
		for _, cb := range m.onChange {
			cb(theme)
		}
	}

	return nil
}

// Watch starts watching for theme file changes.
func (m *Manager) Watch() error {
	m.mu.Lock()
	if m.watching {
		m.mu.Unlock()
		return nil
	}
	m.watching = true
	m.stopWatch = make(chan struct{})
	m.mu.Unlock()

	go m.watchLoop()
	return nil
}

// StopWatch stops watching for theme changes.
func (m *Manager) StopWatch() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.watching {
		close(m.stopWatch)
		m.watching = false
	}
}

// OnChange registers a callback for theme changes.
func (m *Manager) OnChange(callback func(*Theme)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange = append(m.onChange, callback)
}

// SetWatchInterval sets the watch interval.
func (m *Manager) SetWatchInterval(interval time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.watchInterval = interval
}

func (m *Manager) watchLoop() {
	ticker := time.NewTicker(m.watchInterval)
	defer ticker.Stop()

	var lastModTime time.Time

	for {
		select {
		case <-m.stopWatch:
			return
		case <-ticker.C:
			// Check for theme file changes
			if m.themesDir != "" {
				info, err := os.Stat(m.themesDir)
				if err == nil && info.ModTime().After(lastModTime) {
					lastModTime = info.ModTime()
					_ = m.Reload()
				}
			}
		}
	}
}

func (m *Manager) loadFromFile(name string) (*Theme, error) {
	if m.themesDir == "" {
		return nil, ErrThemeNotFound
	}

	path := filepath.Join(m.themesDir, name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		// Try .yml extension
		path = filepath.Join(m.themesDir, name+".yml")
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	var theme Theme
	if err := yaml.Unmarshal(data, &theme); err != nil {
		return nil, ErrInvalidTheme
	}

	if theme.Name == "" {
		theme.Name = name
	}

	return &theme, nil
}

func (m *Manager) loadThemesFromDir() error {
	if m.themesDir == "" {
		return nil
	}

	entries, err := os.ReadDir(m.themesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		themeName := name[:len(name)-len(ext)]
		theme, err := m.loadFromFile(themeName)
		if err != nil {
			continue
		}

		m.themes[themeName] = theme
	}

	return nil
}

func (m *Manager) registerBuiltinThemes() {
	m.themes["default"] = DefaultTheme()
	m.themes["dracula"] = DraculaTheme()
	m.themes["nord"] = NordTheme()
	m.themes["solarized-dark"] = SolarizedDarkTheme()
	m.themes["solarized-light"] = SolarizedLightTheme()
	m.themes["monokai"] = MonokaiTheme()
	m.themes["one-dark"] = OneDarkTheme()
	m.themes["gruvbox"] = GruvboxTheme()
	m.themes["catppuccin"] = CatppuccinTheme()
}

// DefaultTheme returns the default theme.
func DefaultTheme() *Theme {
	return &Theme{
		Name:        "default",
		Description: "Default cbwsh theme",
		IsDark:      true,
		Colors: ColorScheme{
			Background:    "#1a1a2e",
			Foreground:    "#eaeaea",
			Selection:     "#3d3d5c",
			Comment:       "#6c6c9c",
			Keyword:       "#ff79c6",
			String:        "#50fa7b",
			Number:        "#bd93f9",
			Function:      "#8be9fd",
			Variable:      "#f8f8f2",
			Operator:      "#ff79c6",
			Type:          "#8be9fd",
			Primary:       "#bd93f9",
			Secondary:     "#6272a4",
			Success:       "#50fa7b",
			Warning:       "#ffb86c",
			Error:         "#ff5555",
			Info:          "#8be9fd",
			Black:         "#21222c",
			Red:           "#ff5555",
			Green:         "#50fa7b",
			Yellow:        "#f1fa8c",
			Blue:          "#bd93f9",
			Magenta:       "#ff79c6",
			Cyan:          "#8be9fd",
			White:         "#f8f8f2",
			BrightBlack:   "#6272a4",
			BrightRed:     "#ff6e6e",
			BrightGreen:   "#69ff94",
			BrightYellow:  "#ffffa5",
			BrightBlue:    "#d6acff",
			BrightMagenta: "#ff92df",
			BrightCyan:    "#a4ffff",
			BrightWhite:   "#ffffff",
		},
	}
}

// DraculaTheme returns the Dracula theme.
func DraculaTheme() *Theme {
	return &Theme{
		Name:        "dracula",
		Description: "Dracula color scheme",
		Author:      "Zeno Rocha",
		IsDark:      true,
		Colors: ColorScheme{
			Background:    "#282a36",
			Foreground:    "#f8f8f2",
			Selection:     "#44475a",
			Comment:       "#6272a4",
			Keyword:       "#ff79c6",
			String:        "#f1fa8c",
			Number:        "#bd93f9",
			Function:      "#50fa7b",
			Variable:      "#f8f8f2",
			Operator:      "#ff79c6",
			Type:          "#8be9fd",
			Primary:       "#bd93f9",
			Secondary:     "#6272a4",
			Success:       "#50fa7b",
			Warning:       "#ffb86c",
			Error:         "#ff5555",
			Info:          "#8be9fd",
			Black:         "#21222c",
			Red:           "#ff5555",
			Green:         "#50fa7b",
			Yellow:        "#f1fa8c",
			Blue:          "#bd93f9",
			Magenta:       "#ff79c6",
			Cyan:          "#8be9fd",
			White:         "#f8f8f2",
			BrightBlack:   "#6272a4",
			BrightRed:     "#ff6e6e",
			BrightGreen:   "#69ff94",
			BrightYellow:  "#ffffa5",
			BrightBlue:    "#d6acff",
			BrightMagenta: "#ff92df",
			BrightCyan:    "#a4ffff",
			BrightWhite:   "#ffffff",
		},
	}
}

// NordTheme returns the Nord theme.
func NordTheme() *Theme {
	return &Theme{
		Name:        "nord",
		Description: "Arctic, north-bluish color palette",
		IsDark:      true,
		Colors: ColorScheme{
			Background:    "#2e3440",
			Foreground:    "#eceff4",
			Selection:     "#434c5e",
			Comment:       "#616e88",
			Keyword:       "#81a1c1",
			String:        "#a3be8c",
			Number:        "#b48ead",
			Function:      "#88c0d0",
			Variable:      "#eceff4",
			Operator:      "#81a1c1",
			Type:          "#8fbcbb",
			Primary:       "#88c0d0",
			Secondary:     "#81a1c1",
			Success:       "#a3be8c",
			Warning:       "#ebcb8b",
			Error:         "#bf616a",
			Info:          "#88c0d0",
			Black:         "#3b4252",
			Red:           "#bf616a",
			Green:         "#a3be8c",
			Yellow:        "#ebcb8b",
			Blue:          "#81a1c1",
			Magenta:       "#b48ead",
			Cyan:          "#88c0d0",
			White:         "#e5e9f0",
			BrightBlack:   "#4c566a",
			BrightRed:     "#bf616a",
			BrightGreen:   "#a3be8c",
			BrightYellow:  "#ebcb8b",
			BrightBlue:    "#81a1c1",
			BrightMagenta: "#b48ead",
			BrightCyan:    "#8fbcbb",
			BrightWhite:   "#eceff4",
		},
	}
}

// SolarizedDarkTheme returns the Solarized Dark theme.
func SolarizedDarkTheme() *Theme {
	return &Theme{
		Name:        "solarized-dark",
		Description: "Solarized Dark color scheme",
		Author:      "Ethan Schoonover",
		IsDark:      true,
		Colors: ColorScheme{
			Background: "#002b36",
			Foreground: "#839496",
			Selection:  "#073642",
			Comment:    "#586e75",
			Keyword:    "#859900",
			String:     "#2aa198",
			Number:     "#d33682",
			Function:   "#268bd2",
			Variable:   "#839496",
			Operator:   "#859900",
			Type:       "#b58900",
			Primary:    "#268bd2",
			Secondary:  "#586e75",
			Success:    "#859900",
			Warning:    "#b58900",
			Error:      "#dc322f",
			Info:       "#2aa198",
			Black:      "#073642",
			Red:        "#dc322f",
			Green:      "#859900",
			Yellow:     "#b58900",
			Blue:       "#268bd2",
			Magenta:    "#d33682",
			Cyan:       "#2aa198",
			White:      "#eee8d5",
		},
	}
}

// SolarizedLightTheme returns the Solarized Light theme.
func SolarizedLightTheme() *Theme {
	return &Theme{
		Name:        "solarized-light",
		Description: "Solarized Light color scheme",
		Author:      "Ethan Schoonover",
		IsDark:      false,
		Colors: ColorScheme{
			Background: "#fdf6e3",
			Foreground: "#657b83",
			Selection:  "#eee8d5",
			Comment:    "#93a1a1",
			Keyword:    "#859900",
			String:     "#2aa198",
			Number:     "#d33682",
			Function:   "#268bd2",
			Variable:   "#657b83",
			Operator:   "#859900",
			Type:       "#b58900",
			Primary:    "#268bd2",
			Secondary:  "#93a1a1",
			Success:    "#859900",
			Warning:    "#b58900",
			Error:      "#dc322f",
			Info:       "#2aa198",
			Black:      "#002b36",
			Red:        "#dc322f",
			Green:      "#859900",
			Yellow:     "#b58900",
			Blue:       "#268bd2",
			Magenta:    "#d33682",
			Cyan:       "#2aa198",
			White:      "#fdf6e3",
		},
	}
}

// MonokaiTheme returns the Monokai theme.
func MonokaiTheme() *Theme {
	return &Theme{
		Name:        "monokai",
		Description: "Monokai color scheme",
		IsDark:      true,
		Colors: ColorScheme{
			Background: "#272822",
			Foreground: "#f8f8f2",
			Selection:  "#49483e",
			Comment:    "#75715e",
			Keyword:    "#f92672",
			String:     "#e6db74",
			Number:     "#ae81ff",
			Function:   "#a6e22e",
			Variable:   "#f8f8f2",
			Operator:   "#f92672",
			Type:       "#66d9ef",
			Primary:    "#a6e22e",
			Secondary:  "#75715e",
			Success:    "#a6e22e",
			Warning:    "#fd971f",
			Error:      "#f92672",
			Info:       "#66d9ef",
			Black:      "#272822",
			Red:        "#f92672",
			Green:      "#a6e22e",
			Yellow:     "#f4bf75",
			Blue:       "#66d9ef",
			Magenta:    "#ae81ff",
			Cyan:       "#a1efe4",
			White:      "#f8f8f2",
		},
	}
}

// OneDarkTheme returns the One Dark theme.
func OneDarkTheme() *Theme {
	return &Theme{
		Name:        "one-dark",
		Description: "Atom One Dark color scheme",
		IsDark:      true,
		Colors: ColorScheme{
			Background: "#282c34",
			Foreground: "#abb2bf",
			Selection:  "#3e4451",
			Comment:    "#5c6370",
			Keyword:    "#c678dd",
			String:     "#98c379",
			Number:     "#d19a66",
			Function:   "#61afef",
			Variable:   "#e06c75",
			Operator:   "#56b6c2",
			Type:       "#e5c07b",
			Primary:    "#61afef",
			Secondary:  "#5c6370",
			Success:    "#98c379",
			Warning:    "#e5c07b",
			Error:      "#e06c75",
			Info:       "#61afef",
			Black:      "#282c34",
			Red:        "#e06c75",
			Green:      "#98c379",
			Yellow:     "#e5c07b",
			Blue:       "#61afef",
			Magenta:    "#c678dd",
			Cyan:       "#56b6c2",
			White:      "#abb2bf",
		},
	}
}

// GruvboxTheme returns the Gruvbox theme.
func GruvboxTheme() *Theme {
	return &Theme{
		Name:        "gruvbox",
		Description: "Gruvbox color scheme",
		IsDark:      true,
		Colors: ColorScheme{
			Background: "#282828",
			Foreground: "#ebdbb2",
			Selection:  "#504945",
			Comment:    "#928374",
			Keyword:    "#fb4934",
			String:     "#b8bb26",
			Number:     "#d3869b",
			Function:   "#fabd2f",
			Variable:   "#83a598",
			Operator:   "#fe8019",
			Type:       "#8ec07c",
			Primary:    "#fabd2f",
			Secondary:  "#928374",
			Success:    "#b8bb26",
			Warning:    "#fe8019",
			Error:      "#fb4934",
			Info:       "#83a598",
			Black:      "#282828",
			Red:        "#cc241d",
			Green:      "#98971a",
			Yellow:     "#d79921",
			Blue:       "#458588",
			Magenta:    "#b16286",
			Cyan:       "#689d6a",
			White:      "#a89984",
		},
	}
}

// CatppuccinTheme returns the Catppuccin theme.
func CatppuccinTheme() *Theme {
	return &Theme{
		Name:        "catppuccin",
		Description: "Catppuccin Mocha color scheme",
		IsDark:      true,
		Colors: ColorScheme{
			Background: "#1e1e2e",
			Foreground: "#cdd6f4",
			Selection:  "#45475a",
			Comment:    "#6c7086",
			Keyword:    "#cba6f7",
			String:     "#a6e3a1",
			Number:     "#fab387",
			Function:   "#89b4fa",
			Variable:   "#f5e0dc",
			Operator:   "#89dceb",
			Type:       "#f9e2af",
			Primary:    "#89b4fa",
			Secondary:  "#6c7086",
			Success:    "#a6e3a1",
			Warning:    "#f9e2af",
			Error:      "#f38ba8",
			Info:       "#89dceb",
			Black:      "#45475a",
			Red:        "#f38ba8",
			Green:      "#a6e3a1",
			Yellow:     "#f9e2af",
			Blue:       "#89b4fa",
			Magenta:    "#f5c2e7",
			Cyan:       "#94e2d5",
			White:      "#bac2de",
		},
	}
}

// ToLipglossColor converts a color string to lipgloss.Color.
func ToLipglossColor(color string) lipgloss.Color {
	return lipgloss.Color(color)
}

// ApplyTheme applies a theme's colors to create lipgloss styles.
func ApplyTheme(theme *Theme) map[string]lipgloss.Style {
	colors := theme.Colors

	return map[string]lipgloss.Style{
		"background": lipgloss.NewStyle().Background(ToLipglossColor(colors.Background)),
		"foreground": lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Foreground)),
		"primary":    lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Primary)),
		"secondary":  lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Secondary)),
		"success":    lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Success)),
		"warning":    lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Warning)),
		"error":      lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Error)),
		"info":       lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Info)),
		"keyword":    lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Keyword)),
		"string":     lipgloss.NewStyle().Foreground(ToLipglossColor(colors.String)),
		"number":     lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Number)),
		"function":   lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Function)),
		"comment":    lipgloss.NewStyle().Foreground(ToLipglossColor(colors.Comment)),
	}
}
