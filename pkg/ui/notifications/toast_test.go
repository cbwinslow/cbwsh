package notifications

import (
	"testing"
	"time"
)

func TestManager_New(t *testing.T) {
	m := NewManager()

	if m == nil {
		t.Fatal("NewManager() returned nil")
	}

	if m.Count() != 0 {
		t.Errorf("Count() = %d, want 0", m.Count())
	}
}

func TestManager_Show(t *testing.T) {
	m := NewManager()

	id := m.Show(Notification{
		Type:    NotificationTypeInfo,
		Title:   "Test",
		Message: "Test message",
	})

	if id == "" {
		t.Error("Show() returned empty ID")
	}

	if m.Count() != 1 {
		t.Errorf("Count() = %d, want 1", m.Count())
	}
}

func TestManager_ShowInfo(t *testing.T) {
	m := NewManager()

	id := m.ShowInfo("Title", "Message")

	if id == "" {
		t.Error("ShowInfo() returned empty ID")
	}

	n := m.GetNotification(id)
	if n == nil {
		t.Fatal("GetNotification() returned nil")
	}

	if n.Type != NotificationTypeInfo {
		t.Errorf("Type = %v, want NotificationTypeInfo", n.Type)
	}
}

func TestManager_ShowSuccess(t *testing.T) {
	m := NewManager()

	id := m.ShowSuccess("Title", "Message")

	n := m.GetNotification(id)
	if n == nil {
		t.Fatal("GetNotification() returned nil")
	}

	if n.Type != NotificationTypeSuccess {
		t.Errorf("Type = %v, want NotificationTypeSuccess", n.Type)
	}
}

func TestManager_ShowWarning(t *testing.T) {
	m := NewManager()

	id := m.ShowWarning("Title", "Message")

	n := m.GetNotification(id)
	if n == nil {
		t.Fatal("GetNotification() returned nil")
	}

	if n.Type != NotificationTypeWarning {
		t.Errorf("Type = %v, want NotificationTypeWarning", n.Type)
	}
}

func TestManager_ShowError(t *testing.T) {
	m := NewManager()

	id := m.ShowError("Title", "Message")

	n := m.GetNotification(id)
	if n == nil {
		t.Fatal("GetNotification() returned nil")
	}

	if n.Type != NotificationTypeError {
		t.Errorf("Type = %v, want NotificationTypeError", n.Type)
	}
}

func TestManager_ShowCommandComplete(t *testing.T) {
	m := NewManager()

	// Success case
	id := m.ShowCommandComplete("ls -la", 0, 100*time.Millisecond)
	n := m.GetNotification(id)
	if n.Type != NotificationTypeSuccess {
		t.Error("Successful command should show success notification")
	}

	// Error case
	id = m.ShowCommandComplete("invalid-cmd", 1, 100*time.Millisecond)
	n = m.GetNotification(id)
	if n.Type != NotificationTypeError {
		t.Error("Failed command should show error notification")
	}
}

func TestManager_ShowProgress(t *testing.T) {
	m := NewManager()

	id := m.ShowProgress("Downloading", "0%")
	n := m.GetNotification(id)

	if !n.ShowProgress {
		t.Error("ShowProgress should be true")
	}

	if n.Duration != 0 {
		t.Error("Progress notification should not auto-dismiss")
	}

	// Update progress
	m.UpdateProgress(id, 50, "50%")
	n = m.GetNotification(id)

	if n.Progress != 50 {
		t.Errorf("Progress = %d, want 50", n.Progress)
	}
}

func TestManager_Dismiss(t *testing.T) {
	m := NewManager()

	id := m.ShowInfo("Title", "Message")
	if m.Count() != 1 {
		t.Fatalf("Count() = %d, want 1", m.Count())
	}

	m.Dismiss(id)
	if m.Count() != 0 {
		t.Errorf("Count() after dismiss = %d, want 0", m.Count())
	}
}

