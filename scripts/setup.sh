#!/bin/bash

# This script sets up the development environment for ServerCommander.

echo "Setting up ServerCommander development environment..."

# Install dependencies
go mod tidy

# Install any additional system dependencies (e.g., SSH, FTP, SFTP)
# This section can be adjusted depending on the OS you're using.
if [ "$(uname)" == "Darwin" ]; then
  brew install sshpass
elif [ "$(uname)" == "Linux" ]; then
  sudo apt-get update
  sudo apt-get install -y sshpass
fi

echo "Development environment setup complete."
