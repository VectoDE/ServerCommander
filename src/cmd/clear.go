package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"servercommander/src/services"
	"servercommander/src/ui"
)

func ClearConsole() {
	clearCmd, err := getClearCommand()
	if err != nil {
		services.LogToFile(fmt.Sprintf("Error determining clear command: %v", err))
		fmt.Println("Failed to determine the clear command:", err)
		return
	}

	if err := executeClearCommand(clearCmd); err != nil {
		services.LogToFile(fmt.Sprintf("Failed to clear console: %v", err))
		fmt.Println("Failed to clear the console:", err)
		return
	}

	services.LogToFile("Console cleared successfully")

	ui.ApplicationBanner()
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
