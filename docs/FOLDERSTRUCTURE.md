# ServerCommander - Folder Structure

This document explains the folder structure of **ServerCommander** to help contributors and developers understand the project's organization.

```bash
ServerCommander/
├── bin/                    # Compiled binaries for each platform (e.g., .exe, .bin, .app)
│   ├── windows/            # Windows build
│   ├── linux/              # Linux build
│   ├── macos/              # macOS build
│   └── server-commander    # Executable after build
├── cmd/                    # Command-line application entry
│   ├── server-commander/   # Main Go application
│   │   ├── main.go         # Application entry point
│   │   ├── cli.go          # CLI command processing
│   │   ├── ssh.go          # SSH connection handler
│   │   ├── sftp.go         # SFTP handler
│   │   ├── ftp.go          # FTP handler
│   │   ├── process.go      # Process management
│   │   ├── logger.go       # Logging functionality
│   │   ├── status.go       # Server status management
│   │   ├── session.go      # Session management
│   │   ├── file_transfer.go # File transfer logic
│   │   └── utils.go        # General utilities
│   ├── version.go          # Version information
│   └── commands/           # CLI command modules
│       ├── connect.go      # Command to connect to servers
│       ├── upload.go       # Command for file upload
│       ├── download.go     # Command for file download
│       ├── logs.go         # Command for logs
│       ├── status.go       # Command for system status
│       └── process.go      # Command for process management
├── config/                 # Configuration files
│   ├── config.go           # Configuration logic
│   ├── default.yaml        # Default config file
│   └── themes/             # CLI theme configurations
│       ├── dark.yaml       # Dark mode theme
│       ├── light.yaml      # Light mode theme
│       ├── solarized.yaml  # Solarized theme
│       └── custom.yaml     # User-defined themes
├── handlers/               # Protocol handling
│   ├── sshHandler.go       # SSH connection handler
│   ├── sftpHandler.go      # SFTP connection handler
│   ├── ftpHandler.go       # FTP connection handler
│   ├── terminal.go         # Terminal color and formatting
│   ├── authentication.go   # Authentication logic
│   ├── sessionHandler.go   # SSH session management
│   └── errorHandler.go     # Error handling module
├── internal/               # Internal application logic
│   ├── storage/            # Session storage
│   │   ├── session_store.go # Session storage implementation
│   │   ├── json_store.go   # JSON-based session storage
│   │   └── db_store.go     # Database-based session storage
│   ├── process/            # Process handling
│   │   ├── process_list.go # List running processes
│   │   ├── kill.go         # Kill processes
│   │   └── monitor.go      # Monitor system performance
│   ├── file_transfer/      # File transfer operations
│   │   ├── upload.go       # Upload functionality
│   │   ├── download.go     # Download functionality
│   │   ├── progress.go     # File transfer progress bar
│   │   └── checksum.go     # File integrity check
│   ├── utils/              # Utility functions
│   │   ├── logger.go       # Logging utility
│   │   ├── colors.go       # Terminal color support
│   │   ├── paths.go        # Path handling utilities
│   │   ├── ssh_keys.go     # SSH key handling
│   │   └── encryption.go   # Encryption utilities
├── services/               # Core application services
│   ├── sshService.go       # SSH service logic
│   ├── ftpService.go       # FTP service logic
│   ├── sftpService.go      # SFTP service logic
│   ├── sessionService.go   # Session management service
│   ├── loggingService.go   # Logging service
│   ├── processService.go   # Process management service
│   └── monitoringService.go # System monitoring service
├── pkg/                    # Third-party packages and modules
│   ├── termui/             # Terminal UI library
│   ├── crypto/             # Encryption & authentication
│   ├── config/             # Config parsing (Viper)
│   ├── storage/            # Storage layer (BoltDB)
│   ├── logging/            # Advanced logging system
│   ├── ftp/                # FTP handling (go-ftp)
│   ├── sftp/               # SFTP handling (x/crypto/ssh)
│   ├── ssh/                # SSH management
│   └── utils/              # Utility libraries
├── scripts/                # Scripts for automation
│   ├── build.sh            # Build automation script
│   ├── deploy.sh           # Deployment script
│   ├── test.sh             # Run test cases
│   ├── setup.sh            # Environment setup
│   └── cleanup.sh          # Cleanup script
├── tests/                  # Unit & Integration Tests
│   ├── ssh_test.go         # SSH tests
│   ├── ftp_test.go         # FTP tests
│   ├── sftp_test.go        # SFTP tests
│   ├── process_test.go     # Process management tests
│   ├── file_transfer_test.go # File transfer tests
│   ├── session_test.go     # Session handling tests
│   ├── utils_test.go       # Utility function tests
│   └── security_test.go    # Security and authentication tests
├── docs/                   # Documentation
│   ├── README.md           # Main project documentation
│   ├── INSTALL.md          # Installation guide
│   ├── USAGE.md            # Usage guide
│   ├── API.md              # API documentation
│   ├── CONFIGURATION.md    # Configuration details
│   ├── THEMES.md           # Theme customization
│   ├── CONTRIBUTING.md     # Contribution guidelines
│   ├── FOLDERSTRUCTURE.md  # Detailed explanation of the folder structure
│   └── CHANGELOG.md        # Version history and changes
├── .gitignore              # Ignore files for Git
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
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
