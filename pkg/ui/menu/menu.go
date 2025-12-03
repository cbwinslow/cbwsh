// Package menu provides a menu bar UI component for cbwsh.
package menu

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a menu item.
type MenuItem struct {
	Label       string
	Key         key.Binding
	Action      func() tea.Cmd
	SubMenu     *Menu
	Enabled     bool
	Separator   bool
	Description string
}

// Menu represents a dropdown menu.
type Menu struct {
	Label    string
	Items    []MenuItem
	Selected int
	Open     bool
}

// NewMenu creates a new menu.
func NewMenu(label string, items ...MenuItem) *Menu {
	return &Menu{
		Label:    label,
		Items:    items,
		Selected: 0,
	}
}

// IsEnabled indicates if the menu is enabled.
func (m *Menu) IsEnabled() bool {
	return true
}

// AddItem adds an item to the menu.
func (m *Menu) AddItem(item MenuItem) {
	if !item.Separator {
		item.Enabled = true
	}
	m.Items = append(m.Items, item)
}

// AddSeparator adds a separator to the menu.
func (m *Menu) AddSeparator() {
	m.Items = append(m.Items, MenuItem{Separator: true})
}

// SelectNext selects the next enabled item.
func (m *Menu) SelectNext() {
	if len(m.Items) == 0 {
		return
	}

	start := m.Selected
	for {
		m.Selected = (m.Selected + 1) % len(m.Items)
		if !m.Items[m.Selected].Separator && m.Items[m.Selected].Enabled {
			break
		}
		if m.Selected == start {
			break // No enabled items found
		}
	}
}

// SelectPrev selects the previous enabled item.
func (m *Menu) SelectPrev() {
	if len(m.Items) == 0 {
		return
	}

	start := m.Selected
	for {
		m.Selected--
		if m.Selected < 0 {
			m.Selected = len(m.Items) - 1
		}
		if !m.Items[m.Selected].Separator && m.Items[m.Selected].Enabled {
			break
		}
		if m.Selected == start {
			break // No enabled items found
		}
	}
}

// SelectedItem returns the currently selected item.
func (m *Menu) SelectedItem() *MenuItem {
	if len(m.Items) == 0 || m.Selected < 0 || m.Selected >= len(m.Items) {
		return nil
	}
	return &m.Items[m.Selected]
}

// MenuBar represents a menu bar with multiple menus.
type MenuBar struct {
	Menus      []*Menu
	ActiveMenu int
	Open       bool
	Styles     MenuBarStyles
	Width      int
	keyMap     KeyMap
}

// MenuBarStyles defines the styles for the menu bar.
type MenuBarStyles struct {
	Bar             lipgloss.Style
	MenuLabel       lipgloss.Style
	MenuLabelActive lipgloss.Style
	DropdownBorder  lipgloss.Style
	Item            lipgloss.Style
	ItemSelected    lipgloss.Style
	ItemDisabled    lipgloss.Style
	Separator       lipgloss.Style
	Shortcut        lipgloss.Style
}

// DefaultStyles returns default menu bar styles.
func DefaultStyles() MenuBarStyles {
	return MenuBarStyles{
		Bar: lipgloss.NewStyle().
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),
		MenuLabel: lipgloss.NewStyle().
			Padding(0, 1),
		MenuLabelActive: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1),
		DropdownBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Background(lipgloss.Color("236")),
		Item: lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("252")),
		ItemSelected: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("255")),
		ItemDisabled: lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("240")),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		Shortcut: lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Align(lipgloss.Right),
	}
}

// KeyMap defines the key bindings for the menu bar.
type KeyMap struct {
	Toggle key.Binding
	Left   key.Binding
	Right  key.Binding
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Escape key.Binding
	AltF   key.Binding
	AltE   key.Binding
	AltV   key.Binding
	AltH   key.Binding
}

