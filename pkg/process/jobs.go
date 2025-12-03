// Package process provides job control and process management for cbwsh.
package process

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// JobState represents the state of a job.
type JobState int

const (
	// JobStateRunning indicates a running job.
	JobStateRunning JobState = iota
	// JobStateStopped indicates a stopped job.
	JobStateStopped
	// JobStateCompleted indicates a completed job.
	JobStateCompleted
	// JobStateFailed indicates a failed job.
	JobStateFailed
)

// String returns the string representation of the job state.
func (s JobState) String() string {
	switch s {
	case JobStateRunning:
		return "Running"
	case JobStateStopped:
		return "Stopped"
	case JobStateCompleted:
		return "Completed"
	case JobStateFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// Job represents a background job.
type Job struct {
	ID          int
	Command     string
	State       JobState
	PID         int
	StartTime   time.Time
	EndTime     time.Time
	ExitCode    int
	Background  bool
	cmd         *exec.Cmd
	mu          sync.RWMutex
	cancel      context.CancelFunc
	stoppedChan chan struct{}
}

// GetState returns the current job state.
func (j *Job) GetState() JobState {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.State
}

// SetState sets the job state.
func (j *Job) SetState(state JobState) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.State = state
	if state == JobStateCompleted || state == JobStateFailed {
		j.EndTime = time.Now()
	}
}

// GetPID returns the process ID.
func (j *Job) GetPID() int {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.PID
}

// GetExitCode returns the exit code.
func (j *Job) GetExitCode() int {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.ExitCode
}

// Duration returns how long the job has been running.
func (j *Job) Duration() time.Duration {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if !j.EndTime.IsZero() {
		return j.EndTime.Sub(j.StartTime)
	}
	return time.Since(j.StartTime)
}

// String returns a string representation of the job.
func (j *Job) String() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return fmt.Sprintf("[%d] %s %s (PID: %d)", j.ID, j.State.String(), j.Command, j.PID)
}

// JobManager manages background jobs.
type JobManager struct {
	mu        sync.RWMutex
	jobs      map[int]*Job
	nextJobID atomic.Int64
	maxJobs   int
}

// NewJobManager creates a new job manager.
func NewJobManager(maxJobs int) *JobManager {
	if maxJobs <= 0 {
		maxJobs = 100
	}
	return &JobManager{
		jobs:    make(map[int]*Job),
		maxJobs: maxJobs,
	}
}

// StartJob starts a new background job.
func (m *JobManager) StartJob(ctx context.Context, command string, shell string) (*Job, error) {
	m.mu.Lock()

	// Check if we've reached the maximum number of jobs
	activeJobs := 0
	for _, job := range m.jobs {
		if job.GetState() == JobStateRunning || job.GetState() == JobStateStopped {
			activeJobs++
		}
	}
	if activeJobs >= m.maxJobs {
		m.mu.Unlock()
		return nil, errors.New("maximum number of jobs reached")
	}

	jobID := int(m.nextJobID.Add(1))

	jobCtx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(jobCtx, shell, "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Create a new process group
	}

	job := &Job{
		ID:          jobID,
		Command:     command,
		State:       JobStateRunning,
		StartTime:   time.Now(),
		Background:  true,
		cmd:         cmd,
		cancel:      cancel,
		stoppedChan: make(chan struct{}),
	}

	m.jobs[jobID] = job
	m.mu.Unlock()

	if err := cmd.Start(); err != nil {
		m.mu.Lock()
		delete(m.jobs, jobID)
		m.mu.Unlock()
		cancel()
		return nil, fmt.Errorf("failed to start job: %w", err)
	}

	job.mu.Lock()
	job.PID = cmd.Process.Pid
	job.mu.Unlock()

	// Monitor job in background
	go m.monitorJob(job)

	return job, nil
}

func (m *JobManager) monitorJob(job *Job) {
	defer close(job.stoppedChan)

	err := job.cmd.Wait()

	job.mu.Lock()
	defer job.mu.Unlock()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			job.ExitCode = exitErr.ExitCode()
			if exitErr.ExitCode() == -1 {
				// Process was killed
				job.State = JobStateStopped
			} else {
				job.State = JobStateFailed
			}
		} else {
			job.State = JobStateFailed
			job.ExitCode = -1
		}
	} else {
		job.State = JobStateCompleted
		job.ExitCode = 0
	}

	job.EndTime = time.Now()
}