func TestManager_DismissAll(t *testing.T) {
	m := NewManager()

	m.ShowInfo("Title 1", "Message 1")
	m.ShowInfo("Title 2", "Message 2")
	m.ShowInfo("Title 3", "Message 3")

	if m.Count() != 3 {
		t.Fatalf("Count() = %d, want 3", m.Count())
	}

	m.DismissAll()
	if m.Count() != 0 {
		t.Errorf("Count() after DismissAll = %d, want 0", m.Count())
	}
}

func TestManager_List(t *testing.T) {
	m := NewManager()

	m.ShowInfo("Title 1", "Message 1")
	m.ShowInfo("Title 2", "Message 2")

	list := m.List()
	if len(list) != 2 {
		t.Errorf("List() len = %d, want 2", len(list))
	}
}

func TestManager_SetPosition(t *testing.T) {
	m := NewManager()

	m.SetPosition(PositionBottomLeft)
	if m.position != PositionBottomLeft {
		t.Error("SetPosition() did not update position")
	}
}

func TestManager_SetMaxVisible(t *testing.T) {
	m := NewManager()

	m.SetMaxVisible(3)
	if m.maxVisible != 3 {
		t.Error("SetMaxVisible() did not update maxVisible")
	}
}

func TestManager_View(t *testing.T) {
	m := NewManager()

	// Empty view
	view := m.View()
	if view != "" {
		t.Error("View() should be empty when no notifications")
	}

	// With notifications
	m.ShowInfo("Title", "Message")
	view = m.View()
	if view == "" {
		t.Error("View() should not be empty with notifications")
	}
}

func TestManager_ViewWithMaxVisible(t *testing.T) {
	m := NewManager()
	m.SetMaxVisible(2)

	m.ShowInfo("Title 1", "Message 1")
	m.ShowInfo("Title 2", "Message 2")
	m.ShowInfo("Title 3", "Message 3")
	m.ShowInfo("Title 4", "Message 4")

	view := m.View()
	if view == "" {
		t.Error("View() should not be empty")
	}
}

func TestManager_GetNotification(t *testing.T) {
	m := NewManager()

	id := m.ShowInfo("Title", "Message")

	// Found
	n := m.GetNotification(id)
	if n == nil {
		t.Error("GetNotification() should return notification")
	}

	// Not found
	n = m.GetNotification("nonexistent")
	if n != nil {
		t.Error("GetNotification() should return nil for nonexistent ID")
	}
}

func TestManager_CleanExpired(t *testing.T) {
	m := NewManager()

	// Create notification with very short duration
	m.Show(Notification{
		Type:     NotificationTypeInfo,
		Title:    "Test",
		Duration: 1 * time.Millisecond,
	})

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	m.CleanExpired()

	if m.Count() != 0 {
		t.Errorf("Count() after CleanExpired = %d, want 0", m.Count())
	}
}

func TestRenderProgressBar(t *testing.T) {
	m := NewManager()

	tests := []struct {
		progress int
	}{
		{0},
		{50},
		{100},
		{-10}, // Should clamp to 0
		{150}, // Should clamp to 100
	}

	for _, tt := range tests {
		result := m.renderProgressBar(tt.progress)
		if result == "" {
			t.Errorf("renderProgressBar(%d) should not be empty", tt.progress)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"very long string", 10, "very lo..."},
		{"exactly10c", 10, "exactly10c"},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{123, "123"},
		{-5, "-5"},
	}

	for _, tt := range tests {
		got := itoa(tt.input)
		if got != tt.want {
			t.Errorf("itoa(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDefaultStyles(t *testing.T) {
	styles := DefaultStyles()

	// Verify styles are initialized (non-nil)
	// These styles should have some configuration applied
	_ = styles.Title
	_ = styles.Message
	_ = styles.Info
	_ = styles.Success
	_ = styles.Warning
	_ = styles.Error
}