// DefaultKeyMap returns default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Toggle: key.NewBinding(
			key.WithKeys("alt+m", "f10"),
			key.WithHelp("Alt+M", "toggle menu"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "previous menu"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "next menu"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "previous item"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "next item"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "select"),
		),
		Escape: key.NewBinding(
			key.WithKeys("escape"),
			key.WithHelp("Esc", "close menu"),
		),
		AltF: key.NewBinding(
			key.WithKeys("alt+f"),
			key.WithHelp("Alt+F", "File menu"),
		),
		AltE: key.NewBinding(
			key.WithKeys("alt+e"),
			key.WithHelp("Alt+E", "Edit menu"),
		),
		AltV: key.NewBinding(
			key.WithKeys("alt+v"),
			key.WithHelp("Alt+V", "View menu"),
		),
		AltH: key.NewBinding(
			key.WithKeys("alt+h"),
			key.WithHelp("Alt+H", "Help menu"),
		),
	}
}

// NewMenuBar creates a new menu bar.
func NewMenuBar() *MenuBar {
	return &MenuBar{
		Menus:  make([]*Menu, 0),
		Styles: DefaultStyles(),
		keyMap: DefaultKeyMap(),
	}
}

// AddMenu adds a menu to the menu bar.
func (m *MenuBar) AddMenu(menu *Menu) {
	m.Menus = append(m.Menus, menu)
}

// SetWidth sets the width of the menu bar.
func (m *MenuBar) SetWidth(width int) {
	m.Width = width
}

// IsOpen returns whether the menu bar is open.
func (m *MenuBar) IsOpen() bool {
	return m.Open
}

// Toggle toggles the menu bar open/closed.
func (m *MenuBar) Toggle() {
	m.Open = !m.Open
	if m.Open && len(m.Menus) > 0 {
		m.Menus[m.ActiveMenu].Open = true
	} else {
		m.closeAllMenus()
	}
}

// Close closes the menu bar.
func (m *MenuBar) Close() {
	m.Open = false
	m.closeAllMenus()
}

func (m *MenuBar) closeAllMenus() {
	for _, menu := range m.Menus {
		menu.Open = false
	}
}

// NextMenu moves to the next menu.
func (m *MenuBar) NextMenu() {
	if len(m.Menus) == 0 {
		return
	}
	m.closeAllMenus()
	m.ActiveMenu = (m.ActiveMenu + 1) % len(m.Menus)
	if m.Open {
		m.Menus[m.ActiveMenu].Open = true
	}
}

// PrevMenu moves to the previous menu.
func (m *MenuBar) PrevMenu() {
	if len(m.Menus) == 0 {
		return
	}
	m.closeAllMenus()
	m.ActiveMenu--
	if m.ActiveMenu < 0 {
		m.ActiveMenu = len(m.Menus) - 1
	}
	if m.Open {
		m.Menus[m.ActiveMenu].Open = true
	}
}

// SelectMenu selects a menu by index.
func (m *MenuBar) SelectMenu(index int) {
	if index < 0 || index >= len(m.Menus) {
		return
	}
	m.closeAllMenus()
	m.ActiveMenu = index
	m.Open = true
	m.Menus[m.ActiveMenu].Open = true
}

// Update handles input for the menu bar.
func (m *MenuBar) Update(msg tea.Msg) (bool, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check for menu open keys
		switch {
		case key.Matches(msg, m.keyMap.Toggle):
			m.Toggle()
			return true, nil

		case key.Matches(msg, m.keyMap.AltF):
			m.SelectMenuByLabel("File")
			return true, nil

		case key.Matches(msg, m.keyMap.AltE):
			m.SelectMenuByLabel("Edit")
			return true, nil

		case key.Matches(msg, m.keyMap.AltV):
			m.SelectMenuByLabel("View")
			return true, nil

		case key.Matches(msg, m.keyMap.AltH):
			m.SelectMenuByLabel("Help")
			return true, nil
		}

		// Only handle navigation when menu is open
		if !m.Open {
			return false, nil
		}

		switch {
		case key.Matches(msg, m.keyMap.Escape):
			m.Close()
			return true, nil

		case key.Matches(msg, m.keyMap.Left):
			m.PrevMenu()
			return true, nil

		case key.Matches(msg, m.keyMap.Right):
			m.NextMenu()
			return true, nil

		case key.Matches(msg, m.keyMap.Up):
			if m.ActiveMenu >= 0 && m.ActiveMenu < len(m.Menus) {
				m.Menus[m.ActiveMenu].SelectPrev()
			}
			return true, nil

		case key.Matches(msg, m.keyMap.Down):
			if m.ActiveMenu >= 0 && m.ActiveMenu < len(m.Menus) {
				m.Menus[m.ActiveMenu].SelectNext()
			}
			return true, nil

		case key.Matches(msg, m.keyMap.Enter):
			return m.executeSelected()
		}
	}

	return false, nil
}

