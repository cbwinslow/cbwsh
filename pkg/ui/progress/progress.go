// Package progress provides progress bar components for cbwsh.
package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Bar represents a progress bar component.
type Bar struct {
	mu        sync.RWMutex
	progress  progress.Model
	current   int
	total     int
	message   string
	finished  bool
	startTime time.Time
	width     int
}

// NewBar creates a new progress bar.
func NewBar() *Bar {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	return &Bar{
		progress: p,
		width:    40,
	}
}

// NewBarWithColors creates a new progress bar with custom colors.
func NewBarWithColors(colorA, colorB string) *Bar {
	p := progress.New(
		progress.WithGradient(colorA, colorB),
		progress.WithWidth(40),
	)
	return &Bar{
		progress: p,
		width:    40,
	}
}

// Start starts progress tracking.
func (b *Bar) Start(total int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = 0
	b.total = total
	b.finished = false
	b.startTime = time.Now()
}

// Increment increments progress by one.
func (b *Bar) Increment() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.current < b.total {
		b.current++
	}
}

// IncrementBy increments progress by a specific amount.
func (b *Bar) IncrementBy(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current += amount
	if b.current > b.total {
		b.current = b.total
	}
}

// SetMessage sets the progress message.
func (b *Bar) SetMessage(message string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.message = message
}

// Finish completes the progress.
func (b *Bar) Finish() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = b.total
	b.finished = true
}

// View returns the progress view string.
func (b *Bar) View() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.total == 0 {
		return ""
	}

	percent := float64(b.current) / float64(b.total)
	bar := b.progress.ViewAs(percent)

	elapsed := time.Since(b.startTime).Round(time.Second)
	stats := fmt.Sprintf(" %d/%d (%s)", b.current, b.total, elapsed)

	if b.message != "" {
		return fmt.Sprintf("%s %s%s", b.message, bar, stats)
	}
	return bar + stats
}

// Percent returns the current progress percentage.
func (b *Bar) Percent() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.total == 0 {
		return 0
	}
	return float64(b.current) / float64(b.total) * 100
}

// IsFinished returns whether the progress is complete.
func (b *Bar) IsFinished() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.finished
}

// SetWidth sets the progress bar width.
func (b *Bar) SetWidth(width int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.width = width
	b.progress = progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(width),
	)
}

// Model is a bubbletea model for a progress bar.
type Model struct {
	Bar      *Bar
	width    int
	height   int
	quitting bool
}

// NewModel creates a new progress model.
func NewModel() Model {
	return Model{
		Bar:   NewBar(),
		width: 40,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Bar.SetWidth(msg.Width - 20)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return m.Bar.View()
}

// MultiBar manages multiple progress bars.
type MultiBar struct {
	mu   sync.RWMutex
	bars map[string]*Bar
}

// NewMultiBar creates a new multi-bar progress tracker.
func NewMultiBar() *MultiBar {
	return &MultiBar{
		bars: make(map[string]*Bar),
	}
}

// Add adds a new progress bar.
func (m *MultiBar) Add(name string, total int) *Bar {
	m.mu.Lock()
	defer m.mu.Unlock()
	bar := NewBar()
	bar.Start(total)
	m.bars[name] = bar
	return bar
}

// Get returns a progress bar by name.
func (m *MultiBar) Get(name string) (*Bar, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	bar, exists := m.bars[name]
	return bar, exists
}

// Remove removes a progress bar.
func (m *MultiBar) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.bars, name)
}

// View returns the combined view of all progress bars.
func (m *MultiBar) View() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lines []string
	for name, bar := range m.bars {
		lines = append(lines, fmt.Sprintf("%s: %s", name, bar.View()))
	}
	return strings.Join(lines, "\n")
}

// AllFinished returns whether all bars are finished.
func (m *MultiBar) AllFinished() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, bar := range m.bars {
		if !bar.IsFinished() {
			return false
		}
	}
	return true
}

// Spinner represents a spinner component.
type Spinner struct {
	mu      sync.RWMutex
	frames  []string
	current int
	message string
	active  bool
}

// NewSpinner creates a new spinner.
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Start starts the spinner.
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = true
}

// Stop stops the spinner.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = false
}

// SetMessage sets the spinner message.
func (s *Spinner) SetMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

// Tick advances the spinner.
func (s *Spinner) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = (s.current + 1) % len(s.frames)
}

// View returns the spinner view.
func (s *Spinner) View() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.active {
		return ""
	}

	frame := s.frames[s.current]
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	if s.message != "" {
		return style.Render(frame) + " " + s.message
	}
	return style.Render(frame)
}
