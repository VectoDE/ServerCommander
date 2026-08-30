package services

import (
"fmt"
"os"
"path/filepath"
"regexp"
"sync"
"time"
)

// LogLevel defines the severity of a log message
type LogLevel int

const (
LogLevelDebug LogLevel = iota
LogLevelInfo
LogLevelWarn
LogLevelError
)

// Sensitive patterns to redact from logs
var (
passwordPattern   = regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[=:]\s*['"]?[^\s'"]+['"]?`)
tokenPattern      = regexp.MustCompile(`(?i)(token|api[_-]?key|secret[_-]?key)\s*[=:]\s*['"]?[^\s'"]+['"]?`)
privateKeyPattern = regexp.MustCompile(`-----BEGIN (RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----[\s\S]*?-----END (RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----`)
authHeaderPattern = regexp.MustCompile(`(?i)(authorization|auth)\s*[=:]\s*['"]?[^\s'"]+['"]?`)
sessionIDPattern  = regexp.MustCompile(`(?i)(session[_-]?id|sid)\s*[=:]\s*['"]?[^\s'"]+['"]?`)
)

// redactedPlaceholder is the standard placeholder for redacted sensitive data
const redactedPlaceholder = "$1=***REDACTED***"

var (
logMutex sync.Mutex
logFile  *os.File
logErr   error
logLevel LogLevel = LogLevelInfo
)

// SetLogLevel sets the minimum log level to write
func SetLogLevel(level LogLevel) {
logMutex.Lock()
defer logMutex.Unlock()
logLevel = level
}

// getLogLevelString returns string representation of log level
func getLogLevelString(level LogLevel) string {
switch level {
case LogLevelDebug:
return "DEBUG"
case LogLevelInfo:
return "INFO"
case LogLevelWarn:
return "WARN"
case LogLevelError:
return "ERROR"
default:
return "UNKNOWN"
}
}

// sanitizeLogMessage removes sensitive information from log messages
func sanitizeLogMessage(message string) string {
sanitized := message
	sanitized = passwordPattern.ReplaceAllString(sanitized, redactedPlaceholder)
	sanitized = tokenPattern.ReplaceAllString(sanitized, redactedPlaceholder)
sanitized = privateKeyPattern.ReplaceAllString(sanitized, "***PRIVATE_KEY_REDACTED***")
	sanitized = authHeaderPattern.ReplaceAllString(sanitized, redactedPlaceholder)
	sanitized = sessionIDPattern.ReplaceAllString(sanitized, redactedPlaceholder)
return sanitized
}

// Log writes a structured log message with level and sanitization
func Log(level LogLevel, component, message string) {
logMutex.Lock()
defer logMutex.Unlock()

if level < logLevel {
return
}

if logFile == nil {
logFile, logErr = prepareLogFile()
if logErr != nil {
fmt.Printf("ERROR: Unable to initialise log file: %v\n", logErr)
return
}
}

sanitizedMessage := sanitizeLogMessage(message)
logEntry := formatStructuredLog(level, component, sanitizedMessage)

if _, err := logFile.WriteString(logEntry); err != nil {
fmt.Printf("ERROR: Failed to write to log file: %v\n", err)
}
}

// LogDebug logs a debug message
func LogDebug(component, message string) {
Log(LogLevelDebug, component, message)
}

// LogInfo logs an info message
func LogInfo(component, message string) {
Log(LogLevelInfo, component, message)
}

// LogWarn logs a warning message
func LogWarn(component, message string) {
Log(LogLevelWarn, component, message)
}

// LogError logs an error message
func LogError(component, message string) {
Log(LogLevelError, component, message)
}

// LegacyLogToFile maintains backward compatibility while adding sanitization
// Deprecated: Use Log() or LogInfo() instead
func LegacyLogToFile(message string) {
LogInfo("legacy", message)
}

func prepareLogFile() (*os.File, error) {
logDir, err := getLogDir()
if err != nil {
return nil, err
}

logFilePath := filepath.Join(logDir, "servercommander.log")
file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
if err != nil {
return nil, fmt.Errorf("unable to open log file: %w", err)
}

return file, nil
}

func getLogDir() (string, error) {
configDir, err := os.UserConfigDir()
if err != nil {
return "", fmt.Errorf("failed to resolve user config directory: %w", err)
}

logDir := filepath.Join(configDir, "servercommander", "logs")
if err := os.MkdirAll(logDir, 0750); err != nil {
return "", fmt.Errorf("failed to create log directory: %w", err)
}

return logDir, nil
}

// formatStructuredLog creates a structured log entry with timestamp, level, and component
func formatStructuredLog(level LogLevel, component, message string) string {
timestamp := time.Now().UTC().Format(time.RFC3339)
levelStr := getLogLevelString(level)
return fmt.Sprintf("[%s] [%s] [%s] %s\n", timestamp, levelStr, component, message)
}

// formatLogMessage formats the log message with a timestamp (legacy function).
// Deprecated: Use formatStructuredLog instead
func formatLogMessage(message string) string {
return fmt.Sprintf("[%s] %s\n", time.Now().UTC().Format(time.RFC3339), message)
}

// Flush closes the log file, ensuring all buffered data is written
func Flush() {
logMutex.Lock()
defer logMutex.Unlock()

if logFile != nil {
logFile.Sync()
}
}

// Close properly closes the log file
func Close() error {
logMutex.Lock()
defer logMutex.Unlock()

if logFile != nil {
err := logFile.Close()
logFile = nil
return err
}
return nil
}
