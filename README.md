# ServerCommander – The Ultimate SSH & FTP/SFTP Client

## Overview

ServerCommander is a powerful terminal-based client for **SSH**, **FTP**, and **SFTP**, enabling secure and efficient management of **local and remote servers**.

It features a **colorized console**, **stable connections**, and **practical file transfer and server management functionalities** – all packaged as a **standalone executable** for **Windows**, **Linux**, and **macOS**.

## Features

| Feature                   | Description                                                            |
|---------------------------|------------------------------------------------------------------------|
| SSH Support               | Secure connections to remote servers using SSH.                        |
| SFTP Support              | Secure file transfer using SSH File Transfer Protocol.                 |
| FTP Support               | Support for unencrypted and encrypted FTP connections.                 |
| Session Management        | Save and restore server connections.                                   |
| Local & Remote Management | Manage both local systems (Windows, Linux, macOS) and remote servers.  |
| Colorized Terminal        | Improved readability with syntax highlighting and customizable themes. |
| Keyboard Shortcuts        | Efficient navigation with configurable key bindings.                   |
| Authentication            | Supports passowrd and SSH key authentication.                          |
| Logging & Debugging       | Detailed session logs for analysis and troubleshooting.                |
| Cross-Platform            | Runs as .exe (Windows), .bin (Linux) and .app (macOS).                 |

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

### Connecting to Servers

- **Start an SSH connection**:

  ```bash
  server-commander ssh user@host
  ```

- **Open an SFTP session**:

  ```bash
  server-commander sftp user@host
  ```

- **Start an FTP session**:

  ```bash
  server-commander ftp user@host
  ```

### File Transfers

- **Upload a file**:

  ```bash
  server-commander upload /local/file /remote/path
  ```

- **Download a file**:

  ```bash
  server-commander download /remote/file /local/path
  ```

### Retrieving System Information

- **Display server status**:

  ```bash
  server-commander status
  ```

- **Manage running processes**:

  ```bash
  server-commander process-list
  ```

### Logging & Debugging

- **View logs**:

  ```bash
  server-commander logs
  ```

## Architecture & Implementation

- **Programming Language**: Go
- **SSH/SFTP**: ```golang.org/x/crypto/ssh```
- **FTP**: ```github.com/jlaffaye/ftp```
- **Terminal Handling**: ```github.com/muesli/termenv```
- **Config Handling**: ```viper``` for session storage

## Conclusion

**ServerCommander** is a **alternative to PuTTY**, but with **integrated FTP/SFTP support** and packaged as a **standalone executable** for Windows, Linux, and macOS.
It is perfect for **system administrators and developers** who need a **fast and efficient server management tool**.