// StopJob stops a running job.
func (m *JobManager) StopJob(jobID int) error {
	m.mu.RLock()
	job, exists := m.jobs[jobID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job %d not found", jobID)
	}

	state := job.GetState()
	if state != JobStateRunning {
		return fmt.Errorf("job %d is not running (state: %s)", jobID, state.String())
	}

	// Send SIGSTOP to the process group
	if job.cmd.Process != nil {
		pgid, err := syscall.Getpgid(job.cmd.Process.Pid)
		if err == nil {
			_ = syscall.Kill(-pgid, syscall.SIGSTOP)
		}
	}

	job.SetState(JobStateStopped)
	return nil
}

// ContinueJob continues a stopped job.
func (m *JobManager) ContinueJob(jobID int) error {
	m.mu.RLock()
	job, exists := m.jobs[jobID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job %d not found", jobID)
	}

	state := job.GetState()
	if state != JobStateStopped {
		return fmt.Errorf("job %d is not stopped (state: %s)", jobID, state.String())
	}

	// Send SIGCONT to the process group
	if job.cmd.Process != nil {
		pgid, err := syscall.Getpgid(job.cmd.Process.Pid)
		if err == nil {
			_ = syscall.Kill(-pgid, syscall.SIGCONT)
		}
	}

	job.SetState(JobStateRunning)
	return nil
}

// KillJob terminates a job.
func (m *JobManager) KillJob(jobID int) error {
	m.mu.RLock()
	job, exists := m.jobs[jobID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job %d not found", jobID)
	}

	state := job.GetState()
	if state == JobStateCompleted || state == JobStateFailed {
		return nil // Already terminated
	}

	// Cancel the context and kill the process
	if job.cancel != nil {
		job.cancel()
	}

	if job.cmd.Process != nil {
		// Kill the entire process group
		pgid, err := syscall.Getpgid(job.cmd.Process.Pid)
		if err == nil {
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
		} else {
			_ = job.cmd.Process.Kill()
		}
	}

	return nil
}

// GetJob returns a job by ID.
func (m *JobManager) GetJob(jobID int) (*Job, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, exists := m.jobs[jobID]
	return job, exists
}

// ListJobs returns all jobs.
func (m *JobManager) ListJobs() []*Job {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		result = append(result, job)
	}
	return result
}

// ListActiveJobs returns only running or stopped jobs.
func (m *JobManager) ListActiveJobs() []*Job {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Job
	for _, job := range m.jobs {
		state := job.GetState()
		if state == JobStateRunning || state == JobStateStopped {
			result = append(result, job)
		}
	}
	return result
}

// CleanupCompleted removes completed/failed jobs older than the given duration.
func (m *JobManager) CleanupCompleted(maxAge time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	cutoff := time.Now().Add(-maxAge)

	for id, job := range m.jobs {
		state := job.GetState()
		if (state == JobStateCompleted || state == JobStateFailed) && !job.EndTime.IsZero() && job.EndTime.Before(cutoff) {
			delete(m.jobs, id)
			count++
		}
	}

	return count
}

// WaitForJob waits for a job to complete.
func (m *JobManager) WaitForJob(jobID int, timeout time.Duration) error {
	m.mu.RLock()
	job, exists := m.jobs[jobID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job %d not found", jobID)
	}

	select {
	case <-job.stoppedChan:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for job %d", jobID)
	}
}

// ProcessInfo holds information about a process.
type ProcessInfo struct {
	PID        int
	PPID       int
	Command    string
	State      string
	MemoryMB   float64
	CPUPercent float64
	StartTime  time.Time
	User       string
}

// GetProcessInfo returns information about a process.
func GetProcessInfo(pid int) (*ProcessInfo, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("process %d not found: %w", pid, err)
	}

	// Check if process exists by sending signal 0
	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		return nil, fmt.Errorf("process %d not running: %w", pid, err)
	}

	// Basic info - more detailed info would require reading /proc on Linux
	return &ProcessInfo{
		PID: pid,
	}, nil
}

// SendSignal sends a signal to a process.
func SendSignal(pid int, signal syscall.Signal) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process %d not found: %w", pid, err)
	}
	return proc.Signal(signal)
}

// KillProcess kills a process.
func KillProcess(pid int) error {
	return SendSignal(pid, syscall.SIGKILL)
}

// TerminateProcess terminates a process gracefully.
func TerminateProcess(pid int) error {
	return SendSignal(pid, syscall.SIGTERM)
}

// InterruptProcess interrupts a process (Ctrl+C equivalent).
func InterruptProcess(pid int) error {
	return SendSignal(pid, syscall.SIGINT)
}
