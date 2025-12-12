// Package monitor provides shell activity monitoring with AI-powered recommendations.
package monitor

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/ai/ollama"
	"github.com/cbwinslow/cbwsh/pkg/core"
)

// ActivityType represents the type of shell activity.
type ActivityType int

const (
	// ActivityCommand represents a command execution.
	ActivityCommand ActivityType = iota
	// ActivityError represents a command error.
	ActivityError
	// ActivityWorkingDirChange represents a directory change.
	ActivityWorkingDirChange
	// ActivityFileChange represents a file modification.
	ActivityFileChange
)

// String returns the string representation of the activity type.
func (a ActivityType) String() string {
	switch a {
	case ActivityCommand:
		return "command"
	case ActivityError:
		return "error"
	case ActivityWorkingDirChange:
		return "dir_change"
	case ActivityFileChange:
		return "file_change"
	default:
		return "unknown"
	}
}

// Activity represents a shell activity event.
type Activity struct {
	Type      ActivityType
	Timestamp time.Time
	Command   string
	Output    string
	Error     string
	ExitCode  int
	WorkDir   string
	Context   map[string]string
}

// Recommendation represents an AI-generated recommendation.
type Recommendation struct {
	Type      string    // "suggestion", "warning", "info", "tip"
	Title     string
	Message   string
	Timestamp time.Time
	Activity  *Activity
}

// Monitor monitors shell activity and provides AI-powered recommendations.
type Monitor struct {
	mu sync.RWMutex

	// Ollama client
	ollama *ollama.Client

	// Activity tracking
	activities      []Activity
	maxActivities   int
	recommendations []Recommendation
	maxRecommend    int

	// Settings
	enabled           bool
	autoRecommend     bool
	recommendInterval time.Duration
	minActivityGap    time.Duration
	lastActivity      time.Time

	// Callbacks
	onRecommendation func(Recommendation)

	// Context for async operations
	ctx    context.Context
	cancel context.CancelFunc
}

// Config holds monitor configuration.
type Config struct {
	OllamaURL         string
	OllamaModel       string
	MaxActivities     int
	MaxRecommendations int
	AutoRecommend     bool
	RecommendInterval time.Duration
	MinActivityGap    time.Duration
}

// DefaultConfig returns default monitor configuration.
func DefaultConfig() *Config {
	return &Config{
		OllamaURL:          "http://localhost:11434",
		OllamaModel:        "llama2",
		MaxActivities:      100,
		MaxRecommendations: 50,
		AutoRecommend:      true,
		RecommendInterval:  5 * time.Second,
		MinActivityGap:     1 * time.Second,
	}
}

// NewMonitor creates a new activity monitor.
func NewMonitor(cfg *Config) *Monitor {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Monitor{
		ollama:            ollama.NewClient(cfg.OllamaURL, cfg.OllamaModel),
		activities:        make([]Activity, 0, cfg.MaxActivities),
		maxActivities:     cfg.MaxActivities,
		recommendations:   make([]Recommendation, 0, cfg.MaxRecommendations),
		maxRecommend:      cfg.MaxRecommendations,
		enabled:           false,
		autoRecommend:     cfg.AutoRecommend,
		recommendInterval: cfg.RecommendInterval,
		minActivityGap:    cfg.MinActivityGap,
		ctx:               ctx,
		cancel:            cancel,
	}

	return m
}

// Start starts the activity monitor.
func (m *Monitor) Start() {
	m.mu.Lock()
	if m.enabled {
		m.mu.Unlock()
		return
	}
	m.enabled = true
	m.mu.Unlock()

	// Start background recommendation generator if auto-recommend is enabled
	if m.autoRecommend {
		go m.backgroundRecommender()
	}
}

// Stop stops the activity monitor.
func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = false
	if m.cancel != nil {
		m.cancel()
	}
}

// IsEnabled returns whether the monitor is enabled.
func (m *Monitor) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// RecordActivity records a shell activity.
func (m *Monitor) RecordActivity(activity Activity) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.enabled {
		return
	}

	// Check if we should record this activity (throttling)
	now := time.Now()
	if now.Sub(m.lastActivity) < m.minActivityGap {
		return
	}
	m.lastActivity = now

	// Set timestamp if not provided
	if activity.Timestamp.IsZero() {
		activity.Timestamp = now
	}

	// Add to activities
	m.activities = append(m.activities, activity)

	// Trim if exceeds max
	if len(m.activities) > m.maxActivities {
		m.activities = m.activities[len(m.activities)-m.maxActivities:]
	}
}

// RecordCommand is a convenience method to record command execution.
func (m *Monitor) RecordCommand(result *core.CommandResult, workDir string) {
	activity := Activity{
		Type:      ActivityCommand,
		Timestamp: time.Now(),
		Command:   result.Command,
		Output:    result.Output,
		Error:     result.Error,
		ExitCode:  result.ExitCode,
		WorkDir:   workDir,
		Context:   make(map[string]string),
	}

	if result.ExitCode != 0 {
		activity.Type = ActivityError
	}

	m.RecordActivity(activity)
}

