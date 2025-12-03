// Package panes provides terminal pane management for cbwsh.
package panes

import (
	"fmt"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/shell"
	"github.com/google/uuid"
)

// Pane represents a terminal pane with its own shell executor.
type Pane struct {
	mu       sync.RWMutex
	id       string
	title    string
	active   bool
	width    int
	height   int
	executor *shell.Executor
	output   []string
	scrollY  int
}

// NewPane creates a new pane.
func NewPane(shellType core.ShellType) *Pane {
	return &Pane{
		id:       uuid.New().String()[:8],
		title:    "Shell",
		executor: shell.NewExecutor(shellType),
		output:   make([]string, 0),
	}
}

// ID returns the pane identifier.
func (p *Pane) ID() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.id
}

// Title returns the pane title.
func (p *Pane) Title() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.title
}

// SetTitle sets the pane title.
func (p *Pane) SetTitle(title string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.title = title
}

// IsActive returns whether the pane is active.
func (p *Pane) IsActive() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.active
}

// Activate makes this pane active.
func (p *Pane) Activate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active = true
}

// Deactivate makes this pane inactive.
func (p *Pane) Deactivate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.active = false
}

// Width returns the pane width.
func (p *Pane) Width() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.width
}

// Height returns the pane height.
func (p *Pane) Height() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.height
}

// SetSize sets the pane dimensions.
func (p *Pane) SetSize(width, height int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.width = width
	p.height = height
}

// GetExecutor returns the pane's executor.
func (p *Pane) GetExecutor() core.Executor {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.executor
}

// GetShellExecutor returns the pane's shell executor with full functionality.
func (p *Pane) GetShellExecutor() *shell.Executor {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.executor
}

// AppendOutput appends output to the pane.
func (p *Pane) AppendOutput(line string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.output = append(p.output, line)
}

// GetOutput returns the pane output.
func (p *Pane) GetOutput() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]string, len(p.output))
	copy(result, p.output)
	return result
}

// ClearOutput clears the pane output.
func (p *Pane) ClearOutput() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.output = p.output[:0]
}

// ScrollUp scrolls the pane up.
func (p *Pane) ScrollUp(lines int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.scrollY = max(0, p.scrollY-lines)
}

// ScrollDown scrolls the pane down.
func (p *Pane) ScrollDown(lines int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	maxScroll := max(0, len(p.output)-p.height)
	p.scrollY = min(maxScroll, p.scrollY+lines)
}

// GetScrollY returns the current scroll position.
func (p *Pane) GetScrollY() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.scrollY
}

// Manager manages multiple panes.
type Manager struct {
	mu           sync.RWMutex
	panes        map[string]*Pane
	activePaneID string
	layout       core.PaneLayout
	shellType    core.ShellType
}

// NewManager creates a new pane manager.
func NewManager(shellType core.ShellType) *Manager {
	return &Manager{
		panes:     make(map[string]*Pane),
		layout:    core.LayoutSingle,
		shellType: shellType,
	}
}

// Create creates a new pane.
func (m *Manager) Create() (core.Pane, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pane := NewPane(m.shellType)
	m.panes[pane.ID()] = pane

	// If this is the first pane, make it active
	if m.activePaneID == "" {
		pane.Activate()
		m.activePaneID = pane.ID()
	}

	return pane, nil
}

// Close closes a pane by ID.
func (m *Manager) Close(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.panes[id]; !exists {
		return fmt.Errorf("pane not found: %s", id)
	}

	delete(m.panes, id)

	// If we closed the active pane, select another
	if m.activePaneID == id {
		m.activePaneID = ""
		for paneID, pane := range m.panes {
			pane.Activate()
			m.activePaneID = paneID
			break
		}
	}

	return nil
}

// Get returns a pane by ID.
func (m *Manager) Get(id string) (core.Pane, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pane, exists := m.panes[id]
	return pane, exists
}

// GetPane returns a pane by ID with full functionality.
func (m *Manager) GetPane(id string) (*Pane, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pane, exists := m.panes[id]
	return pane, exists
}

// Active returns the active pane.
func (m *Manager) Active() core.Pane {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if pane, exists := m.panes[m.activePaneID]; exists {
		return pane
	}
	return nil
}

// ActivePane returns the active pane with full functionality.
func (m *Manager) ActivePane() *Pane {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if pane, exists := m.panes[m.activePaneID]; exists {
		return pane
	}
	return nil
}

