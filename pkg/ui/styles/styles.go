// Package styles provides UI styling for cbwsh.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a complete color theme.
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	Muted      lipgloss.Color
	Error      lipgloss.Color
	Warning    lipgloss.Color
	Success    lipgloss.Color
	Info       lipgloss.Color
	Border     lipgloss.Color
	Selection  lipgloss.Color
}

// DefaultTheme is the default color theme.
var DefaultTheme = Theme{
	Name:       "default",
	Primary:    lipgloss.Color("86"),  // Cyan
	Secondary:  lipgloss.Color("141"), // Purple
	Accent:     lipgloss.Color("214"), // Orange
	Background: lipgloss.Color("235"), // Dark gray
	Foreground: lipgloss.Color("252"), // Light gray
	Muted:      lipgloss.Color("240"), // Gray
	Error:      lipgloss.Color("196"), // Red
	Warning:    lipgloss.Color("214"), // Orange
	Success:    lipgloss.Color("82"),  // Green
	Info:       lipgloss.Color("39"),  // Blue
	Border:     lipgloss.Color("238"), // Dark gray
	Selection:  lipgloss.Color("236"), // Darker gray
}

// DraculaTheme is the Dracula color theme.
var DraculaTheme = Theme{
	Name:       "dracula",
	Primary:    lipgloss.Color("141"), // Purple
	Secondary:  lipgloss.Color("117"), // Cyan
	Accent:     lipgloss.Color("212"), // Pink
	Background: lipgloss.Color("236"), // Dark
	Foreground: lipgloss.Color("253"), // Light
	Muted:      lipgloss.Color("103"), // Comment gray
	Error:      lipgloss.Color("203"), // Red
	Warning:    lipgloss.Color("215"), // Orange
	Success:    lipgloss.Color("84"),  // Green
	Info:       lipgloss.Color("117"), // Cyan
	Border:     lipgloss.Color("60"),  // Purple border
	Selection:  lipgloss.Color("60"),  // Purple selection
}

// NordTheme is the Nord color theme.
var NordTheme = Theme{
	Name:       "nord",
	Primary:    lipgloss.Color("110"), // Frost blue
	Secondary:  lipgloss.Color("143"), // Aurora green
	Accent:     lipgloss.Color("180"), // Aurora purple
	Background: lipgloss.Color("236"), // Polar night
	Foreground: lipgloss.Color("254"), // Snow storm
	Muted:      lipgloss.Color("103"), // Polar night light
	Error:      lipgloss.Color("167"), // Aurora red
	Warning:    lipgloss.Color("179"), // Aurora yellow
	Success:    lipgloss.Color("143"), // Aurora green
	Info:       lipgloss.Color("110"), // Frost blue
	Border:     lipgloss.Color("60"),  // Polar night border
	Selection:  lipgloss.Color("60"),  // Polar night selection
}

// Themes is a map of available themes.
var Themes = map[string]Theme{
	"default": DefaultTheme,
	"dracula": DraculaTheme,
	"nord":    NordTheme,
}

// GetTheme returns a theme by name.
func GetTheme(name string) Theme {
	if theme, exists := Themes[name]; exists {
		return theme
	}
	return DefaultTheme
}

// Styles provides styled components.
type Styles struct {
	theme Theme

	// Base styles
	App       lipgloss.Style
	Header    lipgloss.Style
	Footer    lipgloss.Style
	StatusBar lipgloss.Style

	// Pane styles
	Pane       lipgloss.Style
	ActivePane lipgloss.Style
	PaneTitle  lipgloss.Style

	// Prompt styles
	Prompt       lipgloss.Style
	PromptSymbol lipgloss.Style
	Input        lipgloss.Style

	// Output styles
	Output  lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Success lipgloss.Style
	Info    lipgloss.Style

	// Command styles
	Command  lipgloss.Style
	Option   lipgloss.Style
	Argument lipgloss.Style
	Variable lipgloss.Style

	// Autocomplete styles
	Suggestion         lipgloss.Style
	SelectedSuggestion lipgloss.Style
	SuggestionDesc     lipgloss.Style

	// Misc styles
	Spinner   lipgloss.Style
	Progress  lipgloss.Style
	Highlight lipgloss.Style
	Muted     lipgloss.Style
	Bold      lipgloss.Style
	Italic    lipgloss.Style
	Underline lipgloss.Style
}

