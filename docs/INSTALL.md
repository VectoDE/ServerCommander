# ServerCommander - Installation Guide

This document provides a step-by-step guide on how to install **ServerCommander** on different operating systems.

## Prerequisites

Before installing ServerCommander, ensure that you have the following dependencies installed:

- **Go 1.20+** (required for building from source) → [Download Go](https://go.dev/doc/install)
- **Git** (for cloning the repository) → [Download Git](https://git-scm.com/downloads)
- **SSH & SFTP Clients** (pre-installed on most systems)
- **FTP Client** (optional, only if using FTP functionality)
- **Make** (for Linux/macOS users)

## Installation Methods

### 1. Download Precompiled Binaries (Recommended)

The easiest way to install ServerCommander is by downloading the latest precompiled binary for your system.

#### Linux & macOS

```bash
curl -L -o server-commander https://github.com/yourrepo/servercommander/releases/latest/download/server-commander-linux
chmod +x server-commander
sudo mv server-commander /usr/local/bin/
```

#### Windows

1. Download ```server-commander.exe``` from the latest [GitHub Releases](https://github.com/VectoDE/ServerCommander/releases).
2. Move the file to a folder in your ```PATH``` (e.g., ```C:\Program Files\ServerCommander```).
3. Add the folder to the system ```PATH``` if necessary.
4. Open ```cmd``` or ```PowerShell``` and run:

    ```bash
    server-commander --version
    ```

### 2. Install via Homebrew (macOS/Linux)

If you are using Homebrew, you can install ServerCommander with:

```bash
brew install yourrepo/tap/server-commander
```

### 3. Build from Source (For Developers)

If you want to build the project from source, follow these steps:

#### Step 1: Clone the Repository

```bash
git clone https://github.com/yourrepo/servercommander.git
cd servercommander
```

#### Step 2: Build the Application

```bash
make build  # For Linux/macOS
make build-windows  # For Windows
```

#### Step 3: Run the Application

```bash
./bin/server-commander --help  # Linux/macOS
bin\server-commander.exe --help  # Windows
```

## Configuration

After installation, you may need to configure ServerCommander to match your requirements.

### 1. Default Configuration File

The default configuration is located at:

- Linux/macOS: ```~/.servercommander/config.yaml```
- Windows: ```C:\Users\<YourUser>\.servercommander\config.yaml```

To create a custom config, copy the default file:

```bash
cp config/default.yaml ~/.servercommander/config.yaml
```

### 2. Environment Variables

You can also configure the application using environment variables:

```bash
export SERVERCOMMANDER_PORT=8080
export SERVERCOMMANDER_LOG_LEVEL=debug
```

## Uninstallation

To remove **ServerCommander**, follow these steps:

### Linux/macOS

```bash
sudo rm -f /usr/local/bin/server-commander
rm -rf ~/.servercommander
```

### Windows

1. Delete ```server-commander.exe``` from its installation directory.
2. Remove the configuration folder: ```C:\Users\<YourUser>\.servercommander```

## Troubleshooting

### 1. Command Not Found

Ensure the binary is in your system ```PATH```.

```bash
echo $PATH  # Linux/macOS
$env:Path   # Windows (PowerShell)
```

### 2. Permission Denied

If you get a permission error, try:

```bash
chmod +x server-commander
sudo ./server-commander
```

### 3. Configuration Not Loading

Make sure the config file is in the correct location and formatted properly.

```bash
cat ~/.servercommander/config.yaml
```

## Need Help?

If you encounter issues, check the [documentation](DOCUMENTATION.md) or open an issue on [GitHub](https://github.com/VectoDE/ServerCommander/issues).
