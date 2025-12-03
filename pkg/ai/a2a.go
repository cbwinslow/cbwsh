// Package ai provides AI integration for cbwsh.
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// A2AMessageType represents the type of A2A protocol message.
type A2AMessageType string

const (
	// A2AMessageTypeQuery represents a query from one agent to another.
	A2AMessageTypeQuery A2AMessageType = "query"
	// A2AMessageTypeResponse represents a response to a query.
	A2AMessageTypeResponse A2AMessageType = "response"
	// A2AMessageTypeEvent represents an event notification.
	A2AMessageTypeEvent A2AMessageType = "event"
	// A2AMessageTypeError represents an error message.
	A2AMessageTypeError A2AMessageType = "error"
)

// A2AMessage represents a message in the A2A protocol.
type A2AMessage struct {
	ID        string                 `json:"id"`
	Type      A2AMessageType         `json:"type"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Payload   string                 `json:"payload"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	ReplyTo   string                 `json:"reply_to,omitempty"`
}

// A2ACapability represents a capability that an agent can offer.
type A2ACapability struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	InputTypes  []string `json:"input_types"`
	OutputType  string   `json:"output_type"`
}

// A2AAgentInfo represents information about an A2A agent.
type A2AAgentInfo struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Capabilities []A2ACapability `json:"capabilities"`
	Version      string          `json:"version"`
}

// A2AHandler handles incoming A2A messages.
type A2AHandler interface {
	// HandleMessage handles an incoming A2A message.
	HandleMessage(ctx context.Context, msg *A2AMessage) (*A2AMessage, error)
	// GetInfo returns information about this agent.
	GetInfo() A2AAgentInfo
}

// A2ARouter routes A2A messages between agents.
type A2ARouter struct {
	mu       sync.RWMutex
	handlers map[string]A2AHandler
	messages chan *A2AMessage
	pending  map[string]chan *A2AMessage
}

// NewA2ARouter creates a new A2A router.
func NewA2ARouter() *A2ARouter {
	return &A2ARouter{
		handlers: make(map[string]A2AHandler),
		messages: make(chan *A2AMessage, 100),
		pending:  make(map[string]chan *A2AMessage),
	}
}

// RegisterHandler registers an A2A handler for an agent.
func (r *A2ARouter) RegisterHandler(agentID string, handler A2AHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[agentID]; exists {
		return fmt.Errorf("handler already registered for agent: %s", agentID)
	}

	r.handlers[agentID] = handler
	return nil
}

// UnregisterHandler removes an A2A handler.
func (r *A2ARouter) UnregisterHandler(agentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[agentID]; !exists {
		return fmt.Errorf("handler not found for agent: %s", agentID)
	}

	delete(r.handlers, agentID)
	return nil
}

// Send sends an A2A message to another agent.
func (r *A2ARouter) Send(ctx context.Context, msg *A2AMessage) (*A2AMessage, error) {
	r.mu.RLock()
	handler, exists := r.handlers[msg.To]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler for agent: %s", msg.To)
	}

	// Set message ID and timestamp if not set
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	return handler.HandleMessage(ctx, msg)
}

// SendAsync sends an A2A message asynchronously.
func (r *A2ARouter) SendAsync(msg *A2AMessage) (string, error) {
	r.mu.Lock()
	_, exists := r.handlers[msg.To]
	if !exists {
		r.mu.Unlock()
		return "", fmt.Errorf("no handler for agent: %s", msg.To)
	}

	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	// Create pending channel for response
	responseChan := make(chan *A2AMessage, 1)
	r.pending[msg.ID] = responseChan
	r.mu.Unlock()

	// Queue message for processing
	select {
	case r.messages <- msg:
		return msg.ID, nil
	default:
		r.mu.Lock()
		delete(r.pending, msg.ID)
		r.mu.Unlock()
		return "", fmt.Errorf("message queue full")
	}
}

// WaitForResponse waits for a response to a message.
func (r *A2ARouter) WaitForResponse(ctx context.Context, msgID string) (*A2AMessage, error) {
	r.mu.RLock()
	ch, exists := r.pending[msgID]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no pending message with ID: %s", msgID)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case response := <-ch:
		r.mu.Lock()
		delete(r.pending, msgID)
		r.mu.Unlock()
		return response, nil
	}
}

// ListAgents returns information about all registered agents.
func (r *A2ARouter) ListAgents() []A2AAgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]A2AAgentInfo, 0, len(r.handlers))
	for _, handler := range r.handlers {
		agents = append(agents, handler.GetInfo())
	}
	return agents
}

// Start starts the router's message processing loop.
func (r *A2ARouter) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-r.messages:
				r.mu.RLock()
				handler, exists := r.handlers[msg.To]
				r.mu.RUnlock()

				if !exists {
					continue
				}

				response, err := handler.HandleMessage(ctx, msg)
				if err != nil {
					response = &A2AMessage{
						ID:        uuid.New().String(),
						Type:      A2AMessageTypeError,
						From:      msg.To,
						To:        msg.From,
						Payload:   err.Error(),
						Timestamp: time.Now(),
						ReplyTo:   msg.ID,
					}
				}

				r.mu.RLock()
				ch, pending := r.pending[msg.ID]
				r.mu.RUnlock()

				if pending && response != nil {
					select {
					case ch <- response:
					default:
					}
				}
			}
		}
	}()
}

