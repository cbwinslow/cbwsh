// Package highlight provides syntax highlighting for cbwsh.
package highlight

import (
	"strings"
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// Highlighter provides syntax highlighting.
type Highlighter struct {
	mu        sync.RWMutex
	styleName string
	style     *chroma.Style
	formatter chroma.Formatter
}

// NewHighlighter creates a new highlighter.
func NewHighlighter() *Highlighter {
	return &Highlighter{
		styleName: "monokai",
		style:     styles.Get("monokai"),
		formatter: formatters.TTY256,
	}
}

// NewHighlighterWithStyle creates a highlighter with a specific style.
func NewHighlighterWithStyle(styleName string) *Highlighter {
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}

	return &Highlighter{
		styleName: styleName,
		style:     style,
		formatter: formatters.TTY256,
	}
}

// Highlight applies syntax highlighting to text.
func (h *Highlighter) Highlight(text, language string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, text)
	if err != nil {
		return text, err
	}

	var buf strings.Builder
	err = h.formatter.Format(&buf, h.style, iterator)
	if err != nil {
		return text, err
	}

	return buf.String(), nil
}

// HighlightCommand highlights a shell command.
func (h *Highlighter) HighlightCommand(command string) (string, error) {
	return h.Highlight(command, "bash")
}

// SetStyle sets the highlighting style.
func (h *Highlighter) SetStyle(styleName string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	style := styles.Get(styleName)
	if style != nil {
		h.styleName = styleName
		h.style = style
	}
}

// GetStyle returns the current style name.
func (h *Highlighter) GetStyle() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.styleName
}

// ListStyles returns all available style names.
func ListStyles() []string {
	return styles.Names()
}

// ShellHighlighter provides specialized shell command highlighting.
type ShellHighlighter struct {
	mu              sync.RWMutex
	commandStyle    lipgloss.Style
	argumentStyle   lipgloss.Style
	optionStyle     lipgloss.Style
	stringStyle     lipgloss.Style
	variableStyle   lipgloss.Style
	operatorStyle   lipgloss.Style
	commentStyle    lipgloss.Style
	errorStyle      lipgloss.Style
	builtinCommands map[string]bool
}

// NewShellHighlighter creates a new shell highlighter.
func NewShellHighlighter() *ShellHighlighter {
	return &ShellHighlighter{
		commandStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("86")),  // Cyan
		argumentStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("252")), // Light gray
		optionStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("214")), // Orange
		stringStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("114")), // Green
		variableStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("141")), // Purple
		operatorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("203")), // Red
		commentStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")), // Dark gray
		errorStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("196")), // Bright red
		builtinCommands: map[string]bool{
			"cd": true, "pwd": true, "echo": true, "exit": true, "export": true,
			"source": true, "alias": true, "unalias": true, "history": true,
			"clear": true, "set": true, "unset": true, "eval": true, "exec": true,
			"return": true, "shift": true, "test": true, "read": true, "printf": true,
			"local": true, "declare": true, "typeset": true, "readonly": true,
			"trap": true, "wait": true, "jobs": true, "fg": true, "bg": true,
			"kill": true, "pushd": true, "popd": true, "dirs": true,
		},
	}
}

// HighlightCommand highlights a shell command with semantic coloring.
func (h *ShellHighlighter) HighlightCommand(command string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if command == "" {
		return ""
	}

	// Handle comments
	if strings.HasPrefix(strings.TrimSpace(command), "#") {
		return h.commentStyle.Render(command)
	}

	tokens := tokenize(command)
	if len(tokens) == 0 {
		return command
	}

	var result strings.Builder
	isFirstWord := true
	inString := false
	stringChar := byte(0)

	for i, token := range tokens {
		// Track string state
		for j := 0; j < len(token); j++ {
			c := token[j]
			if !inString && (c == '"' || c == '\'') {
				inString = true
				stringChar = c
			} else if inString && c == stringChar {
				inString = false
			}
		}

		// Apply highlighting
		highlighted := h.highlightToken(token, isFirstWord, i == 0)

		result.WriteString(highlighted)

		if isFirstWord && token != "" && !isWhitespace(token) {
			isFirstWord = false
		}
	}

	return result.String()
}

