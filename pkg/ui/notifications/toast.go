// Package notifications provides toast notifications for cbwsh.
package notifications

import (
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// NotificationType represents the type of notification.
type NotificationType int

const (
	// NotificationTypeInfo is an informational notification.
	NotificationTypeInfo NotificationType = iota
	// NotificationTypeSuccess is a success notification.
	NotificationTypeSuccess
	// NotificationTypeWarning is a warning notification.
	NotificationTypeWarning
	// NotificationTypeError is an error notification.
	NotificationTypeError
)

// Position represents the notification position.
type Position int

const (
	// PositionTopRight shows notifications in the top right.
	PositionTopRight Position = iota
	// PositionTopLeft shows notifications in the top left.
	PositionTopLeft
	// PositionBottomRight shows notifications in the bottom right.
	PositionBottomRight
	// PositionBottomLeft shows notifications in the bottom left.
	PositionBottomLeft
	// PositionTopCenter shows notifications in the top center.
	PositionTopCenter
	// PositionBottomCenter shows notifications in the bottom center.
	PositionBottomCenter
)

// Notification represents a toast notification.
type Notification struct {
	// ID is the unique identifier.
	ID string
	// Type is the notification type.
	Type NotificationType
	// Title is the notification title.
	Title string
	// Message is the notification message.
	Message string
	// Duration is how long to show the notification.
	Duration time.Duration
	// CreatedAt is when the notification was created.
	CreatedAt time.Time
	// Action is an optional action to run when clicked.
	Action func()
	// ActionLabel is the label for the action button.
	ActionLabel string
	// Dismissible indicates if the notification can be dismissed.
	Dismissible bool
	// Progress is the optional progress value (0-100).
	Progress int
	// ShowProgress indicates whether to show progress bar.
	ShowProgress bool
}

// Styles defines the notification styles.
type Styles struct {
	Container   lipgloss.Style
	Info        lipgloss.Style
	Success     lipgloss.Style
	Warning     lipgloss.Style
	Error       lipgloss.Style
	Title       lipgloss.Style
	Message     lipgloss.Style
	Action      lipgloss.Style
	Progress    lipgloss.Style
	ProgressBar lipgloss.Style
}

// DefaultStyles returns default notification styles.
func DefaultStyles() Styles {
	base := lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Border(lipgloss.RoundedBorder()).
		Width(40)

	return Styles{
		Container: lipgloss.NewStyle(),
		Info: base.
			BorderForeground(lipgloss.Color("62")),
		Success: base.
			BorderForeground(lipgloss.Color("42")),
		Warning: base.
			BorderForeground(lipgloss.Color("214")),
		Error: base.
			BorderForeground(lipgloss.Color("196")),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255")),
		Message: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Action: lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Underline(true),
		Progress: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")),
		ProgressBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")),
	}
}

// Manager manages notifications.
type Manager struct {
	mu            sync.RWMutex
	notifications []*Notification
	maxVisible    int
	position      Position
	styles        Styles
	width         int
	height        int
}

// NewManager creates a new notification manager.
func NewManager() *Manager {
	return &Manager{
		notifications: make([]*Notification, 0),
		maxVisible:    5,
		position:      PositionTopRight,
		styles:        DefaultStyles(),
	}
}

// Show displays a new notification.
func (m *Manager) Show(n Notification) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate ID if not set
	if n.ID == "" {
		n.ID = uuid.New().String()
	}

	// Set created time
	n.CreatedAt = time.Now()

	// Default duration (only if ShowProgress is false - progress notifications don't auto-dismiss)
	if n.Duration == 0 && !n.ShowProgress {
		n.Duration = 5 * time.Second
	}

	// Default dismissible
	if n.Type != NotificationTypeError {
		n.Dismissible = true
	}

	m.notifications = append(m.notifications, &n)

	// Schedule auto-dismiss
	if n.Duration > 0 {
		go func() {
			time.Sleep(n.Duration)
			m.Dismiss(n.ID)
		}()
	}

	return n.ID
}

// ShowInfo shows an info notification.
func (m *Manager) ShowInfo(title, message string) string {
	return m.Show(Notification{
		Type:    NotificationTypeInfo,
		Title:   title,
		Message: message,
	})
}

// ShowSuccess shows a success notification.
func (m *Manager) ShowSuccess(title, message string) string {
	return m.Show(Notification{
		Type:    NotificationTypeSuccess,
		Title:   title,
		Message: message,
	})
}

// ShowWarning shows a warning notification.
func (m *Manager) ShowWarning(title, message string) string {
	return m.Show(Notification{
		Type:    NotificationTypeWarning,
		Title:   title,
		Message: message,
	})
}

