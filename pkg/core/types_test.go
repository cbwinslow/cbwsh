package core_test

import (
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

func TestShellTypeString(t *testing.T) {
	tests := []struct {
		shellType core.ShellType
		expected  string
	}{
		{core.ShellTypeBash, "bash"},
		{core.ShellTypeZsh, "zsh"},
		{core.ShellType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.shellType.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPaneLayoutString(t *testing.T) {
	tests := []struct {
		layout   core.PaneLayout
		expected string
	}{
		{core.LayoutSingle, "single"},
		{core.LayoutHorizontalSplit, "horizontal"},
		{core.LayoutVerticalSplit, "vertical"},
		{core.LayoutGrid, "grid"},
		{core.LayoutCustom, "custom"},
		{core.PaneLayout(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.layout.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPluginTypeString(t *testing.T) {
	tests := []struct {
		pluginType core.PluginType
		expected   string
	}{
		{core.PluginTypeCommand, "command"},
		{core.PluginTypeUI, "ui"},
		{core.PluginTypeHook, "hook"},
		{core.PluginTypeFormatter, "formatter"},
		{core.PluginType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.pluginType.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSSHConnectionStateString(t *testing.T) {
	tests := []struct {
		state    core.SSHConnectionState
		expected string
	}{
		{core.SSHDisconnected, "disconnected"},
		{core.SSHConnecting, "connecting"},
		{core.SSHConnected, "connected"},
		{core.SSHError, "error"},
		{core.SSHConnectionState(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestAIProviderString(t *testing.T) {
	tests := []struct {
		provider core.AIProvider
		expected string
	}{
		{core.AIProviderNone, "none"},
		{core.AIProviderOpenAI, "openai"},
		{core.AIProviderAnthropic, "anthropic"},
		{core.AIProviderGemini, "gemini"},
		{core.AIProviderLocal, "local"},
		{core.AIProvider(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.provider.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