// GetRecentActivities returns recent activities.
func (m *Monitor) GetRecentActivities(count int) []Activity {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if count <= 0 || count > len(m.activities) {
		count = len(m.activities)
	}

	start := len(m.activities) - count
	result := make([]Activity, count)
	copy(result, m.activities[start:])
	return result
}

// GetRecommendations returns all recommendations.
func (m *Monitor) GetRecommendations() []Recommendation {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Recommendation, len(m.recommendations))
	copy(result, m.recommendations)
	return result
}

// ClearRecommendations clears all recommendations.
func (m *Monitor) ClearRecommendations() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recommendations = m.recommendations[:0]
}

// SetOnRecommendation sets the callback for new recommendations.
func (m *Monitor) SetOnRecommendation(callback func(Recommendation)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onRecommendation = callback
}

// addRecommendation adds a recommendation and calls the callback.
func (m *Monitor) addRecommendation(rec Recommendation) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Set timestamp if not provided
	if rec.Timestamp.IsZero() {
		rec.Timestamp = time.Now()
	}

	m.recommendations = append(m.recommendations, rec)

	// Trim if exceeds max
	if len(m.recommendations) > m.maxRecommend {
		m.recommendations = m.recommendations[len(m.recommendations)-m.maxRecommend:]
	}

	// Call callback if set
	if m.onRecommendation != nil {
		go m.onRecommendation(rec)
	}
}

// GenerateRecommendation generates a recommendation based on recent activity.
func (m *Monitor) GenerateRecommendation(ctx context.Context) error {
	// Get recent activities
	activities := m.GetRecentActivities(5)
	if len(activities) == 0 {
		return nil
	}

	// Build context for AI
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Recent shell activity:\n")
	for i, act := range activities {
		contextBuilder.WriteString(fmt.Sprintf("%d. [%s] %s",
			i+1, act.Type.String(), act.Command))
		if act.ExitCode != 0 {
			contextBuilder.WriteString(fmt.Sprintf(" (exit code: %d)", act.ExitCode))
		}
		if act.Error != "" {
			contextBuilder.WriteString(fmt.Sprintf("\n   Error: %s", truncate(act.Error, 100)))
		}
		contextBuilder.WriteString("\n")
	}

	// Add current working directory context
	if len(activities) > 0 {
		contextBuilder.WriteString(fmt.Sprintf("\nWorking directory: %s\n", activities[len(activities)-1].WorkDir))
	}

	prompt := fmt.Sprintf(`%s
Analyze the above shell activity and provide a brief helpful comment, tip, or suggestion. 
Focus on:
- Command improvements or alternatives
- Error explanations and fixes
- Best practices
- Workflow optimization

Keep response under 100 words. Be concise and actionable.`, contextBuilder.String())

	// Query Ollama
	response, err := m.ollama.Generate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate recommendation: %w", err)
	}

	// Create recommendation
	rec := Recommendation{
		Type:      determineRecommendationType(activities),
		Title:     "AI Suggestion",
		Message:   response,
		Timestamp: time.Now(),
	}

	// Add the most recent activity as context
	if len(activities) > 0 {
		lastActivity := activities[len(activities)-1]
		rec.Activity = &lastActivity
	}

	m.addRecommendation(rec)
	return nil
}

// backgroundRecommender periodically generates recommendations.
func (m *Monitor) backgroundRecommender() {
	ticker := time.NewTicker(m.recommendInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.mu.RLock()
			enabled := m.enabled
			activitiesCount := len(m.activities)
			m.mu.RUnlock()

			if !enabled || activitiesCount == 0 {
				continue
			}

			// Generate recommendation with timeout
			ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
			_ = m.GenerateRecommendation(ctx)
			cancel()
		}
	}
}

// determineRecommendationType determines the recommendation type based on activities.
func determineRecommendationType(activities []Activity) string {
	// Check if there are recent errors
	for i := len(activities) - 1; i >= 0 && i >= len(activities)-3; i-- {
		if activities[i].Type == ActivityError {
			return "warning"
		}
	}

	// Check for repeated similar commands (might suggest alias or script)
	if len(activities) >= 3 {
		lastCmd := activities[len(activities)-1].Command
		similarCount := 0
		for i := len(activities) - 2; i >= 0 && i >= len(activities)-5; i-- {
			if similarity(activities[i].Command, lastCmd) > 0.7 {
				similarCount++
			}
		}
		if similarCount >= 2 {
			return "tip"
		}
	}

	return "info"
}

// similarity calculates rough similarity between two strings.
func similarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	// Simple similarity based on common words
	wordsA := strings.Fields(a)
	wordsB := strings.Fields(b)
	if len(wordsA) == 0 || len(wordsB) == 0 {
		return 0.0
	}

	commonCount := 0
	for _, wa := range wordsA {
		for _, wb := range wordsB {
			if wa == wb {
				commonCount++
				break
			}
		}
	}

	return float64(commonCount) / float64(max(len(wordsA), len(wordsB)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
