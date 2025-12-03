// Package logging provides logging infrastructure for cbwsh.
package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Level represents the logging level.
type Level int

const (
	// LevelDebug is for debug messages.
	LevelDebug Level = iota
	// LevelInfo is for informational messages.
	LevelInfo
	// LevelWarn is for warning messages.
	LevelWarn
	// LevelError is for error messages.
	LevelError
	// LevelFatal is for fatal messages.
	LevelFatal
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging capabilities.
type Logger struct {
	mu        sync.RWMutex
	level     Level
	output    io.Writer
	file      *os.File
	formatter Formatter
	fields    map[string]any
}

// Formatter formats log entries.
type Formatter interface {
	Format(entry *Entry) string
}

// Entry represents a log entry.
type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Fields    map[string]any
	Caller    string
}

// TextFormatter formats log entries as text.
type TextFormatter struct {
	TimestampFormat string
	ShowCaller      bool
}

// Format formats a log entry as text.
func (f *TextFormatter) Format(entry *Entry) string {
	timestamp := entry.Timestamp.Format(f.TimestampFormat)

	var caller string
	if f.ShowCaller && entry.Caller != "" {
		caller = fmt.Sprintf(" [%s]", entry.Caller)
	}

	fieldsStr := ""
	if len(entry.Fields) > 0 {
		fieldsStr = " {"
		first := true
		for k, v := range entry.Fields {
			if !first {
				fieldsStr += ", "
			}
			fieldsStr += fmt.Sprintf("%s=%v", k, v)
			first = false
		}
		fieldsStr += "}"
	}

	return fmt.Sprintf("%s %s%s %s%s\n", timestamp, entry.Level.String(), caller, entry.Message, fieldsStr)
}

// JSONFormatter formats log entries as JSON.
type JSONFormatter struct {
	TimestampFormat string
}

// Format formats a log entry as JSON.
func (f *JSONFormatter) Format(entry *Entry) string {
	// Simple JSON formatting without external dependencies
	timestamp := entry.Timestamp.Format(f.TimestampFormat)

	fields := ""
	if len(entry.Fields) > 0 {
		for k, v := range entry.Fields {
			fields += fmt.Sprintf(", %q: %q", k, fmt.Sprintf("%v", v))
		}
	}

	return fmt.Sprintf("{\"timestamp\": %q, \"level\": %q, \"message\": %q%s}\n",
		timestamp, entry.Level.String(), entry.Message, fields)
}

// Option configures a Logger.
type Option func(*Logger)

// WithLevel sets the log level.
func WithLevel(level Level) Option {
	return func(l *Logger) {
		l.level = level
	}
}

// WithOutput sets the output writer.
func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.output = w
	}
}

// WithFormatter sets the formatter.
func WithFormatter(f Formatter) Option {
	return func(l *Logger) {
		l.formatter = f
	}
}

// WithFields sets initial fields.
func WithFields(fields map[string]any) Option {
	return func(l *Logger) {
		for k, v := range fields {
			l.fields[k] = v
		}
	}
}

// New creates a new Logger with the given options.
func New(opts ...Option) *Logger {
	l := &Logger{
		level:  LevelInfo,
		output: os.Stderr,
		formatter: &TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			ShowCaller:      true,
		},
		fields: make(map[string]any),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// SetLevel sets the log level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current log level.
func (l *Logger) GetLevel() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// SetOutput sets the output writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetOutputFile sets output to a file.
func (l *Logger) SetOutputFile(path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open or create log file
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Close previous file if any
	if l.file != nil {
		l.file.Close()
	}

	l.file = file
	l.output = file
	return nil
}

// Close closes the logger and any open files.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		return err
	}
	return nil
}

// WithField returns a new logger with an additional field.
func (l *Logger) WithField(key string, value any) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	fields := make(map[string]any, len(l.fields)+1)
	for k, v := range l.fields {
		fields[k] = v
	}
	fields[key] = value

	return &Logger{
		level:     l.level,
		output:    l.output,
		file:      l.file,
		formatter: l.formatter,
		fields:    fields,
	}
}

// WithFields returns a new logger with additional fields.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	newFields := make(map[string]any, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &Logger{
		level:     l.level,
		output:    l.output,
		file:      l.file,
		formatter: l.formatter,
		fields:    newFields,
	}
}

func (l *Logger) log(level Level, msg string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if level < l.level {
		return
	}

	entry := &Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    l.fields,
		Caller:    getCaller(3),
	}

	output := l.formatter.Format(entry)
	if _, err := l.output.Write([]byte(output)); err != nil {
		// Fallback to stderr if primary output fails
		fmt.Fprintf(os.Stderr, "log write error: %v, message: %s", err, output)
	}
}

func (l *Logger) logf(level Level, format string, args ...any) {
	l.log(level, fmt.Sprintf(format, args...))
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string) {
	l.log(LevelDebug, msg)
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(format string, args ...any) {
	l.logf(LevelDebug, format, args...)
}

// Info logs an info message.
func (l *Logger) Info(msg string) {
	l.log(LevelInfo, msg)
}

// Infof logs a formatted info message.
func (l *Logger) Infof(format string, args ...any) {
	l.logf(LevelInfo, format, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.log(LevelWarn, msg)
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(format string, args ...any) {
	l.logf(LevelWarn, format, args...)
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.log(LevelError, msg)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(format string, args ...any) {
	l.logf(LevelError, format, args...)
}

// Fatal logs a fatal message.
func (l *Logger) Fatal(msg string) {
	l.log(LevelFatal, msg)
}

// Fatalf logs a formatted fatal message.
func (l *Logger) Fatalf(format string, args ...any) {
	l.logf(LevelFatal, format, args...)
}

func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Default logger instance.
var defaultLogger = New()

// SetDefaultLevel sets the default logger level.
func SetDefaultLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetDefaultOutput sets the default logger output.
func SetDefaultOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// Debug logs a debug message using the default logger.
func Debug(msg string) {
	defaultLogger.Debug(msg)
}

// Debugf logs a formatted debug message using the default logger.
func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// Info logs an info message using the default logger.
func Info(msg string) {
	defaultLogger.Info(msg)
}

// Infof logs a formatted info message using the default logger.
func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

// Warn logs a warning message using the default logger.
func Warn(msg string) {
	defaultLogger.Warn(msg)
}

// Warnf logs a formatted warning message using the default logger.
func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

// Error logs an error message using the default logger.
func Error(msg string) {
	defaultLogger.Error(msg)
}

// Errorf logs a formatted error message using the default logger.
func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

// Fatal logs a fatal message using the default logger.
func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

// Fatalf logs a formatted fatal message using the default logger.
func Fatalf(format string, args ...any) {
	defaultLogger.Fatalf(format, args...)
}

// WithField returns a new logger with an additional field using the default logger.
func Field(key string, value any) *Logger {
	return defaultLogger.WithField(key, value)
}

// Fields returns a new logger with additional fields using the default logger.
func Fields(fields map[string]any) *Logger {
	return defaultLogger.WithFields(fields)
}