// BaseA2AAgent provides a base implementation for A2A agents.
type BaseA2AAgent struct {
	info    A2AAgentInfo
	handler func(ctx context.Context, msg *A2AMessage) (*A2AMessage, error)
}

// NewBaseA2AAgent creates a new base A2A agent.
func NewBaseA2AAgent(info A2AAgentInfo) *BaseA2AAgent {
	return &BaseA2AAgent{
		info: info,
	}
}

// SetHandler sets the message handler.
func (a *BaseA2AAgent) SetHandler(handler func(ctx context.Context, msg *A2AMessage) (*A2AMessage, error)) {
	a.handler = handler
}

// HandleMessage handles an incoming A2A message.
func (a *BaseA2AAgent) HandleMessage(ctx context.Context, msg *A2AMessage) (*A2AMessage, error) {
	if a.handler == nil {
		return &A2AMessage{
			ID:        uuid.New().String(),
			Type:      A2AMessageTypeResponse,
			From:      a.info.ID,
			To:        msg.From,
			Payload:   "No handler configured",
			Timestamp: time.Now(),
			ReplyTo:   msg.ID,
		}, nil
	}
	return a.handler(ctx, msg)
}

// GetInfo returns information about this agent.
func (a *BaseA2AAgent) GetInfo() A2AAgentInfo {
	return a.info
}

// A2AMessageEncoder encodes/decodes A2A messages.
type A2AMessageEncoder struct{}

// Encode encodes an A2A message to JSON.
func (e *A2AMessageEncoder) Encode(msg *A2AMessage) ([]byte, error) {
	return json.Marshal(msg)
}

// Decode decodes an A2A message from JSON.
func (e *A2AMessageEncoder) Decode(data []byte) (*A2AMessage, error) {
	var msg A2AMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ShellAssistant is an A2A agent for shell command assistance.
type ShellAssistant struct {
	*BaseA2AAgent
	aiAgent *Agent
}

// NewShellAssistant creates a new shell assistant A2A agent.
func NewShellAssistant(aiAgent *Agent) *ShellAssistant {
	sa := &ShellAssistant{
		BaseA2AAgent: NewBaseA2AAgent(A2AAgentInfo{
			ID:          "shell-assistant",
			Name:        "Shell Assistant",
			Description: "Assists with shell commands and scripting",
			Version:     "1.0.0",
			Capabilities: []A2ACapability{
				{
					Name:        "suggest_command",
					Description: "Suggests a shell command based on a natural language description",
					InputTypes:  []string{"text"},
					OutputType:  "command",
				},
				{
					Name:        "explain_command",
					Description: "Explains what a shell command does",
					InputTypes:  []string{"command"},
					OutputType:  "text",
				},
				{
					Name:        "fix_error",
					Description: "Suggests a fix for a command error",
					InputTypes:  []string{"command", "error"},
					OutputType:  "text",
				},
			},
		}),
		aiAgent: aiAgent,
	}

	sa.SetHandler(sa.handleMessage)
	return sa
}

func (sa *ShellAssistant) handleMessage(ctx context.Context, msg *A2AMessage) (*A2AMessage, error) {
	response, err := sa.processMessage(ctx, msg)
	if err != nil {
		return sa.createErrorResponse(msg, err), nil
	}

	return sa.createSuccessResponse(msg, response), nil
}

func (sa *ShellAssistant) processMessage(ctx context.Context, msg *A2AMessage) (string, error) {
	if sa.aiAgent == nil {
		return "AI agent not configured", nil
	}

	if msg.Type != A2AMessageTypeQuery {
		return "Unknown message type", nil
	}

	return sa.handleQueryCapability(ctx, msg)
}

func (sa *ShellAssistant) handleQueryCapability(ctx context.Context, msg *A2AMessage) (string, error) {
	capability, _ := msg.Metadata["capability"].(string)

	switch capability {
	case "suggest_command":
		return sa.aiAgent.SuggestCommand(ctx, msg.Payload)
	case "explain_command":
		return sa.aiAgent.ExplainCommand(ctx, msg.Payload)
	case "fix_error":
		command, _ := msg.Metadata["command"].(string)
		return sa.aiAgent.FixError(ctx, command, msg.Payload)
	default:
		return sa.aiAgent.Query(ctx, msg.Payload)
	}
}

func (sa *ShellAssistant) createErrorResponse(msg *A2AMessage, err error) *A2AMessage {
	return &A2AMessage{
		ID:        uuid.New().String(),
		Type:      A2AMessageTypeError,
		From:      sa.info.ID,
		To:        msg.From,
		Payload:   err.Error(),
		Timestamp: time.Now(),
		ReplyTo:   msg.ID,
	}
}

func (sa *ShellAssistant) createSuccessResponse(msg *A2AMessage, response string) *A2AMessage {
	return &A2AMessage{
		ID:        uuid.New().String(),
		Type:      A2AMessageTypeResponse,
		From:      sa.info.ID,
		To:        msg.From,
		Payload:   response,
		Timestamp: time.Now(),
		ReplyTo:   msg.ID,
	}
}
