# ServerCommander

**ServerCommander** is a powerful, modular tool designed for managing and controlling local and remote servers. It supports SSH, FTP, SFTP, and more with a sleek, customizable CLI interface.

## Features

- **Multi-Protocol Management**: Supports SSH, FTP, SFTP, and more for seamless server management.
- **Cross-Platform Support**: Manage servers running on Windows, Linux, and macOS.
- **Server Dashboard**: Get real-time system information such as CPU usage, memory, and disk status.
- **File Management**: Upload, download, and manage files on remote servers securely.
- **Task Automation**: Schedule and execute tasks like backups and system maintenance.
- **User & Role Management**: Add and manage users with specific roles and permissions.
- **Security**: Two-Factor Authentication (2FA) and secure connections for enhanced safety.
- **Modular Design**: Easily extendable with additional protocols and server management tools.

## Installation

### Prerequisites

- **Go** (version 1.16 or higher) is required to build the project.
- **Git** to clone the repository.

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/vectode/ServerCommander.git
   cd ServerCommander
   ```
   
2. Build the project:

  ```bash
  go build -o server-commander
  ```

3. Run the application:

  ```bash
  ./server-commander
  ```

## Usage

Once you have the application running, you can interact with it using the following commands:

- **Connect to a server**:

  ```bash
  server-commander connect --host <server-ip> --protocol ssh
  ```
  
- **View server dashboard**:

  ```bash
  server-commander dashboard
  ```

- **Upload a file to a remote server**:

  ```bash
  server-commander upload --host <server-ip> --path /path/to/file
  ```

- **Schedule a task**:

  ```bash
  server-commander schedule --task backup --time "0 3 * * *"
  ```

For a complete list of commands and options, run:

```bash
server-commander --help
```

## Configuration

The application supports various configuration options, which can be adjusted in the ```config/``` directory. You can modify:

- **Themes**: Customize the color scheme of the CLI interface.
- **Server Settings**: Set default configurations for connecting to remote servers.

## Contributing

We welcome contributions to improve **ServerCommander**. If you have suggestions, bug fixes, or new features, feel free to submit a pull request.

### How to Contribute

1. Fork the repository.
2. Create a new branch (```git checkout -b feature/your-feature```).
3. Make your changes.
4. Commit your changes (```git commit -am 'Add new feature'```).
5. Push to the branch (```git push origin feature/your-feature```).
6. Create a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE]() file for details.

## Acknowledgments

- **Go** for providing a powerful and efficient programming language.
- **OpenSSH**, **FTP**, **SFTP** for enabling secure connections and file management.
- **CLI** for enabling a sleek command-line interface.
