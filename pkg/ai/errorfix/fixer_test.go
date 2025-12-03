package errorfix

import (
	"context"
	"regexp"
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// mockAgent implements core.AIAgent for testing.
type mockAgent struct {
	fixResponse string
	err         error
}

func (m *mockAgent) Name() string              { return "mock" }
func (m *mockAgent) Provider() core.AIProvider { return core.AIProviderOpenAI }
func (m *mockAgent) Query(ctx context.Context, prompt string) (string, error) {
	return m.fixResponse, m.err
}

func (m *mockAgent) StreamQuery(ctx context.Context, prompt string) (<-chan string, error) {
	ch := make(chan string, 1)
	ch <- m.fixResponse
	close(ch)
	return ch, m.err
}

func (m *mockAgent) SuggestCommand(ctx context.Context, description string) (string, error) {
	return m.fixResponse, m.err
}

func (m *mockAgent) ExplainCommand(ctx context.Context, command string) (string, error) {
	return "explanation", nil
}

func (m *mockAgent) FixError(ctx context.Context, command string, err string) (string, error) {
	return m.fixResponse, m.err
}

func TestFixer_SuggestFix_PermissionDenied(t *testing.T) {
	fixer := NewFixer(nil)

	suggestion, err := fixer.SuggestFix(context.Background(), "cat /etc/shadow", "permission denied")
	if err != nil {
		t.Fatalf("SuggestFix() error = %v", err)
	}

	if suggestion == nil {
		t.Fatal("SuggestFix() returned nil suggestion")
	}

	if suggestion.SuggestedFix != "sudo cat /etc/shadow" {
		t.Errorf("SuggestFix() = %v, want sudo cat /etc/shadow", suggestion.SuggestedFix)
	}

	if suggestion.Confidence < 0.8 {
		t.Errorf("SuggestFix() confidence = %v, want >= 0.8", suggestion.Confidence)
	}
}

func TestFixer_SuggestFix_CommandNotFound(t *testing.T) {
	fixer := NewFixer(nil)

	suggestion, err := fixer.SuggestFix(context.Background(), "nonexistent", "command not found")
	if err != nil {
		t.Fatalf("SuggestFix() error = %v", err)
	}

	if suggestion == nil {
		t.Fatal("SuggestFix() returned nil suggestion")
	}

	// Should suggest checking if command exists
	if suggestion.SuggestedFix == "" {
		t.Error("SuggestFix() should return a non-empty fix")
	}
}

func TestFixer_SuggestFix_NoSpaceLeft(t *testing.T) {
	fixer := NewFixer(nil)

	suggestion, err := fixer.SuggestFix(context.Background(), "cp large.file /tmp", "no space left on device")
	if err != nil {
		t.Fatalf("SuggestFix() error = %v", err)
	}

	if suggestion == nil {
		t.Fatal("SuggestFix() returned nil suggestion")
	}

	// Should suggest disk usage commands
	if suggestion.SuggestedFix == "" {
		t.Error("SuggestFix() should return a non-empty fix")
	}
}

func TestFixer_SuggestFix_WithAgent(t *testing.T) {
	agent := &mockAgent{fixResponse: "corrected command"}
	fixer := NewFixer(agent)

	// Use an error that doesn't match patterns
	suggestion, err := fixer.SuggestFix(context.Background(), "some command", "unknown error xyz123")
	if err != nil {
		t.Fatalf("SuggestFix() error = %v", err)
	}

	if suggestion == nil {
		t.Fatal("SuggestFix() returned nil suggestion")
	}

	if suggestion.SuggestedFix != "corrected command" {
		t.Errorf("SuggestFix() = %v, want corrected command", suggestion.SuggestedFix)
	}
}

func TestFixer_SuggestFix_NoAgentNoPattern(t *testing.T) {
	fixer := NewFixer(nil)

	_, err := fixer.SuggestFix(context.Background(), "command", "completely unknown error type xyz")
	if err != ErrNoAgent {
		t.Errorf("SuggestFix() error = %v, want %v", err, ErrNoAgent)
	}
}

func TestFixer_SuggestMultipleFixes(t *testing.T) {
	agent := &mockAgent{fixResponse: "ai fix"}
	fixer := NewFixer(agent)

	suggestions, err := fixer.SuggestMultipleFixes(context.Background(), "cat /etc/shadow", "permission denied", 5)
	if err != nil {
		t.Fatalf("SuggestMultipleFixes() error = %v", err)
	}

	if len(suggestions) == 0 {
		t.Fatal("SuggestMultipleFixes() returned no suggestions")
	}

	// First suggestion should be pattern-based
	if suggestions[0].SuggestedFix != "sudo cat /etc/shadow" {
		t.Errorf("First suggestion = %v, want sudo cat /etc/shadow", suggestions[0].SuggestedFix)
	}
}

func TestFixer_AddPattern(t *testing.T) {
	fixer := NewFixer(nil)

	// Add custom pattern
	fixer.AddPattern(ErrorPattern{
		Pattern:     mustCompile(`custom error pattern`),
		Category:    "custom",
		Description: "Custom error handler",
		Fixes:       []string{"custom fix"},
	})

	suggestion, err := fixer.SuggestFix(context.Background(), "cmd", "custom error pattern encountered")
	if err != nil {
		t.Fatalf("SuggestFix() error = %v", err)
	}

	if suggestion == nil {
		t.Fatal("SuggestFix() returned nil suggestion for custom pattern")
	}
}

func TestQuickFix(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		errorMsg string
		wantFix  string
	}{
		{
			name:     "permission denied",
			command:  "cat /etc/shadow",
			errorMsg: "permission denied",
			wantFix:  "sudo cat /etc/shadow",
		},
		{
			name:     "already has sudo",
			command:  "sudo cat /etc/shadow",
			errorMsg: "permission denied",
			wantFix:  "",
		},
		{
			name:     "disk full",
			command:  "cp big.file /tmp",
			errorMsg: "no space left on device",
			wantFix:  "df -h && du -sh * | sort -hr | head -10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QuickFix(tt.command, tt.errorMsg)
			if got != tt.wantFix {
				t.Errorf("QuickFix() = %v, want %v", got, tt.wantFix)
			}
		})
	}
}

func TestCommonTypoFixes(t *testing.T) {
	tests := []struct {
		command string
		wantFix string
	}{
		{"gti status", "git status"},
		{"sl -la", "ls -la"},
		{"cd..", "cd .."},
		{"dcoker ps", "docker ps"},
		{"kubetcl get pods", "kubectl get pods"},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			fixes := commonTypoFixes(tt.command)

			found := false
			for _, fix := range fixes {
				if fix == tt.wantFix {
					found = true
					break
				}
			}

			if !found && tt.wantFix != "" {
				t.Errorf("commonTypoFixes(%s) did not include %s, got %v", tt.command, tt.wantFix, fixes)
			}
		})
	}
}

func TestDefaultPatterns(t *testing.T) {
	patterns := defaultPatterns()

	if len(patterns) == 0 {
		t.Fatal("defaultPatterns() returned empty list")
	}

	// Check for expected patterns
	categories := make(map[string]bool)
	for _, p := range patterns {
		categories[p.Category] = true
	}

	expectedCategories := []string{"permission", "not_found", "file_not_found", "network", "disk", "syntax", "git", "docker"}
	for _, cat := range expectedCategories {
		if !categories[cat] {
			t.Errorf("defaultPatterns() missing category: %s", cat)
		}
	}
}

// Helper to compile regex without error handling in tests
func mustCompile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
