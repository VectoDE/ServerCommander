package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"servercommander/src/console"
)

func clearCommand(args []string) error {
	if console.ClearConsole() {
		console.ApplicationBanner()
		return nil
	}

	clearCmd, err := getClearCommand()
	if err != nil {
		return fmt.Errorf("failed to determine clear command: %w", err)
	}

	if err := executeClearCommand(clearCmd); err != nil {
		return fmt.Errorf("failed to clear the console: %w", err)
	}

	console.ApplicationBanner()
	return nil
}

func getClearCommand() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return "cmd", nil
	case "linux", "darwin":
		return "clear", nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func executeClearCommand(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(command, "/c", "cls")
	} else {
		cmd = exec.Command(command)
	}

	cmd.Stdout = os.Stdout
	return cmd.Run()
}
