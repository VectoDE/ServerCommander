# ServerCommander Architecture Audit Report

**Document Version:** 1.0  
**Audit Date:** 2025-08-30  
**Status:** MAJOR ARCHITECTURAL ISSUES IDENTIFIED

---

## Executive Summary

ServerCommander's current architecture is **minimal and functional** but lacks the structure, patterns, and safeguards necessary for enterprise-grade software. The codebase shows signs of early-stage development with significant technical debt in error handling, testing infrastructure, and architectural boundaries.

### Overall Architecture Rating: **35/100** (Needs Significant Improvement)

| Aspect | Score | Status |
|--------|-------|--------|
| Code Organization | 50/100 | Acceptable |
| Error Handling | 30/100 | Poor |
| Resource Management | 35/100 | Poor |
| Concurrency Safety | NOT VERIFIED | Unknown |
| Testability | 20/100 | Critical |
| Modularity | 40/100 | Fair |
| Documentation | 25/100 | Poor |

---

## Current Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      main.go                                 │
│                    (Entry Point)                             │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   console.Run()                              │
│              (Interactive Console Loop)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                  cmd.Execute()                               │
│            (Command Dispatcher Pattern)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   session    │ │     ssh      │ │     ftp      │
│   Command    │ │   Command    │ │   Command    │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       ▼                ▼                ▼
┌──────────────────────────────────────────────────────────────┐
│                    Services Layer                             │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐              │
│  │   config   │  │    ssh     │  │    ftp     │              │
│  │  Service   │  │  Service   │  │  Service   │              │
│  └────────────┘  └────────────┘  └────────────┘              │
└──────────────────────────────────────────────────────────────┘
```

### Package Structure

```
servercommander/
├── src/
│   ├── main.go                 # Application entry point
│   ├── cmd/                    # Command implementations
│   │   ├── dispatcher.go       # Command registry & execution
│   │   ├── session.go          # Session management commands
│   │   ├── ssh.go              # SSH connection commands
│   │   ├── sftp.go             # SFTP file operations
│   │   ├── ftp.go              # FTP file operations
│   │   ├── help.go             # Help command
│   │   ├── clear.go            # Console clear command
│   │   ├── exit.go             # Exit command
│   │   ├── htop.go             # System monitor
│   │   └── htop_theme.go       # Embedded htop theme
│   ├── console/                # Console UI layer
│   │   └── console.go          # Interactive loop, banner
│   ├── services/               # Business logic
│   │   ├── logger.go           # Logging service
│   │   ├── config/             # Configuration management
│   │   │   ├── sessions.go     # Session CRUD
│   │   │   └── paths.go        # Path resolution
│   │   ├── ssh/                # SSH client wrapper
│   │   │   └── client.go       # System ssh delegation
│   │   └── ftp/                # FTP client implementation
│   │       └── client.go       # Native FTP/FTPS client
│   ├── utils/                  # Utility functions
│   │   ├── colors.go           # ANSI color codes
│   │   ├── prompt.go           # User input helpers
│   │   ├── fileExists.go       # File existence check
│   │   └── usage.go            # Usage formatting
│   └── assets/                 # Embedded resources
│       ├── icon.ico            # Windows icon
│       ├── resource.syso       # Windows resource
│       └── goodbye.mp3         # (Unused?) audio file
├── docs/                       # Documentation
├── scripts/                    # Build scripts
├── build/                      # Compiled binaries
├── go.mod                      # Go module definition
└── LICENSE                     # MIT License
```

---

## Architectural Analysis

### 1. Code Organization

**Strengths:**
- Clear separation between CLI layer (`cmd/`), services (`services/`), and utilities (`utils/`)
- Command pattern implemented for extensibility
- Consistent naming conventions
- Minimal circular dependencies

**Weaknesses:**
- No interfaces defined for services (hard to mock/test)
- Tight coupling between layers
- `utils/` is a catch-all package
- No clear domain model
- Mixed responsibilities in some files

**Recommendation:**
```go
// Define interfaces for testability
type SessionStore interface {
    Get(alias string) (Session, bool)
    List() []Session
    Upsert(session Session) Session
    Remove(alias string) error
    Save() error
}

