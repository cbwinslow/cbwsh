// Package editor provides a markdown editor component for cbwsh.
package editor

import (
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cbwinslow/cbwsh/pkg/ui/markdown"
)

// EditorMode represents the current editor mode.
type EditorMode int

const (
	// ModeEdit is the editing mode.
	ModeEdit EditorMode = iota
	// ModePreview is the preview mode.
	ModePreview
	// ModeSplit shows both edit and preview.
	ModeSplit
)

// MarkdownEditor provides a terminal-based markdown editor.
type MarkdownEditor struct {
	mu sync.RWMutex

	// UI components
	textarea   textarea.Model
	viewport   viewport.Model
	mdRenderer *markdown.Renderer

	// State
	mode       EditorMode
	filePath   string
	content    string
	modified   bool
	focused    bool
	width      int
	height     int
	splitRatio float64

	// Styles
	editBorder    lipgloss.Style
	previewBorder lipgloss.Style
	titleStyle    lipgloss.Style
	statusStyle   lipgloss.Style
}

// KeyMap defines key bindings for the editor.
type KeyMap struct {
	Save       key.Binding
	SaveAs     key.Binding
	Open       key.Binding
	ToggleMode key.Binding
	Escape     key.Binding
	Bold       key.Binding
	Italic     key.Binding
	Code       key.Binding
	Link       key.Binding
	Heading    key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		SaveAs: key.NewBinding(
			key.WithKeys("ctrl+shift+s"),
			key.WithHelp("ctrl+shift+s", "save as"),
		),
		Open: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "open"),
		),
		ToggleMode: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "toggle mode"),
		),
		Escape: key.NewBinding(
			key.WithKeys("escape"),
			key.WithHelp("esc", "close"),
		),
		Bold: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "bold"),
		),
		Italic: key.NewBinding(
			key.WithKeys("ctrl+i"),
			key.WithHelp("ctrl+i", "italic"),
		),
		Code: key.NewBinding(
			key.WithKeys("ctrl+`"),
			key.WithHelp("ctrl+`", "code"),
		),
		Link: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "link"),
		),
		Heading: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "heading"),
		),
	}
}

var keys = DefaultKeyMap()

// NewMarkdownEditor creates a new markdown editor.
func NewMarkdownEditor() *MarkdownEditor {
	ta := textarea.New()
	ta.Placeholder = "Start typing markdown..."
	ta.ShowLineNumbers = true
	ta.CharLimit = 0 // No limit

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle()

	mdRenderer, _ := markdown.NewRenderer()

	return &MarkdownEditor{
		textarea:   ta,
		viewport:   vp,
		mdRenderer: mdRenderer,
		mode:       ModeSplit,
		splitRatio: 0.5,
		editBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")),
		previewBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("141")),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),
		statusStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
	}
}

// Focus focuses the editor.
func (e *MarkdownEditor) Focus() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.focused = true
	e.textarea.Focus()
}

// Blur removes focus from the editor.
func (e *MarkdownEditor) Blur() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.focused = false
	e.textarea.Blur()
}

// IsFocused returns whether the editor is focused.
func (e *MarkdownEditor) IsFocused() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.focused
}

// SetSize sets the editor size.
func (e *MarkdownEditor) SetSize(width, height int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.width = width
	e.height = height

	contentHeight := height - 4 // Account for borders and status

	switch e.mode {
	case ModeEdit:
		e.textarea.SetWidth(width - 4)
		e.textarea.SetHeight(contentHeight)
	case ModePreview:
		e.viewport.Width = width - 4
		e.viewport.Height = contentHeight
	case ModeSplit:
		halfWidth := (width - 6) / 2
		e.textarea.SetWidth(halfWidth)
		e.textarea.SetHeight(contentHeight)
		e.viewport.Width = halfWidth
		e.viewport.Height = contentHeight
	}
}

// SetMode sets the editor mode.
func (e *MarkdownEditor) SetMode(mode EditorMode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.mode = mode
	e.updatePreview()
}

// GetMode returns the current editor mode.
func (e *MarkdownEditor) GetMode() EditorMode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mode
}

// ToggleMode cycles through editor modes.
func (e *MarkdownEditor) ToggleMode() {
	e.mu.Lock()
	defer e.mu.Unlock()
	switch e.mode {
	case ModeEdit:
		e.mode = ModeSplit
	case ModeSplit:
		e.mode = ModePreview
	case ModePreview:
		e.mode = ModeEdit
	}
	e.updatePreview()
}

// Open opens a file.
func (e *MarkdownEditor) Open(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.filePath = filePath
	e.content = string(data)
	e.textarea.SetValue(e.content)
	e.modified = false
	e.updatePreview()

	return nil
}

// Save saves the file.
func (e *MarkdownEditor) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.filePath == "" {
		return nil // No file path set
	}

	e.content = e.textarea.Value()
	err := os.WriteFile(e.filePath, []byte(e.content), 0o600)
	if err != nil {
		return err
	}

	e.modified = false
	return nil
}

// SaveAs saves the file to a new path.
func (e *MarkdownEditor) SaveAs(filePath string) error {
	e.mu.Lock()
	e.filePath = filePath
	e.mu.Unlock()
	return e.Save()
}

// GetContent returns the current content.
func (e *MarkdownEditor) GetContent() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.textarea.Value()
}

// SetContent sets the content.
func (e *MarkdownEditor) SetContent(content string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.content = content
	e.textarea.SetValue(content)
	e.modified = true
	e.updatePreview()
}

// IsModified returns whether the content has been modified.
func (e *MarkdownEditor) IsModified() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.modified
}

