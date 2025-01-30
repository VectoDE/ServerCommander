ServerCommander – The Ultimate SSH & FTP/SFTP Client

1. Overview

ServerCommander is a powerful terminal-based client for SSH, FTP, and SFTP, enabling secure and efficient management of local and remote servers.

It features a colorized console, stable connections, and practical file transfer and server management functionalities – all packaged as a standalone executable for Windows, Linux, and macOS.

2. Features

Feature	Description
SSH Support	Secure connections to remote servers using SSH.
SFTP Support	Secure file transfers using SSH File Transfer Protocol.
FTP Support	Support for unencrypted and encrypted FTP connections.
Session Management	Save and restore server connections.
Local & Remote Management	Manage both local systems (Windows, Linux, macOS) and remote servers.
Colorized Terminal	Improved readability with syntax highlighting and customizable themes.
Keyboard Shortcuts	Efficient navigation with configurable key bindings.
Authentication	Supports password and SSH key authentication.
Logging & Debugging	Detailed session logs for analysis and troubleshooting.
Cross-Platform	Runs as .exe (Windows), .bin (Linux), and .app (macOS).

3. Installation & Build

3.1 Prerequisites
	•	Go (version 1.16 or higher)
	•	Git
	•	SSH client, FTP/SFTP client (optional)

3.2 Build Process

ServerCommander can be compiled for each operating system:

Windows (.exe)

GOOS=windows GOARCH=amd64 go build -o server-commander.exe

Linux (.bin)

GOOS=linux GOARCH=amd64 go build -o server-commander.bin

macOS (.app / .bin)

GOOS=darwin GOARCH=amd64 go build -o server-commander.app

After building, run the executable:

./server-commander

4. Usage

4.1 Connecting to Servers
	•	Start an SSH connection:

server-commander ssh user@host


	•	Open an SFTP session:

server-commander sftp user@host


	•	Start an FTP session:

server-commander ftp user@host



4.2 File Transfers
	•	Upload a file:

server-commander upload /local/file /remote/path


	•	Download a file:

server-commander download /remote/file /local/path



4.3 Retrieving System Information
	•	Display server status:

server-commander status


	•	Manage running processes:

server-commander process-list



4.4 Logging & Debugging
	•	View logs:

server-commander logs

5. Architecture & Implementation
	•	Programming Language: Go
	•	SSH/SFTP: golang.org/x/crypto/ssh
	•	FTP: github.com/jlaffaye/ftp
	•	Terminal Handling: github.com/muesli/termenv
	•	Config Handling: viper for session storage

6. Conclusion

ServerCommander is a 1:1 alternative to PuTTY, but with integrated FTP/SFTP support and packaged as a standalone executable for Windows, Linux, and macOS.
It is perfect for system administrators and developers who need a fast and efficient server management tool.