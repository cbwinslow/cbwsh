package ollama

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		baseURL       string
		model         string
		expectedURL   string
		expectedModel string
	}{
		{
			name:          "default values",
			baseURL:       "",
			model:         "",
			expectedURL:   "http://localhost:11434",
			expectedModel: "llama2",
		},
		{
			name:          "custom values",
			baseURL:       "http://custom:8080",
			model:         "phi3",
			expectedURL:   "http://custom:8080",
			expectedModel: "phi3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, tt.model)
			
			if client.baseURL != tt.expectedURL {
				t.Errorf("expected baseURL %q, got %q", tt.expectedURL, client.baseURL)
			}
			
			if client.model != tt.expectedModel {
				t.Errorf("expected model %q, got %q", tt.expectedModel, client.model)
			}
			
			if client.httpClient == nil {
				t.Error("httpClient should not be nil")
			}
			
			if client.httpClient.Timeout != 60*time.Second {
				t.Errorf("expected timeout 60s, got %v", client.httpClient.Timeout)
			}
		})
	}
}

func TestSetModel(t *testing.T) {
	client := NewClient("", "")
	
	if client.GetModel() != "llama2" {
		t.Errorf("expected initial model 'llama2', got %q", client.GetModel())
	}
	
	client.SetModel("phi3")
	
	if client.GetModel() != "phi3" {
		t.Errorf("expected model 'phi3', got %q", client.GetModel())
	}
}

func TestGenerate_WithInvalidURL(t *testing.T) {
	client := NewClient("http://invalid-host:99999", "llama2")
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	_, err := client.Generate(ctx, "test prompt")
	if err == nil {
		t.Error("expected error with invalid URL, got nil")
	}
}

func TestChat_WithInvalidURL(t *testing.T) {
	client := NewClient("http://invalid-host:99999", "llama2")
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	messages := []ChatMessage{
		{Role: "user", Content: "test"},
	}
	
	_, err := client.Chat(ctx, messages)
	if err == nil {
		t.Error("expected error with invalid URL, got nil")
	}
}

func TestListModels_WithInvalidURL(t *testing.T) {
	client := NewClient("http://invalid-host:99999", "llama2")
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	_, err := client.ListModels(ctx)
	if err == nil {
		t.Error("expected error with invalid URL, got nil")
	}
}

func TestPing_WithInvalidURL(t *testing.T) {
	client := NewClient("http://invalid-host:99999", "llama2")
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	err := client.Ping(ctx)
	if err == nil {
		t.Error("expected error with invalid URL, got nil")
	}
}

func TestChatMessage(t *testing.T) {
	msg := ChatMessage{
		Role:    "user",
		Content: "Hello",
	}
	
	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}
	
	if msg.Content != "Hello" {
		t.Errorf("expected content 'Hello', got %q", msg.Content)
	}
}
