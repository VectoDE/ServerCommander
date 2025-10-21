package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// configRoot resolves (and creates if required) the configuration directory
// used by ServerCommander. The path follows the conventions of the underlying
// operating system.
func configRoot() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user config directory: %w", err)
	}

	target := filepath.Join(configDir, "servercommander")
	if err := os.MkdirAll(target, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", target, err)
	}

	return target, nil
}

func sessionsFile() (string, error) {
	root, err := configRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "sessions.json"), nil
}
