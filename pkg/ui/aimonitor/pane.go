// Package aimonitor provides a dedicated AI monitoring pane for shell activity.
package aimonitor

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cbwinslow/cbwsh/pkg/ai/monitor"
	"github.com/cbwinslow/cbwsh/pkg/ui/markdown"
)

// MonitorPane is a dedicated pane for displaying AI monitoring and recommendations.
type MonitorPane struct {
	mu sync.RWMutex

	// UI components
	viewport   viewport.Model
	mdRenderer *markdown.Renderer

	// State
	monitor  *monitor.Monitor
	visible  bool
	focused  bool
	width    int
	height   int
	position string // "right", "bottom"

	// Styles
	borderStyle     lipgloss.Style
	titleStyle      lipgloss.Style
	infoStyle       lipgloss.Style
	warningStyle    lipgloss.Style
	tipStyle        lipgloss.Style
	suggestionStyle lipgloss.Style
	timestampStyle  lipgloss.Style
	activityStyle   lipgloss.Style
}

// KeyMap defines key bindings for the monitor pane.
type KeyMap struct {
	Toggle     key.Binding
	Clear      key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Refresh    key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Toggle: key.NewBinding(
			key.WithKeys("ctrl+m"),
			key.WithHelp("ctrl+m", "toggle monitor"),
		),
		Clear: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "clear"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("ctrl+up", "pageup"),
			key.WithHelp("pgup", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("ctrl+down", "pagedown"),
			key.WithHelp("pgdn", "scroll down"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "refresh"),
		),
	}
}

var keys = DefaultKeyMap()

// NewMonitorPane creates a new AI monitor pane.
func NewMonitorPane(mon *monitor.Monitor) *MonitorPane {
	vp := viewport.New(40, 20)
	vp.Style = lipgloss.NewStyle()

	mdRenderer, _ := markdown.NewRenderer()

	pane := &MonitorPane{
		viewport:   vp,
		mdRenderer: mdRenderer,
		monitor:    mon,
		visible:    false,
		position:   "right",
		borderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("141")),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("141")).
			Bold(true).
			Padding(0, 1),
		infoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")),
		warningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true),
		tipStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")),
		suggestionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("141")),
		timestampStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
		activityStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("248")),
	}

	// Set callback for new recommendations
	if mon != nil {
		mon.SetOnRecommendation(func(rec monitor.Recommendation) {
			pane.updateContent()
		})
	}

	return pane
}

// Show shows the monitor pane.
func (p *MonitorPane) Show() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.visible = true
	p.updateContent()
}

// Hide hides the monitor pane.
func (p *MonitorPane) Hide() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.visible = false
}

// Toggle toggles the monitor pane visibility.
func (p *MonitorPane) Toggle() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.visible = !p.visible
	if p.visible {
		p.updateContent()
	}
}

// IsVisible returns whether the pane is visible.
func (p *MonitorPane) IsVisible() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.visible
}

// Focus focuses the monitor pane.
func (p *MonitorPane) Focus() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.focused = true
}

// Blur removes focus from the monitor pane.
func (p *MonitorPane) Blur() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.focused = false
}

// IsFocused returns whether the pane is focused.
func (p *MonitorPane) IsFocused() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.focused
}

// SetSize sets the size of the monitor pane.
func (p *MonitorPane) SetSize(width, height int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.width = width
	p.height = height

	// Update viewport size
	contentWidth := width - 4    // Account for borders
	viewportHeight := height - 4 // Account for title and borders

	if viewportHeight < 3 {
		viewportHeight = 3
	}

	p.viewport.Width = contentWidth
	p.viewport.Height = viewportHeight
}

// SetPosition sets the position of the monitor pane.
func (p *MonitorPane) SetPosition(position string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.position = position
}

// Clear clears all recommendations.
func (p *MonitorPane) Clear() {
	if p.monitor != nil {
		p.monitor.ClearRecommendations()
	}
	p.updateContent()
}

