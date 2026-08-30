package ssh

import (
	"bufio"
	"bytes"
	"crypto/md5" // nolint:gosec // Used only for legacy MD5 fingerprint display
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"servercommander/src/services/config"
	"servercommander/src/utils"
)

// Client encapsulates metadata required to spawn SSH processes.
type Client struct {
	session      config.Session
	password     string
	knownHosts   *KnownHostsStore
	strictMode   bool
	sshConfig    *ssh.ClientConfig
	nativeClient bool // true if using system ssh binary
}

// KnownHostsStore manages the known_hosts file for host key verification.
type KnownHostsStore struct {
	filepath string
	hosts    map[string][]ssh.PublicKey
}

// LoadKnownHosts reads the known_hosts file from the standard location or a custom path.
func LoadKnownHosts(customPath string) (*KnownHostsStore, error) {
	var khPath string
	if customPath != "" {
		khPath = customPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		khPath = filepath.Join(homeDir, ".ssh", "known_hosts")
	}

	store := &KnownHostsStore{
		filepath: khPath,
		hosts:    make(map[string][]ssh.PublicKey),
	}

	// File may not exist yet - that's OK
	if _, err := os.Stat(khPath); err == nil {
		if err := store.parse(); err != nil {
			return nil, fmt.Errorf("failed to parse known_hosts: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat known_hosts file: %w", err)
	}

	return store, nil
}

// parse reads and parses the known_hosts file.
func (k *KnownHostsStore) parse() error {
	file, err := os.Open(k.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		hostPatterns := strings.Split(fields[0], ",")
		keyData := fields[2]

		keyBytes, err := base64.StdEncoding.DecodeString(keyData)
		if err != nil {
			continue // Skip malformed entries
		}

		pubKey, err := ssh.ParsePublicKey(keyBytes)
		if err != nil {
			continue
		}

		for _, pattern := range hostPatterns {
			// Handle negated patterns by skipping them for now
			if strings.HasPrefix(pattern, "!") {
				continue
			}
			k.hosts[pattern] = append(k.hosts[pattern], pubKey)
		}
	}

	return scanner.Err()
}

// GetHostKeys returns all known public keys for a host.
func (k *KnownHostsStore) GetHostKeys(host string) []ssh.PublicKey {
	return k.hosts[host]
}

// AddHostKey adds a new host key to the known_hosts store and persists it.
func (k *KnownHostsStore) AddHostKey(host string, key ssh.PublicKey) error {
	keyType := key.Type()
	keyData := base64.StdEncoding.EncodeToString(key.Marshal())

	// Check if already exists
	existing := k.GetHostKeys(host)
	for _, ek := range existing {
		if ek.Type() == keyType && bytes.Equal(ek.Marshal(), key.Marshal()) {
			return nil // Already present
		}
	}

	// Append to file
	file, err := os.OpenFile(k.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open known_hosts for writing: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s %s %s\n", host, keyType, keyData)
	if err != nil {
		return fmt.Errorf("failed to write to known_hosts: %w", err)
	}

	// Update in-memory cache
	k.hosts[host] = append(k.hosts[host], key)
	return nil
}

// HasHostKey checks if a host key is already known.
func (k *KnownHostsStore) HasHostKey(host string, key ssh.PublicKey) bool {
	existing := k.GetHostKeys(host)
	for _, ek := range existing {
		if ek.Type() == key.Type() && bytes.Equal(ek.Marshal(), key.Marshal()) {
			return true
		}
	}
	return false
}

// FormatFingerprint returns a human-readable fingerprint of a public key.
func FormatFingerprint(key ssh.PublicKey) string {
	hash := md5.Sum(key.Marshal()) // nolint:gosec // MD5 is standard for SSH fingerprints
	hexFingerprint := hex.EncodeToString(hash[:])

	// Format as colon-separated pairs
	parts := make([]string, len(hexFingerprint)/2)
	for i := 0; i < len(hexFingerprint); i += 2 {
		parts[i/2] = hexFingerprint[i : i+2]
	}

	return fmt.Sprintf("%s:%s", key.Type(), strings.Join(parts, ":"))
}

// Connect prepares an SSH client for the provided session with proper host key verification.
func Connect(session config.Session, password string, _ []byte) (*Client, error) {
	if session.Protocol != config.ProtocolSSH && session.Protocol != config.ProtocolSFTP {
		return nil, fmt.Errorf("protocol %s cannot be used with SSH", session.Protocol)
	}

	// Load known_hosts
	knownHosts, err := LoadKnownHosts("")
	if err != nil {
		return nil, fmt.Errorf("failed to load known_hosts: %w", err)
	}

	client := &Client{
		session:      session,
		password:     password,
		knownHosts:   knownHosts,
		strictMode:   true, // Default to strict mode
		nativeClient: false,
	}

	// Build SSH client config with host key callback
	callback := client.createHostKeyCallback()
	authMethods := []ssh.AuthMethod{}

	if session.AuthMethod == config.AuthPrivateKey && session.KeyPath != "" {
		key, err := os.ReadFile(session.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	client.sshConfig = &ssh.ClientConfig{
		User:            session.Username,
		Auth:            authMethods,
		HostKeyCallback: callback,
		Timeout:         30 * time.Second,
	}

	return client, nil
}

// createHostKeyCallback returns a callback function for verifying host keys.
func (c *Client) createHostKeyCallback() ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return c.verifyHostKey(hostname, remote, key)
	}
}

// verifyHostKey verifies a host key and handles user prompts for unknown hosts.
func (c *Client) verifyHostKey(hostname string, remote net.Addr, key ssh.PublicKey) error {
	host := extractHost(hostname)
	
	if c.knownHosts.HasHostKey(host, key) {
		return nil
	}

	existingKeys := c.knownHosts.GetHostKeys(host)
	if len(existingKeys) > 0 {
		return c.createHostKeyChangedError(host, key)
	}

	if c.strictMode {
		return c.handleUnknownHost(host, remote, key)
	}

	return nil
}

// extractHost extracts the hostname without port.
func extractHost(hostname string) string {
	if h, _, err := net.SplitHostPort(hostname); err == nil {
		return h
	}
	return hostname
}

// createHostKeyChangedError creates an error for changed host keys (potential MITM).
func (c *Client) createHostKeyChangedError(host string, key ssh.PublicKey) error {
	return fmt.Errorf("WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!\n"+
		"Host: %s\n"+
		"New key type: %s\n"+
		"New fingerprint: %s\n"+
		"This could indicate a man-in-the-middle attack.\n"+
		"Please verify the host key manually and update known_hosts if legitimate.",
		host, key.Type(), FormatFingerprint(key))
}

// handleUnknownHost prompts the user to accept an unknown host key.
func (c *Client) handleUnknownHost(host string, remote net.Addr, key ssh.PublicKey) error {
	fingerprint := FormatFingerprint(key)
	fmt.Printf("\n%s=== SSH Host Key Verification ===%s\n", utils.Cyan, utils.Reset)
	fmt.Printf("The authenticity of host '%s (%s)' can't be established.\n", host, remote.String())
	fmt.Printf("%s key fingerprint: %s\n", key.Type(), fingerprint)
	fmt.Printf("This key is not stored in your known_hosts file.\n\n")

	accept, err := utils.PromptBool(fmt.Sprintf("Are you sure you want to continue connecting and trust this host"), false)
	if err != nil {
		return fmt.Errorf("failed to prompt for host key verification: %w", err)
	}

	if !accept {
		return fmt.Errorf("host key verification declined by user")
	}

	if err := c.knownHosts.AddHostKey(host, key); err != nil {
		fmt.Printf("%sWarning: Failed to save host key to known_hosts: %v%s\n", utils.Yellow, err, utils.Reset)
	} else {
		fmt.Printf("%sHost key added to known_hosts.%s\n\n", utils.Green, utils.Reset)
	}

	return nil
}

// Close is a no-op kept for API compatibility with other services.
func (c *Client) Close() error {
	return nil
}

// InteractiveShell spawns the system ssh command with strict host key checking.
func (c *Client) InteractiveShell() error {
	args := c.buildBaseArgs()
	
	// Enable strict host key checking for system ssh
	args = append([]string{"-o", "StrictHostKeyChecking=yes", "-o", "UserKnownHostsFile=" + c.knownHosts.filepath}, args...)
	
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run executes a remote command via ssh and captures its combined output.
func (c *Client) Run(command string) (string, error) {
	args := append(c.buildBaseArgs(), command)
	args = append([]string{"-o", "StrictHostKeyChecking=yes", "-o", "UserKnownHostsFile=" + c.knownHosts.filepath}, args...)
	
	cmd := exec.Command("ssh", args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	if err := cmd.Run(); err != nil {
		return buffer.String(), fmt.Errorf("remote command failed: %w", err)
	}

	return buffer.String(), nil
}

// RunWithConfig executes a command using the native Go SSH client.
func (c *Client) RunWithConfig(command string) (string, error) {
	addr := fmt.Sprintf("%s:%d", c.session.Host, c.session.Port)
	
	conn, err := ssh.Dial("tcp", addr, c.sshConfig)
	if err != nil {
		return "", fmt.Errorf("failed to establish SSH connection: %w", err)
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer sess.Close()

	output, err := sess.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// Raw exposes the underlying SSH client config for advanced usage.
func (c *Client) Raw() interface{} {
	return c.sshConfig
}

func (c *Client) buildBaseArgs() []string {
	args := []string{"-p", strconv.Itoa(c.session.Port), fmt.Sprintf("%s@%s", c.session.Username, c.session.Host)}
	if c.session.AuthMethod == config.AuthPrivateKey && c.session.KeyPath != "" {
		args = append([]string{"-i", c.session.KeyPath}, args...)
	}
	return args
}
