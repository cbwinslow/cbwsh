// Package posix provides POSIX-compliant system calls and signal handling for cbwsh.
package posix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Signal represents a POSIX signal.
type Signal syscall.Signal

// Common POSIX signals.
const (
	SIGHUP    Signal = Signal(syscall.SIGHUP)
	SIGINT    Signal = Signal(syscall.SIGINT)
	SIGQUIT   Signal = Signal(syscall.SIGQUIT)
	SIGILL    Signal = Signal(syscall.SIGILL)
	SIGTRAP   Signal = Signal(syscall.SIGTRAP)
	SIGABRT   Signal = Signal(syscall.SIGABRT)
	SIGBUS    Signal = Signal(syscall.SIGBUS)
	SIGFPE    Signal = Signal(syscall.SIGFPE)
	SIGKILL   Signal = Signal(syscall.SIGKILL)
	SIGUSR1   Signal = Signal(syscall.SIGUSR1)
	SIGSEGV   Signal = Signal(syscall.SIGSEGV)
	SIGUSR2   Signal = Signal(syscall.SIGUSR2)
	SIGPIPE   Signal = Signal(syscall.SIGPIPE)
	SIGALRM   Signal = Signal(syscall.SIGALRM)
	SIGTERM   Signal = Signal(syscall.SIGTERM)
	SIGCHLD   Signal = Signal(syscall.SIGCHLD)
	SIGCONT   Signal = Signal(syscall.SIGCONT)
	SIGSTOP   Signal = Signal(syscall.SIGSTOP)
	SIGTSTP   Signal = Signal(syscall.SIGTSTP)
	SIGTTIN   Signal = Signal(syscall.SIGTTIN)
	SIGTTOU   Signal = Signal(syscall.SIGTTOU)
	SIGURG    Signal = Signal(syscall.SIGURG)
	SIGXCPU   Signal = Signal(syscall.SIGXCPU)
	SIGXFSZ   Signal = Signal(syscall.SIGXFSZ)
	SIGVTALRM Signal = Signal(syscall.SIGVTALRM)
	SIGPROF   Signal = Signal(syscall.SIGPROF)
	SIGWINCH  Signal = Signal(syscall.SIGWINCH)
	SIGIO     Signal = Signal(syscall.SIGIO)
	SIGSYS    Signal = Signal(syscall.SIGSYS)
)

// String returns the name of the signal.
func (s Signal) String() string {
	names := map[Signal]string{
		SIGHUP:    "SIGHUP",
		SIGINT:    "SIGINT",
		SIGQUIT:   "SIGQUIT",
		SIGILL:    "SIGILL",
		SIGTRAP:   "SIGTRAP",
		SIGABRT:   "SIGABRT",
		SIGBUS:    "SIGBUS",
		SIGFPE:    "SIGFPE",
		SIGKILL:   "SIGKILL",
		SIGUSR1:   "SIGUSR1",
		SIGSEGV:   "SIGSEGV",
		SIGUSR2:   "SIGUSR2",
		SIGPIPE:   "SIGPIPE",
		SIGALRM:   "SIGALRM",
		SIGTERM:   "SIGTERM",
		SIGCHLD:   "SIGCHLD",
		SIGCONT:   "SIGCONT",
		SIGSTOP:   "SIGSTOP",
		SIGTSTP:   "SIGTSTP",
		SIGTTIN:   "SIGTTIN",
		SIGTTOU:   "SIGTTOU",
		SIGURG:    "SIGURG",
		SIGXCPU:   "SIGXCPU",
		SIGXFSZ:   "SIGXFSZ",
		SIGVTALRM: "SIGVTALRM",
		SIGPROF:   "SIGPROF",
		SIGWINCH:  "SIGWINCH",
		SIGIO:     "SIGIO",
		SIGSYS:    "SIGSYS",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return fmt.Sprintf("signal(%d)", int(s))
}

// SignalNumber returns the signal number.
func (s Signal) SignalNumber() int {
	return int(s)
}

// SignalHandler is a function that handles a signal.
type SignalHandler func(sig Signal)

// SignalManager manages signal handling.
type SignalManager struct {
	mu          sync.RWMutex
	handlers    map[Signal][]SignalHandler
	signalChan  chan os.Signal
	stopChan    chan struct{}
	running     bool
	defaultSigs []Signal
}

// NewSignalManager creates a new signal manager.
func NewSignalManager() *SignalManager {
	return &SignalManager{
		handlers:    make(map[Signal][]SignalHandler),
		signalChan:  make(chan os.Signal, 10),
		stopChan:    make(chan struct{}),
		defaultSigs: []Signal{SIGINT, SIGTERM, SIGHUP, SIGQUIT},
	}
}

// RegisterHandler registers a handler for a signal.
func (m *SignalManager) RegisterHandler(sig Signal, handler SignalHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[sig] = append(m.handlers[sig], handler)
}

// UnregisterHandlers removes all handlers for a signal.
func (m *SignalManager) UnregisterHandlers(sig Signal) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, sig)
}

