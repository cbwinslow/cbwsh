// Package autocomplete provides autocompletion for cbwsh.
package autocomplete

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// Completer provides autocompletion suggestions.
type Completer struct {
	mu        sync.RWMutex
	providers []core.CompletionProvider
}

// NewCompleter creates a new completer.
func NewCompleter() *Completer {
	c := &Completer{
		providers: make([]core.CompletionProvider, 0),
	}

	// Add default providers
	c.AddProvider(&CommandProvider{})
	c.AddProvider(&FileProvider{})
	c.AddProvider(&HistoryProvider{})

	return c
}

// Complete returns completion suggestions for the given input.
func (c *Completer) Complete(input string, cursorPos int) ([]core.Suggestion, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var allSuggestions []core.Suggestion

	for _, provider := range c.providers {
		suggestions, err := provider.Provide(input, cursorPos)
		if err != nil {
			continue
		}
		allSuggestions = append(allSuggestions, suggestions...)
	}

	// Sort by relevance (exact prefix matches first)
	sort.Slice(allSuggestions, func(i, j int) bool {
		// Get the word being completed
		word := getWordAtCursor(input, cursorPos)

		iPrefix := strings.HasPrefix(strings.ToLower(allSuggestions[i].Text), strings.ToLower(word))
		jPrefix := strings.HasPrefix(strings.ToLower(allSuggestions[j].Text), strings.ToLower(word))

		if iPrefix != jPrefix {
			return iPrefix
		}

		return allSuggestions[i].Text < allSuggestions[j].Text
	})

	// Limit suggestions
	if len(allSuggestions) > 20 {
		allSuggestions = allSuggestions[:20]
	}

	return allSuggestions, nil
}

// AddProvider adds a completion provider.
func (c *Completer) AddProvider(provider core.CompletionProvider) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.providers = append(c.providers, provider)
}

// RemoveProvider removes a completion provider by name.
func (c *Completer) RemoveProvider(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, provider := range c.providers {
		if provider.Name() == name {
			c.providers = append(c.providers[:i], c.providers[i+1:]...)
			break
		}
	}
}

func getWordAtCursor(input string, cursorPos int) string {
	if cursorPos > len(input) {
		cursorPos = len(input)
	}

	// Find start of word
	start := cursorPos
	for start > 0 && !isDelimiter(input[start-1]) {
		start--
	}

	// Find end of word
	end := cursorPos
	for end < len(input) && !isDelimiter(input[end]) {
		end++
	}

	return input[start:end]
}

func isDelimiter(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '|' || c == '&' || c == ';'
}

// CommandProvider provides command completions.
type CommandProvider struct {
	mu       sync.RWMutex
	commands []string
}

// Name returns the provider name.
func (p *CommandProvider) Name() string {
	return "commands"
}

// Provide returns command suggestions.
func (p *CommandProvider) Provide(input string, cursorPos int) ([]core.Suggestion, error) {
	word := getWordAtCursor(input, cursorPos)
	if word == "" {
		return nil, nil
	}

	// Get commands from PATH
	commands := p.getCommandsFromPath()

	var suggestions []core.Suggestion
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, word) {
			suggestions = append(suggestions, core.Suggestion{
				Text:        cmd,
				Description: "Command",
				Category:    "command",
			})
		}
	}

	return suggestions, nil
}

func (p *CommandProvider) getCommandsFromPath() []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.commands) > 0 {
		return p.commands
	}

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	seen := make(map[string]bool)

	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			if seen[name] {
				continue
			}

			// Check if executable
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.Mode()&0o111 != 0 {
				p.commands = append(p.commands, name)
				seen[name] = true
			}
		}
	}

	sort.Strings(p.commands)
	return p.commands
}

// FileProvider provides file path completions.
type FileProvider struct{}

// Name returns the provider name.
func (p *FileProvider) Name() string {
	return "files"
}

// Provide returns file path suggestions.
func (p *FileProvider) Provide(input string, cursorPos int) ([]core.Suggestion, error) {
	word := getWordAtCursor(input, cursorPos)

	// Only complete if there's a path-like pattern
	if !strings.Contains(word, "/") && !strings.HasPrefix(word, ".") && !strings.HasPrefix(word, "~") {
		return nil, nil
	}

	// Expand home directory
	path := word
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = home + path[1:]
		}
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// If path ends with /, list directory contents
	if strings.HasSuffix(word, "/") {
		dir = path
		base = ""
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	var suggestions []core.Suggestion
	for _, entry := range entries {
		name := entry.Name()
		if base != "" && !strings.HasPrefix(name, base) {
			continue
		}

		fullPath := filepath.Join(dir, name)

		// Convert back to original format
		if strings.HasPrefix(word, "~") {
			home, _ := os.UserHomeDir()
			fullPath = "~" + strings.TrimPrefix(fullPath, home)
		}

		// Add trailing slash for directories
		if entry.IsDir() {
			fullPath += "/"
		}

		description := "File"
		if entry.IsDir() {
			description = "Directory"
		}

		suggestions = append(suggestions, core.Suggestion{
			Text:        fullPath,
			Description: description,
			Category:    "file",
		})
	}

	return suggestions, nil
}

