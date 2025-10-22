# ServerCommander Usage Guide

## Overview

ServerCommander provides a powerful command-line interface (CLI) to manage remote servers over SSH, FTP, and SFTP, transfer files, and monitor server status. Below is a guide on how to use the tool efficiently.

## Installation

Before using ServerCommander, you need to install it. Please refer to the [Installation Guide](INSTALL.md) for details on how to set up ServerCommander on your system.

## Basic Commands

### 1. SSH Connection

To connect to a remote server via SSH:

```bash
server-commander ssh user@host
```

This command initiates an SSH session with the specified server. You will be prompted for the password or key if needed.

### 2. FTP Connection

To connect to a remote server via FTP:

```bash
server-commander ftp user@host
```

This command opens an FTP session with the server, allowing you to transfer files.

### 3. SFTP Connection

For secure file transfers via SFTP:

```bash
server-commander sftp user@host
```

This command uses the SFTP protocol for secure file transfers over the SSH connection.

### 4. List Servers

To list all the servers that are available in your configuration:

```bash
server-commander list
```

This command will show a list of configured servers, including their connection details and status.

### 5. Disconnect from Server

To disconnect from a currently active session:

```bash
server-commander disconnect
```

This command closes the active session and terminates the connection to the server.

## Advanced Commands

### 6. File Upload

To upload a file to a remote server:

```bash
server-commander upload /local/path/to/file /remote/path
```

This command transfers the specified file from your local machine to the specified remote path.

### 7. File Download

To download a file from a remote server:

```bash
server-commander download /remote/path/to/file /local/path
```

This command downloads the specified file from the remote server to your local machine.

### 8. Run a Command on the Remote Server

To run a command remotely via SSH:

```bash
server-commander run "command_to_run"
```

This command allows you to execute a single command on the remote server without opening a full interactive session.

### 9. Server Status

To view the status of a server (CPU, memory, disk usage):

```bash
server-commander status server_id
```

This command provides a real-time status report of the server's resource usage.

### 10. View Logs

To view logs from the server:

```bash
server-commander logs server_id
```

This command shows the most recent logs from the specified server.

### 11. System Monitoring with `htop`

Launch the interactive `htop` monitor using the ServerCommander color palette:

```bash
server-commander htop
```

This opens `htop` in the current terminal session. Press `q` to exit the monitor and return to the ServerCommander console.

## Additional Options

### Help

For more information about any command, use the ```--help``` flag:

```bash
server-commander <command> --help
```

This will display the detailed help for the specified command, including options and syntax.

## Conclusion

ServerCommander is designed to be simple and powerful, giving you full control over remote servers. For more advanced usage and configuration options, please refer to the [API Documentation](API.md) and [Configuration Guide](CONFIGURATION.md).
