package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNetworkError_Error(t *testing.T) {
	err := &NetworkError{
		Op:      "connect",
		Err:     assert.AnError,
		Host:    "example.com",
		Port:    22,
		Retryable: true,
	}
	assert.Contains(t, err.Error(), "network error during connect")
}

func TestAuthError_Error(t *testing.T) {
	err := &AuthError{
		Method: "password",
		Reason: "permission denied",
		Host:   "example.com",
		User:   "admin",
	}
	assert.Contains(t, err.Error(), "authentication failed for admin@example.com")
}

func TestTimeoutError_Error(t *testing.T) {
	err := &TimeoutError{
		Op:      "read",
		Duration: 30 * time.Second,
		Err:     assert.AnError,
	}
	assert.Contains(t, err.Error(), "timeout after 30s during read")
}

func TestTimeoutError_Timeout(t *testing.T) {
	err := &TimeoutError{
		Op:      "read",
		Duration: 30 * time.Second,
		Err:     assert.AnError,
	}
	assert.True(t, err.Timeout())
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:  "hostname",
		Value:  "invalid",
		Reason: "must be a valid domain or IP",
	}
	assert.Contains(t, err.Error(), "validation error for hostname")
}

func TestIsTemporary(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"timeout error", &TimeoutError{Duration: time.Second}, true},
		{"connection reset", assert.AnError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTemporary(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"auth error", NewAuthError("password", "permission denied", "host", "user"), true},
		{"permission denied string", assert.AnError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAuthenticationError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"network error", &NetworkError{Op: "dial"}, true},
		{"generic error", assert.AnError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNetworkError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		attempt   int
		max       int
		expected  bool
	}{
		{"nil error", nil, 1, 3, false},
		{"max attempts reached", assert.AnError, 3, 3, false},
		{"temporary error", &TimeoutError{Duration: time.Second}, 1, 3, true},
		{"auth error should not retry", NewAuthError("password", "denied", "h", "u"), 1, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldRetry(tt.err, tt.attempt, tt.max)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"nil error", nil, ""},
		{"validation error", NewValidationError("port", "99999", "out of range"), "Invalid port"},
		{"auth error permission denied", NewAuthError("password", "permission denied", "h", "u"), "Username or password incorrect"},
		{"auth error invalid key", NewAuthError("key", "invalid key", "h", "u"), "Invalid SSH key provided"},
		{"timeout error", &TimeoutError{Duration: 30 * time.Second}, "Connection timed out after 30s"},
		{"connection refused", assert.AnError, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UserMessage(tt.err)
			if tt.expected != "" {
				assert.Contains(t, result, tt.expected)
			}
		})
	}
}

func TestWrapNetworkError(t *testing.T) {
	err := WrapNetworkError("connect", "example.com", 22, assert.AnError)
	netErr, ok := err.(*NetworkError)
	assert.True(t, ok)
	assert.Equal(t, "connect", netErr.Op)
	assert.Equal(t, "example.com", netErr.Host)
	assert.Equal(t, 22, netErr.Port)
}

func TestNewAuthError(t *testing.T) {
	err := NewAuthError("password", "permission denied", "example.com", "admin")
	authErr, ok := err.(*AuthError)
	assert.True(t, ok)
	assert.Equal(t, "password", authErr.Method)
	assert.Equal(t, "permission denied", authErr.Reason)
	assert.Equal(t, "example.com", authErr.Host)
	assert.Equal(t, "admin", authErr.User)
}

func TestNewTimeoutError(t *testing.T) {
	err := NewTimeoutError("read", 30*time.Second, assert.AnError)
	timeoutErr, ok := err.(*TimeoutError)
	assert.True(t, ok)
	assert.Equal(t, "read", timeoutErr.Op)
	assert.Equal(t, 30*time.Second, timeoutErr.Duration)
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("hostname", "invalid", "must be valid")
	valErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "hostname", valErr.Field)
	assert.Equal(t, "invalid", valErr.Value)
	assert.Equal(t, "must be valid", valErr.Reason)
}

func BenchmarkErrorCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewAuthError("password", "denied", "host", "user")
	}
}

func BenchmarkUserMessage(b *testing.B) {
	err := NewAuthError("password", "permission denied", "host", "user")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = UserMessage(err)
	}
}
