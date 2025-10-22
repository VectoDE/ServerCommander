package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
		return errors.New("htop is not available on Windows. Please install an alternative process monitor.")
	}

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