// NewStyles creates new styles with the given theme.
func NewStyles(theme Theme) *Styles {
	borderStyle := lipgloss.RoundedBorder()

	return &Styles{
		theme: theme,

		// Base styles
		App: lipgloss.NewStyle().
			Background(theme.Background),

		Header: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			Padding(0, 1),

		Footer: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Padding(0, 1),

		StatusBar: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Background(theme.Selection).
			Padding(0, 1),

		// Pane styles
		Pane: lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(theme.Border).
			Padding(0, 1),

		ActivePane: lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(theme.Primary).
			Padding(0, 1),

		PaneTitle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		// Prompt styles
		Prompt: lipgloss.NewStyle().
			Foreground(theme.Primary),

		PromptSymbol: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true),

		Input: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		// Output styles
		Output: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error),

		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success),

		Info: lipgloss.NewStyle().
			Foreground(theme.Info),

		// Command styles
		Command: lipgloss.NewStyle().
			Foreground(theme.Primary),

		Option: lipgloss.NewStyle().
			Foreground(theme.Accent),

		Argument: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		Variable: lipgloss.NewStyle().
			Foreground(theme.Secondary),

		// Autocomplete styles
		Suggestion: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Padding(0, 1),

		SelectedSuggestion: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Background(theme.Selection).
			Padding(0, 1),

		SuggestionDesc: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Italic(true),

		// Misc styles
		Spinner: lipgloss.NewStyle().
			Foreground(theme.Accent),

		Progress: lipgloss.NewStyle().
			Foreground(theme.Primary),

		Highlight: lipgloss.NewStyle().
			Background(theme.Selection),

		Muted: lipgloss.NewStyle().
			Foreground(theme.Muted),

		Bold: lipgloss.NewStyle().
			Bold(true),

		Italic: lipgloss.NewStyle().
			Italic(true),

		Underline: lipgloss.NewStyle().
			Underline(true),
	}
}

// DefaultStyles returns styles with the default theme.
func DefaultStyles() *Styles {
	return NewStyles(DefaultTheme)
}

// SetTheme changes the theme.
func (s *Styles) SetTheme(theme Theme) {
	*s = *NewStyles(theme)
}

// GetTheme returns the current theme.
func (s *Styles) GetTheme() Theme {
	return s.theme
}

// Box creates a box around content.
func (s *Styles) Box(content string, title string, active bool) string {
	style := s.Pane
	if active {
		style = s.ActivePane
	}

	if title != "" {
		return style.
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			Render(s.PaneTitle.Render(title) + "\n" + content)
	}

	return style.Render(content)
}

// PromptLine creates a styled prompt line.
func (s *Styles) PromptLine(cwd, symbol, input string) string {
	return s.Prompt.Render(cwd) + " " +
		s.PromptSymbol.Render(symbol) + " " +
		s.Input.Render(input)
}

// FormatError formats an error message.
func (s *Styles) FormatError(msg string) string {
	return s.Error.Render("✗ " + msg)
}

// FormatWarning formats a warning message.
func (s *Styles) FormatWarning(msg string) string {
	return s.Warning.Render("⚠ " + msg)
}

// FormatSuccess formats a success message.
func (s *Styles) FormatSuccess(msg string) string {
	return s.Success.Render("✓ " + msg)
}

// FormatInfo formats an info message.
func (s *Styles) FormatInfo(msg string) string {
	return s.Info.Render("ℹ " + msg)
}