type SSHClient interface {
    Connect() error
    InteractiveShell() error
    Run(command string) (string, error)
    Close() error
}

type FTPClient interface {
    Connect() error
    Upload(local, remote string) error
    Download(remote, local string) error
    List(path string) ([]Entry, error)
    Close() error
}
```

### 2. Error Handling

**Current State: POOR**

**Issues Identified:**

1. **Inconsistent Error Wrapping:**
```go
// Some places use fmt.Errorf with %w
return nil, fmt.Errorf("unable to stat sessions file: %w", err)

// Others just return bare errors
return fmt.Errorf("unknown ssh action '%s'", action)
```

2. **Error Information Loss:**
```go
// Original error context often lost
if err != nil {
    return nil, err  // No context added
}
```

3. **No Error Types:**
- Cannot programmatically handle different error cases
- No retryable error detection
- No user-friendly vs technical error distinction

**Recommendation:**
```go
// Define error types
var (
    ErrSessionNotFound = errors.New("session not found")
    ErrConnectionFailed = errors.New("connection failed")
    ErrAuthenticationRequired = errors.New("authentication required")
)

// Use consistent wrapping
if err != nil {
    return fmt.Errorf("failed to connect to %s: %w", session.Host, err)
}

// Check for specific errors
if errors.Is(err, ErrSessionNotFound) {
    // Handle not found
}
```

### 3. Resource Management

**Current State: POOR**

**Issues:**

1. **Deferred Cleanup Inconsistent:**
```go
// Good: Proper defer
file, err := os.Open(localPath)
if err != nil {
    return err
}
defer file.Close()

// Bad: Manual cleanup only on success path
dataConn, err := c.openDataConnection(command)
if err != nil {
    return err
}
defer dataConn.Close()  // This IS present, but...

// ...error paths may leak
if err := someOperation(); err != nil {
    return err  // dataConn closed by defer, OK
}
```

2. **Temp File Handling:**
```go
// In sftp.go - temp file created
file, err := os.CreateTemp("", "sftp-batch-*.txt")
// ...
cleanup := func() {
    os.Remove(file.Name())
}
defer cleanup()

