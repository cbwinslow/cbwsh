// Package ai provides AI integration for cbwsh.
package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// Agent represents an AI agent that can assist with shell commands.
type Agent struct {
	mu       sync.RWMutex
	name     string
	provider core.AIProvider
	apiKey   string
	model    string
	enabled  bool
}

// NewAgent creates a new AI agent.
func NewAgent(name string, provider core.AIProvider, apiKey, model string) *Agent {
	return &Agent{
		name:     name,
		provider: provider,
		apiKey:   apiKey,
		model:    model,
		enabled:  true,
	}
}

// Name returns the agent name.
func (a *Agent) Name() string {
	return a.name
}

// Provider returns the AI provider.
func (a *Agent) Provider() core.AIProvider {
	return a.provider
}

// Query sends a query to the AI agent.
func (a *Agent) Query(_ context.Context, prompt string) (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if !a.enabled {
		return "", errors.New("agent is disabled")
	}

	// This is a placeholder - in a real implementation, you would call the actual AI API
	switch a.provider {
	case core.AIProviderOpenAI:
		return a.mockOpenAIQuery(prompt)
	case core.AIProviderAnthropic:
		return a.mockAnthropicQuery(prompt)
	case core.AIProviderGemini:
		return a.mockGeminiQuery(prompt)
	case core.AIProviderLocal:
		return a.mockLocalQuery(prompt)
	default:
		return "", errors.New("no AI provider configured")
	}
}

// StreamQuery sends a query and streams the response.
func (a *Agent) StreamQuery(ctx context.Context, prompt string) (<-chan string, error) {
	a.mu.RLock()
	if !a.enabled {
		a.mu.RUnlock()
		return nil, errors.New("agent is disabled")
	}
	a.mu.RUnlock()

	ch := make(chan string, 10)

	go func() {
		defer close(ch)

		response, err := a.Query(ctx, prompt)
		if err != nil {
			ch <- fmt.Sprintf("Error: %v", err)
			return
		}

		// Simulate streaming by sending word by word
		words := strings.Fields(response)
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			case ch <- word + " ":
			}
		}
	}()

	return ch, nil
}

// SuggestCommand suggests a command based on natural language.
func (a *Agent) SuggestCommand(ctx context.Context, description string) (string, error) {
	prompt := fmt.Sprintf("Suggest a single shell command to: %s. Only output the command, nothing else.", description)
	return a.Query(ctx, prompt)
}

// ExplainCommand explains what a command does.
func (a *Agent) ExplainCommand(ctx context.Context, command string) (string, error) {
	prompt := fmt.Sprintf("Explain what this shell command does: %s", command)
	return a.Query(ctx, prompt)
}

// FixError suggests a fix for an error.
func (a *Agent) FixError(ctx context.Context, command, errMsg string) (string, error) {
	prompt := fmt.Sprintf("The command '%s' failed with error: %s. Suggest a fix.", command, errMsg)
	return a.Query(ctx, prompt)
}

// Enable enables the agent.
func (a *Agent) Enable() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.enabled = true
}

// Disable disables the agent.
func (a *Agent) Disable() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.enabled = false
}

// IsEnabled returns whether the agent is enabled.
func (a *Agent) IsEnabled() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.enabled
}

// SetAPIKey sets the API key.
func (a *Agent) SetAPIKey(apiKey string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.apiKey = apiKey
}

// SetModel sets the model.
func (a *Agent) SetModel(model string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.model = model
}

func (a *Agent) mockOpenAIQuery(prompt string) (string, error) {
	// Mock response for demonstration
	if strings.Contains(strings.ToLower(prompt), "suggest") {
		return "echo 'This is a suggested command'", nil
	}
	if strings.Contains(strings.ToLower(prompt), "explain") {
		return "This command outputs text to the terminal.", nil
	}
	if strings.Contains(strings.ToLower(prompt), "fix") {
		return "Try checking your permissions or syntax.", nil
	}
	return "I can help with shell commands. Please ask me to suggest, explain, or fix commands.", nil
}

func (a *Agent) mockAnthropicQuery(prompt string) (string, error) {
	return a.mockOpenAIQuery(prompt) // Same mock behavior
}

