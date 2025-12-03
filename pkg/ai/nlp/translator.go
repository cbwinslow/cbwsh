// Package nlp provides natural language to command translation for cbwsh.
package nlp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// Common errors.
var (
	ErrNoAgent         = errors.New("no AI agent configured")
	ErrLowConfidence   = errors.New("low confidence in translation")
	ErrEmptyInput      = errors.New("empty input")
	ErrTranslationFail = errors.New("translation failed")
)

// TranslationResult represents the result of translating natural language to a command.
type TranslationResult struct {
	// Command is the translated shell command.
	Command string
	// Explanation describes what the command does.
	Explanation string
	// Confidence is the confidence score (0.0-1.0).
	Confidence float64
	// Alternatives are alternative command suggestions.
	Alternatives []string
	// Warnings contains any warnings about the command.
	Warnings []string
}

// ShellContext provides context for more accurate translations.
type ShellContext struct {
	// CWD is the current working directory.
	CWD string
	// ShellType is the shell type (bash, zsh).
	ShellType core.ShellType
	// EnvVars are relevant environment variables.
	EnvVars map[string]string
	// RecentCommands are recently executed commands.
	RecentCommands []string
	// ProjectType is the detected project type.
	ProjectType string
}

// Translator translates natural language to shell commands.
type Translator interface {
	// Translate converts a natural language description to a shell command.
	Translate(ctx context.Context, description string) (*TranslationResult, error)
	// TranslateWithContext translates with additional shell context.
	TranslateWithContext(ctx context.Context, description string, shellCtx *ShellContext) (*TranslationResult, error)
	// SetMinConfidence sets the minimum confidence threshold.
	SetMinConfidence(confidence float64)
	// GetMinConfidence returns the current minimum confidence.
	GetMinConfidence() float64
}

// AITranslator implements Translator using an AI agent.
type AITranslator struct {
	mu            sync.RWMutex
	agent         core.AIAgent
	minConfidence float64
	cache         map[string]*TranslationResult
	cacheSize     int
}

// NewAITranslator creates a new AI-powered translator.
func NewAITranslator(agent core.AIAgent) *AITranslator {
	return &AITranslator{
		agent:         agent,
		minConfidence: 0.7,
		cache:         make(map[string]*TranslationResult),
		cacheSize:     100,
	}
}

// Translate converts natural language to a shell command.
func (t *AITranslator) Translate(ctx context.Context, description string) (*TranslationResult, error) {
	return t.TranslateWithContext(ctx, description, nil)
}

// TranslateWithContext translates with additional shell context.
func (t *AITranslator) TranslateWithContext(ctx context.Context, description string, shellCtx *ShellContext) (*TranslationResult, error) {
	if description == "" {
		return nil, ErrEmptyInput
	}

	t.mu.RLock()
	agent := t.agent
	t.mu.RUnlock()

	if agent == nil {
		return nil, ErrNoAgent
	}

	// Check cache
	cacheKey := t.cacheKey(description, shellCtx)
	t.mu.RLock()
	if cached, ok := t.cache[cacheKey]; ok {
		t.mu.RUnlock()
		return cached, nil
	}
	t.mu.RUnlock()

	// Build prompt
	prompt := t.buildPrompt(description, shellCtx)

	// Query AI agent
	response, err := agent.SuggestCommand(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTranslationFail, err)
	}

	// Parse response
	result := t.parseResponse(response, description)

	// Validate confidence
	if result.Confidence < t.minConfidence {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Low confidence: %.1f%%", result.Confidence*100))
	}

	// Cache result
	t.cacheResult(cacheKey, result)

	return result, nil
}

// SetMinConfidence sets the minimum confidence threshold.
func (t *AITranslator) SetMinConfidence(confidence float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.minConfidence = confidence
}

// GetMinConfidence returns the current minimum confidence.
func (t *AITranslator) GetMinConfidence() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.minConfidence
}

// SetAgent sets the AI agent.
func (t *AITranslator) SetAgent(agent core.AIAgent) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.agent = agent
}

// ClearCache clears the translation cache.
func (t *AITranslator) ClearCache() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cache = make(map[string]*TranslationResult)
}