// SelectMenuByLabel selects a menu by its label.
func (m *MenuBar) SelectMenuByLabel(label string) {
	for i, menu := range m.Menus {
		if menu.Label == label {
			m.SelectMenu(i)
			return
		}
	}
}

func (m *MenuBar) executeSelected() (bool, tea.Cmd) {
	if m.ActiveMenu < 0 || m.ActiveMenu >= len(m.Menus) {
		return true, nil
	}

	menu := m.Menus[m.ActiveMenu]
	item := menu.SelectedItem()

	if item == nil || !item.Enabled || item.Separator {
		return true, nil
	}

	// If there's a submenu, open it
	if item.SubMenu != nil {
		item.SubMenu.Open = true
		return true, nil
	}

	// Execute the action
	m.Close()
	if item.Action != nil {
		return true, item.Action()
	}

	return true, nil
}

// View renders the menu bar.
func (m *MenuBar) View() string {
	var parts []string

	for i, menu := range m.Menus {
		style := m.Styles.MenuLabel
		if i == m.ActiveMenu && m.Open {
			style = m.Styles.MenuLabelActive
		}
		parts = append(parts, style.Render(menu.Label))
	}

	bar := m.Styles.Bar.Width(m.Width).Render(strings.Join(parts, " "))

	if !m.Open {
		return bar
	}

	// Render dropdown
	dropdown := m.renderDropdown()
	if dropdown != "" {
		return bar + "\n" + dropdown
	}

	return bar
}

func (m *MenuBar) renderDropdown() string {
	if m.ActiveMenu < 0 || m.ActiveMenu >= len(m.Menus) {
		return ""
	}

	menu := m.Menus[m.ActiveMenu]
	if !menu.Open || len(menu.Items) == 0 {
		return ""
	}

	// Calculate dropdown position
	offset := 0
	for i := 0; i < m.ActiveMenu; i++ {
		offset += len(m.Menus[i].Label) + 3 // Label + padding
	}

	// Calculate max width
	maxWidth := 0
	for _, item := range menu.Items {
		if !item.Separator && len(item.Label) > maxWidth {
			maxWidth = len(item.Label)
		}
	}

	// Add shortcut width
	maxWidth += 15

	var renderedLines []string
	for i, item := range menu.Items {
		if item.Separator {
			renderedLines = append(renderedLines, m.Styles.Separator.Render(strings.Repeat("─", maxWidth)))
			continue
		}

		style := m.Styles.Item
		if i == menu.Selected {
			style = m.Styles.ItemSelected
		}
		if !item.Enabled {
			style = m.Styles.ItemDisabled
		}

		label := item.Label
		shortcut := ""
		if item.Key.Enabled() {
			shortcut = item.Key.Help().Key
		}

		// Pad label to maxWidth - shortcut length
		padding := maxWidth - len(label) - len(shortcut)
		if padding < 0 {
			padding = 0
		}

		line := style.Width(maxWidth).Render(label + strings.Repeat(" ", padding) + shortcut)
		renderedLines = append(renderedLines, line)
	}

	dropdown := m.Styles.DropdownBorder.Render(strings.Join(renderedLines, "\n"))

	// Add offset spacing
	if offset > 0 {
		padding := strings.Repeat(" ", offset)
		lines := strings.Split(dropdown, "\n")
		for i := range lines {
			lines[i] = padding + lines[i]
		}
		dropdown = strings.Join(lines, "\n")
	}

	return dropdown
}

