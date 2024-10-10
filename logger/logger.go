package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Logger is an interface for logging
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string)
	Error(msg string)
	DebugArgs(msg string, args ...any)
	InfoArgs(msg string, args ...any)
	SetLogLevel(level string) error
}

// LoggerFactory is responsible for creating loggers
type LoggerFactory struct{}

// NewLogger creates a new logger based on the specified type, format, and level
func (f *LoggerFactory) NewLogger(loggerType string, format string, level string) (Logger, error) {
	if loggerType == "slog" {
		return NewSlogLogger(format, level), nil
	} else if loggerType == "logrus" {
		var logger Logger
		if format == "json" {
			logger = NewLogrusLogger(&log.JSONFormatter{})
		} else {
			logger = NewLogrusLogger(&log.TextFormatter{})
		}

		if err := logger.SetLogLevel(level); err != nil {
			return nil, err
		}
		return logger, nil
	}
	return nil, fmt.Errorf("unknown logger type: %s", loggerType)
}

func createFields(args []any) map[string]any {
	fields := make(map[string]any)
	if len(args)%2 != 0 {
		// Handle odd number of arguments (optional)
		return fields // Or log a warning
	}
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if ok {
			fields[key] = args[i+1]
		}
	}
	return fields
}