// SetActive sets the active pane.
func (m *Manager) SetActive(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pane, exists := m.panes[id]
	if !exists {
		return fmt.Errorf("pane not found: %s", id)
	}

	// Deactivate current pane
	if currentPane, exists := m.panes[m.activePaneID]; exists {
		currentPane.Deactivate()
	}

	// Activate new pane
	pane.Activate()
	m.activePaneID = id

	return nil
}

// List returns all panes.
func (m *Manager) List() []core.Pane {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]core.Pane, 0, len(m.panes))
	for _, pane := range m.panes {
		result = append(result, pane)
	}
	return result
}

// ListPanes returns all panes with full functionality.
func (m *Manager) ListPanes() []*Pane {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Pane, 0, len(m.panes))
	for _, pane := range m.panes {
		result = append(result, pane)
	}
	return result
}

// Layout returns the current layout.
func (m *Manager) Layout() core.PaneLayout {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.layout
}

// SetLayout sets the pane layout.
func (m *Manager) SetLayout(layout core.PaneLayout) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.layout = layout
	return nil
}

// Split splits the active pane.
func (m *Manager) Split(direction core.PaneLayout) (core.Pane, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	newPane := NewPane(m.shellType)
	m.panes[newPane.ID()] = newPane

	// Update layout based on direction
	switch direction {
	case core.LayoutHorizontalSplit:
		m.layout = core.LayoutHorizontalSplit
	case core.LayoutVerticalSplit:
		m.layout = core.LayoutVerticalSplit
	}

	return newPane, nil
}

// NextPane moves focus to the next pane.
func (m *Manager) NextPane() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.panes) <= 1 {
		return nil
	}

	// Get ordered list of pane IDs
	ids := make([]string, 0, len(m.panes))
	for id := range m.panes {
		ids = append(ids, id)
	}

	// Find current position and move to next
	currentIdx := 0
	for i, id := range ids {
		if id == m.activePaneID {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(ids)

	// Deactivate current
	if pane, exists := m.panes[m.activePaneID]; exists {
		pane.Deactivate()
	}

	// Activate next
	nextID := ids[nextIdx]
	m.panes[nextID].Activate()
	m.activePaneID = nextID

	return nil
}

// PrevPane moves focus to the previous pane.
func (m *Manager) PrevPane() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.panes) <= 1 {
		return nil
	}

	// Get ordered list of pane IDs
	ids := make([]string, 0, len(m.panes))
	for id := range m.panes {
		ids = append(ids, id)
	}

	// Find current position and move to previous
	currentIdx := 0
	for i, id := range ids {
		if id == m.activePaneID {
			currentIdx = i
			break
		}
	}

	prevIdx := currentIdx - 1
	if prevIdx < 0 {
		prevIdx = len(ids) - 1
	}

	// Deactivate current
	if pane, exists := m.panes[m.activePaneID]; exists {
		pane.Deactivate()
	}

	// Activate previous
	prevID := ids[prevIdx]
	m.panes[prevID].Activate()
	m.activePaneID = prevID

	return nil
}

// Count returns the number of panes.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.panes)
}

// UpdateAllSizes updates sizes for all panes based on terminal dimensions.
func (m *Manager) UpdateAllSizes(totalWidth, totalHeight int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	paneCount := len(m.panes)
	if paneCount == 0 {
		return
	}

	switch m.layout {
	case core.LayoutSingle:
		for _, pane := range m.panes {
			pane.SetSize(totalWidth, totalHeight)
		}
	case core.LayoutHorizontalSplit:
		paneHeight := totalHeight / paneCount
		for _, pane := range m.panes {
			pane.SetSize(totalWidth, paneHeight)
		}
	case core.LayoutVerticalSplit:
		paneWidth := totalWidth / paneCount
		for _, pane := range m.panes {
			pane.SetSize(paneWidth, totalHeight)
		}
	case core.LayoutGrid:
		cols := 2
		rows := (paneCount + cols - 1) / cols
		paneWidth := totalWidth / cols
		paneHeight := totalHeight / rows
		for _, pane := range m.panes {
			pane.SetSize(paneWidth, paneHeight)
		}
	default:
		for _, pane := range m.panes {
			pane.SetSize(totalWidth, totalHeight)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
