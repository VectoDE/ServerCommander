package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// NetworkError represents a network-related error with additional context.
type NetworkError struct {
	Op      string
	Err     error
	Host    string
	Port    int
	Retryable bool
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error during %s: %v", e.Op, e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// AuthError represents an authentication failure.
type AuthError struct {
	Method string
	Reason string
	Host   string
	User   string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication failed for %s@%s using %s: %s", e.User, e.Host, e.Method, e.Reason)
}

// TimeoutError represents a timeout error with context.
type TimeoutError struct {
	Op      string
	Duration time.Duration
	Err     error
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %v during %s: %v", e.Duration, e.Op, e.Err)
}

func (e *TimeoutError) Unwrap() error {
	return e.Err
}

func (e *TimeoutError) Timeout() bool {
	return true
}

// ValidationError represents an input validation error.
type ValidationError struct {
	Field   string
	Value   string
	Reason  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Reason)
}

// IsTemporary checks if an error is temporary and retrying might help.
func IsTemporary(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Temporary()
	}

	errStr := strings.ToLower(err.Error())
	temporaryKeywords := []string{
		"temporary",
		"timeout",
		"connection reset",
		"broken pipe",
		"no route to host",
		"network unreachable",
		"try again",
	}

	for _, keyword := range temporaryKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// IsAuthenticationError checks if an error is an authentication failure.
func IsAuthenticationError(err error) bool {
	if err == nil {
		return false
	}

	var authErr *AuthError
	if errors.As(err, &authErr) {
		return true
	}

	errStr := strings.ToLower(err.Error())
	authKeywords := []string{
		"permission denied",
		"authentication failed",
		"invalid key",
		"unauthorized",
		"access denied",
		"bad credentials",
	}

	for _, keyword := range authKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// IsNetworkError checks if an error is network-related.
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	var netOpErr *net.OpError
	if errors.As(err, &netOpErr) {
		return true
	}

	errStr := strings.ToLower(err.Error())
	networkKeywords := []string{
		"connection",
		"network",
		"host",
		"dial",
		"read",
		"write",
		"timeout",
		"refused",
		"unreachable",
	}

	for _, keyword := range networkKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// ShouldRetry determines if an operation should be retried based on the error.
func ShouldRetry(err error, attempt int, maxAttempts int) bool {
	if err == nil || attempt >= maxAttempts {
		return false
	}

	return IsTemporary(err) && !IsAuthenticationError(err)
}

// UserMessage returns a user-friendly error message.
func UserMessage(err error) string {
	if err == nil {
		return ""
	}

	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return fmt.Sprintf("Invalid %s: %s", validationErr.Field, validationErr.Reason)
	}

	var authErr *AuthError
	if errors.As(err, &authErr) {
		switch authErr.Reason {
		case "permission denied":
			return "Username or password incorrect"
		case "invalid key":
			return "Invalid SSH key provided"
		case "agent not available":
			return "SSH agent not running"
		default:
			return "Authentication failed"
		}
	}

	var timeoutErr *TimeoutError
	if errors.As(err, &timeoutErr) {
		return fmt.Sprintf("Connection timed out after %v", timeoutErr.Duration)
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return "Connection timed out"
		}
		return "Network error occurred"
	}

	errStr := err.Error()
	if strings.Contains(errStr, "connection refused") {
		return "Connection refused - server may be offline"
	}
	if strings.Contains(errStr, "no such host") {
		return "Host not found - check hostname"
	}
	if strings.Contains(errStr, "network is unreachable") {
		return "Network unreachable - check your connection"
	}

	return errStr
}

// WrapNetworkError wraps a low-level network error with context.
func WrapNetworkError(op string, host string, port int, err error) error {
	retryable := IsTemporary(err)
	return &NetworkError{
		Op:        op,
		Err:       err,
		Host:      host,
		Port:      port,
		Retryable: retryable,
	}
}

// NewAuthError creates a new authentication error.
func NewAuthError(method, reason, host, user string) error {
	return &AuthError{
		Method: method,
		Reason: reason,
		Host:   host,
		User:   user,
	}
}

// NewTimeoutError creates a new timeout error.
func NewTimeoutError(op string, duration time.Duration, err error) error {
	return &TimeoutError{
		Op:       op,
		Duration: duration,
		Err:      err,
	}
}

// NewValidationError creates a new validation error.
func NewValidationError(field, value, reason string) error {
	return &ValidationError{
		Field:  field,
		Value:  value,
		Reason: reason,
	}
}
