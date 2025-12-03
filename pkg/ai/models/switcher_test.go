package models

import (
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

func TestSwitcher_New(t *testing.T) {
	s := NewSwitcher()

	if s == nil {
		t.Fatal("NewSwitcher() returned nil")
	}

	// Should have default models
	models := s.ListModels()
	if len(models) == 0 {
		t.Error("NewSwitcher() should have default models")
	}
}

func TestSwitcher_RegisterModel(t *testing.T) {
	s := NewSwitcher()

	model := ModelInfo{
		ID:       "custom-model",
		Provider: core.AIProviderLocal,
		Name:     "Custom Model",
	}

	s.RegisterModel(model)

	got, ok := s.GetModel("custom-model")
	if !ok {
		t.Fatal("GetModel() should find registered model")
	}

	if got.Name != "Custom Model" {
		t.Errorf("GetModel().Name = %s, want Custom Model", got.Name)
	}
}

func TestSwitcher_UnregisterModel(t *testing.T) {
	s := NewSwitcher()

	s.RegisterModel(ModelInfo{
		ID:       "to-remove",
		Provider: core.AIProviderLocal,
		Name:     "To Remove",
	})

	s.UnregisterModel("to-remove")

	_, ok := s.GetModel("to-remove")
	if ok {
		t.Error("GetModel() should not find unregistered model")
	}
}

func TestSwitcher_ListModels(t *testing.T) {
	s := NewSwitcher()

	models := s.ListModels()
	if len(models) == 0 {
		t.Error("ListModels() should return models")
	}
}

func TestSwitcher_ListModelsByProvider(t *testing.T) {
	s := NewSwitcher()

	// List OpenAI models
	models := s.ListModelsByProvider(core.AIProviderOpenAI)
	for _, m := range models {
		if m.Provider != core.AIProviderOpenAI {
			t.Errorf("ListModelsByProvider() returned model with wrong provider: %v", m.Provider)
		}
	}
}

func TestSwitcher_GetModel(t *testing.T) {
	s := NewSwitcher()

	// Existing model
	_, ok := s.GetModel("gpt-4o")
	if !ok {
		t.Error("GetModel() should find gpt-4o")
	}

	// Non-existing model
	_, ok = s.GetModel("nonexistent")
	if ok {
		t.Error("GetModel() should not find nonexistent model")
	}
}

func TestSwitcher_Switch(t *testing.T) {
	s := NewSwitcher()

	// Set API key for OpenAI
	s.SetAPIKey(core.AIProviderOpenAI, "test-key")

	// Switch to gpt-4o
	err := s.Switch("gpt-4o")
	if err != nil {
		t.Fatalf("Switch() error = %v", err)
	}

	current, err := s.CurrentModel()
	if err != nil {
		t.Fatalf("CurrentModel() error = %v", err)
	}

	if current.ID != "gpt-4o" {
		t.Errorf("CurrentModel().ID = %s, want gpt-4o", current.ID)
	}
}

func TestSwitcher_Switch_NotFound(t *testing.T) {
	s := NewSwitcher()

	err := s.Switch("nonexistent")
	if err != ErrModelNotFound {
		t.Errorf("Switch() error = %v, want ErrModelNotFound", err)
	}
}

func TestSwitcher_Switch_NoAPIKey(t *testing.T) {
	s := NewSwitcher()

	// Don't set API key
	err := s.Switch("gpt-4o")
	if err != ErrProviderNotSet {
		t.Errorf("Switch() error = %v, want ErrProviderNotSet", err)
	}
}

func TestSwitcher_Switch_LocalModel(t *testing.T) {
	s := NewSwitcher()

	// Local models don't need API keys
	err := s.Switch("local-llama")
	if err != nil {
		t.Fatalf("Switch() error = %v", err)
	}

	current, _ := s.CurrentModel()
	if current.ID != "local-llama" {
		t.Errorf("CurrentModel().ID = %s, want local-llama", current.ID)
	}
}

func TestSwitcher_SwitchProvider(t *testing.T) {
	s := NewSwitcher()

	// Set API key
	s.SetAPIKey(core.AIProviderAnthropic, "test-key")

	// Switch to Anthropic
	err := s.SwitchProvider(core.AIProviderAnthropic)
	if err != nil {
		t.Fatalf("SwitchProvider() error = %v", err)
	}

	current, _ := s.CurrentModel()
	if current.Provider != core.AIProviderAnthropic {
		t.Errorf("CurrentModel().Provider = %v, want Anthropic", current.Provider)
	}
}

func TestSwitcher_CurrentModel_NoModel(t *testing.T) {
	s := NewSwitcher()

	_, err := s.CurrentModel()
	if err != ErrModelNotFound {
		t.Errorf("CurrentModel() error = %v, want ErrModelNotFound", err)
	}
}

func TestSwitcher_CurrentModelID(t *testing.T) {
	s := NewSwitcher()
	s.SetAPIKey(core.AIProviderOpenAI, "test-key")
	_ = s.Switch("gpt-4o")

	id := s.CurrentModelID()
	if id != "gpt-4o" {
		t.Errorf("CurrentModelID() = %s, want gpt-4o", id)
	}
}

func TestSwitcher_APIKey(t *testing.T) {
	s := NewSwitcher()

	// Not set
	if s.HasAPIKey(core.AIProviderOpenAI) {
		t.Error("HasAPIKey() should return false initially")
	}

	// Set key
	s.SetAPIKey(core.AIProviderOpenAI, "sk-test")

	if !s.HasAPIKey(core.AIProviderOpenAI) {
		t.Error("HasAPIKey() should return true after SetAPIKey")
	}

	if s.GetAPIKey(core.AIProviderOpenAI) != "sk-test" {
		t.Error("GetAPIKey() should return the set key")
	}
}

func TestSwitcher_OnSwitch(t *testing.T) {
	s := NewSwitcher()
	s.SetAPIKey(core.AIProviderOpenAI, "test-key")

	called := false
	var fromModel, toModel ModelInfo

	s.OnSwitch(func(from, to ModelInfo) {
		called = true
		fromModel = from
		toModel = to
	})

	_ = s.Switch("gpt-4o")

	if !called {
		t.Error("OnSwitch callback should be called")
	}

	if toModel.ID != "gpt-4o" {
		t.Errorf("OnSwitch toModel.ID = %s, want gpt-4o", toModel.ID)
	}

	// Switch again
	_ = s.Switch("gpt-3.5-turbo")

	if fromModel.ID != "gpt-4o" {
		t.Errorf("OnSwitch fromModel.ID = %s, want gpt-4o", fromModel.ID)
	}
}

func TestSwitcher_ListAvailableModels(t *testing.T) {
	s := NewSwitcher()

	// No API keys set - only local models should be available
	available := s.ListAvailableModels()

	for _, m := range available {
		if m.Provider != core.AIProviderLocal && m.RequiresAPIKey {
			t.Errorf("Model %s should not be available without API key", m.ID)
		}
	}

	// Set OpenAI key
	s.SetAPIKey(core.AIProviderOpenAI, "test-key")
	available = s.ListAvailableModels()

	// Should have OpenAI models now
	hasOpenAI := false
	for _, m := range available {
		if m.Provider == core.AIProviderOpenAI {
			hasOpenAI = true
			break
		}
	}

	if !hasOpenAI {
		t.Error("ListAvailableModels() should include OpenAI models after setting API key")
	}
}

func TestSwitcher_ValidateModel(t *testing.T) {
	s := NewSwitcher()

	tests := []struct {
		name    string
		model   ModelInfo
		wantErr error
	}{
		{
			name: "valid model",
			model: ModelInfo{
				ID:       "test",
				Provider: core.AIProviderLocal,
			},
			wantErr: nil,
		},
		{
			name: "missing ID",
			model: ModelInfo{
				Provider: core.AIProviderLocal,
			},
			wantErr: ErrInvalidModel,
		},
		{
			name: "missing provider",
			model: ModelInfo{
				ID:       "test",
				Provider: core.AIProviderNone,
			},
			wantErr: ErrProviderNotSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.ValidateModel(tt.model)
			if err != tt.wantErr {
				t.Errorf("ValidateModel() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultModels(t *testing.T) {
	models := DefaultModels()

	if len(models) == 0 {
		t.Fatal("DefaultModels() returned empty list")
	}

	// Check for expected providers
	providers := make(map[core.AIProvider]bool)
	for _, m := range models {
		providers[m.Provider] = true
	}

	expectedProviders := []core.AIProvider{
		core.AIProviderOpenAI,
		core.AIProviderAnthropic,
		core.AIProviderGemini,
		core.AIProviderLocal,
	}

	for _, p := range expectedProviders {
		if !providers[p] {
			t.Errorf("DefaultModels() missing provider: %v", p)
		}
	}
}

func TestProviderName(t *testing.T) {
	tests := []struct {
		provider core.AIProvider
		want     string
	}{
		{core.AIProviderOpenAI, "OpenAI"},
		{core.AIProviderAnthropic, "Anthropic"},
		{core.AIProviderGemini, "Google Gemini"},
		{core.AIProviderLocal, "Local"},
		{core.AIProviderNone, "Unknown"},
	}

	for _, tt := range tests {
		got := ProviderName(tt.provider)
		if got != tt.want {
			t.Errorf("ProviderName(%v) = %s, want %s", tt.provider, got, tt.want)
		}
	}
}

func TestProviderFromString(t *testing.T) {
	tests := []struct {
		input string
		want  core.AIProvider
	}{
		{"openai", core.AIProviderOpenAI},
		{"anthropic", core.AIProviderAnthropic},
		{"gemini", core.AIProviderGemini},
		{"local", core.AIProviderLocal},
		{"unknown", core.AIProviderNone},
	}

	for _, tt := range tests {
		got := ProviderFromString(tt.input)
		if got != tt.want {
			t.Errorf("ProviderFromString(%s) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
