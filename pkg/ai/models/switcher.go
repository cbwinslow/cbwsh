// Package models provides runtime model switching for cbwsh AI.
package models

import (
	"errors"
	"sync"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// Common errors.
var (
	ErrModelNotFound  = errors.New("model not found")
	ErrProviderNotSet = errors.New("provider not set")
	ErrInvalidModel   = errors.New("invalid model configuration")
)

// ModelInfo contains information about an AI model.
type ModelInfo struct {
	// ID is the unique identifier for the model.
	ID string
	// Provider is the AI provider.
	Provider core.AIProvider
	// Name is the display name.
	Name string
	// Description describes the model.
	Description string
	// MaxTokens is the maximum context length.
	MaxTokens int
	// MaxOutputTokens is the maximum output tokens.
	MaxOutputTokens int
	// Available indicates if the model is currently available.
	Available bool
	// RequiresAPIKey indicates if an API key is needed.
	RequiresAPIKey bool
	// Capabilities lists what the model can do.
	Capabilities []string
	// CostPerToken is the cost per token (for budgeting).
	CostPerToken float64
}

// Switcher manages AI model switching at runtime.
type Switcher struct {
	mu           sync.RWMutex
	models       map[string]ModelInfo
	currentModel string
	apiKeys      map[core.AIProvider]string
	onSwitch     []func(from, to ModelInfo)
}

// NewSwitcher creates a new model switcher.
func NewSwitcher() *Switcher {
	s := &Switcher{
		models:   make(map[string]ModelInfo),
		apiKeys:  make(map[core.AIProvider]string),
		onSwitch: make([]func(from, to ModelInfo), 0),
	}

	// Register default models
	for _, model := range DefaultModels() {
		s.models[model.ID] = model
	}

	return s
}

// RegisterModel registers a new model.
func (s *Switcher) RegisterModel(model ModelInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.models[model.ID] = model
}

// UnregisterModel removes a model.
func (s *Switcher) UnregisterModel(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.models, id)
}

// ListModels returns all registered models.
func (s *Switcher) ListModels() []ModelInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]ModelInfo, 0, len(s.models))
	for _, model := range s.models {
		result = append(result, model)
	}
	return result
}

// ListAvailableModels returns all available models.
func (s *Switcher) ListAvailableModels() []ModelInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []ModelInfo
	for _, model := range s.models {
		if s.isModelAvailable(model) {
			model.Available = true
			result = append(result, model)
		}
	}
	return result
}

// ListModelsByProvider returns models for a specific provider.
func (s *Switcher) ListModelsByProvider(provider core.AIProvider) []ModelInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []ModelInfo
	for _, model := range s.models {
		if model.Provider == provider {
			model.Available = s.isModelAvailable(model)
			result = append(result, model)
		}
	}
	return result
}

// GetModel returns a model by ID.
func (s *Switcher) GetModel(id string) (ModelInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	model, ok := s.models[id]
	if ok {
		model.Available = s.isModelAvailable(model)
	}
	return model, ok
}

// CurrentModel returns the current model.
func (s *Switcher) CurrentModel() (ModelInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.currentModel == "" {
		return ModelInfo{}, ErrModelNotFound
	}

	model, ok := s.models[s.currentModel]
	if !ok {
		return ModelInfo{}, ErrModelNotFound
	}

	model.Available = s.isModelAvailable(model)
	return model, nil
}

// CurrentModelID returns the current model ID.
func (s *Switcher) CurrentModelID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentModel
}

// Switch switches to a different model.
func (s *Switcher) Switch(modelID string) error {
	s.mu.Lock()

	model, ok := s.models[modelID]
	if !ok {
		s.mu.Unlock()
		return ErrModelNotFound
	}

	if !s.isModelAvailable(model) {
		s.mu.Unlock()
		return ErrProviderNotSet
	}

	var oldModel ModelInfo
	if s.currentModel != "" {
		oldModel = s.models[s.currentModel]
	}

	s.currentModel = modelID
	callbacks := s.onSwitch

	s.mu.Unlock()

	// Call callbacks outside the lock
	for _, cb := range callbacks {
		cb(oldModel, model)
	}

	return nil
}

// SwitchProvider switches to the default model for a provider.
func (s *Switcher) SwitchProvider(provider core.AIProvider) error {
	s.mu.RLock()

	// Find the first available model for this provider
	var targetModel string
	for id, model := range s.models {
		if model.Provider == provider && s.isModelAvailable(model) {
			targetModel = id
			break
		}
	}

	s.mu.RUnlock()

	if targetModel == "" {
		return ErrModelNotFound
	}

	return s.Switch(targetModel)
}

// SetAPIKey sets the API key for a provider.
func (s *Switcher) SetAPIKey(provider core.AIProvider, apiKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apiKeys[provider] = apiKey
}

// GetAPIKey gets the API key for a provider.
func (s *Switcher) GetAPIKey(provider core.AIProvider) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apiKeys[provider]
}

// HasAPIKey returns whether an API key is set for a provider.
func (s *Switcher) HasAPIKey(provider core.AIProvider) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apiKeys[provider] != ""
}

// OnSwitch registers a callback for model switches.
func (s *Switcher) OnSwitch(callback func(from, to ModelInfo)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onSwitch = append(s.onSwitch, callback)
}

