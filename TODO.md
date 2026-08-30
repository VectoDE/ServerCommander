# ServerCommander - TODO Roadmap

## Overview

This document contains the complete, actionable roadmap for bringing ServerCommander to production-ready enterprise quality. Each item includes ID, priority, category, problem description, concrete changes required, affected files, dependencies, risk level, required tests, and acceptance criteria.

---

## P0 - Release Blockers (Critical Security & Stability)

### P0-001: Insecure TLS Configuration in FTP Client
- **ID:** P0-001
- **Priority:** P0
- **Kategorie:** Security
- **Problem:** The FTP client uses `InsecureSkipVerify: true` for TLS connections, completely bypassing certificate validation and enabling MITM attacks.
- **Konkrete Änderung:** Implement proper certificate validation with configurable CA bundles. Add hostname verification. Only allow explicit user override via configuration flag with clear warnings.
- **Dateien:** `/workspace/src/services/ftp/client.go` (line 174)
- **Abhängigkeiten:** None
- **Risiko:** Low (breaking change for users relying on self-signed certs without explicit config)
- **Tests:** Unit tests for certificate validation, integration tests with valid/invalid certificates
- **Acceptance Criteria:** 
  - Default behavior validates certificates properly
  - Hostname verification enabled by default
  - Explicit config option to disable verification with warning logged
  - No hardcoded `InsecureSkipVerify: true`

### P0-002: No SSH Host Key Verification
- **ID:** P0-002
- **Priority:** P0
- **Kategorie:** Security
- **Problem:** SSH implementation delegates to system `ssh` binary but provides no known_hosts management, host key verification, or changed-key detection. Users are vulnerable to MITM attacks.
- **Konkrete Änderung:** Implement proper host key verification using known_hosts file. Add fingerprint display and confirmation for unknown hosts. Detect and warn on changed host keys.
- **Dateien:** `/workspace/src/services/ssh/client.go`, `/workspace/src/cmd/ssh.go`
- **Abhängigkeiten:** None
- **Risiko:** Medium (may break existing workflows that rely on accepting all host keys)
- **Tests:** Integration tests with mock SSH servers, known_hosts parsing tests
- **Acceptance Criteria:**
  - known_hosts file support (read/write)
  - Fingerprint display for unknown hosts
  - User prompt for unknown host keys
  - Warning on changed host keys
  - Configurable strictness levels

### P0-003: Password Echo in Terminal
- **ID:** P0-003
- **Priority:** P0
- **Kategorie:** Security
- **Problem:** `PromptPassword` explicitly states "input hidden not supported" and passwords are echoed to terminal, exposing credentials to shoulder surfing and terminal logs.
- **Konkrete Änderung:** Implement proper password masking using platform-specific APIs (golang.org/x/term or similar).
- **Dateien:** `/workspace/src/utils/prompt.go`
- **Abhängigkeiten:** golang.org/x/term package
- **Risiko:** Low
- **Tests:** Manual testing on Windows/Linux/macOS
- **Acceptance Criteria:**
  - Password input not visible in terminal
  - Works on Windows, Linux, macOS
  - Graceful fallback if terminal not available

### P0-004: No Session File Encryption
- **ID:** P0-004
- **Priority:** P0
- **Kategorie:** Security
- **Problem:** Session configurations stored in plaintext JSON at `~/.config/servercommander/sessions.json`. While passwords aren't stored, key paths, usernames, and host information are exposed.
- **Konkrete Änderung:** Encrypt sensitive session data using OS-native secret stores (Windows Credential Manager/DPAPI, macOS Keychain, Linux Secret Service).
- **Dateien:** `/workspace/src/services/config/sessions.go`, `/workspace/src/services/config/paths.go`
- **Abhängigkeiten:** Platform-specific secret store libraries
- **Risiko:** Medium (migration required for existing sessions)
- **Tests:** Migration tests, encryption/decryption tests, cross-platform tests
- **Acceptance Criteria:**
  - Sensitive data encrypted at rest
  - Uses native OS secret stores
  - Backward-compatible migration from plaintext
  - Proper cleanup on session deletion

