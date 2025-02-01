package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func LogToFile(message string) {
	logDir, err := getLogDir()
	if err != nil {
		fmt.Printf("ERROR: Unable to create log directory: %v\n", err)
		return
	}

	logFilePath := filepath.Join(logDir, "console.log")
	file, err := openLogFile(logFilePath)
	if err != nil {
		fmt.Printf("ERROR: Unable to open log file: %v\n", err)
		return
	}
	defer file.Close()

	logMessage := formatLogMessage(message)

	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Printf("ERROR: Failed to write to log file: %v\n", err)
	}
}

func getLogDir() (string, error) {
	tempDir := os.TempDir()
	logDir := filepath.Join(tempDir, "servercommander")

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create log directory: %w", err)
		}
	}
	return logDir, nil
}

func openLogFile(logFilePath string) (*os.File, error) {
	return os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

// formatLogMessage formats the log message with a timestamp.
func formatLogMessage(message string) string {
	return fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}
