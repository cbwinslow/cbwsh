// Package markdown provides markdown rendering for cbwsh.
package markdown

import (
	"sync"

	"github.com/charmbracelet/glamour"
)

// Renderer provides styled markdown rendering.
type Renderer struct {
	mu       sync.RWMutex
	renderer *glamour.TermRenderer
	width    int
	theme    string
}

// Theme represents a markdown rendering theme.
type Theme string

const (
	// ThemeDark is the dark theme.
	ThemeDark Theme = "dark"
	// ThemeLight is the light theme.
	ThemeLight Theme = "light"
	// ThemeDracula is the dracula theme.
	ThemeDracula Theme = "dracula"
	// ThemeTokyoNight is the tokyo-night theme.
	ThemeTokyoNight Theme = "tokyo-night"
	// ThemeNotty is a plain theme without colors.
	ThemeNotty Theme = "notty"
	// ThemeAuto automatically detects the theme.
	ThemeAuto Theme = "auto"
)

// NewRenderer creates a new markdown renderer.
func NewRenderer() (*Renderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		renderer: r,
		width:    80,
		theme:    string(ThemeAuto),
	}, nil
}

// NewRendererWithTheme creates a new markdown renderer with a specific theme.
func NewRendererWithTheme(theme Theme) (*Renderer, error) {
	var opts []glamour.TermRendererOption

	switch theme {
	case ThemeDark:
		opts = append(opts, glamour.WithStylePath("dark"))
	case ThemeLight:
		opts = append(opts, glamour.WithStylePath("light"))
	case ThemeDracula:
		opts = append(opts, glamour.WithStylePath("dracula"))
	case ThemeTokyoNight:
		opts = append(opts, glamour.WithStylePath("tokyo-night"))
	case ThemeNotty:
		opts = append(opts, glamour.WithStylePath("notty"))
	default:
		opts = append(opts, glamour.WithAutoStyle())
	}

	opts = append(opts, glamour.WithWordWrap(80))

	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		renderer: r,
		width:    80,
		theme:    string(theme),
	}, nil
}

// Render renders markdown to styled terminal output.
func (r *Renderer) Render(markdown string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.renderer.Render(markdown)
}

// SetWidth sets the render width.
func (r *Renderer) SetWidth(width int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.width = width
	// Recreate renderer with new width
	newRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err == nil {
		r.renderer = newRenderer
	}
}

// SetTheme sets the rendering theme.
func (r *Renderer) SetTheme(theme string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var opts []glamour.TermRendererOption

	switch Theme(theme) {
	case ThemeDark:
		opts = append(opts, glamour.WithStylePath("dark"))
	case ThemeLight:
		opts = append(opts, glamour.WithStylePath("light"))
	case ThemeDracula:
		opts = append(opts, glamour.WithStylePath("dracula"))
	case ThemeTokyoNight:
		opts = append(opts, glamour.WithStylePath("tokyo-night"))
	case ThemeNotty:
		opts = append(opts, glamour.WithStylePath("notty"))
	default:
		opts = append(opts, glamour.WithAutoStyle())
	}

	opts = append(opts, glamour.WithWordWrap(r.width))

	newRenderer, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return err
	}

	r.renderer = newRenderer
	r.theme = theme
	return nil
}

// GetTheme returns the current theme.
func (r *Renderer) GetTheme() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.theme
}

// GetWidth returns the current width.
func (r *Renderer) GetWidth() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.width
}

// RenderCode renders a code block with syntax highlighting.
func (r *Renderer) RenderCode(code, language string) (string, error) {
	markdown := "```" + language + "\n" + code + "\n```"
	return r.Render(markdown)
}

// RenderTable renders a table.
func (r *Renderer) RenderTable(headers []string, rows [][]string) (string, error) {
	var markdown string

	// Headers
	markdown += "| "
	for _, h := range headers {
		markdown += h + " | "
	}
	markdown += "\n"

	// Separator
	markdown += "| "
	for range headers {
		markdown += "--- | "
	}
	markdown += "\n"

	// Rows
	for _, row := range rows {
		markdown += "| "
		for _, cell := range row {
			markdown += cell + " | "
		}
		markdown += "\n"
	}

	return r.Render(markdown)
}

// RenderList renders a list.
func (r *Renderer) RenderList(items []string, ordered bool) (string, error) {
	var markdown string

	for i, item := range items {
		if ordered {
			markdown += string(rune('1'+i)) + ". " + item + "\n"
		} else {
			markdown += "- " + item + "\n"
		}
	}

	return r.Render(markdown)
}

// RenderHeading renders a heading.
func (r *Renderer) RenderHeading(text string, level int) (string, error) {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}

	hashes := ""
	for i := 0; i < level; i++ {
		hashes += "#"
	}

	markdown := hashes + " " + text + "\n"
	return r.Render(markdown)
}

// RenderBlockquote renders a blockquote.
func (r *Renderer) RenderBlockquote(text string) (string, error) {
	markdown := "> " + text + "\n"
	return r.Render(markdown)
}

// RenderHorizontalRule renders a horizontal rule.
func (r *Renderer) RenderHorizontalRule() (string, error) {
	return r.Render("---\n")
}

// RenderLink renders a link.
func (r *Renderer) RenderLink(text, url string) (string, error) {
	markdown := "[" + text + "](" + url + ")"
	return r.Render(markdown)
}

// RenderBold renders bold text.
func (r *Renderer) RenderBold(text string) (string, error) {
	return r.Render("**" + text + "**")
}

// RenderItalic renders italic text.
func (r *Renderer) RenderItalic(text string) (string, error) {
	return r.Render("*" + text + "*")
}

// RenderInlineCode renders inline code.
func (r *Renderer) RenderInlineCode(code string) (string, error) {
	return r.Render("`" + code + "`")
}
