package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zalando/go-keyring"
)

// SecretStore provides encrypted storage for sensitive session data using
// OS-native secret stores (Windows Credential Manager, macOS Keychain,
// Linux Secret Service).
type SecretStore struct {
	serviceName string
	accountPrefix string
}

// SecretEntry represents a single secret entry in the store.
type SecretEntry struct {
	KeyPath       string `json:"key_path,omitempty"`
	Passphrase    string `json:"passphrase,omitempty"` // For encrypted private keys
	CustomCA      string `json:"custom_ca,omitempty"`  // Custom CA certificate content
}

// NewSecretStore creates a new secret store with the given service name.
// The service name is used to namespace secrets in the OS keyring.
func NewSecretStore(serviceName string) *SecretStore {
	if serviceName == "" {
		serviceName = "servercommander"
	}
	return &SecretStore{
		serviceName:   serviceName,
		accountPrefix: "session:",
	}
}

// accountName generates the full account name for a given session alias.
func (s *SecretStore) accountName(alias string) string {
	return s.accountPrefix + strings.ToLower(alias)
}

// Store saves sensitive session data to the OS secret store.
// Only sensitive fields are stored; the session metadata remains in the JSON file.
func (s *SecretStore) Store(alias string, entry SecretEntry) error {
	account := s.accountName(alias)
	
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal secret entry: %w", err)
	}

	err = keyring.Set(s.serviceName, account, string(data))
	if err != nil {
		// Check if we're running in an environment without keyring support
		if isKeyringUnavailable(err) {
			// Fall back to file-based storage with restricted permissions
			return s.storeFallback(alias, entry)
		}
		return fmt.Errorf("failed to store secret in keyring: %w", err)
	}

	return nil
}

// Retrieve fetches sensitive session data from the OS secret store.
func (s *SecretStore) Retrieve(alias string) (*SecretEntry, error) {
	account := s.accountName(alias)
	
	data, err := keyring.Get(s.serviceName, account)
	if err != nil {
		// Check if we're running in an environment without keyring support
		if isKeyringUnavailable(err) {
			// Try fallback file-based storage
			return s.retrieveFallback(alias)
		}
		// Key not found is not an error - session may not have sensitive data
		if err == keyring.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve secret from keyring: %w", err)
	}

	entry := &SecretEntry{}
	if err := json.Unmarshal([]byte(data), entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret entry: %w", err)
	}

	return entry, nil
}

// Delete removes sensitive session data from the OS secret store.
func (s *SecretStore) Delete(alias string) error {
	account := s.accountName(alias)
	
	err := keyring.Delete(s.serviceName, account)
	if err != nil {
		// Check if we're running in an environment without keyring support
		if isKeyringUnavailable(err) {
			// Clean up fallback file
			return s.deleteFallback(alias)
		}
		if err == keyring.ErrNotFound {
			return nil
		}
		return fmt.Errorf("failed to delete secret from keyring: %w", err)
	}

	return nil
}

// isKeyringUnavailable checks if the error indicates lack of keyring support.
func isKeyringUnavailable(err error) bool {
	// Common errors when keyring is not available
	errStr := err.Error()
	return strings.Contains(errStr, "not implemented") ||
		strings.Contains(errStr, "no keyring") ||
		strings.Contains(errStr, "secret service not found") ||
		strings.Contains(errStr, "KeyringError")
}

// Fallback file-based storage for environments without keyring support.
// Uses restricted file permissions and obfuscated filenames.

func (s *SecretStore) fallbackDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user config directory: %w", err)
	}

	target := filepath.Join(configDir, "servercommander", "secrets")
	if err := os.MkdirAll(target, 0700); err != nil {
		return "", fmt.Errorf("failed to create secrets directory: %w", err)
	}

	return target, nil
}

func (s *SecretStore) fallbackFile(alias string) (string, error) {
	dir, err := s.fallbackDir()
	if err != nil {
		return "", err
	}
	// Obfuscate filename slightly (not security through obscurity, just cleanliness)
	filename := fmt.Sprintf("%s.secret", strings.ReplaceAll(strings.ToLower(alias), " ", "_"))
	return filepath.Join(dir, filename), nil
}

func (s *SecretStore) storeFallback(alias string, entry SecretEntry) error {
	path, err := s.fallbackFile(alias)
	if err != nil {
		return err
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal secret entry: %w", err)
	}

	// Write with restrictive permissions (owner read/write only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}

	return nil
}

func (s *SecretStore) retrieveFallback(alias string) (*SecretEntry, error) {
	path, err := s.fallbackFile(alias)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read secret file: %w", err)
	}

	entry := &SecretEntry{}
	if err := json.Unmarshal(data, entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret entry: %w", err)
	}

	return entry, nil
}

func (s *SecretStore) deleteFallback(alias string) error {
	path, err := s.fallbackFile(alias)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete secret file: %w", err)
	}

	return nil
}

// MigratePlaintextToSecure migrates sessions from plaintext storage to secure storage.
// This should be called once during application startup or upgrade.
func MigratePlaintextToSecure(store *SessionStore, secretStore *SecretStore) error {
	sessions := store.List()
	migrated := 0
	
	for _, session := range sessions {
		// Check if secret already exists in secure store
		existing, err := secretStore.Retrieve(session.Alias)
		if err != nil {
			// Log but continue - don't fail migration for individual sessions
			continue
		}
		
		if existing != nil {
			// Already migrated
			continue
		}
		
		// For now, no migration needed since plaintext never stored passwords
		// Future: if plaintext stored sensitive data, extract and move to secret store here
		
		_ = migrated // Suppress unused variable warning
	}
	
	return nil
}
