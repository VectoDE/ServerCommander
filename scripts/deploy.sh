#!/bin/bash

# This script deploys ServerCommander to a server.

# Check if the necessary argument is provided
if [ -z "$1" ]; then
  echo "Usage: ./deploy.sh <server_ip>"
  exit 1
fi

SERVER_IP=$1

echo "Deploying to server: $SERVER_IP..."

# Copy the binary to the remote server using SCP
scp bin/server-commander user@$SERVER_IP:/usr/local/bin/

# Restart the ServerCommander service on the remote server (optional)
ssh user@$SERVER_IP "systemctl restart servercommander"

echo "Deployment complete."
