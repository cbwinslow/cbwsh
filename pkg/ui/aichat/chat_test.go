package aichat

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cbwinslow/cbwsh/pkg/ai"
)

func TestNewMessage(t *testing.T) {
	t.Parallel()

	pane := NewChatPane(&ai.Manager{})
	pane.AddMessage("user", "Hello, AI!")

	if len(pane.messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(pane.messages))
	}

	msg := pane.messages[0]

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got '%s'", msg.Role)
	}

	if msg.Content != "Hello, AI!" {
		t.Errorf("expected content 'Hello, AI!', got '%s'", msg.Content)
	}

	if msg.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestNewChatPane(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)

	if pane == nil {
		t.Fatal("expected non-nil chat pane")
	}

	if pane.aiManager != aiMgr {
		t.Error("AI manager not set correctly")
	}

	if pane.visible {
		t.Error("pane should not be visible by default")
	}

	if pane.splitRatio <= 0 || pane.splitRatio > 1 {
		t.Error("split ratio should be between 0 and 1")
	}
}

func TestChatPaneSetSize(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)

	pane.SetSize(100, 50)

	if pane.width != 100 {
		t.Errorf("expected width 100, got %d", pane.width)
	}

	if pane.height != 50 {
		t.Errorf("expected height 50, got %d", pane.height)
	}
}

func TestChatPaneToggle(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)

	initialVisible := pane.visible

	pane.Toggle()

	if pane.visible == initialVisible {
		t.Error("toggle should change visibility state")
	}
}

func TestChatPaneFocus(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)

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

func TestChatPaneAddMessage(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)

	if len(pane.messages) != 0 {
		t.Error("chat should start with no messages")
	}

	pane.AddMessage("user", "Test message")

	if len(pane.messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(pane.messages))
	}

	if pane.messages[0].Content != "Test message" {
		t.Errorf("expected message content 'Test message', got '%s'", pane.messages[0].Content)
	}
}

func TestChatPaneClearMessages(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)

	// Add some messages
	pane.AddMessage("user", "Message 1")
	pane.AddMessage("assistant", "Message 2")

	if len(pane.messages) != 2 {
		t.Error("expected 2 messages before clear")
	}

	pane.ClearMessages()

	if len(pane.messages) != 0 {
		t.Errorf("expected 0 messages after clear, got %d", len(pane.messages))
	}
}

func TestChatPaneResize(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)
	pane.SetSize(100, 50)

	initialRatio := pane.splitRatio

	// Increase size
	newRatio := initialRatio + 0.1
	pane.SetSplitRatio(newRatio)

	if pane.splitRatio <= initialRatio {
		t.Error("split ratio should increase")
	}

	// Decrease size
	newRatio = initialRatio - 0.1
	pane.SetSplitRatio(newRatio)

	if pane.splitRatio >= initialRatio {
		t.Error("split ratio should decrease")
	}

	// Ensure bounds
	pane.SetSplitRatio(10.0) // Try to go beyond 1.0
	if pane.splitRatio > 0.9 {
		t.Error("split ratio should be capped at 0.9")
	}

	pane.SetSplitRatio(-10.0) // Try to go below 0
	if pane.splitRatio < 0.1 {
		t.Error("split ratio should have a minimum of 0.1")
	}
}

func TestChatPaneUpdate(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)
	pane.SetSize(100, 50)

	// Test window size message
	updatedPane, _ := pane.Update(tea.WindowSizeMsg{Width: 120, Height: 60})

	if updatedPane.width != 120 {
		t.Errorf("expected width 120 after WindowSizeMsg, got %d", updatedPane.width)
	}
}

func TestChatPaneView(t *testing.T) {
	t.Parallel()

	aiMgr := &ai.Manager{}
	pane := NewChatPane(aiMgr)
	pane.SetSize(100, 50)

	// When not visible
	rendered := pane.View()
	if rendered != "" {
		t.Error("view should return empty string when not visible")
	}

	// When visible
	pane.Show()
	rendered = pane.View()
	if rendered == "" {
		t.Error("view should return non-empty string when visible")
	}
}

func TestDefaultChatKeyMap(t *testing.T) {
	t.Parallel()

	km := DefaultKeyMap()

	if km.Send.Enabled() {
		// Just checking initialization
	}

	if km.Clear.Enabled() {
		// Just checking initialization
	}
}
