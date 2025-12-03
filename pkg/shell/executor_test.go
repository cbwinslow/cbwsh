package shell_test

import (
	"context"
	"testing"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/shell"
)

func TestNewExecutor(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)
	if exec == nil {
		t.Fatal("expected non-nil executor")
	}

	if exec.GetShellType() != core.ShellTypeBash {
		t.Errorf("expected bash shell type")
	}

	wd := exec.GetWorkingDirectory()
	if wd == "" {
		t.Error("expected non-empty working directory")
	}
}

func TestExecuteSimpleCommand(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)
	ctx := context.Background()

	result, err := exec.Execute(ctx, "echo hello")
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}

	if result.Output != "hello\n" {
		t.Errorf("expected 'hello\\n', got '%s'", result.Output)
	}
}

func TestExecuteFailingCommand(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)
	ctx := context.Background()

	result, err := exec.Execute(ctx, "exit 1")
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
}

func TestSetGetShellType(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)

	err := exec.SetShellType(core.ShellTypeZsh)
	if err != nil {
		t.Fatalf("set shell type failed: %v", err)
	}

	if exec.GetShellType() != core.ShellTypeZsh {
		t.Error("expected zsh shell type")
	}
}

func TestSetGetWorkingDirectory(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)

	err := exec.SetWorkingDirectory("/tmp")
	if err != nil {
		t.Fatalf("set working directory failed: %v", err)
	}

	wd := exec.GetWorkingDirectory()
	if wd != "/tmp" {
		t.Errorf("expected /tmp, got %s", wd)
	}
}

func TestSetWorkingDirectoryInvalid(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)

	err := exec.SetWorkingDirectory("/nonexistent/path")
	if err == nil {
		t.Error("expected error for invalid directory")
	}
}

func TestSetGetEnvironment(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)

	env := map[string]string{
		"TEST_VAR": "test_value",
	}

	err := exec.SetEnvironment(env)
	if err != nil {
		t.Fatalf("set environment failed: %v", err)
	}

	result := exec.GetEnvironment()
	if result["TEST_VAR"] != "test_value" {
		t.Errorf("expected test_value, got %s", result["TEST_VAR"])
	}
}

func TestAliases(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)

	exec.SetAlias("ll", "ls -la")

	aliases := exec.GetAliases()
	if aliases["ll"] != "ls -la" {
		t.Errorf("expected 'ls -la', got '%s'", aliases["ll"])
	}

	exec.RemoveAlias("ll")
	aliases = exec.GetAliases()
	if _, exists := aliases["ll"]; exists {
		t.Error("expected alias to be removed")
	}
}

func TestHistory(t *testing.T) {
	history := shell.NewHistory(100, "/tmp/test_history")

	history.Add("command1")
	history.Add("command2")
	history.Add("command3")

	all := history.All()
	if len(all) != 3 {
		t.Errorf("expected 3 commands, got %d", len(all))
	}

	// Navigate history
	cmd, ok := history.Previous()
	if !ok || cmd != "command3" {
		t.Errorf("expected command3, got %s", cmd)
	}

	cmd, ok = history.Previous()
	if !ok || cmd != "command2" {
		t.Errorf("expected command2, got %s", cmd)
	}

	cmd, ok = history.Next()
	if !ok || cmd != "command3" {
		t.Errorf("expected command3, got %s", cmd)
	}

	// Reset
	history.Reset()
	_, ok = history.Next()
	if ok {
		t.Error("expected no more history after reset")
	}
}

func TestHistorySearch(t *testing.T) {
	history := shell.NewHistory(100, "/tmp/test_history")

	history.Add("git status")
	history.Add("git commit")
	history.Add("ls -la")
	history.Add("git push")

	results := history.Search("git")
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestHistoryNoDuplicates(t *testing.T) {
	history := shell.NewHistory(100, "/tmp/test_history")

	history.Add("command1")
	history.Add("command1")
	history.Add("command1")

	all := history.All()
	if len(all) != 1 {
		t.Errorf("expected 1 command (no duplicates), got %d", len(all))
	}
}

func TestHistoryMaxSize(t *testing.T) {
	history := shell.NewHistory(3, "/tmp/test_history")

	history.Add("command1")
	history.Add("command2")
	history.Add("command3")
	history.Add("command4")

	all := history.All()
	if len(all) != 3 {
		t.Errorf("expected 3 commands (max size), got %d", len(all))
	}

	if all[0] != "command2" {
		t.Errorf("expected command2 as first, got %s", all[0])
	}
}

func TestHistoryClear(t *testing.T) {
	history := shell.NewHistory(100, "/tmp/test_history")

	history.Add("command1")
	history.Add("command2")

	history.Clear()

	all := history.All()
	if len(all) != 0 {
		t.Errorf("expected empty history, got %d commands", len(all))
	}
}

func TestExecuteAsync(t *testing.T) {
	exec := shell.NewExecutor(core.ShellTypeBash)
	ctx := context.Background()

	resultCh, err := exec.ExecuteAsync(ctx, "echo async")
	if err != nil {
		t.Fatalf("execute async failed: %v", err)
	}

	select {
	case result := <-resultCh:
		if result.ExitCode != 0 {
			t.Errorf("expected exit code 0, got %d", result.ExitCode)
		}
		if result.Output != "async\n" {
			t.Errorf("expected 'async\\n', got '%s'", result.Output)
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for async result")
	}
}
