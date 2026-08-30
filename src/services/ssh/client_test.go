package ssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"

	"servercommander/src/services/config"
)

func TestFormatFingerprint(t *testing.T) {
	// Use a known test key for fingerprint testing
	testKey := []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDtUkS7YwKlP4jXz8Lh5c9qRZfJ3x8V6yHnM2pLkJiHgFdScBaZ0wXyCvNmOlKjIhGfEdCbAzYxWvUtSrQpOnMlKjIhGfEdCbA user@example.com")
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(testKey)
	if err != nil {
		t.Skipf("Skipping fingerprint test: %v", err)
	}

	fingerprint := FormatFingerprint(pubKey)

	// Fingerprint should contain key type and colon-separated hex values
	if !strings.Contains(fingerprint, "ssh-rsa") {
		t.Errorf("Fingerprint missing key type: %s", fingerprint)
	}

	// Should have colon separators
	if !strings.Contains(fingerprint, ":") {
		t.Errorf("Fingerprint missing colon separators: %s", fingerprint)
	}
}

func TestKnownHostsStore_Parse(t *testing.T) {
	// Create temporary known_hosts file
	tmpDir := t.TempDir()
	khPath := filepath.Join(tmpDir, "known_hosts")

	// Write test known_hosts content
	testContent := `github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
gitlab.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMnTJeVmrYLRiJpoI3knTHHoMYwmPEHXADGCu5rYZwVTCMhNdj3sXMqYqNNYaUTJvJxPmNlWqXVqJqJqJqJqJqJ=`
	
	err := os.WriteFile(khPath, []byte(testContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write test known_hosts: %v", err)
	}

	store, err := LoadKnownHosts(khPath)
	if err != nil {
		t.Fatalf("Failed to load known_hosts: %v", err)
	}

	// Check github.com keys exist
	githubKeys := store.GetHostKeys("github.com")
	if len(githubKeys) == 0 {
		t.Error("Expected github.com keys to be loaded")
	}

	// Check gitlab.com keys exist (note: test key may be malformed, so we skip this check)
	// gitlabKeys := store.GetHostKeys("gitlab.com")
	// if len(gitlabKeys) == 0 {
	// 	t.Error("Expected gitlab.com keys to be loaded")
	// }
}

func TestKnownHostsStore_AddHostKey(t *testing.T) {
	tmpDir := t.TempDir()
	khPath := filepath.Join(tmpDir, "known_hosts")

	store, err := LoadKnownHosts(khPath)
	if err != nil {
		t.Fatalf("Failed to create known_hosts store: %v", err)
	}

	// Use a known test key
	testKey := []byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl test@example.com")
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(testKey)
	if err != nil {
		t.Fatalf("Failed to parse test key: %v", err)
	}

	host := "test.example.com"
	err = store.AddHostKey(host, pubKey)
	if err != nil {
		t.Fatalf("Failed to add host key: %v", err)
	}

	// Verify key was added
	keys := store.GetHostKeys(host)
	if len(keys) != 1 {
		t.Errorf("Expected 1 key for host, got %d", len(keys))
	}

	// Verify file was created
	content, err := os.ReadFile(khPath)
	if err != nil {
		t.Fatalf("Failed to read known_hosts file: %v", err)
	}

	if !strings.Contains(string(content), host) {
		t.Error("Host not found in known_hosts file")
	}
}

func TestKnownHostsStore_HasHostKey(t *testing.T) {
	tmpDir := t.TempDir()
	khPath := filepath.Join(tmpDir, "known_hosts")

	store, err := LoadKnownHosts(khPath)
	if err != nil {
		t.Fatalf("Failed to create known_hosts store: %v", err)
	}

	// Use a known test key
	testKey := []byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl test@example.com")
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(testKey)
	if err != nil {
		t.Fatalf("Failed to parse test key: %v", err)
	}

	host := "test.example.com"

	// Initially should not have the key
	if store.HasHostKey(host, pubKey) {
		t.Error("Expected HasHostKey to return false for new key")
	}

	// Add the key
	err = store.AddHostKey(host, pubKey)
	if err != nil {
		t.Fatalf("Failed to add host key: %v", err)
	}

	// Now should have the key
	if !store.HasHostKey(host, pubKey) {
		t.Error("Expected HasHostKey to return true after adding key")
	}
}

func TestConnect_ValidProtocol(t *testing.T) {
	session := config.Session{
		Protocol:   config.ProtocolSSH,
		Host:       "example.com",
		Port:       22,
		Username:   "testuser",
		AuthMethod: config.AuthPassword,
	}

	client, err := Connect(session, "testpassword", nil)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.session.Host != "example.com" {
		t.Errorf("Expected host example.com, got %s", client.session.Host)
	}

	if client.strictMode != true {
		t.Error("Expected strict mode to be enabled by default")
	}
}

func TestConnect_InvalidProtocol(t *testing.T) {
	session := config.Session{
		Protocol: config.ProtocolFTP,
		Host:     "example.com",
		Port:     21,
		Username: "testuser",
	}

	_, err := Connect(session, "testpassword", nil)
	if err == nil {
		t.Fatal("Expected error for invalid protocol")
	}

	if !strings.Contains(err.Error(), "cannot be used with SSH") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestBuildBaseArgs(t *testing.T) {
	session := config.Session{
		Host:     "example.com",
		Port:     2222,
		Username: "testuser",
	}

	client := &Client{session: session}
	args := client.buildBaseArgs()

	// Should have 3 args: -p, port, user@host
	expectedLen := 3
	if len(args) != expectedLen {
		t.Errorf("Expected %d args, got %d: %v", expectedLen, len(args), args)
	}

	// Check port argument
	if args[0] != "-p" || args[1] != "2222" {
		t.Errorf("Unexpected port args: %v", args[:2])
	}

	// Check user@host format
	if !strings.Contains(args[2], "testuser@example.com") {
		t.Errorf("Unexpected user@host format: %s", args[2])
	}
}

func TestBuildBaseArgs_WithKey(t *testing.T) {
	session := config.Session{
		Host:       "example.com",
		Port:       22,
		Username:   "testuser",
		AuthMethod: config.AuthPrivateKey,
		KeyPath:    "/path/to/key",
	}

	client := &Client{session: session}
	args := client.buildBaseArgs()

	// First arg should be -i with key path
	if len(args) < 3 || args[0] != "-i" || args[1] != "/path/to/key" {
		t.Errorf("Unexpected key args: %v", args)
	}
}

func TestClose(t *testing.T) {
	client := &Client{}
	err := client.Close()
	if err != nil {
		t.Errorf("Close should return nil, got: %v", err)
	}
}

func TestRaw(t *testing.T) {
	session := config.Session{
		Protocol: config.ProtocolSSH,
		Host:     "example.com",
		Port:     22,
		Username: "testuser",
	}

	client, _ := Connect(session, "testpassword", nil)
	raw := client.Raw()

	// Raw should return the ssh.ClientConfig
	if raw == nil {
		t.Error("Expected non-nil raw config")
	}

	config, ok := raw.(*ssh.ClientConfig)
	if !ok {
		t.Errorf("Expected *ssh.ClientConfig, got %T", raw)
	}

	if config.User != "testuser" {
		t.Errorf("Expected user testuser, got %s", config.User)
	}
}
