package menu_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"

	"github.com/cbwinslow/cbwsh/pkg/ui/menu"
)

func TestNewMenu(t *testing.T) {
	t.Parallel()

	m := menu.NewMenu("File")
	if m == nil {
		t.Fatal("expected non-nil menu")
	}

	if m.Label != "File" {
		t.Errorf("expected label 'File', got '%s'", m.Label)
	}
}

func TestMenuAddItem(t *testing.T) {
	t.Parallel()

	m := menu.NewMenu("File")
	m.AddItem(menu.MenuItem{
		Label:   "New",
		Enabled: true,
	})

	if len(m.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(m.Items))
	}

	if m.Items[0].Label != "New" {
		t.Errorf("expected item label 'New', got '%s'", m.Items[0].Label)
	}
}

func TestMenuAddSeparator(t *testing.T) {
	t.Parallel()

	m := menu.NewMenu("File")
	m.AddItem(menu.MenuItem{Label: "New", Enabled: true})
	m.AddSeparator()
	m.AddItem(menu.MenuItem{Label: "Exit", Enabled: true})

	if len(m.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(m.Items))
	}

	if !m.Items[1].Separator {
		t.Error("expected second item to be separator")
	}
}

func TestMenuSelectNext(t *testing.T) {
	t.Parallel()

	m := menu.NewMenu("File",
		menu.MenuItem{Label: "New", Enabled: true},
		menu.MenuItem{Label: "Open", Enabled: true},
		menu.MenuItem{Label: "Save", Enabled: true},
	)

	if m.Selected != 0 {
		t.Errorf("expected initial selection 0, got %d", m.Selected)
	}

	m.SelectNext()
	if m.Selected != 1 {
		t.Errorf("expected selection 1, got %d", m.Selected)
	}

	m.SelectNext()
	if m.Selected != 2 {
		t.Errorf("expected selection 2, got %d", m.Selected)
	}

	m.SelectNext()
	if m.Selected != 0 {
		t.Errorf("expected selection to wrap to 0, got %d", m.Selected)
	}
}

func TestMenuSelectPrev(t *testing.T) {
	t.Parallel()

	m := menu.NewMenu("File",
		menu.MenuItem{Label: "New", Enabled: true},
		menu.MenuItem{Label: "Open", Enabled: true},
		menu.MenuItem{Label: "Save", Enabled: true},
	)

	m.SelectPrev()
	if m.Selected != 2 {
		t.Errorf("expected selection to wrap to 2, got %d", m.Selected)
	}

	m.SelectPrev()
	if m.Selected != 1 {
		t.Errorf("expected selection 1, got %d", m.Selected)
	}
}

func TestMenuSelectSkipsSeparator(t *testing.T) {
	t.Parallel()

	m := menu.NewMenu("File",
		menu.MenuItem{Label: "New", Enabled: true},
		menu.MenuItem{Separator: true},
		menu.MenuItem{Label: "Exit", Enabled: true},
	)

	m.SelectNext()
	if m.Selected != 2 {
		t.Errorf("expected selection to skip separator to 2, got %d", m.Selected)
	}
}

func TestMenuSelectedItem(t *testing.T) {
	t.Parallel()

	m := menu.NewMenu("File",
		menu.MenuItem{Label: "New", Enabled: true},
		menu.MenuItem{Label: "Open", Enabled: true},
	)

	item := m.SelectedItem()
	if item == nil {
		t.Fatal("expected non-nil selected item")
	}

	if item.Label != "New" {
		t.Errorf("expected selected item 'New', got '%s'", item.Label)
	}
}

func TestNewMenuBar(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	if bar == nil {
		t.Fatal("expected non-nil menu bar")
	}
}

func TestMenuBarAddMenu(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	bar.AddMenu(menu.NewMenu("File"))
	bar.AddMenu(menu.NewMenu("Edit"))

	if len(bar.Menus) != 2 {
		t.Errorf("expected 2 menus, got %d", len(bar.Menus))
	}
}

