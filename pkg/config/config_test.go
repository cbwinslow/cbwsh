package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/config"
	"github.com/cbwinslow/cbwsh/pkg/core"
)

func TestDefault(t *testing.T) {
	cfg := config.Default()

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if cfg.Shell.DefaultShell != core.ShellTypeBash {
		t.Errorf("expected bash shell, got %v", cfg.Shell.DefaultShell)
	}

	if cfg.Shell.HistorySize != 10000 {
		t.Errorf("expected history size 10000, got %d", cfg.Shell.HistorySize)
	}

	if cfg.UI.Theme != "default" {
		t.Errorf("expected default theme, got %s", cfg.UI.Theme)
	}

	if !cfg.UI.ShowStatusBar {
		t.Error("expected status bar to be shown")
	}

	if !cfg.UI.EnableAnimations {
		t.Error("expected animations to be enabled")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create and save config
	cfg := config.Default()
	cfg.UI.Theme = "dracula"
	cfg.Shell.HistorySize = 5000

	err := cfg.Save(configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load config
	loaded, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.UI.Theme != "dracula" {
		t.Errorf("expected theme dracula, got %s", loaded.UI.Theme)
	}

	if loaded.Shell.HistorySize != 5000 {
		t.Errorf("expected history size 5000, got %d", loaded.Shell.HistorySize)
	}
}

func TestLoadNonExistent(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("expected no error for nonexistent file, got %v", err)
	}

	// Should return default config
	if cfg.Shell.DefaultShell != core.ShellTypeBash {
		t.Errorf("expected bash shell, got %v", cfg.Shell.DefaultShell)
	}
}

func TestGetSetShellConfig(t *testing.T) {
	cfg := config.Default()

	shellCfg := cfg.GetShellConfig()
	if shellCfg.DefaultShell != core.ShellTypeBash {
		t.Errorf("expected bash shell")
	}

	shellCfg.DefaultShell = core.ShellTypeZsh
	cfg.SetShellConfig(shellCfg)

	updated := cfg.GetShellConfig()
	if updated.DefaultShell != core.ShellTypeZsh {
		t.Errorf("expected zsh shell after update")
	}
}

func TestGetSetUIConfig(t *testing.T) {
	cfg := config.Default()

	uiCfg := cfg.GetUIConfig()
	if uiCfg.Theme != "default" {
		t.Errorf("expected default theme")
	}

	uiCfg.Theme = "nord"
	cfg.SetUIConfig(uiCfg)

	updated := cfg.GetUIConfig()
	if updated.Theme != "nord" {
		t.Errorf("expected nord theme after update")
	}
}

func TestGetSetAIConfig(t *testing.T) {
	cfg := config.Default()

	aiCfg := cfg.GetAIConfig()
	if aiCfg.Provider != core.AIProviderNone {
		t.Errorf("expected no AI provider")
	}

	aiCfg.Provider = core.AIProviderOpenAI
	aiCfg.Model = "gpt-4"
	cfg.SetAIConfig(aiCfg)

	updated := cfg.GetAIConfig()
	if updated.Provider != core.AIProviderOpenAI {
		t.Errorf("expected OpenAI provider after update")
	}
	if updated.Model != "gpt-4" {
		t.Errorf("expected gpt-4 model after update")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.yaml")

	cfg := config.Default()
	err := cfg.Save(configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}
