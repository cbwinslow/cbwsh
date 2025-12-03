// Package aichat provides a resizable AI chat pane component for cbwsh.
package aichat

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cbwinslow/cbwsh/pkg/ai"
	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/ui/markdown"
)

// Message represents a chat message.
type Message struct {
	Role      string // "user" or "assistant"
	Content   string
	Timestamp time.Time
}

// ChatPane provides a resizable AI chat interface.
type ChatPane struct {
	mu sync.RWMutex

	// UI components
	viewport   viewport.Model
	textarea   textarea.Model
	mdRenderer *markdown.Renderer

	// State
	messages   []Message
	focused    bool
	loading    bool
	width      int
	height     int
	splitRatio float64 // Ratio of chat pane to total width (0.0 to 1.0)
	visible    bool

	// AI
	aiManager *ai.Manager

	// Styles
	userStyle      lipgloss.Style
	assistantStyle lipgloss.Style
	borderStyle    lipgloss.Style
	titleStyle     lipgloss.Style
}

// KeyMap defines key bindings for the chat pane.
type KeyMap struct {
	Send       key.Binding
	Clear      key.Binding
	Escape     key.Binding
	Resize     key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Send: key.NewBinding(
			key.WithKeys("ctrl+enter", "ctrl+s"),
			key.WithHelp("ctrl+enter", "send message"),
		),
		Clear: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "clear chat"),
		),
		Escape: key.NewBinding(
			key.WithKeys("escape"),
			key.WithHelp("esc", "unfocus"),
		),
		Resize: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "resize"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("ctrl+up", "pageup"),
			key.WithHelp("pgup", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("ctrl+down", "pagedown"),
			key.WithHelp("pgdn", "scroll down"),
		),
	}
}

var keys = DefaultKeyMap()

// NewChatPane creates a new chat pane.
func NewChatPane(aiManager *ai.Manager) *ChatPane {
	ta := textarea.New()
	ta.Placeholder = "Ask AI a question..."
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.CharLimit = 5000
	ta.SetHeight(3)

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle()

	mdRenderer, _ := markdown.NewRenderer()

	return &ChatPane{
		viewport:   vp,
		textarea:   ta,
		mdRenderer: mdRenderer,
		messages:   make([]Message, 0),
		splitRatio: 0.3, // Default 30% width
		visible:    false,
		aiManager:  aiManager,
		userStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),
		assistantStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("141")),
		borderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Padding(0, 1),
	}
}

// Focus focuses the chat pane.
func (c *ChatPane) Focus() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.focused = true
	c.textarea.Focus()
}

// Blur removes focus from the chat pane.
func (c *ChatPane) Blur() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.focused = false
	c.textarea.Blur()
}

// IsFocused returns whether the chat pane is focused.
func (c *ChatPane) IsFocused() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.focused
}

// Show shows the chat pane.
func (c *ChatPane) Show() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visible = true
}

// Hide hides the chat pane.
func (c *ChatPane) Hide() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visible = false
}

// IsVisible returns whether the chat pane is visible.
func (c *ChatPane) IsVisible() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.visible
}

// Toggle toggles the chat pane visibility.
func (c *ChatPane) Toggle() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visible = !c.visible
}

// SetSize sets the size of the chat pane.
func (c *ChatPane) SetSize(width, height int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.width = width
	c.height = height

	// Update component sizes
	contentWidth := width - 4    // Account for borders
	viewportHeight := height - 8 // Account for input and borders

	if viewportHeight < 5 {
		viewportHeight = 5
	}

	c.viewport.Width = contentWidth
	c.viewport.Height = viewportHeight
	c.textarea.SetWidth(contentWidth)
}

// SetSplitRatio sets the split ratio for the chat pane.
func (c *ChatPane) SetSplitRatio(ratio float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ratio < 0.1 {
		ratio = 0.1
	}
	if ratio > 0.9 {
		ratio = 0.9
	}
	c.splitRatio = ratio
}

// GetSplitRatio returns the current split ratio.
func (c *ChatPane) GetSplitRatio() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.splitRatio
}

// GetWidth returns the calculated width based on split ratio.
func (c *ChatPane) GetWidth(totalWidth int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return int(float64(totalWidth) * c.splitRatio)
}

// ClearMessages clears all chat messages.
func (c *ChatPane) ClearMessages() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = c.messages[:0]
	c.viewport.SetContent("")
}

// AddMessage adds a message to the chat.
func (c *ChatPane) AddMessage(role, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	c.updateViewportContent()
}

func (c *ChatPane) updateViewportContent() {
	var content strings.Builder

	for _, msg := range c.messages {
		if msg.Role == "user" {
			content.WriteString(c.userStyle.Render("You: "))
			content.WriteString(msg.Content)
		} else {
			content.WriteString(c.assistantStyle.Render("AI: "))
			// Render markdown for assistant messages
			if c.mdRenderer != nil {
				rendered, err := c.mdRenderer.Render(msg.Content)
				if err == nil {
					content.WriteString(rendered)
				} else {
					content.WriteString(msg.Content)
				}
			} else {
				content.WriteString(msg.Content)
			}
		}
		content.WriteString("\n\n")
	}

	c.viewport.SetContent(content.String())
	c.viewport.GotoBottom()
}

