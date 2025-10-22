package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"servercommander/src/utils"
)

func init() {
	RegisterCommand("htop", "Launch the system monitor with ServerCommander colors", htopCommand)
}

func htopCommand(args []string) error {
	if err := ensureUsage(args, 0, 0, "htop"); err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		return runWindowsProcessMonitor()
	}

	return runHtopBinary()
}

func prepareHtopTheme() (string, func(), error) {
	if len(htopTheme) == 0 {
		return "", nil, nil
	}

	dir, err := os.MkdirTemp("", "servercommander-htop")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temporary configuration directory: %w", err)
	}

	path := filepath.Join(dir, "htoprc")
	if err := os.WriteFile(path, htopTheme, 0o600); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to prepare htop theme: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(dir)
	}

	return path, cleanup, nil
}

func runHtopBinary() error {
	binary, err := exec.LookPath("htop")
	if err != nil {
		return fmt.Errorf("htop executable not found in PATH: %w", err)
	}

	themePath, cleanup, err := prepareHtopTheme()
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	fmt.Printf("%sLaunching htop with ServerCommander theme. Press 'q' to return to the console.%s\n", utils.Green, utils.Reset)

	cmd := exec.Command(binary)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if themePath != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("HTOPRC=%s", themePath))
	}

	return cmd.Run()
}

func runWindowsProcessMonitor() error {
	// Prefer a native htop installation if available.
	if binary, err := exec.LookPath("htop.exe"); err == nil {
		return runBinaryWithTheme(binary)
	}
	if binary, err := exec.LookPath("htop"); err == nil {
		return runBinaryWithTheme(binary)
	}

	// Fall back to PowerShell's process listing with paging.
	powershellCandidates := []string{"powershell.exe", "powershell", "pwsh.exe", "pwsh"}
	var psBinary string
	for _, candidate := range powershellCandidates {
		if binary, err := exec.LookPath(candidate); err == nil {
			psBinary = binary
			break
		}
	}
	if psBinary == "" {
		return errors.New("neither htop nor PowerShell were found in PATH. Please install htop or ensure PowerShell is available")
	}

	fmt.Printf("%sLaunching PowerShell process monitor (htop not found). Use Space/PageUp/PageDown to navigate and Ctrl+C to exit.%s\n", utils.Green, utils.Reset)

	script := strings.Join([]string{
		"$ErrorActionPreference = 'Stop'",
		"Write-Host 'Process list sorted by CPU usage. Press Ctrl+C to exit paging.'",
		"Get-Process | Sort-Object -Property CPU -Descending | Select-Object -Property Id, ProcessName, CPU, PM, WS, StartTime -ErrorAction SilentlyContinue | Format-Table -AutoSize | Out-Host -Paging",
	}, "; ")

	cmd := exec.Command(psBinary, "-NoProfile", "-Command", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runBinaryWithTheme(binary string) error {
	themePath, cleanup, err := prepareHtopTheme()
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	fmt.Printf("%sLaunching htop with ServerCommander theme. Press 'q' to return to the console.%s\n", utils.Green, utils.Reset)

	cmd := exec.Command(binary)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if themePath != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("HTOPRC=%s", themePath))
	}

	return cmd.Run()
}
