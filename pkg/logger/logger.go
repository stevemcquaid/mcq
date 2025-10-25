package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger provides structured logging functionality
type Logger struct {
	config *LoggerConfig
}

// NewLogger creates a new Logger instance
func NewLogger() *Logger {
	return &Logger{
		config: &LoggerConfig{},
	}
}

// Setup configures the logger based on verbosity level
func (l *Logger) Setup(verbosityLevel int) {
	level, exists := VerbosityLevels[verbosityLevel]
	if !exists {
		level = LevelOff
	}

	slogLevel := l.convertToSlogLevel(level)

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // Remove timestamp
			}
			return a
		},
	})

	l.config.Logger = slog.New(handler)
}

// convertToSlogLevel converts our LogLevel to slog.Level
func (l *Logger) convertToSlogLevel(level LogLevel) slog.Level {
	switch level {
	case LevelOff:
		return slog.LevelError - 1
	case LevelBasic:
		return slog.LevelInfo
	case LevelDetailed:
		return slog.LevelDebug
	case LevelVerbose:
		return slog.LevelDebug - 1
	default:
		return slog.LevelError - 1
	}
}

// Basic logs basic process information
func (l *Logger) Basic(msg string, args ...interface{}) {
	if l.config.Logger != nil {
		l.config.Logger.Info(msg, args...)
	}
}

// Detailed logs detailed information
func (l *Logger) Detailed(msg string, args ...interface{}) {
	if l.config.Logger != nil {
		l.config.Logger.Debug(msg, args...)
	}
}

// Verbose logs verbose information
func (l *Logger) Verbose(msg string, args ...interface{}) {
	if l.config.Logger != nil {
		l.config.Logger.Log(context.Background(), l.convertToSlogLevel(LevelVerbose), msg, args...)
	}
}

// Error logs error information
func (l *Logger) Error(operation string, err error) {
	if l.config.Logger != nil {
		l.config.Logger.Error("Operation failed", "operation", operation, "error", err)
	}
}

// Info logs info level messages
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.config.Logger != nil {
		l.config.Logger.Info(msg, args...)
	}
}

// Debug logs debug level messages
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.config.Logger != nil {
		l.config.Logger.Debug(msg, args...)
	}
}

// Warn logs warning level messages
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.config.Logger != nil {
		l.config.Logger.Warn(msg, args...)
	}
}

// GetConfig returns the logger configuration
func (l *Logger) GetConfig() *LoggerConfig {
	return l.config
}

// SetConfig sets the logger configuration
func (l *Logger) SetConfig(config *LoggerConfig) {
	l.config = config
}
