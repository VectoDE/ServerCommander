# ServerCommander Documentation

## Introduction

ServerCommander is a **cross-platform** command-line tool for managing **SSH, FTP, and SFTP** connections, transferring files, and monitoring system status efficiently.

This documentation provides an in-depth look at:

- Installation & setup
- Configuration options
- Supported commands
- API reference

## Documentation Structure

| Section                                   | Description                               |
|-------------------------------------------|-------------------------------------------|
| [Installation](INSTALL.md)                | How to install and build the application  |
| [Usage](USAGE.md)                         | Detailed CLI command guide                |
| [API Reference](API.md)                   | API structure and available functions     |
| [Configuration](CONFIGURATION.md)         | How to configure settings and themes      |
| [Themes](THEMES.md)                       | Customizing the UI theme                  |
| [Folder Structure](FOLDERSTRUCTURE.md)    | Explanation of the project structure      |
| [Changelog](CHANGELOG.md)                 | Version history and new features          |

## Quick Start

To start using ServerCommander, follow these steps:

1. **Install dependencies** (Go, Git, SSH client)
2. **Build the project** using:

   ```bash
   go build -o server-commander
   ```

3. **Run the CLI tool**:

    ```bash
    ./server-commander
    ```

4. **Connect to a server**:

```bash
server-commander ssh user@host
```

For detailed command usage, refer to [USAGE](USAGE.md).

## Additional Resources

Official Repository: [GitHub Link](https://github.com/VectoDE/ServerCommander)
License: [MIT License](https://github.com/VectoDE/ServerCommander/blob/main/LICENSE)
Contact: [support@hauknetz.de](mailto:support@hauknetz.de)