// HistoryProvider provides history-based completions.
type HistoryProvider struct {
	mu      sync.RWMutex
	history []string
}

// Name returns the provider name.
func (p *HistoryProvider) Name() string {
	return "history"
}

// Provide returns history-based suggestions.
func (p *HistoryProvider) Provide(input string, cursorPos int) ([]core.Suggestion, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if input == "" {
		return nil, nil
	}

	var suggestions []core.Suggestion
	seen := make(map[string]bool)

	for i := len(p.history) - 1; i >= 0; i-- {
		cmd := p.history[i]
		if strings.HasPrefix(cmd, input) && !seen[cmd] {
			suggestions = append(suggestions, core.Suggestion{
				Text:        cmd,
				Description: "History",
				Category:    "history",
			})
			seen[cmd] = true
		}

		if len(suggestions) >= 5 {
			break
		}
	}

	return suggestions, nil
}

// AddHistory adds a command to history.
func (p *HistoryProvider) AddHistory(command string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.history = append(p.history, command)
}

// SetHistory sets the history.
func (p *HistoryProvider) SetHistory(history []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.history = make([]string, len(history))
	copy(p.history, history)
}

// AliasProvider provides alias completions.
type AliasProvider struct {
	mu      sync.RWMutex
	aliases map[string]string
}

// NewAliasProvider creates a new alias provider.
func NewAliasProvider() *AliasProvider {
	return &AliasProvider{
		aliases: make(map[string]string),
	}
}

// Name returns the provider name.
func (p *AliasProvider) Name() string {
	return "aliases"
}

// Provide returns alias suggestions.
func (p *AliasProvider) Provide(input string, cursorPos int) ([]core.Suggestion, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	word := getWordAtCursor(input, cursorPos)

	var suggestions []core.Suggestion
	for alias, expansion := range p.aliases {
		if strings.HasPrefix(alias, word) {
			suggestions = append(suggestions, core.Suggestion{
				Text:        alias,
				Description: "Alias: " + expansion,
				Category:    "alias",
			})
		}
	}

	return suggestions, nil
}

// SetAliases sets the aliases.
func (p *AliasProvider) SetAliases(aliases map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.aliases = make(map[string]string)
	for k, v := range aliases {
		p.aliases[k] = v
	}
}

// EnvProvider provides environment variable completions.
type EnvProvider struct{}

// Name returns the provider name.
func (p *EnvProvider) Name() string {
	return "environment"
}

// Provide returns environment variable suggestions.
func (p *EnvProvider) Provide(input string, cursorPos int) ([]core.Suggestion, error) {
	word := getWordAtCursor(input, cursorPos)

	// Only complete if starting with $
	if !strings.HasPrefix(word, "$") {
		return nil, nil
	}

	prefix := strings.TrimPrefix(word, "$")

	var suggestions []core.Suggestion
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		value := parts[1]

		if strings.HasPrefix(name, prefix) {
			// Truncate long values
			if len(value) > 30 {
				value = value[:27] + "..."
			}

			suggestions = append(suggestions, core.Suggestion{
				Text:        "$" + name,
				Description: value,
				Category:    "env",
			})
		}
	}

	return suggestions, nil
}

// BuiltinProvider provides shell builtin completions.
type BuiltinProvider struct{}

// Name returns the provider name.
func (p *BuiltinProvider) Name() string {
	return "builtins"
}

var builtins = []struct {
	name string
	desc string
}{
	{"cd", "Change directory"},
	{"pwd", "Print working directory"},
	{"echo", "Display text"},
	{"exit", "Exit shell"},
	{"export", "Set environment variable"},
	{"unset", "Unset variable"},
	{"alias", "Create alias"},
	{"unalias", "Remove alias"},
	{"source", "Execute file"},
	{".", "Execute file"},
	{"history", "Show command history"},
	{"clear", "Clear screen"},
	{"help", "Show help"},
}

// Provide returns builtin suggestions.
func (p *BuiltinProvider) Provide(input string, cursorPos int) ([]core.Suggestion, error) {
	word := getWordAtCursor(input, cursorPos)

	var suggestions []core.Suggestion
	for _, builtin := range builtins {
		if strings.HasPrefix(builtin.name, word) {
			suggestions = append(suggestions, core.Suggestion{
				Text:        builtin.name,
				Description: builtin.desc,
				Category:    "builtin",
			})
		}
	}

	return suggestions, nil
}
