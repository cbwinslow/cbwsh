// Package errorfix provides error fix suggestions for cbwsh.
package errorfix

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// Common errors.
var (
	ErrNoAgent      = errors.New("no AI agent configured")
	ErrNoSuggestion = errors.New("no fix suggestion available")
)

// ErrorPattern represents a known error pattern.
type ErrorPattern struct {
	// Pattern is the regex pattern to match.
	Pattern *regexp.Regexp
	// Category is the error category.
	Category string
	// Description describes the error.
	Description string
	// Fixes are potential fixes.
	Fixes []string
}

// FixSuggestion represents a fix suggestion for an error.
type FixSuggestion struct {
	// OriginalCommand is the command that failed.
	OriginalCommand string
	// OriginalError is the error message.
	OriginalError string
	// SuggestedFix is the suggested fix command.
	SuggestedFix string
	// Explanation explains why this fix should work.
	Explanation string
	// Confidence is the confidence score (0.0-1.0).
	Confidence float64
	// AlternativeFixes are other potential fixes.
	AlternativeFixes []string
}

// Fixer provides error fix suggestions.
type Fixer struct {
	mu       sync.RWMutex
	agent    core.AIAgent
	patterns []ErrorPattern
}

// NewFixer creates a new error fixer.
func NewFixer(agent core.AIAgent) *Fixer {
	f := &Fixer{
		agent:    agent,
		patterns: defaultPatterns(),
	}
	return f
}

// SuggestFix suggests a fix for a failed command.
func (f *Fixer) SuggestFix(ctx context.Context, command, errorMsg string) (*FixSuggestion, error) {
	// First try pattern matching for common errors
	if suggestion := f.patternMatch(command, errorMsg); suggestion != nil {
		return suggestion, nil
	}

	// Fall back to AI agent
	f.mu.RLock()
	agent := f.agent
	f.mu.RUnlock()

	if agent == nil {
		return nil, ErrNoAgent
	}

	fix, err := agent.FixError(ctx, command, errorMsg)
	if err != nil {
		return nil, fmt.Errorf("AI fix failed: %w", err)
	}

	return &FixSuggestion{
		OriginalCommand: command,
		OriginalError:   errorMsg,
		SuggestedFix:    fix,
		Explanation:     "AI-suggested fix based on error analysis",
		Confidence:      0.7,
	}, nil
}

// SuggestMultipleFixes returns multiple fix suggestions.
func (f *Fixer) SuggestMultipleFixes(ctx context.Context, command, errorMsg string, maxSuggestions int) ([]*FixSuggestion, error) {
	var suggestions []*FixSuggestion

	// Pattern-based suggestions
	if suggestion := f.patternMatch(command, errorMsg); suggestion != nil {
		suggestions = append(suggestions, suggestion)

		// Add alternatives from pattern match
		for _, alt := range suggestion.AlternativeFixes {
			if len(suggestions) >= maxSuggestions {
				break
			}
			suggestions = append(suggestions, &FixSuggestion{
				OriginalCommand: command,
				OriginalError:   errorMsg,
				SuggestedFix:    alt,
				Explanation:     "Alternative fix based on error pattern",
				Confidence:      0.6,
			})
		}
	}

	// AI-based suggestion if we have room
	if len(suggestions) < maxSuggestions {
		f.mu.RLock()
		agent := f.agent
		f.mu.RUnlock()

		if agent != nil {
			fix, err := agent.FixError(ctx, command, errorMsg)
			if err == nil && fix != "" {
				// Check if this fix is already suggested
				isDuplicate := false
				for _, s := range suggestions {
					if s.SuggestedFix == fix {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					suggestions = append(suggestions, &FixSuggestion{
						OriginalCommand: command,
						OriginalError:   errorMsg,
						SuggestedFix:    fix,
						Explanation:     "AI-suggested fix",
						Confidence:      0.7,
					})
				}
			}
		}
	}

	if len(suggestions) == 0 {
		return nil, ErrNoSuggestion
	}

	return suggestions, nil
}

// AddPattern adds a custom error pattern.
func (f *Fixer) AddPattern(pattern ErrorPattern) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.patterns = append(f.patterns, pattern)
}

