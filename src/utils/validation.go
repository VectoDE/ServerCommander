package utils

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidateHostname prüft ob ein Hostname oder IP-Adresse gültig ist
// Verhindert Command Injection und Path Traversal in Hostnames
func ValidateHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname darf nicht leer sein")
	}

	// Prüfe auf gefährliche Zeichen für Command Injection
	dangerousChars := []string{
		";", "|", "&", "$", "`", "\\", "(", ")", "<", ">",
		"[", "]", "{", "}", "!", "'", "\"", "\n", "\r",
	}
	for _, char := range dangerousChars {
		if strings.Contains(hostname, char) {
			return fmt.Errorf("hostname enthält ungültige Zeichen: %s", char)
		}
	}

	// Prüfe auf Path Traversal Versuche
	if strings.Contains(hostname, "..") || strings.Contains(hostname, "/") {
		return fmt.Errorf("hostname darf keine Pfadnavigation enthalten")
	}

	// Versuche als IP-Adresse zu parsen (IPv4 oder IPv6)
	if net.ParseIP(hostname) != nil {
		return nil // Gültige IP
	}

	// Validiere als Domain-Name
	// RFC 1035 + RFC 1123 compliant regex
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if domainRegex.MatchString(hostname) {
		return nil
	}

	// Spezialfall: localhost
	if hostname == "localhost" {
		return nil
	}

	return fmt.Errorf("ungültiger Hostname: %s", hostname)
}

// ValidateUsername prüft ob ein Benutzername sicher ist
// Verhindert Command Injection in Benutzernamen
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username darf nicht leer sein")
	}

	if len(username) > 64 {
		return fmt.Errorf("username ist zu lang (max 64 Zeichen)")
	}

	// Gefährliche Zeichen für Command Injection
	dangerousChars := []string{
		";", "|", "&", "$", "`", "\\", "(", ")", "<", ">",
		"[", "]", "{", "}", "!", "'", "\"", "\n", "\r", " ",
	}
	for _, char := range dangerousChars {
		if strings.Contains(username, char) {
			return fmt.Errorf("username enthält ungültige Zeichen: %s", char)
		}
	}

	// Nur alphanumerische Zeichen, Unterstrich, Bindestrich, Punkt erlaubt
	validRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if validRegex.MatchString(username) {
		return nil
	}

	return fmt.Errorf("username enthält ungültige Zeichen")
}

// ValidatePort prüft ob eine Portnummer im gültigen Bereich liegt
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port muss zwischen 1 und 65535 liegen (aktuell: %d)", port)
	}
	return nil
}

// ValidateFilePath prüft einen Dateipfad auf Sicherheit
// Verhindert Path Traversal Attacks
func ValidateFilePath(path string, allowRelative bool) error {
	if path == "" {
		return fmt.Errorf("pfad darf nicht leer sein")
	}

	// Prüfe auf absolute Path Traversal Versuche
	cleanPath := filepath.Clean(path)
	if strings.HasPrefix(cleanPath, "..") && !allowRelative {
		return fmt.Errorf("pfad darf nicht außerhalb des erlaubten Bereichs liegen")
	}

	// Gefährliche Zeichen
	dangerousChars := []string{"\x00", ";", "|", "&", "$", "`"}
	for _, char := range dangerousChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("pfad enthält ungültige Zeichen")
		}
	}

	return nil
}

// SanitizeString entfernt gefährliche Steuerzeichen aus einem String
func SanitizeString(input string) string {
	// Entferne Null-Bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Entferne andere gefährliche Steuerzeichen (außer \n und \t für normale Textausgabe)
	sanitized := strings.Map(func(r rune) rune {
		if r >= 0 && r < 32 && r != '\n' && r != '\r' && r != '\t' {
			return -1 // Entferne das Zeichen
		}
		return r
	}, input)

	return sanitized
}

// ContainsInjectionChars prüft ob ein String Command-Injection-Zeichen enthält
func ContainsInjectionChars(input string) bool {
	dangerousChars := []string{
		";", "|", "&", "$", "`", "\\", "(", ")", "<", ">",
		"[", "]", "{", "}", "!", "'", "\"", "\n", "\r", "\x00",
	}
	for _, char := range dangerousChars {
		if strings.Contains(input, char) {
			return true
		}
	}
	return false
}

// ValidateURL prüft ob eine URL sicher und gültig ist
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL darf nicht leer sein")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("ungültige URL: %w", err)
	}

	// Nur erlaubte Schemes
	allowedSchemes := map[string]bool{
		"http":  true,
		"https": true,
		"ftp":   true,
		"sftp":  true,
		"ssh":   true,
	}
	if !allowedSchemes[parsed.Scheme] {
		return fmt.Errorf("nicht unterstütztes Protokoll: %s", parsed.Scheme)
	}

	// Host validieren
	if parsed.Host != "" {
		host := parsed.Hostname()
		if host != "" {
			if err := ValidateHostname(host); err != nil {
				return fmt.Errorf("ungültiger Host in URL: %w", err)
			}
		}
	}

	return nil
}