// SendMessage sends a message to the AI and adds the response.
func (c *ChatPane) SendMessage(ctx context.Context) tea.Cmd {
	c.mu.Lock()
	message := strings.TrimSpace(c.textarea.Value())
	if message == "" {
		c.mu.Unlock()
		return nil
	}

	c.textarea.Reset()
	c.loading = true
	c.mu.Unlock()

	// Add user message
	c.AddMessage("user", message)

	return func() tea.Msg {
		var response string
		var err error

		if c.aiManager != nil {
			response, err = c.aiManager.Query(ctx, message)
		} else {
			response = "AI is not configured. Please set up an AI provider."
		}

		if err != nil {
			response = "Error: " + err.Error()
		}

		return aiResponseMsg{response: response}
	}
}

type aiResponseMsg struct {
	response string
}

// Update handles messages for the chat pane.
func (c *ChatPane) Update(msg tea.Msg) (*ChatPane, tea.Cmd) {
	c.mu.Lock()
	if !c.visible {
		c.mu.Unlock()
		return c, nil
	}
	c.mu.Unlock()

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		c.mu.RLock()
		focused := c.focused
		c.mu.RUnlock()

		if !focused {
			return c, nil
		}

		switch {
		case key.Matches(msg, keys.Send):
			return c, c.SendMessage(context.Background())

		case key.Matches(msg, keys.Clear):
			c.ClearMessages()
			return c, nil

		case key.Matches(msg, keys.Escape):
			c.Blur()
			return c, nil

		case key.Matches(msg, keys.ScrollUp):
			c.mu.Lock()
			for i := 0; i < 3; i++ {
				c.viewport, _ = c.viewport.Update(tea.KeyMsg{Type: tea.KeyUp})
			}
			c.mu.Unlock()
			return c, nil

		case key.Matches(msg, keys.ScrollDown):
			c.mu.Lock()
			for i := 0; i < 3; i++ {
				c.viewport, _ = c.viewport.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			c.mu.Unlock()
			return c, nil
		}

	case aiResponseMsg:
		c.mu.Lock()
		c.loading = false
		c.mu.Unlock()
		c.AddMessage("assistant", msg.response)
		return c, nil

	case tea.WindowSizeMsg:
		c.SetSize(c.GetWidth(msg.Width), msg.Height-4)
		return c, nil
	}

	// Update textarea
	c.mu.Lock()
	var cmd tea.Cmd
	c.textarea, cmd = c.textarea.Update(msg)
	c.mu.Unlock()
	cmds = append(cmds, cmd)

	// Update viewport
	c.mu.Lock()
	c.viewport, cmd = c.viewport.Update(msg)
	c.mu.Unlock()
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

// View renders the chat pane.
func (c *ChatPane) View() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.visible {
		return ""
	}

	var content strings.Builder

	// Title
	title := c.titleStyle.Render("🤖 AI Assistant")
	if c.loading {
		title += " (thinking...)"
	}
	content.WriteString(title)
	content.WriteString("\n")

	// Messages viewport
	content.WriteString(c.viewport.View())
	content.WriteString("\n")

	// Input area
	inputStyle := lipgloss.NewStyle()
	if c.focused {
		inputStyle = inputStyle.BorderForeground(lipgloss.Color("86"))
	}
	content.WriteString(inputStyle.Render(c.textarea.View()))

	// Wrap in border
	borderStyle := c.borderStyle
	if c.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("86"))
	}

	return borderStyle.
		Width(c.width - 2).
		Height(c.height - 2).
		Render(content.String())
}

// GetMessages returns all chat messages.
func (c *ChatPane) GetMessages() []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]Message, len(c.messages))
	copy(result, c.messages)
	return result
}

// SetAIManager sets the AI manager.
func (c *ChatPane) SetAIManager(manager *ai.Manager) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.aiManager = manager
}

// Conversation represents a saved conversation.
type Conversation struct {
	ID        string
	Title     string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConversationManager manages multiple conversations.
type ConversationManager struct {
	mu            sync.RWMutex
	conversations map[string]*Conversation
	activeID      string
}

// NewConversationManager creates a new conversation manager.
func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		conversations: make(map[string]*Conversation),
	}
}

// Create creates a new conversation.
func (cm *ConversationManager) Create(title string) *Conversation {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	id := generateID()
	conv := &Conversation{
		ID:        id,
		Title:     title,
		Messages:  make([]Message, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cm.conversations[id] = conv
	if cm.activeID == "" {
		cm.activeID = id
	}

	return conv
}

// Get returns a conversation by ID.
func (cm *ConversationManager) Get(id string) (*Conversation, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conv, exists := cm.conversations[id]
	return conv, exists
}

// Delete removes a conversation.
func (cm *ConversationManager) Delete(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.conversations, id)
	if cm.activeID == id {
		cm.activeID = ""
		for newID := range cm.conversations {
			cm.activeID = newID
			break
		}
	}
	return nil
}

// List returns all conversations.
func (cm *ConversationManager) List() []*Conversation {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make([]*Conversation, 0, len(cm.conversations))
	for _, conv := range cm.conversations {
		result = append(result, conv)
	}
	return result
}

// Active returns the active conversation.
func (cm *ConversationManager) Active() *Conversation {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.conversations[cm.activeID]
}

// SetActive sets the active conversation.
func (cm *ConversationManager) SetActive(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.conversations[id]; !exists {
		return core.ErrNotFound
	}
	cm.activeID = id
	return nil
}

// AddMessage adds a message to a conversation.
func (cm *ConversationManager) AddMessage(id string, msg Message) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conv, exists := cm.conversations[id]
	if !exists {
		return core.ErrNotFound
	}

	conv.Messages = append(conv.Messages, msg)
	conv.UpdatedAt = time.Now()
	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
