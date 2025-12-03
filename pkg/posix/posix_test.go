package posix_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/posix"
)

func TestSignalString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		signal   posix.Signal
		expected string
	}{
		{posix.SIGHUP, "SIGHUP"},
		{posix.SIGINT, "SIGINT"},
		{posix.SIGTERM, "SIGTERM"},
		{posix.SIGKILL, "SIGKILL"},
		{posix.SIGCHLD, "SIGCHLD"},
		{posix.Signal(999), "signal(999)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			result := tt.signal.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSignalNumber(t *testing.T) {
	t.Parallel()

	sig := posix.SIGTERM
	num := sig.SignalNumber()
	if num != 15 {
		t.Errorf("expected SIGTERM to be 15, got %d", num)
	}
}

func TestNewSignalManager(t *testing.T) {
	t.Parallel()

	manager := posix.NewSignalManager()
	if manager == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestSignalManagerRegisterHandler(t *testing.T) {
	t.Parallel()

	manager := posix.NewSignalManager()
	called := false

	manager.RegisterHandler(posix.SIGUSR1, func(sig posix.Signal) {
		called = true
	})

	// Just verify registration doesn't panic
	manager.UnregisterHandlers(posix.SIGUSR1)

	if called {
		t.Error("handler should not have been called yet")
	}
}

func TestGetProcessInfo(t *testing.T) {
	t.Parallel()

	info, err := posix.GetProcessInfo()
	if err != nil {
		t.Fatalf("failed to get process info: %v", err)
	}

	if info.PID <= 0 {
		t.Errorf("expected positive PID, got %d", info.PID)
	}

	if info.UID < 0 {
		t.Errorf("expected non-negative UID, got %d", info.UID)
	}
}

func TestPipe(t *testing.T) {
	t.Parallel()

	read, write, err := posix.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	defer posix.Close(read)
	defer posix.Close(write)

	if read < 0 || write < 0 {
		t.Error("expected valid file descriptors")
	}
}

func TestDup(t *testing.T) {
	t.Parallel()

	// Create a temp file
	f, err := os.CreateTemp("", "posix_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	fd := posix.FileDescriptor(f.Fd())
	newFD, err := posix.Dup(fd)
	if err != nil {
		t.Fatalf("failed to dup: %v", err)
	}
	defer posix.Close(newFD)

	if newFD <= 0 {
		t.Errorf("expected valid file descriptor, got %d", newFD)
	}
}

func TestGetResourceLimit(t *testing.T) {
	t.Parallel()

	limit, err := posix.GetResourceLimit(posix.ResourceNoFile)
	if err != nil {
		t.Fatalf("failed to get resource limit: %v", err)
	}

	if limit.Current == 0 && limit.Maximum == 0 {
		t.Error("expected non-zero limits")
	}
}

func TestChangeDirectory(t *testing.T) {
	t.Parallel()

	// Save current directory
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(orig) }()

	err = posix.ChangeDirectory("/tmp")
	if err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get new working directory: %v", err)
	}

	if wd != "/tmp" {
		t.Errorf("expected /tmp, got %s", wd)
	}
}

func TestSetUmask(t *testing.T) {
	t.Parallel()

	oldMask := posix.SetUmask(0o022)
	newMask := posix.SetUmask(oldMask)

	if newMask != 0o022 {
		t.Errorf("expected umask 022, got %o", newMask)
	}
}

func TestTimer(t *testing.T) {
	t.Parallel()

	called := make(chan bool, 1)
	timer := posix.NewTimer(50*time.Millisecond, func() {
		called <- true
	})

	timer.Start()

	select {
	case <-called:
		// Success
	case <-time.After(200 * time.Millisecond):
		t.Error("timer callback was not called")
	}
}

func TestTimerStop(t *testing.T) {
	t.Parallel()

	called := make(chan bool, 1)
	timer := posix.NewTimer(100*time.Millisecond, func() {
		called <- true
	})

	timer.Start()
	timer.Stop()

	select {
	case <-called:
		t.Error("timer callback should not have been called after stop")
	case <-time.After(200 * time.Millisecond):
		// Success - callback was not called
	}
}

func TestTimerReset(t *testing.T) {
	t.Parallel()

	called := make(chan bool, 1)
	timer := posix.NewTimer(50*time.Millisecond, func() {
		called <- true
	})

	timer.Start()
	timer.Reset(10 * time.Millisecond)

	select {
	case <-called:
		// Success
	case <-time.After(200 * time.Millisecond):
		t.Error("timer callback was not called after reset")
	}
}

func TestGetEnvironment(t *testing.T) {
	t.Parallel()

	env := posix.GetEnvironment()
	if len(env) == 0 {
		t.Error("expected non-empty environment")
	}

	// Should have at least PATH
	hasPath := false
	for _, v := range env {
		if v.Name == "PATH" {
			hasPath = true
			break
		}
	}

	if !hasPath {
		t.Error("expected PATH in environment")
	}
}

func TestSetEnvironmentVariable(t *testing.T) {
	varName := "CBWSH_TEST_VAR"
	varValue := "test_value"

	err := posix.SetEnvironmentVariable(varName, varValue)
	if err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}

	if os.Getenv(varName) != varValue {
		t.Errorf("expected %s, got %s", varValue, os.Getenv(varName))
	}

	err = posix.UnsetEnvironmentVariable(varName)
	if err != nil {
		t.Fatalf("failed to unset environment variable: %v", err)
	}

	if os.Getenv(varName) != "" {
		t.Error("expected empty value after unset")
	}
}

func TestSignalManagerWait(t *testing.T) {
	t.Parallel()

	manager := posix.NewSignalManager()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := manager.Wait(ctx, posix.SIGUSR1)
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestSendSignal(t *testing.T) {
	t.Parallel()

	// Send signal 0 (no-op, just check if process exists) to self
	err := posix.Send(os.Getpid(), posix.Signal(0))
	if err != nil {
		t.Errorf("failed to send signal 0 to self: %v", err)
	}
}

func TestGetProcessGroup(t *testing.T) {
	t.Parallel()

	pgid, err := posix.GetProcessGroup(os.Getpid())
	if err != nil {
		t.Fatalf("failed to get process group: %v", err)
	}

	if pgid <= 0 {
		t.Errorf("expected positive PGID, got %d", pgid)
	}
}
