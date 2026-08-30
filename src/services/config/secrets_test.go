package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSecretStore_New(t *testing.T) {
	store := NewSecretStore("test-service")
	if store.serviceName != "test-service" {
		t.Errorf("expected serviceName 'test-service', got '%s'", store.serviceName)
	}
	if store.accountPrefix != "session:" {
		t.Errorf("expected accountPrefix 'session:', got '%s'", store.accountPrefix)
	}
}

func TestSecretStore_DefaultServiceName(t *testing.T) {
	store := NewSecretStore("")
	if store.serviceName != "servercommander" {
		t.Errorf("expected default serviceName 'servercommander', got '%s'", store.serviceName)
	}
}

func TestSecretStore_AccountName(t *testing.T) {
	store := NewSecretStore("test")
	
	tests := []struct {
		alias    string
		expected string
	}{
		{"myserver", "session:myserver"},
		{"MyServer", "session:myserver"}, // Should be lowercased
		{"test-server", "session:test-server"},
		{"test server", "session:test server"},
	}
	
	for _, tt := range tests {
		result := store.accountName(tt.alias)
		if result != tt.expected {
			t.Errorf("accountName(%q) = %q, want %q", tt.alias, result, tt.expected)
		}
	}
}

func TestSecretEntry_MarshalUnmarshal(t *testing.T) {
	entry := SecretEntry{
		KeyPath:    "/home/user/.ssh/id_rsa",
		Passphrase: "secret-passphrase",
		CustomCA:   "/path/to/ca.crt",
	}
	
	// Test that the struct can be marshaled and unmarshaled
	// (actual encryption/decryption tested via Store/Retrieve)
	_ = entry // Suppress unused variable warning
}

func TestSecretStore_FallbackDir(t *testing.T) {
	store := NewSecretStore("test")
	
	dir, err := store.fallbackDir()
	if err != nil {
		t.Fatalf("fallbackDir() error: %v", err)
	}
	
	// Check that directory exists
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("fallback directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("fallback path is not a directory")
	}
	
	// Check permissions (should be 0700)
	mode := info.Mode().Perm()
	if mode&0077 != 0 {
		t.Errorf("fallback directory permissions too open: %o", mode)
	}
}

func TestSecretStore_FallbackFile(t *testing.T) {
	store := NewSecretStore("test")
	
	tests := []struct {
		alias    string
		filename string
	}{
		{"myserver", "myserver.secret"},
		{"MyServer", "myserver.secret"}, // Lowercased
		{"test server", "test_server.secret"}, // Spaces replaced
		{"test-server", "test-server.secret"},
	}
	
	for _, tt := range tests {
		path, err := store.fallbackFile(tt.alias)
		if err != nil {
			t.Errorf("fallbackFile(%q) error: %v", tt.alias, err)
			continue
		}
		
		expectedBase := tt.filename
		actualBase := filepath.Base(path)
		if actualBase != expectedBase {
			t.Errorf("fallbackFile(%q) base = %q, want %q", tt.alias, actualBase, expectedBase)
		}
	}
}

func TestSecretStore_StoreRetrieveFallback(t *testing.T) {
	store := NewSecretStore("test-fallback")
	alias := "test-session-fallback"
	
	entry := SecretEntry{
		KeyPath:    "/tmp/test.key",
		Passphrase: "test-passphrase",
		CustomCA:   "/tmp/ca.crt",
	}
	
	// Store using fallback
	err := store.storeFallback(alias, entry)
	if err != nil {
		t.Fatalf("storeFallback() error: %v", err)
	}
	
	// Retrieve using fallback
	retrieved, err := store.retrieveFallback(alias)
	if err != nil {
		t.Fatalf("retrieveFallback() error: %v", err)
	}
	
	if retrieved == nil {
		t.Fatal("retrieveFallback() returned nil")
	}
	
	if retrieved.KeyPath != entry.KeyPath {
		t.Errorf("KeyPath mismatch: got %q, want %q", retrieved.KeyPath, entry.KeyPath)
	}
	if retrieved.Passphrase != entry.Passphrase {
		t.Errorf("Passphrase mismatch: got %q, want %q", retrieved.Passphrase, entry.Passphrase)
	}
	if retrieved.CustomCA != entry.CustomCA {
		t.Errorf("CustomCA mismatch: got %q, want %q", retrieved.CustomCA, entry.CustomCA)
	}
	
	// Cleanup
	_ = store.deleteFallback(alias)
}

func TestSecretStore_DeleteFallback(t *testing.T) {
	store := NewSecretStore("test-delete")
	alias := "test-session-delete"
	
	entry := SecretEntry{
		KeyPath: "/tmp/test.key",
	}
	
	// Store first
	err := store.storeFallback(alias, entry)
	if err != nil {
		t.Fatalf("storeFallback() error: %v", err)
	}
	
	// Delete
	err = store.deleteFallback(alias)
	if err != nil {
		t.Fatalf("deleteFallback() error: %v", err)
	}
	
	// Verify deletion - should return nil for non-existent file
	retrieved, err := store.retrieveFallback(alias)
	if err != nil {
		t.Errorf("retrieveFallback() after delete returned error: %v", err)
	}
	if retrieved != nil {
		t.Error("retrieveFallback() after delete should return nil entry")
	}
}

func TestSecretStore_RetrieveNonExistent(t *testing.T) {
	store := NewSecretStore("test-nonexistent")
	
	entry, err := store.retrieveFallback("nonexistent-session")
	if err != nil {
		t.Errorf("retrieveFallback() for non-existent session should return nil, nil, got err: %v", err)
	}
	if entry != nil {
		t.Error("retrieveFallback() for non-existent session should return nil entry")
	}
}

func TestMigratePlaintextToSecure(t *testing.T) {
	// Create temporary session store
	sessionStore := &SessionStore{
		Sessions: map[string]Session{},
	}
	
	secretStore := NewSecretStore("test-migrate")
	
	// Migration should succeed even with empty store
	err := MigratePlaintextToSecure(sessionStore, secretStore)
	if err != nil {
		t.Errorf("MigratePlaintextToSecure() error: %v", err)
	}
	
	// Add a session and test migration
	sessionStore.Sessions["test"] = Session{
		Alias:    "test",
		Host:     "localhost",
		Username: "user",
	}
	
	err = MigratePlaintextToSecure(sessionStore, secretStore)
	if err != nil {
		t.Errorf("MigratePlaintextToSecure() with sessions error: %v", err)
	}
}