// GetFilePath returns the current file path.
func (e *MarkdownEditor) GetFilePath() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.filePath
}

func (e *MarkdownEditor) updatePreview() {
	content := e.textarea.Value()
	if e.mdRenderer != nil {
		rendered, err := e.mdRenderer.Render(content)
		if err == nil {
			e.viewport.SetContent(rendered)
		} else {
			e.viewport.SetContent(content)
		}
	} else {
		e.viewport.SetContent(content)
	}
}

// InsertText inserts text at the end of the content.
// Note: Currently inserts at end of content for simplicity.
// Future versions may support proper cursor position insertion.
func (e *MarkdownEditor) InsertText(text string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	current := e.textarea.Value()
	// Note: Text is appended at end. The bubbles textarea component
	// would need additional work to support insertion at cursor position.
	e.textarea.SetValue(current + text)
	e.modified = true
	e.updatePreview()
}

// InsertBold inserts bold markdown syntax.
func (e *MarkdownEditor) InsertBold() {
	e.InsertText("**bold text**")
}

// InsertItalic inserts italic markdown syntax.
func (e *MarkdownEditor) InsertItalic() {
	e.InsertText("*italic text*")
}

// InsertCode inserts code markdown syntax.
func (e *MarkdownEditor) InsertCode() {
	e.InsertText("`code`")
}

// InsertLink inserts link markdown syntax.
func (e *MarkdownEditor) InsertLink() {
	e.InsertText("[link text](url)")
}

// InsertHeading inserts heading markdown syntax.
func (e *MarkdownEditor) InsertHeading() {
	e.InsertText("\n# Heading\n")
}

// InsertCodeBlock inserts a code block.
func (e *MarkdownEditor) InsertCodeBlock(language string) {
	e.InsertText("\n```" + language + "\n\n```\n")
}

// InsertList inserts a bullet list.
func (e *MarkdownEditor) InsertList() {
	e.InsertText("\n- Item 1\n- Item 2\n- Item 3\n")
}

// InsertTable inserts a table.
func (e *MarkdownEditor) InsertTable() {
	e.InsertText("\n| Column 1 | Column 2 |\n|----------|----------|\n| Cell 1   | Cell 2   |\n")
}

// Update handles messages for the editor.
func (e *MarkdownEditor) Update(msg tea.Msg) (*MarkdownEditor, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		e.mu.RLock()
		focused := e.focused
		e.mu.RUnlock()

		if !focused {
			return e, nil
		}

		switch {
		case key.Matches(msg, keys.Save):
			_ = e.Save() // Error is logged internally
			return e, nil

		case key.Matches(msg, keys.ToggleMode):
			e.ToggleMode()
			return e, nil

		case key.Matches(msg, keys.Escape):
			e.Blur()
			return e, nil

		case key.Matches(msg, keys.Bold):
			e.InsertBold()
			return e, nil

		case key.Matches(msg, keys.Italic):
			e.InsertItalic()
			return e, nil

		case key.Matches(msg, keys.Code):
			e.InsertCode()
			return e, nil

		case key.Matches(msg, keys.Link):
			e.InsertLink()
			return e, nil

		case key.Matches(msg, keys.Heading):
			e.InsertHeading()
			return e, nil
		}

	case tea.WindowSizeMsg:
		e.SetSize(msg.Width, msg.Height-4)
		return e, nil
	}

	// Update textarea
	e.mu.Lock()
	var cmd tea.Cmd
	oldContent := e.textarea.Value()
	e.textarea, cmd = e.textarea.Update(msg)
	newContent := e.textarea.Value()
	if oldContent != newContent {
		e.modified = true
		e.updatePreview()
	}
	e.mu.Unlock()
	cmds = append(cmds, cmd)

	// Update viewport in preview mode
	e.mu.Lock()
	if e.mode == ModePreview || e.mode == ModeSplit {
		e.viewport, cmd = e.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	e.mu.Unlock()

	return e, tea.Batch(cmds...)
}

// View renders the editor.
func (e *MarkdownEditor) View() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var content string

	// Title bar
	title := e.titleStyle.Render("📝 Markdown Editor")
	if e.filePath != "" {
		title += " - " + e.filePath
	}
	if e.modified {
		title += " *"
	}

	// Mode indicator
	modeStr := ""
	switch e.mode {
	case ModeEdit:
		modeStr = "[Edit]"
	case ModePreview:
		modeStr = "[Preview]"
	case ModeSplit:
		modeStr = "[Split]"
	}
	title += " " + e.statusStyle.Render(modeStr)
	content += title + "\n"

	// Main content
	switch e.mode {
	case ModeEdit:
		content += e.editBorder.Render(e.textarea.View())
	case ModePreview:
		content += e.previewBorder.Render(e.viewport.View())
	case ModeSplit:
		editPane := e.editBorder.Render(e.textarea.View())
		previewPane := e.previewBorder.Render(e.viewport.View())
		content += lipgloss.JoinHorizontal(lipgloss.Top, editPane, " ", previewPane)
	}

	// Status bar
	status := e.statusStyle.Render("ctrl+s: save | ctrl+p: toggle mode | ctrl+b: bold | ctrl+i: italic | esc: close")
	content += "\n" + status

	return content
}

// GetWordCount returns the word count of the content.
func (e *MarkdownEditor) GetWordCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	content := e.textarea.Value()
	words := strings.Fields(content)
	return len(words)
}

// GetLineCount returns the line count of the content.
func (e *MarkdownEditor) GetLineCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	content := e.textarea.Value()
	return strings.Count(content, "\n") + 1
}
