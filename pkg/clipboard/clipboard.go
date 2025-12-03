// Package clipboard provides cross-platform clipboard support for cbwsh.
package clipboard

import (
	"errors"
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

// Common errors.
var (
	ErrClipboardUnavailable = errors.New("clipboard unavailable")
	ErrEmptyClipboard       = errors.New("clipboard is empty")
)

// Entry represents a clipboard history entry.
type Entry struct {
	// Content is the clipboard content.
	Content string
	// Timestamp is when the entry was added.
	Timestamp time.Time
	// Source indicates where the content came from.
	Source string
}

// Manager manages clipboard operations with history.
type Manager struct {
	mu         sync.RWMutex
	history    []Entry
	maxHistory int
	enabled    bool
}

// NewManager creates a new clipboard manager.
func NewManager() *Manager {
	return &Manager{
		history:    make([]Entry, 0),
		maxHistory: 100,
		enabled:    true,
	}
}

// Read reads text from the clipboard.
func (m *Manager) Read() (string, error) {
	if !m.enabled {
		return "", ErrClipboardUnavailable
	}

	content, err := clipboard.ReadAll()
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", ErrEmptyClipboard
	}

	return content, nil
}

// Write writes text to the clipboard.
func (m *Manager) Write(text string) error {
	if !m.enabled {
		return ErrClipboardUnavailable
	}

	if err := clipboard.WriteAll(text); err != nil {
		return err
	}

	// Add to history
	m.addToHistory(text, "write")

	return nil
}

// WriteAndTrack writes to clipboard and tracks the operation.
func (m *Manager) WriteAndTrack(text, source string) error {
	if !m.enabled {
		return ErrClipboardUnavailable
	}

	if err := clipboard.WriteAll(text); err != nil {
		return err
	}

	m.addToHistory(text, source)

	return nil
}

// History returns the clipboard history.
func (m *Manager) History() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Entry, len(m.history))
	copy(result, m.history)
	return result
}

// HistoryLast returns the last n history entries.
func (m *Manager) HistoryLast(n int) []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if n >= len(m.history) {
		result := make([]Entry, len(m.history))
		copy(result, m.history)
		return result
	}

	start := len(m.history) - n
	result := make([]Entry, n)
	copy(result, m.history[start:])
	return result
}

// ClearHistory clears the clipboard history.
func (m *Manager) ClearHistory() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = make([]Entry, 0)
}

// SetMaxHistory sets the maximum history size.
func (m *Manager) SetMaxHistory(max int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxHistory = max

	// Trim if necessary
	if len(m.history) > max {
		m.history = m.history[len(m.history)-max:]
	}
}

// GetMaxHistory returns the maximum history size.
func (m *Manager) GetMaxHistory() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxHistory
}

// Enable enables clipboard operations.
func (m *Manager) Enable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = true
}

// Disable disables clipboard operations.
func (m *Manager) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = false
}

// IsEnabled returns whether clipboard is enabled.
func (m *Manager) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// Copy is a convenience function to copy text to clipboard.
func (m *Manager) Copy(text string) error {
	return m.Write(text)
}

// Paste is a convenience function to paste from clipboard.
func (m *Manager) Paste() (string, error) {
	return m.Read()
}

func (m *Manager) addToHistory(content, source string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := Entry{
		Content:   content,
		Timestamp: time.Now(),
		Source:    source,
	}

	// Don't add duplicate consecutive entries
	if len(m.history) > 0 {
		last := m.history[len(m.history)-1]
		if last.Content == content {
			return
		}
	}

	m.history = append(m.history, entry)

	// Trim if exceeds max
	if len(m.history) > m.maxHistory {
		m.history = m.history[1:]
	}
}

// Available checks if clipboard is available on the system.
func Available() bool {
	// Try to read from clipboard to check availability
	_, err := clipboard.ReadAll()
	// clipboard.ReadAll returns an empty string and nil error on empty clipboard
	// but returns an error if clipboard is unavailable
	return err == nil || err.Error() == "exit status 1" // xclip returns 1 on empty
}

// CopyText is a standalone function to copy text to clipboard.
func CopyText(text string) error {
	return clipboard.WriteAll(text)
}

// PasteText is a standalone function to paste text from clipboard.
func PasteText() (string, error) {
	return clipboard.ReadAll()
}