### P0-005: Command Injection Risk in SSH/SFTP Delegation
- **ID:** P0-005
- **Priority:** P0
- **Kategorie:** Security
- **Problem:** External commands (`ssh`, `sftp`) are invoked with user-controlled arguments. While basic sanitization exists, there's potential for argument injection through malformed session configurations.
- **Konkrete Änderung:** Validate and sanitize all inputs before passing to external commands. Use argument arrays instead of string concatenation where possible.
- **Dateien:** `/workspace/src/services/ssh/client.go`, `/workspace/src/cmd/sftp.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Fuzzing tests with malicious inputs, penetration testing
- **Acceptance Criteria:**
  - All external command inputs validated
  - No shell interpolation used
  - Input length limits enforced
  - Special characters properly escaped

### P0-006: Missing go.sum File
- **ID:** P0-006
- **Priority:** P0
- **Kategorie:** Supply Chain Security
- **Problem:** Repository has no `go.sum` file, making dependency verification impossible and enabling supply chain attacks through modified dependencies.
- **Konkrete Änderung:** Run `go mod tidy` to generate proper go.sum. Pin all dependency versions. Add CI check for go.sum presence.
- **Dateien:** `/workspace/go.mod`, `/workspace/go.sum` (create)
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** CI pipeline validation
- **Acceptance Criteria:**
  - go.sum file present and committed
  - All dependencies pinned
  - CI fails if go.sum missing or modified unexpectedly

### P0-007: No Test Coverage
- **ID:** P0-007
- **Priority:** P0
- **Kategorie:** Testing
- **Problem:** Zero test files exist in the entire codebase. No unit tests, integration tests, or security tests. Changes cannot be verified safely.
- **Konkrete Änderung:** Create comprehensive test suite starting with critical security components (FTP TLS, SSH connection, session management).
- **Dateien:** All packages need `*_test.go` files
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** N/A (this IS the testing infrastructure)
- **Acceptance Criteria:**
  - Minimum 60% code coverage
  - All security-critical functions tested
  - CI runs tests on every commit
  - Tests cover error paths

---

## P1 - Before Public Beta

### P1-001: Go Version Mismatch
- **ID:** P1-001
- **Priority:** P1
- **Kategorie:** Infrastructure
- **Problem:** `go.mod` specifies Go 1.23, but CI uses 1.23.5 and many systems have older versions. Build scripts may fail.
- **Konkrete Änderung:** Standardize on specific Go version across all configs. Update documentation with minimum version requirements.
- **Dateien:** `/workspace/go.mod`, `/workspace/.github/workflows/main.yml`, `/workspace/scripts/build.sh`, `/workspace/scripts/build.bat`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Build tests on multiple Go versions
- **Acceptance Criteria:**
  - Consistent Go version specification
  - Documentation updated
  - Build succeeds on specified version

### P1-002: Documentation-Inconsistency
- **ID:** P1-002
- **Priority:** P1
- **Kategorie:** Documentation
- **Problem:** Documentation describes features that don't exist (REST API, config.yaml, themes, process management). Actual implementation is CLI-only with JSON sessions.
- **Konkrete Änderung:** Either implement described features or update documentation to match actual functionality. Remove misleading API docs.
- **Dateien:** `/workspace/docs/API.md`, `/workspace/docs/CONFIGURATION.md`, `/workspace/docs/THEMES.md`, `/workspace/docs/USAGE.md`
- **Abhängigkeiten:** Feature implementation decisions
- **Risiko:** High (breaking changes to expected functionality)
- **Tests:** Documentation review
- **Acceptance Criteria:**
  - Documentation matches implementation
  - No references to non-existent features
  - Clear indication of planned vs implemented features

### P1-003: No Error Handling for Network Failures
- **ID:** P1-003
- **Priority:** P1
- **Kategorie:** Resilience
- **Problem:** Network operations lack timeout handling, retry logic, and graceful degradation. Connection failures may hang indefinitely.
- **Konkrete Änderung:** Add context-based timeouts, exponential backoff retries, and clear error messages for network failures.
- **Dateien:** `/workspace/src/services/ftp/client.go`, `/workspace/src/services/ssh/client.go`, `/workspace/src/cmd/sftp.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Network failure simulation tests
- **Acceptance Criteria:**
  - All network operations have timeouts
  - Retry logic with backoff
  - Clear error messages
  - No indefinite hangs

### P1-004: Incomplete FTP Implementation
- **ID:** P1-004
- **Priority:** P1
- **Kategorie:** Features
- **Problem:** FTP client lacks essential features: no resume support, no bandwidth limiting, no transfer queue, no checksum verification, incomplete MLSD parsing.
- **Konkrete Änderung:** Implement missing FTP features or clearly document limitations. Consider using established FTP library instead of custom implementation.
- **Dateien:** `/workspace/src/services/ftp/client.go`
- **Abhängigkeiten:** May require external library
- **Risiko:** Medium
- **Tests:** FTP interoperability tests
- **Acceptance Criteria:**
  - Resume interrupted transfers
  - Progress reporting
  - Checksum verification option
  - Better directory listing parsing

