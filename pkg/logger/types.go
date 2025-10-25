package logger

import "log/slog"

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Logger *slog.Logger
}

// LogLevel represents different verbosity levels
type LogLevel int

const (
	LevelOff      LogLevel = iota // No logging output
	LevelBasic                    // Essential process information
	LevelDetailed                 // API details and processing summaries
	LevelVerbose                  // All details including streaming chunks
)

// VerbosityLevels maps integer verbosity to LogLevel
var VerbosityLevels = map[int]LogLevel{
	0: LevelOff,
	1: LevelBasic,
	2: LevelDetailed,
	3: LevelVerbose,
}