// ValidateModel validates a model configuration.
func (s *Switcher) ValidateModel(model ModelInfo) error {
	if model.ID == "" {
		return ErrInvalidModel
	}
	if model.Provider == core.AIProviderNone {
		return ErrProviderNotSet
	}
	return nil
}

func (s *Switcher) isModelAvailable(model ModelInfo) bool {
	// Local models don't need API keys
	if model.Provider == core.AIProviderLocal {
		return true
	}

	// Check if API key is set for the provider
	if model.RequiresAPIKey {
		return s.apiKeys[model.Provider] != ""
	}

	return true
}

// DefaultModels returns the default model configurations.
func DefaultModels() []ModelInfo {
	return []ModelInfo{
		// OpenAI Models
		{
			ID:              "gpt-4o",
			Provider:        core.AIProviderOpenAI,
			Name:            "GPT-4o",
			Description:     "Most capable OpenAI model",
			MaxTokens:       128000,
			MaxOutputTokens: 4096,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code", "analysis", "vision"},
		},
		{
			ID:              "gpt-4o-mini",
			Provider:        core.AIProviderOpenAI,
			Name:            "GPT-4o Mini",
			Description:     "Fast, affordable model for simple tasks",
			MaxTokens:       128000,
			MaxOutputTokens: 4096,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code"},
		},
		{
			ID:              "gpt-4-turbo",
			Provider:        core.AIProviderOpenAI,
			Name:            "GPT-4 Turbo",
			Description:     "Fast GPT-4 with vision capabilities",
			MaxTokens:       128000,
			MaxOutputTokens: 4096,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code", "analysis", "vision"},
		},
		{
			ID:              "gpt-3.5-turbo",
			Provider:        core.AIProviderOpenAI,
			Name:            "GPT-3.5 Turbo",
			Description:     "Fast model for simple tasks",
			MaxTokens:       16385,
			MaxOutputTokens: 4096,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code"},
		},

		// Anthropic Models
		{
			ID:              "claude-3-5-sonnet",
			Provider:        core.AIProviderAnthropic,
			Name:            "Claude 3.5 Sonnet",
			Description:     "Most intelligent Claude model",
			MaxTokens:       200000,
			MaxOutputTokens: 8192,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code", "analysis", "vision"},
		},
		{
			ID:              "claude-3-opus",
			Provider:        core.AIProviderAnthropic,
			Name:            "Claude 3 Opus",
			Description:     "Powerful model for complex tasks",
			MaxTokens:       200000,
			MaxOutputTokens: 4096,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code", "analysis", "vision"},
		},
		{
			ID:              "claude-3-sonnet",
			Provider:        core.AIProviderAnthropic,
			Name:            "Claude 3 Sonnet",
			Description:     "Balanced performance and cost",
			MaxTokens:       200000,
			MaxOutputTokens: 4096,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code", "analysis"},
		},
		{
			ID:              "claude-3-haiku",
			Provider:        core.AIProviderAnthropic,
			Name:            "Claude 3 Haiku",
			Description:     "Fast and affordable",
			MaxTokens:       200000,
			MaxOutputTokens: 4096,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code"},
		},

		// Google Models
		{
			ID:              "gemini-pro",
			Provider:        core.AIProviderGemini,
			Name:            "Gemini Pro",
			Description:     "Google's advanced AI model",
			MaxTokens:       30720,
			MaxOutputTokens: 2048,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code", "analysis"},
		},
		{
			ID:              "gemini-1.5-pro",
			Provider:        core.AIProviderGemini,
			Name:            "Gemini 1.5 Pro",
			Description:     "Long context Google model",
			MaxTokens:       1000000,
			MaxOutputTokens: 8192,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code", "analysis", "vision"},
		},
		{
			ID:              "gemini-1.5-flash",
			Provider:        core.AIProviderGemini,
			Name:            "Gemini 1.5 Flash",
			Description:     "Fast Google model",
			MaxTokens:       1000000,
			MaxOutputTokens: 8192,
			RequiresAPIKey:  true,
			Capabilities:    []string{"chat", "code"},
		},

		// Local Models
		{
			ID:              "local-llama",
			Provider:        core.AIProviderLocal,
			Name:            "Local LLaMA",
			Description:     "Run LLaMA locally",
			MaxTokens:       4096,
			MaxOutputTokens: 2048,
			RequiresAPIKey:  false,
			Capabilities:    []string{"chat", "code"},
		},
		{
			ID:              "local-mistral",
			Provider:        core.AIProviderLocal,
			Name:            "Local Mistral",
			Description:     "Run Mistral locally",
			MaxTokens:       8192,
			MaxOutputTokens: 2048,
			RequiresAPIKey:  false,
			Capabilities:    []string{"chat", "code"},
		},
	}
}

// ProviderName returns a human-readable name for a provider.
func ProviderName(provider core.AIProvider) string {
	switch provider {
	case core.AIProviderOpenAI:
		return "OpenAI"
	case core.AIProviderAnthropic:
		return "Anthropic"
	case core.AIProviderGemini:
		return "Google Gemini"
	case core.AIProviderLocal:
		return "Local"
	default:
		return "Unknown"
	}
}

// ProviderFromString converts a string to AIProvider.
func ProviderFromString(s string) core.AIProvider {
	switch s {
	case "openai":
		return core.AIProviderOpenAI
	case "anthropic":
		return core.AIProviderAnthropic
	case "gemini":
		return core.AIProviderGemini
	case "local":
		return core.AIProviderLocal
	default:
		return core.AIProviderNone
	}
}
