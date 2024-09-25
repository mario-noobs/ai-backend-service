// slog_logger.go
package logger

import (
	"golang.org/x/exp/slog"
	"os"
)

// SlogLogger is an implementation of Logger using Slog
type SlogLogger struct {
	logger *slog.Logger
	level  slog.Level
}

// NewSlogLogger creates a new SlogLogger with a specified format and level
func NewSlogLogger(format string, level string) *SlogLogger {
	var handler slog.Handler

	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, nil) // Use nil for default options
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil) // Use nil for default options
	}

	// Set the log level
	lvl := parseLogLevel(level)
	logger := slog.New(handler)

	return &SlogLogger{logger: logger, level: lvl}
}

// parseLogLevel converts a string to slog.Level
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to Info if invalid level is provided
	}
}

// SetLogLevel sets the log level for Slog (not commonly used, but kept for interface)
func (l *SlogLogger) SetLogLevel(level string) error {
	l.level = parseLogLevel(level)
	return nil
}

// Debug logs a debug message with fields
func (l *SlogLogger) Debug(msg string, fields map[string]interface{}) {
	if l.level <= slog.LevelDebug {
		l.logger.Debug(msg, fields)
	}
}

// Info logs an info message with fields
func (l *SlogLogger) Info(msg string, fields map[string]interface{}) {
	if l.level <= slog.LevelInfo {
		l.logger.Info(msg, fields)
	}
}

// Warn logs a warning message with fields
func (l *SlogLogger) Warn(msg string, fields map[string]interface{}) {
	if l.level <= slog.LevelWarn {
		l.logger.Warn(msg, fields)
	}
}

// Error logs an error message with fields
func (l *SlogLogger) Error(msg string, fields map[string]interface{}) {
	if l.level <= slog.LevelError {
		l.logger.Error(msg, fields)
	}
}

// Info logs an info message with key-value pairs
func (l *SlogLogger) InfoArgs(msg string, args ...any) {
	if l.level <= slog.LevelInfo {
		fields := createFields(args)
		l.logger.Info(msg, fields)
	}
}

// Debug logs a debug message with key-value pairs
func (l *SlogLogger) DebugArgs(msg string, args ...any) {
	if l.level <= slog.LevelDebug {
		fields := createFields(args)
		l.logger.Debug(msg, fields)
	}
}