func TestMenuBarToggle(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	bar.AddMenu(menu.NewMenu("File"))

	if bar.IsOpen() {
		t.Error("expected menu bar to be closed initially")
	}

	bar.Toggle()
	if !bar.IsOpen() {
		t.Error("expected menu bar to be open after toggle")
	}

	bar.Toggle()
	if bar.IsOpen() {
		t.Error("expected menu bar to be closed after second toggle")
	}
}

func TestMenuBarNextPrevMenu(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	bar.AddMenu(menu.NewMenu("File"))
	bar.AddMenu(menu.NewMenu("Edit"))
	bar.AddMenu(menu.NewMenu("View"))
	bar.Toggle()

	if bar.ActiveMenu != 0 {
		t.Errorf("expected active menu 0, got %d", bar.ActiveMenu)
	}

	bar.NextMenu()
	if bar.ActiveMenu != 1 {
		t.Errorf("expected active menu 1, got %d", bar.ActiveMenu)
	}

	bar.NextMenu()
	if bar.ActiveMenu != 2 {
		t.Errorf("expected active menu 2, got %d", bar.ActiveMenu)
	}

	bar.NextMenu()
	if bar.ActiveMenu != 0 {
		t.Errorf("expected active menu to wrap to 0, got %d", bar.ActiveMenu)
	}

	bar.PrevMenu()
	if bar.ActiveMenu != 2 {
		t.Errorf("expected active menu to wrap to 2, got %d", bar.ActiveMenu)
	}
}

func TestMenuBarClose(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	bar.AddMenu(menu.NewMenu("File"))
	bar.Toggle()

	if !bar.IsOpen() {
		t.Error("expected menu bar to be open")
	}

	bar.Close()
	if bar.IsOpen() {
		t.Error("expected menu bar to be closed after Close()")
	}
}

func TestMenuBarSelectMenu(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	bar.AddMenu(menu.NewMenu("File"))
	bar.AddMenu(menu.NewMenu("Edit"))
	bar.AddMenu(menu.NewMenu("View"))

	bar.SelectMenu(2)
	if bar.ActiveMenu != 2 {
		t.Errorf("expected active menu 2, got %d", bar.ActiveMenu)
	}
	if !bar.IsOpen() {
		t.Error("expected menu bar to be open after SelectMenu")
	}
}

func TestMenuBarSelectMenuByLabel(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	bar.AddMenu(menu.NewMenu("File"))
	bar.AddMenu(menu.NewMenu("Edit"))
	bar.AddMenu(menu.NewMenu("View"))

	bar.SelectMenuByLabel("Edit")
	if bar.ActiveMenu != 1 {
		t.Errorf("expected active menu 1 (Edit), got %d", bar.ActiveMenu)
	}
}

func TestMenuBarView(t *testing.T) {
	t.Parallel()

	bar := menu.NewMenuBar()
	bar.AddMenu(menu.NewMenu("File"))
	bar.AddMenu(menu.NewMenu("Edit"))
	bar.SetWidth(80)

	view := bar.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestDefaultStyles(t *testing.T) {
	t.Parallel()

	styles := menu.DefaultStyles()
	// Just verify it doesn't panic
	_ = styles.Bar.Render("Test")
}

func TestDefaultKeyMap(t *testing.T) {
	t.Parallel()

	keyMap := menu.DefaultKeyMap()
	if !keyMap.Toggle.Enabled() {
		t.Error("expected Toggle key to be enabled")
	}
}

func TestCreateDefaultMenus(t *testing.T) {
	t.Parallel()

	menus := menu.CreateDefaultMenus()
	if len(menus) != 4 {
		t.Errorf("expected 4 default menus, got %d", len(menus))
	}

	labels := []string{"File", "Edit", "View", "Help"}
	for i, m := range menus {
		if m.Label != labels[i] {
			t.Errorf("expected menu %d to be '%s', got '%s'", i, labels[i], m.Label)
		}
	}
}

func TestMenuItemWithKey(t *testing.T) {
	t.Parallel()

	item := menu.MenuItem{
		Label:   "New",
		Key:     key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("Ctrl+N", "new")),
		Enabled: true,
	}

	if !item.Key.Enabled() {
		t.Error("expected key binding to be enabled")
	}
}