// SetAgent sets the AI agent.
func (f *Fixer) SetAgent(agent core.AIAgent) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.agent = agent
}

func (f *Fixer) patternMatch(command, errorMsg string) *FixSuggestion {
	f.mu.RLock()
	patterns := f.patterns
	f.mu.RUnlock()

	errorLower := strings.ToLower(errorMsg)
	commandLower := strings.ToLower(command)

	for _, pattern := range patterns {
		if pattern.Pattern.MatchString(errorLower) {
			fixes := f.generateFixes(pattern, command, commandLower, errorMsg)
			if len(fixes) == 0 {
				continue
			}

			suggestion := &FixSuggestion{
				OriginalCommand: command,
				OriginalError:   errorMsg,
				SuggestedFix:    fixes[0],
				Explanation:     pattern.Description,
				Confidence:      0.85,
			}

			if len(fixes) > 1 {
				suggestion.AlternativeFixes = fixes[1:]
			}

			return suggestion
		}
	}

	return nil
}

func (f *Fixer) generateFixes(pattern ErrorPattern, command, commandLower, _ string) []string {
	var fixes []string

	for _, fix := range pattern.Fixes {
		generated := fix

		// Replace placeholders
		if strings.Contains(fix, "${command}") {
			generated = strings.ReplaceAll(generated, "${command}", command)
		}

		// Handle specific patterns
		switch pattern.Category {
		case "permission":
			if strings.HasPrefix(commandLower, "sudo ") {
				// Already has sudo, suggest other options
				generated = strings.TrimPrefix(command, "sudo ")
			} else {
				generated = "sudo " + command
			}
		case "not_found":
			// Extract command name
			parts := strings.Fields(command)
			if len(parts) > 0 {
				cmdName := parts[0]
				fixes = append(fixes, fmt.Sprintf("which %s", cmdName))
				fixes = append(fixes, fmt.Sprintf("type %s", cmdName))
				generated = fmt.Sprintf("command -v %s || echo 'Command not found'", cmdName)
			}
		case "typo":
			// Suggest common typo corrections
			typoFixes := commonTypoFixes(command)
			fixes = append(fixes, typoFixes...)
			if len(fixes) > 0 {
				generated = fixes[0]
				fixes = fixes[1:]
			}
		}

		if generated != "" && generated != command {
			fixes = append([]string{generated}, fixes...)
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, fix := range fixes {
		if !seen[fix] {
			seen[fix] = true
			unique = append(unique, fix)
		}
	}

	return unique
}

func defaultPatterns() []ErrorPattern {
	return []ErrorPattern{
		{
			Pattern:     regexp.MustCompile(`permission denied`),
			Category:    "permission",
			Description: "Permission denied - try with elevated privileges",
			Fixes:       []string{"sudo ${command}"},
		},
		{
			Pattern:     regexp.MustCompile(`command not found`),
			Category:    "not_found",
			Description: "Command not found - check if installed or path",
			Fixes:       []string{"which ${cmd}", "type ${cmd}"},
		},
		{
			Pattern:     regexp.MustCompile(`no such file or directory`),
			Category:    "file_not_found",
			Description: "File or directory not found - check path",
			Fixes:       []string{"ls -la", "pwd"},
		},
		{
			Pattern:     regexp.MustCompile(`connection refused`),
			Category:    "network",
			Description: "Connection refused - check if service is running",
			Fixes:       []string{"ping localhost", "netstat -tuln"},
		},
		{
			Pattern:     regexp.MustCompile(`disk.*full|no space left`),
			Category:    "disk",
			Description: "Disk full - free up space",
			Fixes:       []string{"df -h", "du -sh * | sort -hr | head -10"},
		},
		{
			Pattern:     regexp.MustCompile(`timeout|timed out`),
			Category:    "timeout",
			Description: "Operation timed out - increase timeout or check connectivity",
			Fixes:       []string{"ping -c 3 ${host}"},
		},
		{
			Pattern:     regexp.MustCompile(`syntax error`),
			Category:    "syntax",
			Description: "Syntax error - check command syntax",
			Fixes:       []string{"man ${cmd}", "${cmd} --help"},
		},
		{
			Pattern:     regexp.MustCompile(`port.*already.*use|address already in use`),
			Category:    "port",
			Description: "Port already in use - find and stop the process",
			Fixes:       []string{"lsof -i :${port}", "netstat -tuln | grep ${port}"},
		},
		{
			Pattern:     regexp.MustCompile(`authentication fail|auth.*fail|invalid.*password`),
			Category:    "auth",
			Description: "Authentication failed - check credentials",
			Fixes:       []string{},
		},
		{
			Pattern:     regexp.MustCompile(`cannot.*find.*module|module not found`),
			Category:    "module",
			Description: "Module not found - install dependencies",
			Fixes:       []string{"npm install", "pip install", "go mod tidy"},
		},
		{
			Pattern:     regexp.MustCompile(`git.*not a git repository`),
			Category:    "git",
			Description: "Not a git repository - initialize or navigate to one",
			Fixes:       []string{"git init", "cd .."},
		},
		{
			Pattern:     regexp.MustCompile(`docker.*daemon.*not running`),
			Category:    "docker",
			Description: "Docker daemon not running - start Docker",
			Fixes:       []string{"sudo systemctl start docker", "sudo service docker start"},
		},
		{
			Pattern:     regexp.MustCompile(`kubectl.*connection.*refused|kubernetes.*unreachable`),
			Category:    "kubernetes",
			Description: "Kubernetes cluster unreachable - check connection",
			Fixes:       []string{"kubectl cluster-info", "kubectl config current-context"},
		},
	}
}

func commonTypoFixes(command string) []string {
	typos := map[string]string{
		"gti":     "git",
		"got":     "git",
		"gi":      "git",
		"sl":      "ls",
		"cd..":    "cd ..",
		"cd-":     "cd -",
		"sudp":    "sudo",
		"suod":    "sudo",
		"mkdri":   "mkdir",
		"mkdier":  "mkdir",
		"mkadir":  "mkdir",
		"cta":     "cat",
		"tial":    "tail",
		"haed":    "head",
		"grpe":    "grep",
		"gerp":    "grep",
		"nmp":     "npm",
		"dcoker":  "docker",
		"dokcer":  "docker",
		"kubetcl": "kubectl",
		"kubeclt": "kubectl",
	}

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	var fixes []string
	for typo, correct := range typos {
		if strings.HasPrefix(parts[0], typo) {
			fixed := correct + strings.TrimPrefix(parts[0], typo)
			if len(parts) > 1 {
				fixed += " " + strings.Join(parts[1:], " ")
			}
			fixes = append(fixes, fixed)
		}
	}

	return fixes
}

// QuickFix provides a quick fix for common errors.
func QuickFix(command, errorMsg string) string {
	errorLower := strings.ToLower(errorMsg)

	// Permission denied -> add sudo
	if strings.Contains(errorLower, "permission denied") {
		if !strings.HasPrefix(command, "sudo ") {
			return "sudo " + command
		}
	}

	// Command not found -> suggest installation
	if strings.Contains(errorLower, "command not found") {
		parts := strings.Fields(command)
		if len(parts) > 0 {
			return fmt.Sprintf("which %s || echo 'Try: apt install %s'", parts[0], parts[0])
		}
	}

	// No space left -> show disk usage
	if strings.Contains(errorLower, "no space left") || strings.Contains(errorLower, "disk full") {
		return "df -h && du -sh * | sort -hr | head -10"
	}

	return ""
}
