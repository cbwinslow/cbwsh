package process_test

import (
	"context"
	"testing"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/process"
)

func TestJobStateString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state    process.JobState
		expected string
	}{
		{process.JobStateRunning, "Running"},
		{process.JobStateStopped, "Stopped"},
		{process.JobStateCompleted, "Completed"},
		{process.JobStateFailed, "Failed"},
		{process.JobState(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNewJobManager(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	if manager == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestJobManagerStartJob(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	job, err := manager.StartJob(ctx, "sleep 0.1", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	if job.ID <= 0 {
		t.Error("expected positive job ID")
	}

	if job.Command != "sleep 0.1" {
		t.Errorf("expected command 'sleep 0.1', got '%s'", job.Command)
	}

	// Wait for completion
	time.Sleep(200 * time.Millisecond)

	if job.GetState() != process.JobStateCompleted {
		t.Errorf("expected completed state, got %s", job.GetState())
	}
}

func TestJobManagerListJobs(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	_, err := manager.StartJob(ctx, "sleep 0.1", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	jobs := manager.ListJobs()
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
}

func TestJobManagerGetJob(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	job, err := manager.StartJob(ctx, "sleep 0.1", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	found, exists := manager.GetJob(job.ID)
	if !exists {
		t.Fatal("expected job to exist")
	}

	if found.ID != job.ID {
		t.Errorf("expected job ID %d, got %d", job.ID, found.ID)
	}

	_, exists = manager.GetJob(9999)
	if exists {
		t.Error("expected job 9999 to not exist")
	}
}

func TestJobManagerKillJob(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	job, err := manager.StartJob(ctx, "sleep 10", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	// Give the process time to start
	time.Sleep(50 * time.Millisecond)

	err = manager.KillJob(job.ID)
	if err != nil {
		t.Fatalf("failed to kill job: %v", err)
	}

	// Wait for the job to be killed
	time.Sleep(100 * time.Millisecond)

	state := job.GetState()
	if state != process.JobStateCompleted && state != process.JobStateFailed && state != process.JobStateStopped {
		t.Errorf("expected non-running state after kill, got %s", state)
	}
}

func TestJobManagerMaxJobs(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(2)
	ctx := context.Background()

	_, err1 := manager.StartJob(ctx, "sleep 10", "/bin/bash")
	if err1 != nil {
		t.Fatalf("failed to start job 1: %v", err1)
	}

	_, err2 := manager.StartJob(ctx, "sleep 10", "/bin/bash")
	if err2 != nil {
		t.Fatalf("failed to start job 2: %v", err2)
	}

	_, err3 := manager.StartJob(ctx, "sleep 10", "/bin/bash")
	if err3 == nil {
		t.Error("expected error when exceeding max jobs")
	}
}

func TestJobDuration(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	job, err := manager.StartJob(ctx, "sleep 0.1", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	// Check duration while running
	time.Sleep(50 * time.Millisecond)
	duration := job.Duration()
	if duration < 50*time.Millisecond {
		t.Errorf("expected duration >= 50ms, got %v", duration)
	}

	// Wait for completion
	time.Sleep(100 * time.Millisecond)

	finalDuration := job.Duration()
	if finalDuration < 100*time.Millisecond {
		t.Errorf("expected final duration >= 100ms, got %v", finalDuration)
	}
}

func TestJobString(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	job, err := manager.StartJob(ctx, "echo test", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	str := job.String()
	if str == "" {
		t.Error("expected non-empty string representation")
	}

	// Clean up
	time.Sleep(100 * time.Millisecond)
}

func TestListActiveJobs(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	// Start a long-running job
	_, err := manager.StartJob(ctx, "sleep 1", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	// Start a quick job
	_, err = manager.StartJob(ctx, "echo done", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start quick job: %v", err)
	}

	// Wait for quick job to complete
	time.Sleep(100 * time.Millisecond)

	activeJobs := manager.ListActiveJobs()
	// At least the long-running job should still be active
	if len(activeJobs) < 1 {
		t.Errorf("expected at least 1 active job, got %d", len(activeJobs))
	}
}

func TestCleanupCompleted(t *testing.T) {
	t.Parallel()

	manager := process.NewJobManager(10)
	ctx := context.Background()

	_, err := manager.StartJob(ctx, "echo done", "/bin/bash")
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	// Wait for job to complete
	time.Sleep(100 * time.Millisecond)

	// Cleanup with 0 max age should remove completed jobs
	removed := manager.CleanupCompleted(0)
	if removed != 1 {
		t.Errorf("expected 1 job removed, got %d", removed)
	}

	jobs := manager.ListJobs()
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs after cleanup, got %d", len(jobs))
	}
}
