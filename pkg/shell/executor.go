// Package shell provides shell execution functionality for cbwsh.
package shell

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

// Executor implements the core.Executor interface for running shell commands.
type Executor struct {
	mu         sync.RWMutex
	shellType  core.ShellType
	workingDir string
	env        map[string]string
	currentCmd *exec.Cmd
	aliases    map[string]string
}

// NewExecutor creates a new shell executor.
func NewExecutor(shellType core.ShellType) *Executor {
	wd, _ := os.Getwd()
	return &Executor{
		shellType:  shellType,
		workingDir: wd,
		env:        make(map[string]string),
		aliases:    make(map[string]string),
	}
}

// Execute runs a command and returns the result.
func (e *Executor) Execute(ctx context.Context, command string) (*core.CommandResult, error) {
	e.mu.Lock()
	command = e.expandAliases(command)

	startTime := time.Now()

	shell := e.getShellPath()
	cmd := exec.CommandContext(ctx, shell, "-c", command)
	cmd.Dir = e.workingDir
	cmd.Env = e.buildEnv()
	e.currentCmd = cmd
	e.mu.Unlock()

	output, err := cmd.CombinedOutput()

	e.mu.Lock()
	e.currentCmd = nil
	e.mu.Unlock()

	result := &core.CommandResult{
		Command:  command,
		Duration: time.Since(startTime).Milliseconds(),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			result.Error = string(output)
		} else {
			result.ExitCode = -1
			result.Error = err.Error()
		}
	} else {
		result.Output = string(output)
		result.ExitCode = 0
	}

	return result, nil
}

// ExecuteAsync runs a command asynchronously and streams the output.
func (e *Executor) ExecuteAsync(ctx context.Context, command string) (<-chan *core.CommandResult, error) {
	resultCh := make(chan *core.CommandResult, 1)

	go func() {
		defer close(resultCh)

		e.mu.Lock()
		command = e.expandAliases(command)
		startTime := time.Now()

		shell := e.getShellPath()
		cmd := exec.CommandContext(ctx, shell, "-c", command)
		cmd.Dir = e.workingDir
		cmd.Env = e.buildEnv()
		e.currentCmd = cmd
		e.mu.Unlock()

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			resultCh <- &core.CommandResult{
				Command:  command,
				Error:    err.Error(),
				ExitCode: -1,
			}
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			resultCh <- &core.CommandResult{
				Command:  command,
				Error:    err.Error(),
				ExitCode: -1,
			}
			return
		}

		if err := cmd.Start(); err != nil {
			resultCh <- &core.CommandResult{
				Command:  command,
				Error:    err.Error(),
				ExitCode: -1,
			}
			return
		}

		var outputBuf, errorBuf strings.Builder
		outputDone := make(chan struct{})
		errorDone := make(chan struct{})

		go func() {
			defer close(outputDone)
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				outputBuf.WriteString(scanner.Text())
				outputBuf.WriteString("\n")
			}
		}()

		go func() {
			defer close(errorDone)
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				errorBuf.WriteString(scanner.Text())
				errorBuf.WriteString("\n")
			}
		}()

		<-outputDone
		<-errorDone

		err = cmd.Wait()

		e.mu.Lock()
		e.currentCmd = nil
		e.mu.Unlock()

		result := &core.CommandResult{
			Command:  command,
			Output:   outputBuf.String(),
			Error:    errorBuf.String(),
			Duration: time.Since(startTime).Milliseconds(),
		}

		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				result.ExitCode = exitErr.ExitCode()
			} else {
				result.ExitCode = -1
			}
		} else {
			result.ExitCode = 0
		}

		resultCh <- result
	}()

	return resultCh, nil
}

// ExecuteInteractive runs an interactive command with PTY support.
func (e *Executor) ExecuteInteractive(ctx context.Context, command string, stdin io.Reader, stdout, stderr io.Writer) error {
	e.mu.Lock()
	command = e.expandAliases(command)
	shell := e.getShellPath()
	cmd := exec.CommandContext(ctx, shell, "-c", command)
	cmd.Dir = e.workingDir
	cmd.Env = e.buildEnv()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	e.currentCmd = cmd
	e.mu.Unlock()

	err := cmd.Run()

	e.mu.Lock()
	e.currentCmd = nil
	e.mu.Unlock()

	return err
}

// Interrupt stops the currently running command.
func (e *Executor) Interrupt() error {
	e.mu.RLock()
	cmd := e.currentCmd
	e.mu.RUnlock()

	if cmd != nil && cmd.Process != nil {
		return cmd.Process.Kill()
	}
	return nil
}

