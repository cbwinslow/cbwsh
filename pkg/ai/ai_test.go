package ai_test

import (
	"context"
	"testing"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/ai"
	"github.com/cbwinslow/cbwsh/pkg/core"
)

func TestA2ARouter(t *testing.T) {
	t.Parallel()

	router := ai.NewA2ARouter()
	if router == nil {
		t.Fatal("expected non-nil router")
	}

	// Create a test handler
	testAgent := ai.NewBaseA2AAgent(ai.A2AAgentInfo{
		ID:          "test-agent",
		Name:        "Test Agent",
		Description: "A test agent",
		Version:     "1.0.0",
	})

	testAgent.SetHandler(func(_ context.Context, msg *ai.A2AMessage) (*ai.A2AMessage, error) {
		return &ai.A2AMessage{
			ID:        "response-id",
			Type:      ai.A2AMessageTypeResponse,
			From:      "test-agent",
			To:        msg.From,
			Payload:   "Hello from test agent!",
			Timestamp: time.Now(),
			ReplyTo:   msg.ID,
		}, nil
	})

	// Register handler
	err := router.RegisterHandler("test-agent", testAgent)
	if err != nil {
		t.Fatalf("failed to register handler: %v", err)
	}

	// Try to register again - should fail
	err = router.RegisterHandler("test-agent", testAgent)
	if err == nil {
		t.Error("expected error when registering duplicate handler")
	}

	// List agents
	agents := router.ListAgents()
	if len(agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(agents))
	}

	// Send a message
	msg := &ai.A2AMessage{
		Type:    ai.A2AMessageTypeQuery,
		From:    "sender",
		To:      "test-agent",
		Payload: "Hello!",
	}

	response, err := router.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	if response.Payload != "Hello from test agent!" {
		t.Errorf("unexpected response: %s", response.Payload)
	}

	// Unregister handler
	err = router.UnregisterHandler("test-agent")
	if err != nil {
		t.Fatalf("failed to unregister handler: %v", err)
	}

	// Try to unregister again - should fail
	err = router.UnregisterHandler("test-agent")
	if err == nil {
		t.Error("expected error when unregistering non-existent handler")
	}
}

func TestA2AMessageTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		msgType  ai.A2AMessageType
		expected string
	}{
		{ai.A2AMessageTypeQuery, "query"},
		{ai.A2AMessageTypeResponse, "response"},
		{ai.A2AMessageTypeEvent, "event"},
		{ai.A2AMessageTypeError, "error"},
	}

	for _, tt := range tests {
		if string(tt.msgType) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.msgType)
		}
	}
}

func TestA2AMessageEncoder(t *testing.T) {
	t.Parallel()

	encoder := &ai.A2AMessageEncoder{}

	msg := &ai.A2AMessage{
		ID:        "test-id",
		Type:      ai.A2AMessageTypeQuery,
		From:      "sender",
		To:        "receiver",
		Payload:   "Hello!",
		Timestamp: time.Now(),
	}

	// Encode
	data, err := encoder.Encode(msg)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	// Decode
	decoded, err := encoder.Decode(data)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if decoded.ID != msg.ID {
		t.Errorf("expected ID %s, got %s", msg.ID, decoded.ID)
	}
	if decoded.Payload != msg.Payload {
		t.Errorf("expected payload %s, got %s", msg.Payload, decoded.Payload)
	}
}

func TestShellAssistant(t *testing.T) {
	t.Parallel()

	// Create an AI agent
	agent := ai.NewAgent("test", core.AIProviderLocal, "", "")

	// Create shell assistant
	sa := ai.NewShellAssistant(agent)
	if sa == nil {
		t.Fatal("expected non-nil shell assistant")
	}

	// Get info
	info := sa.GetInfo()
	if info.ID != "shell-assistant" {
		t.Errorf("expected ID 'shell-assistant', got %s", info.ID)
	}

	if len(info.Capabilities) == 0 {
		t.Error("expected at least one capability")
	}

	// Handle a message
	msg := &ai.A2AMessage{
		ID:      "test-msg",
		Type:    ai.A2AMessageTypeQuery,
		From:    "user",
		To:      "shell-assistant",
		Payload: "How do I list files?",
	}

	response, err := sa.HandleMessage(context.Background(), msg)
	if err != nil {
		t.Fatalf("failed to handle message: %v", err)
	}

	if response.Type != ai.A2AMessageTypeResponse {
		t.Errorf("expected response type, got %s", response.Type)
	}
}

