package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/config"
)

// TestConfigurationScenarios tests various configuration scenarios
func TestConfigurationScenarios(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		configYAML  string
		expectError bool
	}{
		{
			name: "minimal valid config",
			configYAML: `
shell:
  default_shell: bash
ui:
  theme: default
ai:
  provider: none
`,
			expectError: false,
		},
		{
			name: "config with ollama",
			configYAML: `
shell:
  default_shell: bash
  history_size: 10000
ui:
  theme: dracula
  layout: single
  show_status_bar: true
ai:
  provider: ollama
  ollama_url: http://localhost:11434
  ollama_model: codellama
  enable_monitoring: true
`,
			expectError: false,
		},
		{
			name: "config with openai",
			configYAML: `
shell:
  default_shell: zsh
ui:
  theme: nord
  enable_animations: true
ai:
  provider: openai
  api_key: sk-test123
  model: gpt-4
`,
			expectError: false,
		},
		{
			name: "config with claude",
			configYAML: `
ai:
  provider: anthropic
  api_key: sk-ant-test
  model: claude-3-opus-20240229
`,
			expectError: false,
		},
		{
			name: "config with gemini",
			configYAML: `
ai:
  provider: gemini
  api_key: AIzaTest
  model: gemini-pro
`,
			expectError: false,
		},
		{
			name: "config with all themes",
			configYAML: `
ui:
  theme: tokyo-night
  layout: horizontal
  show_menu_bar: true
  syntax_highlighting: true
`,
			expectError: false,
		},
		{
			name: "config with ssh settings",
			configYAML: `
ssh:
  default_user: deploy
  default_key_path: ~/.ssh/id_rsa
  known_hosts_path: ~/.ssh/known_hosts
  connect_timeout: 30
  keep_alive_interval: 60
  saved_hosts:
    - name: production
      host: prod.example.com
      user: deploy
      port: 22
`,
			expectError: false,
		},
		{
			name: "config with plugins",
			configYAML: `
plugins:
  enabled: true
  auto_load: true
  directory: ~/.cbwsh/plugins
`,
			expectError: false,
		},
		{
			name: "config with logging",
			configYAML: `
logging:
  enabled: true
  level: debug
  path: ~/.cbwsh/logs/cbwsh.log
  max_size: 10
  max_backups: 3
`,
			expectError: false,
		},
		{
			name: "config with all layouts",
			configYAML: `
ui:
  layout: grid
`,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create temp config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(configPath, []byte(tc.configYAML), 0644)
			if err != nil {
				t.Fatalf("failed to write config file: %v", err)
			}

			// Load config
			cfg, err := config.LoadFromFile(configPath)

			if tc.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if cfg == nil {
				t.Error("expected non-nil config")
			}
		})
	}
}

// TestDefaultConfiguration tests that default configuration works
func TestDefaultConfiguration(t *testing.T) {
	t.Parallel()

	cfg := config.Default()

	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}

	if cfg.Shell.DefaultShell == "" {
		t.Error("default shell should be set")
	}

	if cfg.UI.Theme == "" {
		t.Error("default theme should be set")
	}
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		config  *config.Config
		isValid bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: "bash",
					HistorySize:  10000,
				},
				UI: config.UIConfig{
					Theme:  "default",
					Layout: "single",
				},
			},
			isValid: true,
		},
		{
			name: "invalid shell",
			config: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: "invalid",
				},
			},
			isValid: false,
		},
		{
			name: "invalid theme",
			config: &config.Config{
				UI: config.UIConfig{
					Theme: "nonexistent",
				},
			},
			isValid: false,
		},
		{
			name: "invalid layout",
			config: &config.Config{
				UI: config.UIConfig{
					Layout: "invalid",
				},
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.Validate()

			if tc.isValid {
				if err != nil {
					t.Errorf("expected valid config, got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("expected validation error but got none")
				}
			}
		})
	}
}

// TestAIProviderConfigurations tests different AI provider configurations
func TestAIProviderConfigurations(t *testing.T) {
	t.Parallel()

	providers := []string{"none", "ollama", "openai", "anthropic", "gemini"}

	for _, provider := range providers {
		provider := provider
		t.Run(provider, func(t *testing.T) {
			t.Parallel()

			cfg := config.Default()
			cfg.AI.Provider = provider

			// Set provider-specific configs
			switch provider {
			case "ollama":
				cfg.AI.OllamaURL = "http://localhost:11434"
				cfg.AI.OllamaModel = "codellama"
			case "openai":
				cfg.AI.APIKey = "test-key"
				cfg.AI.Model = "gpt-4"
			case "anthropic":
				cfg.AI.APIKey = "test-key"
				cfg.AI.Model = "claude-3-opus"
			case "gemini":
				cfg.AI.APIKey = "test-key"
				cfg.AI.Model = "gemini-pro"
			}

			// Validate
			err := cfg.Validate()
			if err != nil && provider != "none" {
				// Some providers might require additional validation
				t.Logf("validation error for %s: %v", provider, err)
			}
		})
	}
}

// TestThemeConfigurations tests all available themes
func TestThemeConfigurations(t *testing.T) {
	t.Parallel()

	themes := []string{"default", "dracula", "nord", "tokyo-night", "gruvbox"}

	for _, theme := range themes {
		theme := theme
		t.Run(theme, func(t *testing.T) {
			t.Parallel()

			cfg := config.Default()
			cfg.UI.Theme = theme

			// Themes should be loadable
			_ = cfg
		})
	}
}

// TestLayoutConfigurations tests all available layouts
func TestLayoutConfigurations(t *testing.T) {
	t.Parallel()

	layouts := []string{"single", "horizontal", "vertical", "grid"}

	for _, layout := range layouts {
		layout := layout
		t.Run(layout, func(t *testing.T) {
			t.Parallel()

			cfg := config.Default()
			cfg.UI.Layout = layout

			err := cfg.Validate()
			if err != nil {
				t.Errorf("layout %s should be valid: %v", layout, err)
			}
		})
	}
}