// CreateDefaultMenus creates standard File, Edit, View, and Help menus.
func CreateDefaultMenus() []*Menu {
	fileMenu := NewMenu("File",
		MenuItem{
			Label:       "New Window",
			Key:         key.NewBinding(key.WithKeys("ctrl+shift+n"), key.WithHelp("Ctrl+Shift+N", "new window")),
			Description: "Open a new window",
			Enabled:     true,
		},
		MenuItem{
			Label:       "New Tab",
			Key:         key.NewBinding(key.WithKeys("ctrl+t"), key.WithHelp("Ctrl+T", "new tab")),
			Description: "Open a new tab",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "Open...",
			Key:         key.NewBinding(key.WithKeys("ctrl+o"), key.WithHelp("Ctrl+O", "open")),
			Description: "Open a file",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Save Session",
			Key:         key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("Ctrl+S", "save")),
			Description: "Save the current session",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "Preferences",
			Key:         key.NewBinding(key.WithKeys("ctrl+,"), key.WithHelp("Ctrl+,", "preferences")),
			Description: "Open preferences",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "Exit",
			Key:         key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("Ctrl+Q", "quit")),
			Description: "Exit the application",
			Enabled:     true,
		},
	)

	editMenu := NewMenu("Edit",
		MenuItem{
			Label:       "Copy",
			Key:         key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("Ctrl+C", "copy")),
			Description: "Copy selection",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Paste",
			Key:         key.NewBinding(key.WithKeys("ctrl+v"), key.WithHelp("Ctrl+V", "paste")),
			Description: "Paste from clipboard",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Select All",
			Key:         key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("Ctrl+A", "select all")),
			Description: "Select all text",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "Find",
			Key:         key.NewBinding(key.WithKeys("ctrl+f"), key.WithHelp("Ctrl+F", "find")),
			Description: "Find text",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "Clear Screen",
			Key:         key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("Ctrl+L", "clear")),
			Description: "Clear the screen",
			Enabled:     true,
		},
	)

	viewMenu := NewMenu("View",
		MenuItem{
			Label:       "Split Horizontal",
			Key:         key.NewBinding(key.WithKeys("ctrl+-"), key.WithHelp("Ctrl+-", "split horizontal")),
			Description: "Split pane horizontally",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Split Vertical",
			Key:         key.NewBinding(key.WithKeys("ctrl+\\"), key.WithHelp("Ctrl+\\", "split vertical")),
			Description: "Split pane vertically",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "Toggle Sidebar",
			Key:         key.NewBinding(key.WithKeys("ctrl+b"), key.WithHelp("Ctrl+B", "sidebar")),
			Description: "Toggle sidebar visibility",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Toggle Status Bar",
			Key:         key.NewBinding(key.WithKeys("ctrl+shift+b"), key.WithHelp("Ctrl+Shift+B", "status bar")),
			Description: "Toggle status bar visibility",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "Zoom In",
			Key:         key.NewBinding(key.WithKeys("ctrl++"), key.WithHelp("Ctrl++", "zoom in")),
			Description: "Increase font size",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Zoom Out",
			Key:         key.NewBinding(key.WithKeys("ctrl+minus"), key.WithHelp("Ctrl+-", "zoom out")),
			Description: "Decrease font size",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Reset Zoom",
			Key:         key.NewBinding(key.WithKeys("ctrl+0"), key.WithHelp("Ctrl+0", "reset zoom")),
			Description: "Reset font size",
			Enabled:     true,
		},
	)

	helpMenu := NewMenu("Help",
		MenuItem{
			Label:       "Keyboard Shortcuts",
			Key:         key.NewBinding(key.WithKeys("ctrl+?"), key.WithHelp("Ctrl+?", "shortcuts")),
			Description: "Show keyboard shortcuts",
			Enabled:     true,
		},
		MenuItem{
			Label:       "Documentation",
			Key:         key.NewBinding(key.WithKeys("f1"), key.WithHelp("F1", "docs")),
			Description: "Open documentation",
			Enabled:     true,
		},
		MenuItem{Separator: true},
		MenuItem{
			Label:       "About",
			Description: "About cbwsh",
			Enabled:     true,
		},
	)

	return []*Menu{fileMenu, editMenu, viewMenu, helpMenu}
}