func TestAgentWithGeminiProvider(t *testing.T) {
	t.Parallel()

	agent := ai.NewAgent("gemini-test", core.AIProviderGemini, "fake-key", "gemini-pro")

	if agent.Provider() != core.AIProviderGemini {
		t.Errorf("expected Gemini provider, got %s", agent.Provider())
	}

	// Query should work (mock response)
	response, err := agent.Query(context.Background(), "Hello")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}

	if response == "" {
		t.Error("expected non-empty response")
	}
}

func TestAgentManager(t *testing.T) {
	t.Parallel()

	manager := ai.NewManager()

	// Create agents
	agent1 := ai.NewAgent("agent1", core.AIProviderOpenAI, "", "")
	agent2 := ai.NewAgent("agent2", core.AIProviderGemini, "", "")

	// Register agents
	if err := manager.RegisterAgent(agent1); err != nil {
		t.Fatalf("failed to register agent1: %v", err)
	}
	if err := manager.RegisterAgent(agent2); err != nil {
		t.Fatalf("failed to register agent2: %v", err)
	}

	// List agents
	agents := manager.ListAgents()
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}

	// Get agent
	agent, exists := manager.GetAgent("agent1")
	if !exists {
		t.Error("expected to find agent1")
	}
	if agent.Name() != "agent1" {
		t.Errorf("expected name 'agent1', got %s", agent.Name())
	}

	// Set active agent
	if err := manager.SetActiveAgent("agent2"); err != nil {
		t.Fatalf("failed to set active agent: %v", err)
	}

	active := manager.ActiveAgent()
	if active.Name() != "agent2" {
		t.Errorf("expected active agent 'agent2', got %s", active.Name())
	}

	// Unregister agent
	if err := manager.UnregisterAgent("agent1"); err != nil {
		t.Fatalf("failed to unregister agent: %v", err)
	}

	agents = manager.ListAgents()
	if len(agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(agents))
	}
}

func TestToolRegistry(t *testing.T) {
	t.Parallel()

	registry := ai.NewToolRegistry()

	// Register a tool
	tool := &ai.Tool{
		Name:        "test-tool",
		Description: "A test tool",
		Handler: func(_ context.Context, args map[string]string) (string, error) {
			return "Tool executed with: " + args["input"], nil
		},
	}

	if err := registry.Register(tool); err != nil {
		t.Fatalf("failed to register tool: %v", err)
	}

	// Try to register again - should fail
	if err := registry.Register(tool); err == nil {
		t.Error("expected error when registering duplicate tool")
	}

	// Get tool
	found, exists := registry.Get("test-tool")
	if !exists {
		t.Error("expected to find tool")
	}
	if found.Name != "test-tool" {
		t.Errorf("expected name 'test-tool', got %s", found.Name)
	}

	// List tools
	tools := registry.List()
	if len(tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(tools))
	}

	// Execute tool
	result, err := registry.Execute(context.Background(), "test-tool", map[string]string{"input": "hello"})
	if err != nil {
		t.Fatalf("failed to execute tool: %v", err)
	}
	if result != "Tool executed with: hello" {
		t.Errorf("unexpected result: %s", result)
	}

	// Unregister tool
	if err := registry.Unregister("test-tool"); err != nil {
		t.Fatalf("failed to unregister tool: %v", err)
	}

	// Try to unregister again - should fail
	if err := registry.Unregister("test-tool"); err == nil {
		t.Error("expected error when unregistering non-existent tool")
	}
}
