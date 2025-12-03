package nlp

import (
	"context"
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// mockAgent implements core.AIAgent for testing.
type mockAgent struct {
	response string
	err      error
}

func (m *mockAgent) Name() string              { return "mock" }
func (m *mockAgent) Provider() core.AIProvider { return core.AIProviderOpenAI }
func (m *mockAgent) Query(ctx context.Context, prompt string) (string, error) {
	return m.response, m.err
}

func (m *mockAgent) StreamQuery(ctx context.Context, prompt string) (<-chan string, error) {
	ch := make(chan string, 1)
	ch <- m.response
	close(ch)
	return ch, m.err
}

func (m *mockAgent) SuggestCommand(ctx context.Context, description string) (string, error) {
	return m.response, m.err
}

func (m *mockAgent) ExplainCommand(ctx context.Context, command string) (string, error) {
	return "explanation", nil
}

func (m *mockAgent) FixError(ctx context.Context, command string, err string) (string, error) {
	return "fix", nil
}

func TestAITranslator_Translate(t *testing.T) {
	tests := []struct {
		name        string
		description string
		agentResp   string
		wantCommand string
		wantErr     bool
	}{
		{
			name:        "find large files",
			description: "find files larger than 100MB",
			agentResp:   "find . -type f -size +100M",
			wantCommand: "find . -type f -size +100M",
			wantErr:     false,
		},
		{
			name:        "list hidden files",
			description: "show all hidden files",
			agentResp:   "ls -la",
			wantCommand: "ls -la",
			wantErr:     false,
		},
		{
			name:        "with code blocks",
			description: "list files",
			agentResp:   "```bash\nls -la\n```",
			wantCommand: "ls -la",
			wantErr:     false,
		},
		{
			name:        "empty description",
			description: "",
			agentResp:   "",
			wantCommand: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &mockAgent{response: tt.agentResp}
			translator := NewAITranslator(agent)

			result, err := translator.Translate(context.Background(), tt.description)

			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.Command != tt.wantCommand {
				t.Errorf("Translate() command = %v, want %v", result.Command, tt.wantCommand)
			}
		})
	}
}

func TestAITranslator_TranslateWithContext(t *testing.T) {
	agent := &mockAgent{response: "ls -la"}
	translator := NewAITranslator(agent)

	ctx := &ShellContext{
		CWD:       "/home/user",
		ShellType: core.ShellTypeZsh,
	}

	result, err := translator.TranslateWithContext(context.Background(), "list files", ctx)
	if err != nil {
		t.Fatalf("TranslateWithContext() error = %v", err)
	}

	if result.Command != "ls -la" {
		t.Errorf("TranslateWithContext() command = %v, want %v", result.Command, "ls -la")
	}
}

func TestAITranslator_Confidence(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		wantMinConf float64
	}{
		{
			name:        "valid command",
			response:    "ls -la",
			wantMinConf: 0.5,
		},
		{
			name:        "short command",
			response:    "ls",
			wantMinConf: 0.5,
		},
		{
			name:        "placeholder response",
			response:    "echo 'This is a suggested command'",
			wantMinConf: 0.3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &mockAgent{response: tt.response}
			translator := NewAITranslator(agent)

			result, err := translator.Translate(context.Background(), "test")
			if err != nil {
				t.Fatalf("Translate() error = %v", err)
			}

			if result.Confidence < tt.wantMinConf {
				t.Errorf("Translate() confidence = %v, want >= %v", result.Confidence, tt.wantMinConf)
			}
		})
	}
}

func TestAITranslator_MinConfidence(t *testing.T) {
	agent := &mockAgent{response: "ls"}
	translator := NewAITranslator(agent)

	// Default should be 0.7
	if got := translator.GetMinConfidence(); got != 0.7 {
		t.Errorf("GetMinConfidence() = %v, want %v", got, 0.7)
	}

	// Set new value
	translator.SetMinConfidence(0.5)
	if got := translator.GetMinConfidence(); got != 0.5 {
		t.Errorf("GetMinConfidence() after set = %v, want %v", got, 0.5)
	}
}

func TestAITranslator_NoAgent(t *testing.T) {
	translator := NewAITranslator(nil)

	_, err := translator.Translate(context.Background(), "list files")
	if err != ErrNoAgent {
		t.Errorf("Translate() error = %v, want %v", err, ErrNoAgent)
	}
}

func TestAITranslator_Cache(t *testing.T) {
	agent := &mockAgent{response: "ls -la"}

	// Create translator for cache testing
	translator := &AITranslator{
		agent:         agent,
		minConfidence: 0.7,
		cache:         make(map[string]*TranslationResult),
		cacheSize:     100,
	}

	// First call should hit the agent
	_, err := translator.Translate(context.Background(), "list files")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}

	// Second call with same input should use cache
	result, err := translator.Translate(context.Background(), "list files")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}

	if result.Command != "ls -la" {
		t.Errorf("Cached Translate() command = %v, want %v", result.Command, "ls -la")
	}

	// Clear cache
	translator.ClearCache()

	// Next call should hit agent again
	_, err = translator.Translate(context.Background(), "list files")
	if err != nil {
		t.Fatalf("Translate() after clear error = %v", err)
	}
}

func TestCommonTranslations(t *testing.T) {
	// Verify common translations are defined
	expected := map[string]bool{
		"find large files":       true,
		"list hidden files":      true,
		"show disk usage":        true,
		"show running processes": true,
	}

	for phrase := range expected {
		if _, ok := CommonTranslations[phrase]; !ok {
			t.Errorf("CommonTranslations missing phrase: %s", phrase)
		}
	}
}
