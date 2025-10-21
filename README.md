# ServerCommander – The Ultimate SSH & FTP/SFTP Client

## Overview

ServerCommander is a powerful terminal-based client for **SSH**, **FTP**, and **SFTP**, enabling secure and efficient management of **local and remote servers**.

It features a **colorized console**, **stable connections**, and **practical file transfer and server management functionalities** – all packaged as a **standalone executable** for **Windows**, **Linux**, and **macOS**.

## Features

| Feature                   | Description                                                                                 |
|---------------------------|---------------------------------------------------------------------------------------------|
| Interactive shell         | Built-in command dispatcher with contextual help and graceful shutdown handling.           |
| Session management        | Persist SSH/SFTP/FTP endpoints in the user configuration directory for effortless reuse.   |
| SSH connectivity          | Launch interactive shells or run one-off commands through the system `ssh` client.         |
| SFTP file operations      | Upload, download, and list remote files via the `sftp` client with friendly output.        |
| FTP/FTPS support          | Native passive-mode FTP client with optional explicit TLS for secure transfers.            |
| Secure authentication     | Supports password and private-key based authentication (keys are never stored).            |
| Structured logging        | Execution traces are written to `~/.config/servercommander/logs/servercommander.log`.      |
| Cross-platform            | Designed for Windows, macOS, and Linux with no CGO dependencies.                           |

## Installation & Build

### Prerequisites

- Go (version 1.16 or higher)
- Git
- SSH client, FTP/SFTP client (optional)

### Build Process

ServerCommander can be compiled for each operating system:

### Windows (.exe)

```bash
go build -o "build/ServerCommander.exe" ./src/main.go
```

### Linux (.bin)

```bash
go build -o "build/ServerCommander" ./src/main.go
```

### macOS (.app / .bin)

```bash
go build -o "build/ServerCommander.app" ./src/main.go
```

After building, run the executable:

```bash
./server-commander
```

## Usage

Run the binary and use the interactive prompt:

```bash
./server-commander
```

Key commands inside the shell:

| Command                               | Description                                                                 |
|---------------------------------------|-----------------------------------------------------------------------------|
| `session add <alias>`                 | Create or update a stored session (you will be guided through the fields).  |
| `session list`                        | Display all saved sessions.                                                 |
| `session show <alias>`                | Show session details.                                                       |
| `session remove <alias>`              | Delete a stored session.                                                    |
| `connect <alias>`                     | Start an interactive SSH shell for the given session.                       |
| `ssh exec <alias> <command>`          | Execute a single command via SSH and print the output.                      |
| `sftp list <alias> [remote-path]`     | List remote files using SFTP.                                               |
| `sftp upload <alias> <local> <remote>`| Upload a file via SFTP.                                                     |
| `sftp download <alias> <remote> <local>`| Download a file via SFTP.                                                |
| `ftp list <alias> [remote-path]`      | List remote files using the built-in FTP client.                            |
| `ftp upload <alias> <local> <remote>` | Upload a file via FTP/FTPS.                                                 |
| `ftp download <alias> <remote> <local>`| Download a file via FTP/FTPS.                                             |
| `help`                                | Print the command catalogue.                                                |
| `clear`                               | Clear the terminal and reprint the banner.                                  |
| `exit`                                | Gracefully shut down ServerCommander.                                       |

### Requirements

- Go 1.20+ for building (the module targets Go 1.23).
- OpenSSH client tools (`ssh` and `sftp`) available in `PATH` for SSH/SFTP features.
- Remote FTP servers must allow passive connections when using the FTP client.
- Password prompts in the interactive shell are echoed. Prefer SSH keys or run the system tools directly if hidden input is required.

## Architecture & Implementation

- **Programming Language**: Go
- **Command layer**: Custom dispatcher with pluggable command registration.
- **Session storage**: JSON file stored under `~/.config/servercommander/sessions.json`.
- **SSH/SFTP**: Delegates to the system OpenSSH tooling for reliability and feature parity.
- **FTP/FTPS**: In-house passive-mode client supporting optional explicit TLS.
- **Logging**: Structured text logs in the user configuration directory.

## Conclusion

**ServerCommander** is a **alternative to PuTTY**, but with **integrated FTP/SFTP support** and packaged as a **standalone executable** for Windows, Linux, and macOS.
It is perfect for **system administrators and developers** who need a **fast and efficient server management tool**.
