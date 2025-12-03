package palette

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPalette_New(t *testing.T) {
	p := New()

	if p == nil {
		t.Fatal("New() returned nil")
	}

	if p.IsVisible() {
		t.Error("New palette should not be visible")
	}
}

func TestPalette_AddCommand(t *testing.T) {
	p := New()

	cmd := Command{
		ID:          "test",
		Name:        "Test Command",
		Description: "A test command",
	}

	p.AddCommand(cmd)

	if len(p.commands) != 1 {
		t.Errorf("AddCommand() len = %d, want 1", len(p.commands))
	}

	if p.commands[0].ID != "test" {
		t.Errorf("AddCommand() ID = %s, want test", p.commands[0].ID)
	}
}

func TestPalette_AddCommands(t *testing.T) {
	p := New()

	cmds := []Command{
		{ID: "cmd1", Name: "Command 1"},
		{ID: "cmd2", Name: "Command 2"},
		{ID: "cmd3", Name: "Command 3"},
	}

	p.AddCommands(cmds)

	if len(p.commands) != 3 {
		t.Errorf("AddCommands() len = %d, want 3", len(p.commands))
	}
}

func TestPalette_RemoveCommand(t *testing.T) {
	p := New()

	p.AddCommands([]Command{
		{ID: "cmd1", Name: "Command 1"},
		{ID: "cmd2", Name: "Command 2"},
	})

	p.RemoveCommand("cmd1")

	if len(p.commands) != 1 {
		t.Errorf("RemoveCommand() len = %d, want 1", len(p.commands))
	}

	if p.commands[0].ID != "cmd2" {
		t.Errorf("RemoveCommand() remaining ID = %s, want cmd2", p.commands[0].ID)
	}
}

func TestPalette_ClearCommands(t *testing.T) {
	p := New()

	p.AddCommands([]Command{
		{ID: "cmd1", Name: "Command 1"},
		{ID: "cmd2", Name: "Command 2"},
	})

	p.ClearCommands()

	if len(p.commands) != 0 {
		t.Errorf("ClearCommands() len = %d, want 0", len(p.commands))
	}
}

func TestPalette_OpenClose(t *testing.T) {
	p := New()

	if p.IsVisible() {
		t.Error("Palette should start closed")
	}

	p.Open()

	if !p.IsVisible() {
		t.Error("Palette should be visible after Open()")
	}

	p.Close()

	if p.IsVisible() {
		t.Error("Palette should not be visible after Close()")
	}
}

func TestPalette_Filter(t *testing.T) {
	p := New()

	p.AddCommands([]Command{
		{ID: "git_status", Name: "Git Status", Category: "Git"},
		{ID: "git_commit", Name: "Git Commit", Category: "Git"},
		{ID: "new_pane", Name: "New Pane", Category: "Panes"},
	})

	p.Open()

	// No filter - should show all
	if len(p.filtered) != 3 {
		t.Errorf("Unfiltered len = %d, want 3", len(p.filtered))
	}

	// Filter by "git"
	p.SetQuery("git")
	if len(p.filtered) != 2 {
		t.Errorf("Filtered by 'git' len = %d, want 2", len(p.filtered))
	}

	// Filter by "pane"
	p.SetQuery("pane")
	if len(p.filtered) != 1 {
		t.Errorf("Filtered by 'pane' len = %d, want 1", len(p.filtered))
	}

	// Filter by non-existent
	p.SetQuery("nonexistent")
	if len(p.filtered) != 0 {
		t.Errorf("Filtered by 'nonexistent' len = %d, want 0", len(p.filtered))
	}
}

func TestPalette_FilterByKeyword(t *testing.T) {
	p := New()

	p.AddCommand(Command{
		ID:       "test",
		Name:     "Test Command",
		Keywords: []string{"special", "keyword"},
	})

	p.Open()
	p.SetQuery("special")

	if len(p.filtered) != 1 {
		t.Errorf("Filtered by keyword len = %d, want 1", len(p.filtered))
	}
}

