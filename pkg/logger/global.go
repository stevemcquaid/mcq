package logger

// Global logger instance for backward compatibility
var globalLogger = NewLogger()

// SetupLogger configures the global logger based on verbosity level
func SetupLogger(verbosityLevel int) {
	globalLogger.Setup(verbosityLevel)
}

// LogBasic logs basic process information using the global logger
func LogBasic(msg string, args ...interface{}) {
	globalLogger.Basic(msg, args...)
}

// LogDetailed logs detailed information using the global logger
func LogDetailed(msg string, args ...interface{}) {
	globalLogger.Detailed(msg, args...)
}

// LogVerbose logs verbose information using the global logger
func LogVerbose(msg string, args ...interface{}) {
	globalLogger.Verbose(msg, args...)
}

// LogError logs error information using the global logger
func LogError(operation string, err error) {
	globalLogger.Error(operation, err)
}

// LogInfo logs info level messages using the global logger
func LogInfo(msg string, args ...interface{}) {
	globalLogger.Info(msg, args...)
}

// LogDebug logs debug level messages using the global logger
func LogDebug(msg string, args ...interface{}) {
	globalLogger.Debug(msg, args...)
}

// LogWarn logs warning level messages using the global logger
func LogWarn(msg string, args ...interface{}) {
	globalLogger.Warn(msg, args...)
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return globalLogger
}
