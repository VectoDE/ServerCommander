# ServerCommander - Folder Structure

This document explains the folder structure of **ServerCommander** to help contributors and developers understand the project's organization.

```bash
ServerCommander/
├── bin/                    # Compiled binaries for each platform
│   ├── windows/            # Windows build
│   ├── linux/              # Linux build
│   ├── macos/              # macOS build
│   └── server-commander    # Executable after build
├── src/                    # Source code directory
│   ├── main.go             # Application entry point (main Go file)
│   ├── cmd/                # Command-line application entry
│   │   ├── server-commander/   # Main Go application logic
│   │   │   ├── cli.go          # CLI command processing
│   │   │   ├── server.go       # Handles server connection logic (SSH, FTP, SFTP)
│   │   │   ├── process.go      # Process management (start, stop, list)
│   │   │   ├── file.go         # File operations (upload, download, transfer)
│   │   │   └── utils.go        # General utilities
│   │   └── version.go           # Version information
│   ├── config/               # Configuration files
│   │   ├── config.go           # Configuration loading and parsing logic
│   │   └── themes/             # Theme configurations
│   │       ├── dark.yaml       # Dark mode theme
│   │       └── light.yaml      # Light mode theme
│   ├── internal/               # Internal application logic
│   │   ├── storage/            # Session storage (in-memory, JSON, DB)
│   │   │   ├── memory_store.go # In-memory session storage
│   │   │   ├── json_store.go   # JSON-based session storage
│   │   │   └── db_store.go     # Database-based session storage
│   │   ├── service/            # Core services (file transfers, SSH sessions)
│   │   │   ├── ssh_service.go  # SSH service logic
│   │   │   ├── ftp_service.go  # FTP service logic
│   │   │   ├── sftp_service.go # SFTP service logic
│   │   │   ├── session_service.go # Session management service
│   │   │   ├── process_service.go # Process management service
│   │   │   └── logging_service.go # Logging service
│   │   └── utils/              # Utility functions (logging, encryption)
│   │       ├── logger.go       # Logger utility
│   │       ├── encryption.go   # Encryption utilities
│   │       ├── colors.go       # Terminal color utilities
│   │       ├── paths.go        # Path handling utilities
│   │       └── ssh_keys.go     # SSH key handling utilities
│   ├── pkg/                    # Third-party packages and modules
│   │   ├── termui/             # Terminal UI library
│   │   │   ├── termui.go       # Terminal UI handling
│   │   ├── crypto/             # Encryption & authentication
│   │   │   ├── encryption.go   # Encryption logic
│   │   ├── ssh/                # SSH connection handling
│   │   │   ├── ssh.go          # SSH connection handling
│   │   ├── ftp/                # FTP handling
│   │   │   ├── ftp.go          # FTP connection handling
│   │   ├── sftp/               # SFTP handling
│   │   │   ├── sftp.go         # SFTP connection handling
│   │   └── logging/            # Logging system
│   │       ├── logger.go       # Logging utility
│   └── scripts/                # Automation scripts
│       ├── build.sh            # Build automation script
│       ├── deploy.sh           # Deployment script
│       ├── setup.sh            # Environment setup
│       └── cleanup.sh          # Cleanup script
├── tests/                  # Unit & Integration Tests
│   ├── server_test.go      # Tests for server connections (SSH, FTP, SFTP)
│   ├── process_test.go     # Tests for process management (start, stop, list)
│   ├── file_test.go        # Tests for file transfers (upload, download)
│   ├── session_test.go     # Tests for session handling
│   ├── utils_test.go       # Tests for utility functions
│   └── logging_test.go     # Tests for logging functionality
├── docs/                   # Documentation
│   ├── README.md           # Main project documentation
│   ├── INSTALL.md          # Installation guide
│   ├── USAGE.md            # Usage guide
│   ├── API.md              # API documentation
│   ├── CONFIGURATION.md    # Configuration details
│   ├── THEMES.md           # Theme customization guide
│   ├── CONTRIBUTING.md     # Contribution guidelines
│   ├── FOLDERSTRUCTURE.md  # Folder structure explanation
│   └── CHANGELOG.md        # Version history and changes
├── .gitignore              # Git ignore settings
├── go.mod                  # Go module definition
├── LICENSE                 # Project license (MIT)
├── Makefile                # Build automation for all OS
└── README.md               # Main project documentation
```

## Explanation of Key Directories

### **bin/**

- Contains compiled binary files for different operating systems.

### **cmd/**

- Houses the **main entry point** of the application, including command parsing.

### **config/**

- Stores configuration settings such as **default.yaml** and CLI themes.

### **handlers/**

- Manages different protocol **(SSH, FTP, SFTP)** functionalities.

### **internal/**

- Contains internal utilities, including **encryption, validation, and system monitoring**.

### **services/**

- Core services like **file transfers, SSH execution, and logging**.

### **utils/**

- Utility functions for **logging, error handling, and CLI interaction**.

### **tests/**

- Contains unit and integration tests for **all major features**.

### **scripts/**

- Automation scripts for **installation, building, and cleanup**.

### **docs/**

- Documentation related to **installation, usage, and contribution**.

## How to Navigate the Project

- If you want to **modify the core application**, check `cmd/` and `services/`.
- If you want to **extend SSH or FTP functionality**, look into `handlers/`.
- If you want to **change configuration or themes**, check `config/`.
- If you want to **contribute to testing**, modify `tests/`.
- If you want to **understand the folder structure**, refer to `docs/FOLDERSTRUCTURE.md` (this file).

## Contribution

If you want to contribute, check out [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
