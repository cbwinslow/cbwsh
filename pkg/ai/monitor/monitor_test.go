package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/core"
)

func TestNewMonitor(t *testing.T) {
	cfg := DefaultConfig()
	monitor := NewMonitor(cfg)

	if monitor == nil {
		t.Fatal("NewMonitor returned nil")
	}

	if monitor.ollama == nil {
		t.Error("ollama client should not be nil")
	}

	if monitor.enabled {
		t.Error("monitor should start disabled")
	}

	if len(monitor.activities) != 0 {
		t.Errorf("expected 0 activities, got %d", len(monitor.activities))
	}

	if len(monitor.recommendations) != 0 {
		t.Errorf("expected 0 recommendations, got %d", len(monitor.recommendations))
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OllamaURL != "http://localhost:11434" {
		t.Errorf("expected default URL 'http://localhost:11434', got %q", cfg.OllamaURL)
	}

	if cfg.OllamaModel != "llama2" {
		t.Errorf("expected default model 'llama2', got %q", cfg.OllamaModel)
	}

	if cfg.MaxActivities != 100 {
		t.Errorf("expected max activities 100, got %d", cfg.MaxActivities)
	}

	if cfg.AutoRecommend != true {
		t.Error("expected auto recommend to be true")
	}
}

func TestMonitor_StartStop(t *testing.T) {
	monitor := NewMonitor(nil)

	if monitor.IsEnabled() {
		t.Error("monitor should start disabled")
	}

	monitor.Start()

	if !monitor.IsEnabled() {
		t.Error("monitor should be enabled after Start")
	}

	monitor.Stop()

	if monitor.IsEnabled() {
		t.Error("monitor should be disabled after Stop")
	}
}

func TestMonitor_RecordActivity(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoRecommend = false // Disable auto-recommend for testing
	monitor := NewMonitor(cfg)
	monitor.Start()

	activity := Activity{
		Type:      ActivityCommand,
		Command:   "ls -la",
		Output:    "files...",
		ExitCode:  0,
		WorkDir:   "/home/test",
		Timestamp: time.Now(),
	}

	monitor.RecordActivity(activity)

	activities := monitor.GetRecentActivities(10)
	if len(activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(activities))
	}

	if activities[0].Command != "ls -la" {
		t.Errorf("expected command 'ls -la', got %q", activities[0].Command)
	}
}

func TestMonitor_RecordCommand(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoRecommend = false
	monitor := NewMonitor(cfg)
	monitor.Start()

	result := &core.CommandResult{
		Command:  "echo test",
		Output:   "test",
		ExitCode: 0,
	}

	monitor.RecordCommand(result, "/home/test")

	activities := monitor.GetRecentActivities(10)
	if len(activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(activities))
	}

	if activities[0].Type != ActivityCommand {
		t.Errorf("expected type ActivityCommand, got %v", activities[0].Type)
	}
}

func TestMonitor_RecordCommandWithError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoRecommend = false
	monitor := NewMonitor(cfg)
	monitor.Start()

	result := &core.CommandResult{
		Command:  "invalid-command",
		Error:    "command not found",
		ExitCode: 127,
	}

	monitor.RecordCommand(result, "/home/test")

	activities := monitor.GetRecentActivities(10)
	if len(activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(activities))
	}

	if activities[0].Type != ActivityError {
		t.Errorf("expected type ActivityError, got %v", activities[0].Type)
	}
}

func TestMonitor_MaxActivities(t *testing.T) {
	cfg := &Config{
		OllamaURL:          "http://localhost:11434",
		OllamaModel:        "llama2",
		MaxActivities:      5,
		MaxRecommendations: 50,
		AutoRecommend:      false,
		RecommendInterval:  5 * time.Second,
		MinActivityGap:     0, // No throttling for test
	}
	monitor := NewMonitor(cfg)
	monitor.Start()

	// Add more than max activities
	for i := 0; i < 10; i++ {
		activity := Activity{
			Type:      ActivityCommand,
			Command:   "test",
			Timestamp: time.Now(),
		}
		monitor.RecordActivity(activity)
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}

	activities := monitor.GetRecentActivities(100)
	if len(activities) > 5 {
		t.Errorf("expected max 5 activities, got %d", len(activities))
	}
}

func TestMonitor_ClearRecommendations(t *testing.T) {
	monitor := NewMonitor(nil)

	// Add a recommendation directly
	monitor.addRecommendation(Recommendation{
		Type:    "info",
		Title:   "Test",
		Message: "Test message",
	})

	recs := monitor.GetRecommendations()
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}

	monitor.ClearRecommendations()

	recs = monitor.GetRecommendations()
	if len(recs) != 0 {
		t.Errorf("expected 0 recommendations after clear, got %d", len(recs))
	}
}

func TestMonitor_OnRecommendationCallback(t *testing.T) {
	monitor := NewMonitor(nil)

	called := false
	var receivedRec Recommendation

	monitor.SetOnRecommendation(func(rec Recommendation) {
		called = true
		receivedRec = rec
	})

	testRec := Recommendation{
		Type:    "info",
		Title:   "Test",
		Message: "Test message",
	}

	monitor.addRecommendation(testRec)

	// Give callback time to execute
	time.Sleep(10 * time.Millisecond)

	if !called {
		t.Error("expected callback to be called")
	}

	if receivedRec.Title != "Test" {
		t.Errorf("expected title 'Test', got %q", receivedRec.Title)
	}
}

func TestActivityType_String(t *testing.T) {
	tests := []struct {
		activityType ActivityType
		expected     string
	}{
		{ActivityCommand, "command"},
		{ActivityError, "error"},
		{ActivityWorkingDirChange, "dir_change"},
		{ActivityFileChange, "file_change"},
		{ActivityType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.activityType.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetermineRecommendationType(t *testing.T) {
	// Test with errors
	activities := []Activity{
		{Type: ActivityCommand, Command: "ls"},
		{Type: ActivityError, Command: "invalid"},
	}

	recType := determineRecommendationType(activities)
	if recType != "warning" {
		t.Errorf("expected 'warning' for error activity, got %q", recType)
	}

	// Test with normal commands
	activities = []Activity{
		{Type: ActivityCommand, Command: "ls"},
		{Type: ActivityCommand, Command: "pwd"},
	}

	recType = determineRecommendationType(activities)
	if recType != "info" {
		t.Errorf("expected 'info' for normal commands, got %q", recType)
	}
}

func TestSimilarity(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected float64
	}{
		{"ls -la", "ls -la", 1.0},
		{"ls -la", "ls -l", 0.5}, // 1 common word out of 2
		{"", "", 1.0},            // Empty strings are equal
		{"ls", "pwd", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			result := similarity(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected similarity %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenerateRecommendation_NoActivities(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AutoRecommend = false
	monitor := NewMonitor(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Should not error with no activities
	err := monitor.GenerateRecommendation(ctx)
	if err != nil {
		t.Errorf("expected no error with no activities, got %v", err)
	}
}
