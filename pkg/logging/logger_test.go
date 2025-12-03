package logging_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/logging"
)

func TestLevelString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level    logging.Level
		expected string
	}{
		{logging.LevelDebug, "DEBUG"},
		{logging.LevelInfo, "INFO"},
		{logging.LevelWarn, "WARN"},
		{logging.LevelError, "ERROR"},
		{logging.LevelFatal, "FATAL"},
		{logging.Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	logger := logging.New()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if logger.GetLevel() != logging.LevelInfo {
		t.Errorf("expected default level to be INFO")
	}
}

func TestLoggerWithLevel(t *testing.T) {
	t.Parallel()

	logger := logging.New(logging.WithLevel(logging.LevelDebug))
	if logger.GetLevel() != logging.LevelDebug {
		t.Errorf("expected DEBUG level")
	}
}

func TestLoggerWithOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := logging.New(logging.WithOutput(&buf))

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected output to contain INFO, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := logging.New(
		logging.WithOutput(&buf),
		logging.WithLevel(logging.LevelWarn),
	)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("debug message should be filtered")
	}
	if strings.Contains(output, "info message") {
		t.Error("info message should be filtered")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("warn message should not be filtered")
	}
	if !strings.Contains(output, "error message") {
		t.Error("error message should not be filtered")
	}
}

func TestLoggerWithField(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := logging.New(logging.WithOutput(&buf))

	fieldLogger := logger.WithField("key", "value")
	fieldLogger.Info("message with field")

	output := buf.String()
	if !strings.Contains(output, "key=value") {
		t.Errorf("expected field in output, got: %s", output)
	}
}

func TestLoggerWithFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := logging.New(logging.WithOutput(&buf))

	fieldLogger := logger.WithFields(map[string]any{
		"key1": "value1",
		"key2": "value2",
	})
	fieldLogger.Info("message with fields")

	output := buf.String()
	if !strings.Contains(output, "key1=value1") {
		t.Errorf("expected key1 in output, got: %s", output)
	}
	if !strings.Contains(output, "key2=value2") {
		t.Errorf("expected key2 in output, got: %s", output)
	}
}

func TestTextFormatter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	formatter := &logging.TextFormatter{
		TimestampFormat: "2006-01-02",
		ShowCaller:      false,
	}
	logger := logging.New(
		logging.WithOutput(&buf),
		logging.WithFormatter(formatter),
	)

	logger.Info("formatted message")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected INFO in output, got: %s", output)
	}
}

func TestJSONFormatter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	formatter := &logging.JSONFormatter{
		TimestampFormat: "2006-01-02",
	}
	logger := logging.New(
		logging.WithOutput(&buf),
		logging.WithFormatter(formatter),
	)

	logger.Info("json message")

	output := buf.String()
	if !strings.Contains(output, `"level": "INFO"`) {
		t.Errorf("expected JSON format in output, got: %s", output)
	}
	if !strings.Contains(output, `"message": "json message"`) {
		t.Errorf("expected message in JSON output, got: %s", output)
	}
}

func TestSetLevel(t *testing.T) {
	t.Parallel()

	logger := logging.New()
	logger.SetLevel(logging.LevelError)

	if logger.GetLevel() != logging.LevelError {
		t.Errorf("expected ERROR level after SetLevel")
	}
}

func TestFormattedLogging(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := logging.New(logging.WithOutput(&buf))

	logger.Infof("formatted %s %d", "test", 123)

	output := buf.String()
	if !strings.Contains(output, "formatted test 123") {
		t.Errorf("expected formatted message, got: %s", output)
	}
}