func (h *ShellHighlighter) highlightToken(token string, isFirstWord, isVeryFirst bool) string {
	if isWhitespace(token) {
		return token
	}

	// Comments
	if strings.HasPrefix(token, "#") {
		return h.commentStyle.Render(token)
	}

	// Strings
	if (strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"")) ||
		(strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'")) {
		return h.stringStyle.Render(token)
	}

	// Variables
	if strings.HasPrefix(token, "$") {
		return h.variableStyle.Render(token)
	}

	// Operators and redirections
	if isOperator(token) {
		return h.operatorStyle.Render(token)
	}

	// Options (flags)
	if strings.HasPrefix(token, "-") {
		return h.optionStyle.Render(token)
	}

	// Commands (first word)
	if isFirstWord && isVeryFirst {
		if h.builtinCommands[token] {
			return h.commandStyle.Bold(true).Render(token)
		}
		return h.commandStyle.Render(token)
	}

	// Arguments
	return h.argumentStyle.Render(token)
}

func tokenize(command string) []string {
	var tokens []string
	var current strings.Builder
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(command); i++ {
		c := command[i]

		if !inString && (c == '"' || c == '\'') {
			inString = true
			stringChar = c
			current.WriteByte(c)
		} else if inString && c == stringChar {
			current.WriteByte(c)
			inString = false
		} else if !inString && (c == ' ' || c == '\t') {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(c))
		} else if !inString && isOperatorChar(c) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			// Handle multi-character operators
			if i+1 < len(command) && isOperatorChar(command[i+1]) {
				tokens = append(tokens, string(c)+string(command[i+1]))
				i++
			} else {
				tokens = append(tokens, string(c))
			}
		} else {
			current.WriteByte(c)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func isWhitespace(s string) bool {
	for _, c := range s {
		if c != ' ' && c != '\t' && c != '\n' {
			return false
		}
	}
	return true
}

func isOperator(s string) bool {
	operators := []string{
		"|", "||", "&", "&&", ";", ";;",
		">", ">>", "<", "<<", "<<<",
		"2>", "2>>", "&>", "&>>",
		"$(", "(", ")", "{", "}", "[", "]", "[[", "]]",
	}
	for _, op := range operators {
		if s == op {
			return true
		}
	}
	return false
}

func isOperatorChar(c byte) bool {
	return c == '|' || c == '&' || c == ';' || c == '>' || c == '<' ||
		c == '(' || c == ')' || c == '{' || c == '}' || c == '[' || c == ']'
}

// ConditionalHighlighter applies conditional highlighting.
type ConditionalHighlighter struct {
	mu           sync.RWMutex
	conditions   map[string]lipgloss.Style
	defaultStyle lipgloss.Style
}

// NewConditionalHighlighter creates a new conditional highlighter.
func NewConditionalHighlighter() *ConditionalHighlighter {
	return &ConditionalHighlighter{
		conditions:   make(map[string]lipgloss.Style),
		defaultStyle: lipgloss.NewStyle(),
	}
}

// AddCondition adds a highlighting condition.
func (h *ConditionalHighlighter) AddCondition(pattern string, style lipgloss.Style) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conditions[pattern] = style
}

// RemoveCondition removes a highlighting condition.
func (h *ConditionalHighlighter) RemoveCondition(pattern string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conditions, pattern)
}

// SetDefaultStyle sets the default style.
func (h *ConditionalHighlighter) SetDefaultStyle(style lipgloss.Style) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.defaultStyle = style
}

// Highlight applies conditional highlighting.
func (h *ConditionalHighlighter) Highlight(text string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for pattern, style := range h.conditions {
		if strings.Contains(text, pattern) {
			return style.Render(text)
		}
	}

	return h.defaultStyle.Render(text)
}

// ErrorHighlighter highlights error messages.
type ErrorHighlighter struct {
	errorStyle   lipgloss.Style
	warningStyle lipgloss.Style
	successStyle lipgloss.Style
}

// NewErrorHighlighter creates a new error highlighter.
func NewErrorHighlighter() *ErrorHighlighter {
	return &ErrorHighlighter{
		errorStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
		warningStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("214")),
		successStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("82")),
	}
}

// Highlight highlights based on content.
func (h *ErrorHighlighter) Highlight(text string) string {
	lower := strings.ToLower(text)

	if strings.Contains(lower, "error") || strings.Contains(lower, "fail") ||
		strings.Contains(lower, "fatal") || strings.Contains(lower, "exception") {
		return h.errorStyle.Render(text)
	}

	if strings.Contains(lower, "warning") || strings.Contains(lower, "warn") {
		return h.warningStyle.Render(text)
	}

	if strings.Contains(lower, "success") || strings.Contains(lower, "ok") ||
		strings.Contains(lower, "done") || strings.Contains(lower, "passed") {
		return h.successStyle.Render(text)
	}

	return text
}
