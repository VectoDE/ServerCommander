package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logOnce sync.Once
	logFile *os.File
	logErr  error
)

// LogToFile appends the provided message to the application log. The file is
// lazily opened to avoid unnecessary IO during tests and keeps the
// initialisation cost minimal for short-lived commands.
func LogToFile(message string) {
	logOnce.Do(func() {
		logFile, logErr = prepareLogFile()
	})

	if logErr != nil {
		fmt.Printf("ERROR: Unable to initialise log file: %v\n", logErr)
		return
	}

	if _, err := logFile.WriteString(formatLogMessage(message)); err != nil {
		fmt.Printf("ERROR: Failed to write to log file: %v\n", err)
	}
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

// formatLogMessage formats the log message with a timestamp.
func formatLogMessage(message string) string {
	return fmt.Sprintf("[%s] %s\n", time.Now().UTC().Format(time.RFC3339), message)
}