func (t *AITranslator) buildPrompt(description string, shellCtx *ShellContext) string {
	var sb strings.Builder

	sb.WriteString(description)

	if shellCtx != nil {
		if shellCtx.CWD != "" {
			sb.WriteString(fmt.Sprintf(" (current directory: %s)", shellCtx.CWD))
		}
		if shellCtx.ShellType != core.ShellTypeBash {
			sb.WriteString(fmt.Sprintf(" (shell: %s)", shellCtx.ShellType.String()))
		}
		if shellCtx.ProjectType != "" {
			sb.WriteString(fmt.Sprintf(" (project type: %s)", shellCtx.ProjectType))
		}
	}

	return sb.String()
}

func (t *AITranslator) parseResponse(response, originalDescription string) *TranslationResult {
	// Clean up the response
	command := strings.TrimSpace(response)

	// Remove markdown code blocks if present
	command = strings.TrimPrefix(command, "```bash\n")
	command = strings.TrimPrefix(command, "```sh\n")
	command = strings.TrimPrefix(command, "```\n")
	command = strings.TrimSuffix(command, "\n```")
	command = strings.TrimSuffix(command, "```")
	command = strings.TrimSpace(command)

	// Take only the first line if multiple lines
	if idx := strings.Index(command, "\n"); idx > 0 {
		command = command[:idx]
	}

	result := &TranslationResult{
		Command:    command,
		Confidence: t.estimateConfidence(command, originalDescription),
	}

	// Generate explanation
	result.Explanation = t.generateExplanation(command)

	return result
}

func (t *AITranslator) estimateConfidence(command, _ string) float64 {
	// Basic confidence estimation based on command structure
	confidence := 0.8

	// Lower confidence for empty or very short commands
	if len(command) < 3 {
		confidence -= 0.3
	}

	// Lower confidence for commands with echo only (might be placeholder)
	if strings.HasPrefix(command, "echo ") && strings.Contains(command, "suggested") {
		confidence -= 0.4
	}

	// Higher confidence for common commands
	commonPrefixes := []string{"ls", "cd", "grep", "find", "cat", "rm", "mv", "cp", "mkdir", "git", "docker", "kubectl"}
	for _, prefix := range commonPrefixes {
		if strings.HasPrefix(command, prefix+" ") || command == prefix {
			confidence += 0.1
			break
		}
	}

	// Clamp to valid range
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

func (t *AITranslator) generateExplanation(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}

	// Basic explanations for common commands
	explanations := map[string]string{
		"ls":      "List directory contents",
		"cd":      "Change directory",
		"grep":    "Search for patterns in files",
		"find":    "Search for files and directories",
		"cat":     "Display file contents",
		"rm":      "Remove files or directories",
		"mv":      "Move or rename files",
		"cp":      "Copy files or directories",
		"mkdir":   "Create directories",
		"git":     "Git version control command",
		"docker":  "Docker container command",
		"kubectl": "Kubernetes command",
	}

	if explanation, ok := explanations[parts[0]]; ok {
		return explanation
	}

	return fmt.Sprintf("Execute %s command", parts[0])
}

func (t *AITranslator) cacheKey(description string, shellCtx *ShellContext) string {
	key := strings.ToLower(strings.TrimSpace(description))
	if shellCtx != nil && shellCtx.ShellType != core.ShellTypeBash {
		key += "|" + shellCtx.ShellType.String()
	}
	return key
}

func (t *AITranslator) cacheResult(key string, result *TranslationResult) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Simple LRU-like eviction
	if len(t.cache) >= t.cacheSize {
		// Delete a random entry
		for k := range t.cache {
			delete(t.cache, k)
			break
		}
	}

	t.cache[key] = result
}

// CommonTranslations provides pre-defined translations for common phrases.
var CommonTranslations = map[string]string{
	"find large files":           "find . -type f -size +100M",
	"list hidden files":          "ls -la",
	"show disk usage":            "df -h",
	"show folder sizes":          "du -sh *",
	"find text in files":         "grep -r 'text' .",
	"count lines in file":        "wc -l filename",
	"show running processes":     "ps aux",
	"kill process":               "kill -9 PID",
	"show network connections":   "netstat -tuln",
	"download file":              "curl -O URL",
	"compress folder":            "tar -czvf archive.tar.gz folder/",
	"extract archive":            "tar -xzvf archive.tar.gz",
	"find and replace in files":  "sed -i 's/old/new/g' file",
	"show environment variables": "env",
	"show current directory":     "pwd",
}
