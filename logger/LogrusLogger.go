package logger

import (
	log "github.com/sirupsen/logrus"
)

// LogrusLogger is an implementation of Logger using Logrus
type LogrusLogger struct {
	logger *log.Logger
}

// NewLogrusLogger creates a new LogrusLogger with a specified formatter
func NewLogrusLogger(formatter log.Formatter) *LogrusLogger {
	logger := log.New()
	logger.SetFormatter(formatter)
	return &LogrusLogger{logger: logger}
}

// SetLogLevel sets the log level for Logrus
func (l *LogrusLogger) SetLogLevel(level string) error {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	l.logger.SetLevel(lvl)
	return nil
}

// Info logs an info message with fields
func (l *LogrusLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.WithFields(log.Fields(fields)).Info(msg)
}

// Warn logs a warning message with fields
func (l *LogrusLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.WithFields(log.Fields(fields)).Warn(msg)
}

// Error logs an error message with fields
func (l *LogrusLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.WithFields(log.Fields(fields)).Error(msg)
}

// Debug logs a debug message with fields
func (l *LogrusLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.WithFields(log.Fields(fields)).Debug(msg)
}

// Info logs an info message with key-value pairs
func (l *LogrusLogger) InfoArgs(msg string, args ...any) {
	fields := createFields(args)
	l.logger.WithFields(log.Fields(fields)).Info(msg)
}

// Debug logs a debug message with key-value pairs
func (l *LogrusLogger) DebugArgs(msg string, args ...any) {
	fields := createFields(args)
	l.logger.WithFields(log.Fields(fields)).Debug(msg)
}
