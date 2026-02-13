package autocomplete

import (
	"testing"
)

func TestNewCompleter(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	if completer == nil {
		t.Fatal("expected non-nil completer")
	}

	// Should have default providers
	if len(completer.providers) == 0 {
		t.Error("expected default providers to be added")
	}
}

func TestCompleterBasicCompletion(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	// Test basic command completion
	suggestions, err := completer.Complete("ls", 2)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should have at least some suggestions (depends on the system)
	_ = suggestions
}

func TestCompleterFileCompletion(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	// Test file path completion with current directory
	suggestions, err := completer.Complete("./", 2)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Depending on the implementation, this might return files in current directory
	// We just check that it doesn't panic
	_ = suggestions
}

func TestCompleterCommandCompletion(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	// Test various command patterns
	testCases := []string{
		"git",
		"cd",
		"mkdir",
		"echo",
	}

	for _, tc := range testCases {
		suggestions, err := completer.Complete(tc, len(tc))
		if err != nil {
			t.Errorf("unexpected error for '%s': %v", tc, err)
		}
		// Just ensure no panic occurs
		_ = suggestions
	}
}

func TestCompleterEmptyInput(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	suggestions, err := completer.Complete("", 0)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Empty input might return common commands or nothing
	// Just ensure it doesn't panic
	_ = suggestions
}

func TestCompleterWithOptions(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	// Test completion with command options
	testCases := []struct {
		input     string
		cursorPos int
	}{
		{"git -", 5},
		{"ls -l", 5},
		{"grep --", 7},
	}

	for _, tc := range testCases {
		suggestions, err := completer.Complete(tc.input, tc.cursorPos)
		if err != nil {
			t.Errorf("unexpected error for '%s': %v", tc.input, err)
		}
		_ = suggestions
	}
}

func TestCompleterAddProvider(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()
	initialCount := len(completer.providers)

	// Add a custom provider (we'll use an existing one for testing)
	completer.AddProvider(&CommandProvider{})

	if len(completer.providers) != initialCount+1 {
		t.Errorf("expected %d providers, got %d", initialCount+1, len(completer.providers))
	}
}

func TestCompleterRemoveProvider(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	// Try to remove a provider by name
	completer.RemoveProvider("command")

	// Provider count should change (if the provider existed)
	// This test mainly ensures the method doesn't panic
}

func TestCompleterCursorPosition(t *testing.T) {
	t.Parallel()

	completer := NewCompleter()

	testCases := []struct {
		input     string
		cursorPos int
	}{
		{"git status", 3},   // cursor at end of 'git'
		{"git status", 11},  // cursor at end of line
		{"git status", 0},   // cursor at start
		{"git status", 100}, // cursor beyond end (should be handled gracefully)
	}

	for _, tc := range testCases {
		suggestions, err := completer.Complete(tc.input, tc.cursorPos)
		if err != nil {
			t.Errorf("unexpected error for cursor at %d: %v", tc.cursorPos, err)
		}
		_ = suggestions
	}
}