### P1-005: No Logging Sanitization
- **ID:** P1-005
- **Priority:** P1
- **Kategorie:** Security
- **Problem:** Logger writes command outputs and potentially sensitive data to log file without sanitization. Passwords, tokens, or terminal content could be logged.
- **Konkrete Änderung:** Implement structured logging with automatic redaction of sensitive patterns. Add log levels and rotation.
- **Dateien:** `/workspace/src/services/logger.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Log output analysis
- **Acceptance Criteria:**
  - No passwords in logs
  - No private keys in logs
  - Configurable log levels
  - Log rotation implemented

### P1-006: Race Condition in Prompt Reader
- **ID:** P1-006
- **Priority:** P1
- **Kategorie:** Code Quality
- **Problem:** `SetPromptReader` uses mutex but `readLine` acquires RLock while the setter uses Lock. Potential race during concurrent access.
- **Konkrete Änderung:** Review and fix concurrent access patterns. Consider if mutex is needed at all for typical usage.
- **Dateien:** `/workspace/src/utils/prompt.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Race detector tests (`go test -race`)
- **Acceptance Criteria:**
  - No race conditions detected
  - Thread-safe prompt handling

### P1-007: Missing Context Propagation
- **ID:** P1-007
- **Priority:** P1
- **Kategorie:** Code Quality
- **Problem:** No context.Context usage throughout codebase. Cannot cancel operations, no timeout propagation, resource leaks possible.
- **Konkrete Änderung:** Add context.Context to all long-running operations. Support cancellation and timeouts.
- **Dateien:** All service and command files
- **Abhängigkeiten:** None
- **Risiko:** Medium (API changes)
- **Tests:** Cancellation tests
- **Acceptance Criteria:**
  - Context passed through call chain
  - Operations respect cancellation
  - Resources cleaned up on cancel

### P1-008: No Resource Cleanup on Errors
- **ID:** P1-008
- **Priority:** P1
- **Kategorie:** Code Quality
- **Problem:** Several code paths don't properly clean up resources (temp files, connections) when errors occur mid-operation.
- **Konkrete Änderung:** Use defer consistently for cleanup. Implement proper error handling chains.
- **Dateien:** `/workspace/src/cmd/sftp.go`, `/workspace/src/services/ftp/client.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Resource leak detection
- **Acceptance Criteria:**
  - No temp file leaks
  - Connections always closed
  - Memory properly freed

---

## P2 - Before Stable Release

### P2-001: Add go.sum and Dependency Audit
- **ID:** P2-001
- **Priority:** P2
- **Kategorie:** Supply Chain
- **Problem:** No dependency tracking beyond go.mod. Need full SBOM and vulnerability scanning.
- **Konkrete Änderung:** Generate SBOM, run govulncheck, document all transitive dependencies.
- **Dateien:** `/workspace/go.sum`, new SBOM file
- **Abhängigkeiten:** P0-006
- **Risiko:** Low
- **Tests:** Dependency scan in CI
- **Acceptance Criteria:**
  - Complete SBOM generated
  - No known vulnerabilities
  - All licenses documented

### P2-002: Implement Proper Terminal Handling
- **ID:** P2-002
- **Priority:** P2
- **Kategorie:** UX
- **Problem:** Basic ANSI escape sequences only. No proper terminal emulation, resize handling, or Unicode support verification.
- **Konkrete Änderung:** Evaluate terminal library (e.g., github.com/nsf/termbox-go or similar). Implement proper terminal handling.
- **Dateien:** `/workspace/src/console/console.go`, `/workspace/src/utils/colors.go`
- **Abhängigkeiten:** Terminal library
- **Risiko:** Medium
- **Tests:** Terminal compatibility tests
- **Acceptance Criteria:**
  - Proper terminal resize handling
  - Unicode/CJK support
  - Mouse support (optional)
  - Scrollback buffer

### P2-003: Session Import/Export
- **ID:** P2-003
- **Priority:** P2
- **Kategorie:** Features
- **Problem:** No way to backup or transfer sessions between installations.
- **Konkrete Änderung:** Add export/import commands for sessions. Support encrypted exports for sensitive data.
- **Dateien:** `/workspace/src/cmd/session.go`
- **Abhängigkeiten:** P0-004
- **Risiko:** Low
- **Tests:** Import/export round-trip tests
- **Acceptance Criteria:**
  - Export sessions to file
  - Import sessions from file
  - Optional encryption
  - Validation on import

### P2-004: Configuration File Support
- **ID:** P2-004
- **Priority:** P2
- **Kategorie:** Features
- **Problem:** No configuration file support despite documentation claiming YAML config exists. All settings must be set per-session.
- **Konkrete Änderung:** Implement actual configuration file support matching documentation OR update documentation to remove false claims.
- **Dateien:** New config package, update all commands
- **Abhängigkeiten:** Decision on feature scope
- **Risiko:** Medium
- **Tests:** Configuration loading tests
- **Acceptance Criteria:**
  - Config file parsed correctly
  - Defaults applied properly
  - CLI overrides work
  - Environment variable support

### P2-005: Windows Console Compatibility
- **ID:** P2-005
- **Priority:** P2
- **Kategorie:** Platform Support
- **Problem:** Windows console handling is minimal. ANSI colors may not work on older Windows versions. Clear console falls back to cmd.
- **Konkrete Änderung:** Enable virtual terminal processing on Windows 10+. Provide legacy fallback.
- **Dateien:** `/workspace/src/console/console.go`, `/workspace/src/utils/colors.go`
- **Abhängigkeiten:** Windows API calls
- **Risiko:** Low
- **Tests:** Windows console tests
- **Acceptance Criteria:**
  - Colors work on Windows 10+
  - Graceful degradation on older Windows
  - Proper console mode restoration

### P2-006: Proper Exit Handling
- **ID:** P2-006
- **Priority:** P2
- **Kategorie:** Code Quality
- **Problem:** `exitCommand` calls `os.Exit(0)` directly, bypassing defer statements and cleanup. Resources may leak.
- **Konkrete Änderung:** Refactor exit flow to use proper cleanup. Return error to main instead of direct exit.
- **Dateien:** `/workspace/src/cmd/exit.go`, `/workspace/src/main.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Cleanup verification tests
- **Acceptance Criteria:**
  - All defers execute on exit
  - Connections closed properly
  - Temp files cleaned up