// Start starts the signal manager.
func (m *SignalManager) Start(signals ...Signal) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return errors.New("signal manager already running")
	}
	m.running = true
	m.mu.Unlock()

	if len(signals) == 0 {
		signals = m.defaultSigs
	}

	// Convert to os.Signal
	osSigs := make([]os.Signal, len(signals))
	for i, sig := range signals {
		osSigs[i] = syscall.Signal(sig)
	}

	signal.Notify(m.signalChan, osSigs...)

	go m.processSignals()

	return nil
}

// Stop stops the signal manager.
func (m *SignalManager) Stop() {
	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		return
	}
	m.running = false
	m.mu.Unlock()

	close(m.stopChan)
	signal.Stop(m.signalChan)
}

func (m *SignalManager) processSignals() {
	for {
		select {
		case <-m.stopChan:
			return
		case sig := <-m.signalChan:
			m.mu.RLock()
			handlers := m.handlers[Signal(sig.(syscall.Signal))]
			m.mu.RUnlock()

			for _, handler := range handlers {
				handler(Signal(sig.(syscall.Signal)))
			}
		}
	}
}

// Wait waits for a signal from the specified set.
func (m *SignalManager) Wait(ctx context.Context, signals ...Signal) (Signal, error) {
	sigChan := make(chan os.Signal, 1)
	defer signal.Stop(sigChan)

	osSigs := make([]os.Signal, len(signals))
	for i, sig := range signals {
		osSigs[i] = syscall.Signal(sig)
	}
	signal.Notify(sigChan, osSigs...)

	select {
	case sig := <-sigChan:
		return Signal(sig.(syscall.Signal)), nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// Send sends a signal to a process.
func Send(pid int, sig Signal) error {
	return syscall.Kill(pid, syscall.Signal(sig))
}

// SendToGroup sends a signal to a process group.
func SendToGroup(pgid int, sig Signal) error {
	return syscall.Kill(-pgid, syscall.Signal(sig))
}

// FileDescriptor represents a file descriptor.
type FileDescriptor int

// Standard file descriptors.
const (
	Stdin  FileDescriptor = 0
	Stdout FileDescriptor = 1
	Stderr FileDescriptor = 2
)

// Dup duplicates a file descriptor.
func Dup(fd FileDescriptor) (FileDescriptor, error) {
	newFD, err := syscall.Dup(int(fd))
	return FileDescriptor(newFD), err
}

// Dup2 duplicates a file descriptor to a specific number.
func Dup2(oldFD, newFD FileDescriptor) error {
	return syscall.Dup2(int(oldFD), int(newFD))
}

// Close closes a file descriptor.
func Close(fd FileDescriptor) error {
	return syscall.Close(int(fd))
}

// Pipe creates a pipe.
func Pipe() (read, write FileDescriptor, err error) {
	var fds [2]int
	err = syscall.Pipe(fds[:])
	return FileDescriptor(fds[0]), FileDescriptor(fds[1]), err
}

// ProcessInfo contains information about a process.
type ProcessInfo struct {
	PID        int
	PPID       int
	PGID       int
	SID        int
	UID        int
	GID        int
	EUID       int
	EGID       int
	Umask      int
	WorkingDir string
}

// GetProcessInfo returns information about the current process.
func GetProcessInfo() (*ProcessInfo, error) {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	pgid, _ := syscall.Getpgid(0)
	// Get session ID - not available on all platforms
	sid := 0

	return &ProcessInfo{
		PID:        os.Getpid(),
		PPID:       os.Getppid(),
		PGID:       pgid,
		SID:        sid,
		UID:        os.Getuid(),
		GID:        os.Getgid(),
		EUID:       os.Geteuid(),
		EGID:       os.Getegid(),
		WorkingDir: wd,
	}, nil
}

// Fork forks the current process (wrapper around exec.Command for safety).
// Note: Direct fork is not available in Go; use exec.Command for child processes.
func Fork() error {
	return errors.New("direct fork not available in Go; use exec.Command")
}

// Exec replaces the current process with a new one.
func Exec(path string, args []string, env []string) error {
	return syscall.Exec(path, args, env)
}

// SetProcessGroup sets the process group ID.
func SetProcessGroup(pid, pgid int) error {
	return syscall.Setpgid(pid, pgid)
}

// GetProcessGroup gets the process group ID.
func GetProcessGroup(pid int) (int, error) {
	return syscall.Getpgid(pid)
}

// CreateSession creates a new session.
func CreateSession() (int, error) {
	return syscall.Setsid()
}

// SetUmask sets the file mode creation mask.
func SetUmask(mask int) int {
	return syscall.Umask(mask)
}

// ChangeDirectory changes the current working directory.
func ChangeDirectory(path string) error {
	return syscall.Chdir(path)
}

// ChangeRoot changes the root directory.
func ChangeRoot(path string) error {
	return syscall.Chroot(path)
}

// ResourceLimit holds resource limit values.
type ResourceLimit struct {
	Current uint64
	Maximum uint64
}

// ResourceType represents a resource type for limits.
type ResourceType int

// Resource types.
const (
	ResourceCPU    ResourceType = ResourceType(syscall.RLIMIT_CPU)
	ResourceFSize  ResourceType = ResourceType(syscall.RLIMIT_FSIZE)
	ResourceData   ResourceType = ResourceType(syscall.RLIMIT_DATA)
	ResourceStack  ResourceType = ResourceType(syscall.RLIMIT_STACK)
	ResourceCore   ResourceType = ResourceType(syscall.RLIMIT_CORE)
	ResourceNoFile ResourceType = ResourceType(syscall.RLIMIT_NOFILE)
	ResourceAS     ResourceType = ResourceType(syscall.RLIMIT_AS)
)

// GetResourceLimit gets a resource limit.
func GetResourceLimit(resource ResourceType) (*ResourceLimit, error) {
	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(int(resource), &rlimit)
	if err != nil {
		return nil, err
	}
	return &ResourceLimit{
		Current: rlimit.Cur,
		Maximum: rlimit.Max,
	}, nil
}

// SetResourceLimit sets a resource limit.
func SetResourceLimit(resource ResourceType, limit *ResourceLimit) error {
	rlimit := syscall.Rlimit{
		Cur: limit.Current,
		Max: limit.Maximum,
	}
	return syscall.Setrlimit(int(resource), &rlimit)
}

// Timer represents a POSIX timer.
type Timer struct {
	duration time.Duration
	callback func()
	timer    *time.Timer
	mu       sync.Mutex
	stopped  bool
}

// NewTimer creates a new timer.
func NewTimer(duration time.Duration, callback func()) *Timer {
	return &Timer{
		duration: duration,
		callback: callback,
	}
}

// Start starts the timer.
func (t *Timer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.timer != nil {
		return
	}

	t.stopped = false
	t.timer = time.AfterFunc(t.duration, func() {
		t.mu.Lock()
		stopped := t.stopped
		t.mu.Unlock()

		if !stopped && t.callback != nil {
			t.callback()
		}
	})
}

// Stop stops the timer.
func (t *Timer) Stop() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.stopped = true
	if t.timer != nil {
		return t.timer.Stop()
	}
	return false
}

// Reset resets the timer.
func (t *Timer) Reset(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.duration = duration
	t.stopped = false
	if t.timer != nil {
		t.timer.Reset(duration)
	}
}

// ExitStatus represents a process exit status.
type ExitStatus struct {
	Code       int
	Signal     Signal
	Signaled   bool
	CoreDumped bool
}

// ParseWaitStatus parses a wait status.
func ParseWaitStatus(status syscall.WaitStatus) *ExitStatus {
	return &ExitStatus{
		Code:       status.ExitStatus(),
		Signal:     Signal(status.Signal()),
		Signaled:   status.Signaled(),
		CoreDumped: status.CoreDump(),
	}
}

// WaitPID waits for a specific process.
func WaitPID(pid int, options int) (*ExitStatus, error) {
	var status syscall.WaitStatus
	_, err := syscall.Wait4(pid, &status, options, nil)
	if err != nil {
		return nil, err
	}
	return ParseWaitStatus(status), nil
}

// Wait options.
const (
	WaitNoHang    = syscall.WNOHANG
	WaitUntraced  = syscall.WUNTRACED
	WaitContinued = syscall.WCONTINUED
)

// EnvironmentVariable represents an environment variable.
type EnvironmentVariable struct {
	Name  string
	Value string
}

// GetEnvironment returns all environment variables.
func GetEnvironment() []EnvironmentVariable {
	environ := os.Environ()
	result := make([]EnvironmentVariable, 0, len(environ))

	for _, env := range environ {
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				result = append(result, EnvironmentVariable{
					Name:  env[:i],
					Value: env[i+1:],
				})
				break
			}
		}
	}

	return result
}

// SetEnvironmentVariable sets an environment variable.
func SetEnvironmentVariable(name, value string) error {
	return os.Setenv(name, value)
}

// UnsetEnvironmentVariable unsets an environment variable.
func UnsetEnvironmentVariable(name string) error {
	return os.Unsetenv(name)
}

// ClearEnvironment clears all environment variables.
func ClearEnvironment() {
	os.Clearenv()
}
