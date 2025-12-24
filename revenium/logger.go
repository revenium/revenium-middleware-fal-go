package revenium

import (
	"log"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var (
	currentLogLevel = LogLevelInfo
	logger          = log.New(os.Stdout, "[Revenium] ", log.LstdFlags)
)

// InitializeLogger initializes the logger with the configured log level
func InitializeLogger() {
	levelStr := os.Getenv("REVENIUM_LOG_LEVEL")
	if levelStr == "" {
		levelStr = "INFO"
	}

	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		currentLogLevel = LogLevelDebug
	case "INFO":
		currentLogLevel = LogLevelInfo
	case "WARN", "WARNING":
		currentLogLevel = LogLevelWarn
	case "ERROR":
		currentLogLevel = LogLevelError
	default:
		currentLogLevel = LogLevelInfo
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelDebug {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelInfo {
		logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelWarn {
		logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if currentLogLevel <= LogLevelError {
		logger.Printf("[ERROR] "+format, v...)
	}
}

// SetLogLevel sets the current log level
func SetLogLevel(level LogLevel) {
	currentLogLevel = level
}

// GetLogLevel returns the current log level
func GetLogLevel() LogLevel {
	return currentLogLevel
}

// LogLevelFromString converts a string to a LogLevel
func LogLevelFromString(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelInfo
	case "WARN", "WARNING":
		return LogLevelWarn
	case "ERROR":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "INFO"
	}
}

// logRequest logs an HTTP request for debugging
func logRequest(method, url string, headers map[string]string) {
	Debug("HTTP %s %s", method, url)
	if currentLogLevel <= LogLevelDebug {
		for k, v := range headers {
			// Don't log full API keys
			if k == "Authorization" || k == "x-api-key" {
				Debug("  %s: [REDACTED]", k)
			} else {
				Debug("  %s: %s", k, v)
			}
		}
	}
}

// logResponse logs an HTTP response for debugging
func logResponse(statusCode int, body string) {
	Debug("HTTP Response: %d", statusCode)
	if currentLogLevel <= LogLevelDebug && body != "" {
		// Truncate long responses
		if len(body) > 500 {
			Debug("  Body: %s... (truncated)", body[:500])
		} else {
			Debug("  Body: %s", body)
		}
	}
}

// logError logs an error with context
func logError(context string, err error) {
	Error("%s: %v", context, err)
}

// logMeteringPayload logs a metering payload for debugging
func logMeteringPayload(payload interface{}) {
	if currentLogLevel <= LogLevelDebug {
		Debug("Metering payload: %+v", payload)
	}
}