// But what if cleanup fails? No error reported.
```

3. **No Connection Pooling:**
- New connection for every operation
- No keepalive between commands
- Performance impact

**Recommendation:**
```go
// Use context for cancellation and timeouts
func (c *Client) Upload(ctx context.Context, localPath, remotePath string) error {
    // Check context before starting
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // Use context-aware operations
    file, err := os.Open(localPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // ... rest of operation
}
```

### 4. Concurrency Safety

**Status: NOT VERIFIED**

**Observations:**

1. **Prompt Reader Mutex:**
```go
var (
    readerMu sync.RWMutex
    reader   = bufio.NewReader(os.Stdin)
)

func SetPromptReader(r io.Reader) {
    readerMu.Lock()
    reader = bufio.NewReader(r)
    readerMu.Unlock()
}

func readLine() (string, error) {
    readerMu.RLock()  // RLock for read
    active := reader
    readerMu.RUnlock()
    return active.ReadString('\n')
}
```

**Potential Issue:** RLock during read, but ReadString may block indefinitely, holding RLock.

2. **No Goroutines Detected:**
- Entirely single-threaded
- No parallel operations
- Safe but limits performance

3. **No Race Condition Testing:**
```bash
go test -race ./...
# Never run
```

**Recommendation:**
- Run race detector on all tests
- Document thread-safety guarantees
- Consider worker pools for file transfers

### 5. Testability

**Current State: CRITICAL**

**Issues:**

1. **No Interfaces:**
- Cannot mock external dependencies
- Cannot test without real filesystem/network
- All functions call real implementations

2. **Global State:**
```go
var commandRegistry = map[string]CommandDescriptor{}

func RegisterCommand(name, description string, handler CommandHandler) {
    // Modifies global state
    commandRegistry[key] = ...
}
```

3. **Hard Dependencies:**
- Direct calls to `os.Stdin`, `os.Stdout`
- Direct file system access
- Direct network calls

**Recommendation:**
```go
// Dependency injection
type CommandContext struct {
    Stdin  io.Reader
    Stdout io.Writer
    Stderr io.Writer
    Config ConfigService
    Sessions SessionStore
}

func Execute(ctx *CommandContext, input string) error {
    // Now testable with mocked context
}
```

### 6. Modularity

**Current State: FAIR**

**Strengths:**
- Clear package boundaries
- Limited cross-package dependencies
- Single responsibility per package (mostly)

**Weaknesses:**
- No plugin architecture
- Cannot extend without modifying source
- No versioned APIs between packages
- Tightly coupled to CLI interface

**Recommendation:**
Consider plugin architecture for enterprise features:
```go
type Plugin interface {
    Name() string
    Version() string
    RegisterCommands(registry CommandRegistry) error
    Initialize(config Config) error
    Shutdown() error
}
```

---

## Data Flow Analysis

### Session Management Flow

```
User Input → cmd.sessionCommand → config.LoadSessions() 
                                      ↓
                                 Read sessions.json
                                      ↓
                                 Parse JSON
                                      ↓
                                 Return SessionStore
                                      ↓
cmd.sessionAdd → store.Upsert() → store.Save() → Write sessions.json
```

**Issues:**
- No atomic writes (corruption risk)
- No file locking (concurrent access risk)
- No schema versioning (migration risk)

### SSH Connection Flow

```
User Input → cmd.connectCommand → promptPassword() 
                                     ↓
                            sshservice.Connect()
                                     ↓
                            buildBaseArgs()
                                     ↓
                            exec.Command("ssh", args...)
                                     ↓
                            cmd.Run() (blocks)
```

**Issues:**
- No timeout on connection
- No host key verification
- Password echoed in terminal
- No error recovery

### FTP Transfer Flow

```
User Input → cmd.ftpCommand → withFTPClient()
                                   ↓
                          ftpservice.Connect()
                                   ↓
                          net.DialTimeout()
                                   ↓
                          TLS handshake (INSECURE)
                                   ↓
                          Authentication
                                   ↓
                          PASV mode setup
                                   ↓
                          Data transfer
                                   ↓
                          client.Close()
```

**Issues:**
- TLS validation DISABLED
- No resume support
- No progress reporting
- No checksum verification

---

## Trust Boundaries

```
┌─────────────────────────────────────────────────────────────┐
│                     TRUSTED ZONE                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              ServerCommander Process                 │   │
│  │                                                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐          │   │
│  │  │  Memory  │  │  Config  │  │   Logs   │          │   │
│  │  │  (Secrets)│  │  Files   │  │  Files   │          │   │
│  │  └──────────┘  └──────────┘  └──────────┘          │   │
│  └─────────────────────────────────────────────────────┘   │
│                          │                                  │
│              ═══════════════════════                        │
│              ║  TRUST BOUNDARY  ║                           │
│              ═══════════════════════                        │
│                          │                                  │
└──────────────────────────┼──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│  Network      │  │  Filesystem   │  │  External     │
│  (UNTRUSTED)  │  │  (PARTIAL)    │  │  Commands     │
│               │  │               │  │  (PARTIAL)    │
│ • MITM Risk   │  │ • Path Traversal│  │ • Injection  │
│ • Bad Certs   │  │ • Permissions │  │ • Args       │
│ • Bad Hosts   │  │ • Symlinks    │  │ • Exit Codes │
└───────────────┘  └───────────────┘  └───────────────┘
```

---

## Technical Debt Inventory

### Critical Debt

| ID | Area | Description | Impact | Effort |
|----|------|-------------|--------|--------|
| TD-001 | Security | No TLS validation | CRITICAL | Low |
| TD-002 | Security | No SSH host key check | CRITICAL | Low |
| TD-003 | Security | Password echo | HIGH | Low |
| TD-004 | Testing | Zero tests | HIGH | High |
| TD-005 | Supply Chain | No go.sum | CRITICAL | Low |

### High Priority Debt

| ID | Area | Description | Impact | Effort |
|----|------|-------------|--------|--------|
| TD-006 | Architecture | No interfaces | HIGH | Medium |
| TD-007 | Error Handling | Inconsistent wrapping | MEDIUM | Low |
| TD-008 | Resources | Incomplete cleanup | MEDIUM | Low |
| TD-009 | Config | No atomic writes | MEDIUM | Low |
| TD-010 | Docs | Outdated documentation | MEDIUM | Medium |

### Medium Priority Debt

| ID | Area | Description | Impact | Effort |
|----|------|-------------|--------|--------|
| TD-011 | Performance | No connection pooling | LOW | Medium |
| TD-012 | UX | Basic terminal handling | LOW | High |
| TD-013 | Features | Missing import/export | LOW | Low |
| TD-014 | Platform | Windows compatibility gaps | MEDIUM | Medium |
| TD-015 | Resilience | No signal handling | LOW | Low |

---

## Target Architecture

### Proposed Layered Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │     CLI     │  │     GUI     │  │     API     │         │
│  │  (Current)  │  │  (Future)   │  │  (Future)   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Command Handlers                        │   │
│  │         (Use Cases / Application Services)           │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Session  │  │Transfer  │  │   Auth   │  │  Server  │   │
│  │  Entity  │  │  Entity  │  │  Entity  │  │  Entity  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Domain Services                          │  │
│  │         (Business Logic, Validation)                  │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   SSH    │  │   FTP    │  │  Config  │  │  Logger  │   │
│  │ Adapter  │  │ Adapter  │  │ Adapter  │  │ Adapter  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Key Architectural Improvements

1. **Dependency Injection:**
```go
type Application struct {
    sessionStore SessionStore
    sshFactory   SSHClientFactory
    ftpFactory   FTPClientFactory
    config       Config
    logger       Logger
}

func NewApplication(deps Dependencies) *Application {
    return &Application{...}
}
```

2. **Interface-Based Design:**
```go
type TransferService interface {
    Upload(ctx context.Context, spec TransferSpec) error
    Download(ctx context.Context, spec TransferSpec) error
    List(ctx context.Context, path string) ([]FileEntry, error)
}
```

3. **Event-Driven Architecture (Future):**
```go
type EventBus interface {
    Subscribe(eventType string, handler EventHandler)
    Publish(event Event)
}

// Events: SessionCreated, TransferStarted, TransferCompleted, ConnectionFailed
```

---

## Migration Strategy

### Phase 1: Foundation (Weeks 1-2)
- [ ] Add interfaces for all services
- [ ] Implement dependency injection
- [ ] Fix critical security issues
- [ ] Add basic tests

### Phase 2: Hardening (Weeks 3-4)
- [ ] Improve error handling
- [ ] Add context propagation
- [ ] Implement resource cleanup
- [ ] Add integration tests

### Phase 3: Enhancement (Weeks 5-8)
- [ ] Refactor to layered architecture
- [ ] Add configuration file support
- [ ] Implement session encryption
- [ ] Performance optimization

### Phase 4: Enterprise (Months 3-6)
- [ ] Plugin architecture
- [ ] API layer
- [ ] GUI frontend (optional)
- [ ] Advanced features

---

## Conclusion

ServerCommander's architecture is **functional but fragile**. The current implementation works for basic use cases but lacks the robustness, testability, and extensibility required for enterprise software.

**Key Recommendations:**

1. **Immediate:** Fix P0 security findings (TLS, SSH, passwords)
2. **Short-term:** Add interfaces and dependency injection
3. **Medium-term:** Refactor to layered architecture
4. **Long-term:** Consider plugin system and API layer

**Risk Assessment:**
- **Current:** High risk of bugs, security issues, difficult maintenance
- **After Phase 1:** Medium risk, testable, more maintainable
- **After Phase 3:** Low risk, robust, extensible

---

**END OF ARCHITECTURE AUDIT REPORT**
