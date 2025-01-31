# ServerCommander - Configuration Guide

This document provides details on how to configure ServerCommander to suit your needs. It includes explanations of configuration files, available options, and customization settings.

## Configuration File Location

By default, ServerCommander uses a YAML configuration file located at:
    ```bash
    ~/.servercommander/config.yaml  # Linux & macOS
    %USERPROFILE%\.servercommander\config.yaml  # Windows
    ```

Alternatively, you can specify a custom configuration file using the ```--config``` flag:
    ```bash
    server-commander --config /path/to/custom-config.yaml
    ```

## Default Configuration File Structure

```bash
server:
  host: 0.0.0.0
  port: 22
  default_protocol: ssh
  timeout: 30  # Timeout for connections in seconds

authentication:
  use_key_auth: true  # Enable SSH key authentication
  private_key_path: ~/.ssh/id_rsa  # Path to private key
  allow_passwords: false  # Allow password authentication

theme:
  color_scheme: dark  # Available options: dark, light, solarized, custom
  custom_theme_path: ~/.servercommander/themes/custom.yaml

logging:
  enable: true
  log_file: ~/.servercommander/servercommander.log
  level: info  # Available options: debug, info, warn, error

session:
  save_sessions: true  # Save session history
  session_timeout: 600  # Timeout before session closes (in seconds)
  max_sessions: 5  # Maximum concurrent sessions

ftp:
  enable: true
  default_port: 21
  passive_mode: true
  max_transfer_speed: 1024  # KB/s

sftp:
  enable: true
  default_port: 22
  max_transfer_speed: 1024  # KB/s

updates:
  auto_check: true  # Automatically check for updates
  notify: true  # Notify user when updates are available
```

## Detailed Configuration Options

### 1. Server Settings (```server:```)

- ```host```: Defines the server address (default: ```0.0.0.0```).
- ```port```: Sets the port for SSH and SFTP connections (default: ```22```).
- ```default_protocol```: Chooses between ```ssh```, ```sftp```, or ```ftp``` (default: ```ssh```).
- ```timeout```: Specifies connection timeout in seconds (default: ```30```).

### 2. Authentication (```authentication:```)

- ```use_key_auth```: Enables SSH key authentication (default: ```true```).
- ```private_key_path```: Defines the private SSH key location.
- ```allow_passwords```: Allows password authentication if set to ```true``` (default: ```false```).

### 3. Theme (```theme:```)

- ```color_scheme```: Available options: ```dark```, ```light```, ```solarized```, ```custom```.
- ```custom_theme_path```: Path to a custom YAML theme file.

### 4. Logging (```logging:```)

- ```enable```: Enables logging (default: ```true```).
- ```log_file```: Specifies log file location.
- ```level```: Log verbosity: ```debug```, ```info```, ```warn```, ```error```.

### 5. Session Management (```session:```)

- ```save_sessions```: Saves previous session history.
- ```session_timeout```: Timeout before session closes (in seconds).
- ```max_sessions```: Maximum concurrent open sessions.

### 6. FTP Configuration (```ftp:```)

- ```enable```: Enables FTP functionality (default: ```true```).
- ```default_port```: Sets default FTP port (default: ```21```).
- ```passive_mode```: Enables FTP passive mode.
- ```max_transfer_speed```: Limits transfer speed (default: ```1024 KB/s```).

### 7. SFTP Configuration (```sftp:```)

- ```enable```: Enables SFTP functionality (default: ```true```).
- ```default_port```: Sets default SFTP port (default: ```22```).
- ```max_transfer_speed```: Limits transfer speed (default: ```1024 KB/s```).

### 8. Updates (updates:)

- ```auto_check```: Automatically checks for software updates.
- ```notify```: Notifies the user when updates are available.

## Modifying Configuration

### 1. Editing the Config File Manually

Edit the ```config.yaml``` file using any text editor:
    ```bash
    nano ~/.servercommander/config.yaml
    ```

### 2. Using Command-Line Flags

Some options can be overridden using CLI flags:
    ```bash
    server-commander --port 2222 --default_protocol ftp
    ```

### 3. Using Environment Variables

You can override settings using environment variables:
    ```bash
    export SERVERCOMMANDER_PORT=2222
    export SERVERCOMMANDER_DEFAULT_PROTOCOL=ftp
    ```

## Resetting Configuration

To reset the configuration to default values:
    ```bash
    server-commander --reset-config
    ```

This will restore the default ```config.yaml``` file.

## Conclusion

This guide covers the essential configuration settings for ServerCommander. Adjust these parameters to optimize security, logging, sessions, and transfer protocols according to your needs. For more information, refer to the [API documentation](API.md) or check the [Contributing Guide](CONTRIBUTING.md).