// updateContent updates the viewport content with recommendations.
func (p *MonitorPane) updateContent() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.monitor == nil {
		p.viewport.SetContent("AI Monitor not configured")
		return
	}

	recommendations := p.monitor.GetRecommendations()

	var content strings.Builder

	if len(recommendations) == 0 {
		content.WriteString(p.infoStyle.Render("No recommendations yet.\n\n"))
		content.WriteString(p.activityStyle.Render("The AI will analyze your shell activity and provide suggestions."))
	} else {
		// Show recommendations in reverse order (newest first)
		for i := len(recommendations) - 1; i >= 0; i-- {
			rec := recommendations[i]

			// Recommendation type indicator
			var typeStyle lipgloss.Style
			var typeIcon string
			switch rec.Type {
			case "warning":
				typeStyle = p.warningStyle
				typeIcon = "⚠️ "
			case "tip":
				typeStyle = p.tipStyle
				typeIcon = "💡 "
			case "suggestion":
				typeStyle = p.suggestionStyle
				typeIcon = "✨ "
			default:
				typeStyle = p.infoStyle
				typeIcon = "ℹ️ "
			}

			// Header with timestamp
			timestamp := rec.Timestamp.Format("15:04:05")
			header := typeStyle.Render(typeIcon+rec.Title) + " " +
				p.timestampStyle.Render(timestamp)
			content.WriteString(header)
			content.WriteString("\n")

			// Message
			if p.mdRenderer != nil {
				rendered, err := p.mdRenderer.Render(rec.Message)
				if err == nil {
					content.WriteString(rendered)
				} else {
					content.WriteString(rec.Message)
				}
			} else {
				content.WriteString(rec.Message)
			}
			content.WriteString("\n")

			// Related activity context (if available)
			if rec.Activity != nil {
				activityInfo := p.activityStyle.Render(fmt.Sprintf("↳ %s", rec.Activity.Command))
				content.WriteString(activityInfo)
				content.WriteString("\n")
			}

			content.WriteString("\n")
			if i > 0 {
				content.WriteString(strings.Repeat("─", p.viewport.Width))
				content.WriteString("\n\n")
			}
		}
	}

	p.viewport.SetContent(content.String())
	p.viewport.GotoBottom()
}

// Update handles messages for the monitor pane.
func (p *MonitorPane) Update(msg tea.Msg) (*MonitorPane, tea.Cmd) {
	p.mu.RLock()
	if !p.visible {
		p.mu.RUnlock()
		return p, nil
	}
	focused := p.focused
	p.mu.RUnlock()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !focused {
			return p, nil
		}

		switch {
		case key.Matches(msg, keys.Clear):
			p.Clear()
			return p, nil

		case key.Matches(msg, keys.Refresh):
			p.updateContent()
			return p, nil

		case key.Matches(msg, keys.ScrollUp):
			p.mu.Lock()
			for i := 0; i < 3; i++ {
				p.viewport, _ = p.viewport.Update(tea.KeyMsg{Type: tea.KeyUp})
			}
			p.mu.Unlock()
			return p, nil

		case key.Matches(msg, keys.ScrollDown):
			p.mu.Lock()
			for i := 0; i < 3; i++ {
				p.viewport, _ = p.viewport.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			p.mu.Unlock()
			return p, nil
		}

	case tea.WindowSizeMsg:
		// Size will be set externally by the main app
		return p, nil

	case recommendationMsg:
		// New recommendation received
		p.updateContent()
		return p, nil
	}

	// Update viewport
	p.mu.Lock()
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	p.mu.Unlock()

	return p, cmd
}

// recommendationMsg signals a new recommendation.
type recommendationMsg struct{}

// View renders the monitor pane.
func (p *MonitorPane) View() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.visible {
		return ""
	}

	var content strings.Builder

	// Title bar
	monitorStatus := "●"
	if p.monitor != nil && p.monitor.IsEnabled() {
		monitorStatus = p.infoStyle.Render("●") // Active
	} else {
		monitorStatus = p.timestampStyle.Render("○") // Inactive
	}

	title := p.titleStyle.Render("🤖 AI Monitor ") + monitorStatus
	content.WriteString(title)
	content.WriteString("\n")

	// Viewport content
	content.WriteString(p.viewport.View())

	// Help text at bottom
	if p.focused {
		helpText := p.timestampStyle.Render("ctrl+l: clear | ctrl+r: refresh | pgup/pgdn: scroll")
		content.WriteString("\n")
		content.WriteString(helpText)
	}

	// Wrap in border
	borderStyle := p.borderStyle
	if p.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("86"))
	}

	return borderStyle.
		Width(p.width - 2).
		Height(p.height - 2).
		Render(content.String())
}

// GetWidth returns the current width.
func (p *MonitorPane) GetWidth() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.width
}

// GetHeight returns the current height.
func (p *MonitorPane) GetHeight() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.height
}

// SetMonitor sets the activity monitor.
func (p *MonitorPane) SetMonitor(mon *monitor.Monitor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.monitor = mon

	// Set callback for new recommendations
	if mon != nil {
		mon.SetOnRecommendation(func(rec monitor.Recommendation) {
			p.updateContent()
		})
	}
}

// Tick returns a command that updates the pane periodically.
func Tick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return recommendationMsg{}
	})
}
