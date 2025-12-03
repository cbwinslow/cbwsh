package clipboard

import (
	"testing"
	"time"
)

func TestManager_New(t *testing.T) {
	m := NewManager()

	if m == nil {
		t.Fatal("NewManager() returned nil")
	}

	if !m.IsEnabled() {
		t.Error("NewManager() should be enabled by default")
	}

	if m.GetMaxHistory() != 100 {
		t.Errorf("GetMaxHistory() = %d, want 100", m.GetMaxHistory())
	}
}

func TestManager_History(t *testing.T) {
	m := NewManager()

	// Initially empty
	history := m.History()
	if len(history) != 0 {
		t.Errorf("History() len = %d, want 0", len(history))
	}

	// Add entries via addToHistory
	m.addToHistory("entry1", "test")
	m.addToHistory("entry2", "test")
	m.addToHistory("entry3", "test")

	history = m.History()
	if len(history) != 3 {
		t.Errorf("History() len = %d, want 3", len(history))
	}

	if history[0].Content != "entry1" {
		t.Errorf("History()[0].Content = %s, want entry1", history[0].Content)
	}
}

func TestManager_HistoryLast(t *testing.T) {
	m := NewManager()

	m.addToHistory("entry1", "test")
	m.addToHistory("entry2", "test")
	m.addToHistory("entry3", "test")
	m.addToHistory("entry4", "test")
	m.addToHistory("entry5", "test")

	// Get last 3
	last := m.HistoryLast(3)
	if len(last) != 3 {
		t.Errorf("HistoryLast(3) len = %d, want 3", len(last))
	}

	if last[0].Content != "entry3" {
		t.Errorf("HistoryLast(3)[0].Content = %s, want entry3", last[0].Content)
	}

	// Get more than available
	last = m.HistoryLast(10)
	if len(last) != 5 {
		t.Errorf("HistoryLast(10) len = %d, want 5", len(last))
	}
}

func TestManager_ClearHistory(t *testing.T) {
	m := NewManager()

	m.addToHistory("entry1", "test")
	m.addToHistory("entry2", "test")

	m.ClearHistory()

	history := m.History()
	if len(history) != 0 {
		t.Errorf("History() len after clear = %d, want 0", len(history))
	}
}

func TestManager_SetMaxHistory(t *testing.T) {
	m := NewManager()

	m.SetMaxHistory(3)
	if m.GetMaxHistory() != 3 {
		t.Errorf("GetMaxHistory() = %d, want 3", m.GetMaxHistory())
	}

	// Add more than max
	m.addToHistory("entry1", "test")
	m.addToHistory("entry2", "test")
	m.addToHistory("entry3", "test")
	m.addToHistory("entry4", "test")
	m.addToHistory("entry5", "test")

	history := m.History()
	if len(history) != 3 {
		t.Errorf("History() len = %d, want 3 (max)", len(history))
	}

	// Oldest entries should be removed
	if history[0].Content != "entry3" {
		t.Errorf("History()[0].Content = %s, want entry3", history[0].Content)
	}
}

func TestManager_EnableDisable(t *testing.T) {
	m := NewManager()

	if !m.IsEnabled() {
		t.Error("Should be enabled by default")
	}

	m.Disable()
	if m.IsEnabled() {
		t.Error("Should be disabled after Disable()")
	}

	m.Enable()
	if !m.IsEnabled() {
		t.Error("Should be enabled after Enable()")
	}
}

func TestManager_DisabledOperations(t *testing.T) {
	m := NewManager()
	m.Disable()

	_, err := m.Read()
	if err != ErrClipboardUnavailable {
		t.Errorf("Read() when disabled error = %v, want ErrClipboardUnavailable", err)
	}

	err = m.Write("test")
	if err != ErrClipboardUnavailable {
		t.Errorf("Write() when disabled error = %v, want ErrClipboardUnavailable", err)
	}
}

func TestManager_NoDuplicateConsecutive(t *testing.T) {
	m := NewManager()

	m.addToHistory("same", "test")
	m.addToHistory("same", "test")
	m.addToHistory("same", "test")

	history := m.History()
	if len(history) != 1 {
		t.Errorf("History() len = %d, want 1 (no duplicates)", len(history))
	}
}

func TestManager_DifferentEntries(t *testing.T) {
	m := NewManager()

	m.addToHistory("entry1", "test")
	m.addToHistory("entry2", "test")
	m.addToHistory("entry1", "test") // Same as first but not consecutive

	history := m.History()
	if len(history) != 3 {
		t.Errorf("History() len = %d, want 3", len(history))
	}
}

func TestEntry(t *testing.T) {
	entry := Entry{
		Content:   "test content",
		Timestamp: time.Now(),
		Source:    "test",
	}

	if entry.Content != "test content" {
		t.Errorf("Entry.Content = %s, want 'test content'", entry.Content)
	}

	if entry.Source != "test" {
		t.Errorf("Entry.Source = %s, want 'test'", entry.Source)
	}

	if entry.Timestamp.IsZero() {
		t.Error("Entry.Timestamp should not be zero")
	}
}

func TestManager_CopyPaste(t *testing.T) {
	m := NewManager()

	// Copy is alias for Write
	// Paste is alias for Read
	// These will actually try clipboard operations, so we just check they don't panic
	// In a real environment without clipboard access, these would return errors

	// Just verify methods exist and can be called
	_ = m.Copy
	_ = m.Paste
}

// Note: Actual clipboard read/write tests are skipped because they require
// a working display/clipboard on the system. These tests focus on the
// history and management functionality.

func TestManager_WriteAndTrack(t *testing.T) {
	m := NewManager()

	// This will attempt actual clipboard write, which may fail in CI
	// But we can verify the tracking logic by checking history after enabling disabled
	m.Disable()
	err := m.WriteAndTrack("test", "custom-source")
	if err != ErrClipboardUnavailable {
		t.Errorf("WriteAndTrack() when disabled should return ErrClipboardUnavailable")
	}
}