// SetShellType sets the shell type (bash/zsh).
func (e *Executor) SetShellType(shellType core.ShellType) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.shellType = shellType
	return nil
}

// GetShellType returns the current shell type.
func (e *Executor) GetShellType() core.ShellType {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.shellType
}

// SetWorkingDirectory sets the current working directory.
func (e *Executor) SetWorkingDirectory(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.workingDir = path
	return nil
}

// GetWorkingDirectory returns the current working directory.
func (e *Executor) GetWorkingDirectory() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.workingDir
}

// SetEnvironment sets environment variables.
func (e *Executor) SetEnvironment(env map[string]string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for k, v := range env {
		e.env[k] = v
	}
	return nil
}

// GetEnvironment returns current environment variables.
func (e *Executor) GetEnvironment() map[string]string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make(map[string]string, len(e.env))
	for k, v := range e.env {
		result[k] = v
	}
	return result
}

// SetAlias sets a command alias.
func (e *Executor) SetAlias(name, command string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.aliases[name] = command
}

// RemoveAlias removes a command alias.
func (e *Executor) RemoveAlias(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.aliases, name)
}

// GetAliases returns all aliases.
func (e *Executor) GetAliases() map[string]string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make(map[string]string, len(e.aliases))
	for k, v := range e.aliases {
		result[k] = v
	}
	return result
}

func (e *Executor) getShellPath() string {
	switch e.shellType {
	case core.ShellTypeBash:
		return "/bin/bash"
	case core.ShellTypeZsh:
		return "/bin/zsh"
	default:
		return "/bin/sh"
	}
}

func (e *Executor) buildEnv() []string {
	baseEnv := os.Environ()
	result := make([]string, 0, len(baseEnv)+len(e.env))
	result = append(result, baseEnv...)
	for k, v := range e.env {
		result = append(result, k+"="+v)
	}
	return result
}

func (e *Executor) expandAliases(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return command
	}

	if expanded, ok := e.aliases[parts[0]]; ok {
		parts[0] = expanded
		return strings.Join(parts, " ")
	}

	return command
}

// History manages command history.
type History struct {
	mu       sync.RWMutex
	commands []string
	position int
	maxSize  int
	filePath string
}

// NewHistory creates a new history manager.
func NewHistory(maxSize int, filePath string) *History {
	return &History{
		commands: make([]string, 0, maxSize),
		position: -1,
		maxSize:  maxSize,
		filePath: filePath,
	}
}

// Add adds a command to history.
func (h *History) Add(command string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Don't add duplicates of the last command
	if len(h.commands) > 0 && h.commands[len(h.commands)-1] == command {
		return
	}

	h.commands = append(h.commands, command)
	if len(h.commands) > h.maxSize {
		h.commands = h.commands[1:]
	}
	h.position = len(h.commands)
}

// Previous returns the previous command in history.
func (h *History) Previous() (string, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.position > 0 {
		h.position--
		return h.commands[h.position], true
	}
	return "", false
}

// Next returns the next command in history.
func (h *History) Next() (string, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.position < len(h.commands)-1 {
		h.position++
		return h.commands[h.position], true
	}
	h.position = len(h.commands)
	return "", false
}

// Reset resets the history position to the end.
func (h *History) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.position = len(h.commands)
}

// Search searches history for commands matching the query.
func (h *History) Search(query string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var results []string
	for i := len(h.commands) - 1; i >= 0; i-- {
		if strings.Contains(h.commands[i], query) {
			results = append(results, h.commands[i])
		}
	}
	return results
}

// Load loads history from file.
func (h *History) Load() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := os.ReadFile(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			h.commands = append(h.commands, line)
		}
	}

	if len(h.commands) > h.maxSize {
		h.commands = h.commands[len(h.commands)-h.maxSize:]
	}

	h.position = len(h.commands)
	return nil
}

// Save saves history to file.
func (h *History) Save() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	dir := strings.TrimSuffix(h.filePath, "/"+strings.Split(h.filePath, "/")[len(strings.Split(h.filePath, "/"))-1])
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	data := strings.Join(h.commands, "\n")
	return os.WriteFile(h.filePath, []byte(data), 0o600)
}

// All returns all commands in history.
func (h *History) All() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make([]string, len(h.commands))
	copy(result, h.commands)
	return result
}

// Clear clears all history.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.commands = h.commands[:0]
	h.position = 0
}
