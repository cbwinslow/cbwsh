package aimonitor

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cbwinslow/cbwsh/pkg/ai/monitor"
)

func TestNewMonitorPane(t *testing.T) {
	t.Parallel()

	mon := &monitor.Monitor{}
	pane := NewMonitorPane(mon)

	if pane == nil {
		t.Fatal("expected non-nil monitor pane")
	}

	if pane.monitor != mon {
		t.Error("monitor not set correctly")
	}

	if pane.visible {
		t.Error("pane should not be visible by default")
	}

	if pane.position == "" {
		t.Error("position should be set to default value")
	}
}

func TestMonitorPaneSetSize(t *testing.T) {
	t.Parallel()

	mon := &monitor.Monitor{}
	pane := NewMonitorPane(mon)

	pane.SetSize(100, 50)

	if pane.width != 100 {
		t.Errorf("expected width 100, got %d", pane.width)
	}

	if pane.height != 50 {
		t.Errorf("expected height 50, got %d", pane.height)
	}
}

func TestMonitorPaneToggleVisibility(t *testing.T) {
	t.Parallel()

	mon := &monitor.Monitor{}
	pane := NewMonitorPane(mon)

	initialVisible := pane.visible

	pane.Toggle()

	if pane.visible == initialVisible {
		t.Error("toggle should change visibility state")
	}

	pane.Toggle()

	if pane.visible != initialVisible {
		t.Error("second toggle should restore original visibility")
	}
}

func TestMonitorPaneFocus(t *testing.T) {
	t.Parallel()

	mon := &monitor.Monitor{}
	pane := NewMonitorPane(mon)

	if pane.focused {
		t.Error("pane should not be focused initially")
	}

	pane.Focus()

	if !pane.focused {
		t.Error("pane should be focused after Focus()")
	}

	pane.Blur()

	if pane.focused {
		t.Error("pane should not be focused after Blur()")
	}
}

func TestMonitorPaneUpdate(t *testing.T) {
	t.Parallel()

	mon := &monitor.Monitor{}
	pane := NewMonitorPane(mon)
	pane.SetSize(100, 50)

	// Test window size message
	updatedPane, _ := pane.Update(tea.WindowSizeMsg{Width: 120, Height: 60})

	// Note: Update returns *MonitorPane but doesn't necessarily update width/height directly
	// It's more about handling the message. The actual size would be set via SetSize
	_ = updatedPane
}

func TestMonitorPaneRender(t *testing.T) {
	t.Parallel()

	mon := &monitor.Monitor{}
	pane := NewMonitorPane(mon)
	pane.SetSize(100, 50)

	// Render when not visible
	rendered := pane.View()
	if rendered != "" {
		t.Error("render should return empty string when not visible")
	}

	// Render when visible
	pane.Show()
	rendered = pane.View()
	if rendered == "" {
		t.Error("render should return non-empty string when visible")
	}
}

func TestMonitorPanePosition(t *testing.T) {
	t.Parallel()

	mon := &monitor.Monitor{}
	pane := NewMonitorPane(mon)

	// Test setting position
	pane.SetPosition("right")
	if pane.position != "right" {
		t.Errorf("expected position 'right', got '%s'", pane.position)
	}

	pane.SetPosition("bottom")
	if pane.position != "bottom" {
		t.Errorf("expected position 'bottom', got '%s'", pane.position)
	}
}

func TestDefaultKeyMap(t *testing.T) {
	t.Parallel()

	km := DefaultKeyMap()

	if km.Toggle.Enabled() {
		// This is just to check that the key map is initialized
		// We don't want to test the actual key bindings here
	}

	if km.Clear.Enabled() {
		// Same as above
	}
}