// ShowError shows an error notification.
func (m *Manager) ShowError(title, message string) string {
	return m.Show(Notification{
		Type:     NotificationTypeError,
		Title:    title,
		Message:  message,
		Duration: 10 * time.Second,
	})
}

// ShowCommandComplete shows a notification for command completion.
func (m *Manager) ShowCommandComplete(command string, exitCode int, duration time.Duration) string {
	var notifType NotificationType
	var title string

	if exitCode == 0 {
		notifType = NotificationTypeSuccess
		title = "Command Completed"
	} else {
		notifType = NotificationTypeError
		title = "Command Failed"
	}

	return m.Show(Notification{
		Type:    notifType,
		Title:   title,
		Message: truncate(command, 30) + " (" + duration.Round(time.Millisecond).String() + ")",
	})
}

// ShowProgress shows a progress notification.
func (m *Manager) ShowProgress(title, message string) string {
	return m.Show(Notification{
		Type:         NotificationTypeInfo,
		Title:        title,
		Message:      message,
		ShowProgress: true,
		Progress:     0,
		Duration:     0, // Don't auto-dismiss
	})
}

// UpdateProgress updates a progress notification.
func (m *Manager) UpdateProgress(id string, progress int, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, n := range m.notifications {
		if n.ID == id {
			n.Progress = progress
			if message != "" {
				n.Message = message
			}
			break
		}
	}
}

// Dismiss removes a notification.
func (m *Manager) Dismiss(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, n := range m.notifications {
		if n.ID == id {
			m.notifications = append(m.notifications[:i], m.notifications[i+1:]...)
			break
		}
	}
}

// DismissAll removes all notifications.
func (m *Manager) DismissAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifications = make([]*Notification, 0)
}

// List returns all active notifications.
func (m *Manager) List() []*Notification {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Notification, len(m.notifications))
	copy(result, m.notifications)
	return result
}

// Count returns the number of active notifications.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.notifications)
}

// SetPosition sets the notification position.
func (m *Manager) SetPosition(pos Position) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.position = pos
}

// SetMaxVisible sets the maximum visible notifications.
func (m *Manager) SetMaxVisible(max int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxVisible = max
}

// SetStyles sets the notification styles.
func (m *Manager) SetStyles(styles Styles) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.styles = styles
}

// SetSize sets the manager size.
func (m *Manager) SetSize(width, height int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.width = width
	m.height = height
}

// View renders all visible notifications.
func (m *Manager) View() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.notifications) == 0 {
		return ""
	}

	var rendered []string
	count := 0

	for _, n := range m.notifications {
		if count >= m.maxVisible {
			break
		}
		rendered = append(rendered, m.renderNotification(n))
		count++
	}

	// Show remaining count
	remaining := len(m.notifications) - m.maxVisible
	if remaining > 0 {
		rendered = append(rendered, m.styles.Message.Render(
			"... and "+itoa(remaining)+" more",
		))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rendered...)
	return m.styles.Container.Render(content)
}

func (m *Manager) renderNotification(n *Notification) string {
	// Get style based on type
	var style lipgloss.Style
	switch n.Type {
	case NotificationTypeSuccess:
		style = m.styles.Success
	case NotificationTypeWarning:
		style = m.styles.Warning
	case NotificationTypeError:
		style = m.styles.Error
	default:
		style = m.styles.Info
	}

	// Build content
	var content string

	// Title
	if n.Title != "" {
		content = m.styles.Title.Render(n.Title)
	}

	// Message
	if n.Message != "" {
		if content != "" {
			content += "\n"
		}
		content += m.styles.Message.Render(n.Message)
	}

	// Progress bar
	if n.ShowProgress {
		if content != "" {
			content += "\n"
		}
		content += m.renderProgressBar(n.Progress)
	}

	// Action
	if n.ActionLabel != "" {
		if content != "" {
			content += "\n"
		}
		content += m.styles.Action.Render("[" + n.ActionLabel + "]")
	}

	return style.Render(content)
}

func (m *Manager) renderProgressBar(progress int) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	width := 30
	filled := width * progress / 100
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	return m.styles.ProgressBar.Render(bar) + " " + m.styles.Progress.Render(itoa(progress)+"%")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}

	var digits []byte
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	return string(digits)
}

// GetNotification returns a notification by ID.
func (m *Manager) GetNotification(id string) *Notification {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, n := range m.notifications {
		if n.ID == id {
			return n
		}
	}
	return nil
}

// CleanExpired removes expired notifications.
func (m *Manager) CleanExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	var active []*Notification

	for _, n := range m.notifications {
		if n.Duration == 0 || now.Before(n.CreatedAt.Add(n.Duration)) {
			active = append(active, n)
		}
	}

	m.notifications = active
}