func (a *Agent) mockGeminiQuery(prompt string) (string, error) {
	return a.mockOpenAIQuery(prompt) // Same mock behavior
}

func (a *Agent) mockLocalQuery(prompt string) (string, error) {
	return a.mockOpenAIQuery(prompt) // Same mock behavior
}

// Manager manages multiple AI agents.
type Manager struct {
	mu     sync.RWMutex
	agents map[string]*Agent
	active string
}

// NewManager creates a new AI manager.
func NewManager() *Manager {
	return &Manager{
		agents: make(map[string]*Agent),
	}
}

// RegisterAgent registers an AI agent.
func (m *Manager) RegisterAgent(agent *Agent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[agent.Name()]; exists {
		return fmt.Errorf("agent already registered: %s", agent.Name())
	}

	m.agents[agent.Name()] = agent

	// Set as active if first agent
	if m.active == "" {
		m.active = agent.Name()
	}

	return nil
}

// UnregisterAgent removes an AI agent.
func (m *Manager) UnregisterAgent(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[name]; !exists {
		return fmt.Errorf("agent not found: %s", name)
	}

	delete(m.agents, name)

	if m.active == name {
		m.active = ""
		for agentName := range m.agents {
			m.active = agentName
			break
		}
	}

	return nil
}

// GetAgent returns an agent by name.
func (m *Manager) GetAgent(name string) (*Agent, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agent, exists := m.agents[name]
	return agent, exists
}

// ActiveAgent returns the active agent.
func (m *Manager) ActiveAgent() *Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.agents[m.active]
}

// SetActiveAgent sets the active agent.
func (m *Manager) SetActiveAgent(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[name]; !exists {
		return fmt.Errorf("agent not found: %s", name)
	}

	m.active = name
	return nil
}

// ListAgents returns all registered agents.
func (m *Manager) ListAgents() []*Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Agent, 0, len(m.agents))
	for _, agent := range m.agents {
		result = append(result, agent)
	}
	return result
}

// Query sends a query to the active agent.
func (m *Manager) Query(ctx context.Context, prompt string) (string, error) {
	agent := m.ActiveAgent()
	if agent == nil {
		return "", errors.New("no active AI agent")
	}
	return agent.Query(ctx, prompt)
}

// SuggestCommand suggests a command using the active agent.
func (m *Manager) SuggestCommand(ctx context.Context, description string) (string, error) {
	agent := m.ActiveAgent()
	if agent == nil {
		return "", errors.New("no active AI agent")
	}
	return agent.SuggestCommand(ctx, description)
}

// ExplainCommand explains a command using the active agent.
func (m *Manager) ExplainCommand(ctx context.Context, command string) (string, error) {
	agent := m.ActiveAgent()
	if agent == nil {
		return "", errors.New("no active AI agent")
	}
	return agent.ExplainCommand(ctx, command)
}

// FixError suggests a fix using the active agent.
func (m *Manager) FixError(ctx context.Context, command, errMsg string) (string, error) {
	agent := m.ActiveAgent()
	if agent == nil {
		return "", errors.New("no active AI agent")
	}
	return agent.FixError(ctx, command, errMsg)
}

// Tool represents an AI-powered tool.
type Tool struct {
	Name        string
	Description string
	Handler     func(ctx context.Context, args map[string]string) (string, error)
}

// ToolRegistry manages AI tools.
type ToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]*Tool
}

// NewToolRegistry creates a new tool registry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*Tool),
	}
}

// Register registers a tool.
func (r *ToolRegistry) Register(tool *Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool already registered: %s", tool.Name)
	}

	r.tools[tool.Name] = tool
	return nil
}

// Unregister removes a tool.
func (r *ToolRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[name]; !exists {
		return fmt.Errorf("tool not found: %s", name)
	}

	delete(r.tools, name)
	return nil
}

// Get returns a tool by name.
func (r *ToolRegistry) Get(name string) (*Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, exists := r.tools[name]
	return tool, exists
}

// List returns all registered tools.
func (r *ToolRegistry) List() []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool)
	}
	return result
}

// Execute executes a tool by name.
func (r *ToolRegistry) Execute(ctx context.Context, name string, args map[string]string) (string, error) {
	r.mu.RLock()
	tool, exists := r.tools[name]
	r.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	return tool.Handler(ctx, args)
}