### P2-007: Signal Handling
- **ID:** P2-007
- **Priority:** P2
- **Kategorie:** Resilience
- **Problem:** No signal handling for SIGINT, SIGTERM. Abrupt termination may leave resources open.
- **Konkrete Änderung:** Add signal handlers for graceful shutdown. Clean up resources on interrupt.
- **Dateien:** `/workspace/src/main.go`, `/workspace/src/console/console.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Signal handling tests
- **Acceptance Criteria:**
  - Ctrl+C handled gracefully
  - Resources cleaned up
  - Goodbye message shown

### P2-008: Path Traversal Prevention
- **ID:** P2-008
- **Priority:** P2
- **Kategorie:** Security
- **Problem:** File operations don't validate paths for traversal attacks. Malicious remote paths could write outside intended directories.
- **Konkrete Änderung:** Validate and sanitize all file paths. Prevent path traversal using filepath.Clean and prefix checks.
- **Dateien:** `/workspace/src/cmd/sftp.go`, `/workspace/src/cmd/ftp.go`, `/workspace/src/services/ftp/client.go`
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Path traversal attack tests
- **Acceptance Criteria:**
  - No path traversal possible
  - Paths validated before use
  - Clear error on invalid paths

---

## P3 - Enterprise/Advanced Features

### P3-001: License API Integration
- **ID:** P3-001
- **Priority:** P3
- **Kategorie:** Enterprise
- **Problem:** No licensing system implemented for commercial distribution.
- **Konkrete Änderung:** Integrate license API with offline support, device binding, and entitlement validation.
- **Dateien:** New license package
- **Abhängigkeiten:** License API specification
- **Risiko:** Medium
- **Tests:** License validation tests
- **Acceptance Criteria:**
  - Online validation
  - Offline grace period
  - Device binding
  - Entitlement checks

### P3-002: RBAC/Policies
- **ID:** P3-002
- **Priority:** P3
- **Kategorie:** Enterprise
- **Problem:** No role-based access control or policy enforcement.
- **Konkrete Änderung:** Implement policy engine for restricting connections, protocols, and actions.
- **Dateien:** New policy package
- **Abhängigkeiten:** P3-001
- **Risiko:** Medium
- **Tests:** Policy enforcement tests
- **Acceptance Criteria:**
  - Role definitions
  - Policy enforcement
  - Audit logging

### P3-003: Central Session Management
- **ID:** P3-003
- **Priority:** P3
- **Kategorie:** Enterprise
- **Problem:** Sessions stored locally only. No centralized management for teams.
- **Konkrete Änderung:** Add optional cloud sync for sessions with encryption.
- **Dateien:** New sync package
- **Abhängigkeiten:** P0-004
- **Risiko:** Medium
- **Tests:** Sync conflict resolution tests
- **Acceptance Criteria:**
  - Encrypted sync
  - Conflict resolution
  - Team sharing

### P3-004: Audit Logging
- **ID:** P3-004
- **Priority:** P3
- **Kategorie:** Enterprise
- **Problem:** No audit trail for compliance requirements.
- **Konkrete Änderung:** Implement structured audit logging with tamper-evident storage.
- **Dateien:** Enhanced logger, new audit package
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Audit log integrity tests
- **Acceptance Criteria:**
  - All actions logged
  - Tamper-evident
  - Searchable format

### P3-005: SSO Integration
- **ID:** P3-005
- **Priority:** P3
- **Kategorie:** Enterprise
- **Problem:** No single sign-on support for enterprise authentication.
- **Konkrete Änderung:** Add SAML/OIDC support for enterprise authentication.
- **Dateien:** New auth package
- **Abhängigkeiten:** P3-001
- **Risiko:** Medium
- **Tests:** SSO integration tests
- **Acceptance Criteria:**
  - SAML support
  - OIDC support
  - Fallback to local auth

---

## P4 - Future/Optional Enhancements

### P4-001: GUI Frontend
- **ID:** P4-001
- **Priority:** P4
- **Kategorie:** Features
- **Problem:** CLI-only interface limits adoption by users preferring graphical tools.
- **Konkrete Änderung:** Evaluate GUI frameworks (fyne, walk, gotk3) for optional graphical interface.
- **Dateien:** New gui package
- **Abhängigkeiten:** GUI framework selection
- **Risiko:** High
- **Tests:** GUI testing framework
- **Acceptance Criteria:**
  - Cross-platform GUI
  - Feature parity with CLI
  - Accessibility support

### P4-002: Plugin System
- **ID:** P4-002
- **Priority:** P4
- **Kategorie:** Architecture
- **Problem:** Extensibility limited to source modifications.
- **Konkrete Änderung:** Design plugin architecture for third-party extensions.
- **Dateien:** New plugin package
- **Abhängigkeiten:** None
- **Risiko:** High
- **Tests:** Plugin isolation tests
- **Acceptance Criteria:**
  - Secure plugin loading
  - API stability
  - Sandboxing

### P4-003: Additional Protocols
- **ID:** P4-003
- **Priority:** P4
- **Kategorie:** Features
- **Problem:** Limited to SSH/SFTP/FTP. Missing WebDAV, SCP, RDP, VNC.
- **Konkrete Änderung:** Evaluate and implement additional protocols based on user demand.
- **Dateien:** New protocol packages
- **Abhängigkeiten:** None
- **Risiko:** Medium
- **Tests:** Protocol interoperability tests
- **Acceptance Criteria:**
  - Protocol-specific tests
  - Security reviews
  - Documentation

### P4-004: Internationalization
- **ID:** P4-004
- **Priority:** P4
- **Kategorie:** UX
- **Problem:** English/German only, hardcoded strings.
- **Konkrete Änderung:** Implement i18n framework with translation files.
- **Dateien:** All UI strings extracted
- **Abhängigkeiten:** i18n library
- **Risiko:** Low
- **Tests:** Translation completeness tests
- **Acceptance Criteria:**
  - All strings externalized
  - Multiple languages supported
  - RTL support (if needed)

### P4-005: Performance Optimization
- **ID:** P4-005
- **Priority:** P4
- **Kategorie:** Performance
- **Problem:** No benchmarks or performance profiling done.
- **Konkrete Änderung:** Add benchmarks, profile hot paths, optimize based on data.
- **Dateien:** Benchmark files
- **Abhängigkeiten:** None
- **Risiko:** Low
- **Tests:** Performance regression tests
- **Acceptance Criteria:**
  - Benchmarks established
  - Performance targets met
  - No regressions

---

## Summary Statistics

| Priority | Count | Status |
|----------|-------|--------|
| P0 | 7 | BLOCKING |
| P1 | 8 | HIGH |
| P2 | 8 | MEDIUM |
| P3 | 5 | LOW |
| P4 | 5 | OPTIONAL |

**Total Items:** 33

---

## Critical Path to Production

1. **Immediate (Week 1-2):** Complete all P0 items
2. **Beta Preparation (Week 3-4):** Complete P1 items
3. **Stable Preparation (Week 5-8):** Complete P2 items
4. **Enterprise Features (Month 3+):** Implement P3 based on customer needs
5. **Future Enhancements:** P4 items as resources allow

---

## Notes

- All security findings MUST be addressed before any public release
- Documentation must accurately reflect implemented features
- Test coverage is non-negotiable for production readiness
- Supply chain security requires immediate attention
- Platform compatibility must be verified on Windows, Linux, and macOS