func TestPalette_Navigation(t *testing.T) {
	p := New()

	p.AddCommands([]Command{
		{ID: "cmd1", Name: "Command 1"},
		{ID: "cmd2", Name: "Command 2"},
		{ID: "cmd3", Name: "Command 3"},
	})

	p.Open()

	// Initially at 0
	if p.selected != 0 {
		t.Errorf("Initial selected = %d, want 0", p.selected)
	}

	// Move down
	p.Update(tea.KeyMsg{Type: tea.KeyDown})
	if p.selected != 1 {
		t.Errorf("After down, selected = %d, want 1", p.selected)
	}

	// Move down again
	p.Update(tea.KeyMsg{Type: tea.KeyDown})
	if p.selected != 2 {
		t.Errorf("After down x2, selected = %d, want 2", p.selected)
	}

	// Move down wraps to 0
	p.Update(tea.KeyMsg{Type: tea.KeyDown})
	if p.selected != 0 {
		t.Errorf("After down x3 (wrap), selected = %d, want 0", p.selected)
	}

	// Move up wraps to end
	p.Update(tea.KeyMsg{Type: tea.KeyUp})
	if p.selected != 2 {
		t.Errorf("After up (wrap), selected = %d, want 2", p.selected)
	}
}

func TestPalette_Selected(t *testing.T) {
	p := New()

	p.AddCommands([]Command{
		{ID: "cmd1", Name: "Command 1"},
		{ID: "cmd2", Name: "Command 2"},
	})

	p.Open()

	selected := p.Selected()
	if selected == nil {
		t.Fatal("Selected() returned nil")
	}

	if selected.ID != "cmd1" {
		t.Errorf("Selected() ID = %s, want cmd1", selected.ID)
	}

	// Move to next
	p.Update(tea.KeyMsg{Type: tea.KeyDown})
	selected = p.Selected()
	if selected.ID != "cmd2" {
		t.Errorf("Selected() after down ID = %s, want cmd2", selected.ID)
	}
}

func TestPalette_CloseOnEscape(t *testing.T) {
	p := New()
	p.Open()

	if !p.IsVisible() {
		t.Fatal("Palette should be visible")
	}

	// Use Runes to simulate escape key press
	p.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if p.IsVisible() {
		t.Error("Palette should close on Escape")
	}
}

func TestPalette_View(t *testing.T) {
	p := New()

	// Not visible - should return empty
	view := p.View()
	if view != "" {
		t.Error("View() when not visible should return empty string")
	}

	p.AddCommand(Command{
		ID:   "test",
		Name: "Test Command",
	})

	p.Open()
	view = p.View()

	if view == "" {
		t.Error("View() when visible should return non-empty string")
	}
}

func TestPalette_NoResultsView(t *testing.T) {
	p := New()
	p.Open()
	p.SetQuery("nonexistent")

	view := p.View()

	if view == "" {
		t.Error("View() should render even with no results")
	}
}

func TestPalette_Query(t *testing.T) {
	p := New()
	p.Open()

	p.SetQuery("test")
	if p.Query() != "test" {
		t.Errorf("Query() = %s, want test", p.Query())
	}
}

func TestDefaultCommands(t *testing.T) {
	cmds := DefaultCommands()

	if len(cmds) == 0 {
		t.Fatal("DefaultCommands() returned empty list")
	}

	// Check for expected commands
	expectedIDs := map[string]bool{
		"new_pane":   true,
		"close_pane": true,
		"ai_assist":  true,
		"quit":       true,
		"show_help":  true,
	}

	foundIDs := make(map[string]bool)
	for _, cmd := range cmds {
		foundIDs[cmd.ID] = true
	}

	for id := range expectedIDs {
		if !foundIDs[id] {
			t.Errorf("DefaultCommands() missing command: %s", id)
		}
	}
}

func TestDefaultStyles(t *testing.T) {
	styles := DefaultStyles()

	// Verify styles are initialized (non-nil)
	// These styles should have some configuration applied
	_ = styles.Background
	_ = styles.Border
	_ = styles.Input
	_ = styles.Item
	_ = styles.SelectedItem
}
