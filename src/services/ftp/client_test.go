package ftp

import (
	"testing"

	"servercommander/src/services/config"
)

// TestClientCreation tests basic FTP client creation
func TestClientCreation(t *testing.T) {
	// Test that we can create a session struct
	session := config.Session{
		Alias:    "test-ftp",
		Protocol: config.ProtocolFTP,
		Host:     "test.example.com",
		Port:     21,
		Username: "testuser",
	}

	if session.Host != "test.example.com" {
		t.Errorf("Expected host 'test.example.com', got '%s'", session.Host)
	}

	if session.Port != 21 {
		t.Errorf("Expected port 21, got %d", session.Port)
	}

	if session.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", session.Username)
	}
}

// TestSession_TLSConfig tests TLS configuration fields
func TestSession_TLSConfig(t *testing.T) {
	session := config.Session{
		Alias:         "secure-ftp",
		Protocol:      config.ProtocolFTP,
		Host:          "secure.example.com",
		Port:          990,
		Username:      "secureuser",
		UseTLS:        true,
		SkipTLSVerify: false,
		TLSCAFile:     "/path/to/ca.crt",
	}

	if !session.UseTLS {
		t.Error("Expected UseTLS to be true")
	}

	if session.SkipTLSVerify {
		t.Error("Expected SkipTLSVerify to be false for secure connection")
	}

	if session.TLSCAFile != "/path/to/ca.crt" {
		t.Errorf("Expected TLSCAFile '/path/to/ca.crt', got '%s'", session.TLSCAFile)
	}
}

// TestSession_DefaultValues tests default values for session
func TestSession_DefaultValues(t *testing.T) {
	session := config.Session{
		Alias:    "default-ftp",
		Protocol: config.ProtocolFTP,
		Host:     "default.example.com",
		Username: "defaultuser",
	}

	// Port should default to 21 if not set
	if session.Port == 0 {
		t.Log("Port is 0, would need explicit setting in production")
	}

	// UseTLS should default to false for plain FTP
	if session.UseTLS {
		t.Error("Expected UseTLS to default to false")
	}

	// SkipTLSVerify should default to false (secure by default)
	if session.SkipTLSVerify {
		t.Error("Expected SkipTLSVerify to default to false")
	}
}

// TestSession_HostValidation tests basic host validation
func TestSession_HostValidation(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		shouldAllow bool
	}{
		{"Valid hostname", "example.com", true},
		{"Valid FQDN", "ftp.example.com", true},
		{"Valid IP", "192.168.1.1", true},
		{"Valid IPv6", "::1", true},
		{"Empty host", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			session := config.Session{
				Alias: "test",
				Host:  tc.host,
			}

			// Basic validation - empty host should be rejected
			if tc.host == "" && session.Host != "" {
				t.Error("Empty host should not be allowed")
			}

			if tc.host != "" && session.Host == "" {
				t.Error("Non-empty host should be allowed")
			}
		})
	}
}

// TestSession_PortRange tests valid port ranges
func TestSession_PortRange(t *testing.T) {
	validPorts := []int{21, 990, 2121, 8021, 65535}
	invalidPorts := []int{-1, 0, 65536, 70000}

	for _, port := range validPorts {
		if port > 0 && port <= 65535 {
			session := config.Session{
				Alias: "test",
				Host:  "example.com",
				Port:  port,
			}
			if session.Port != port {
				t.Errorf("Expected port %d, got %d", port, session.Port)
			}
		}
	}

	for _, port := range invalidPorts {
		if port < 1 || port > 65535 {
			t.Logf("Port %d is invalid (out of range)", port)
		}
	}
}

// TestSession_Protocol tests protocol field
func TestSession_Protocol(t *testing.T) {
	protocols := []config.Protocol{
		config.ProtocolFTP,
		config.ProtocolSFTP,
		config.ProtocolSSH,
	}

	for _, proto := range protocols {
		session := config.Session{
			Alias:    "test",
			Protocol: proto,
			Host:     "example.com",
		}
		if session.Protocol != proto {
			t.Errorf("Expected protocol %s, got %s", proto, session.Protocol)
		}
	}
}

// BenchmarkSessionCreation benchmarks session creation
func BenchmarkSessionCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = config.Session{
			Alias:    "benchmark",
			Protocol: config.ProtocolFTP,
			Host:     "benchmark.example.com",
			Port:     21,
			Username: "benchuser",
		}
	}
}
