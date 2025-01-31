#!/bin/bash

# This script builds the ServerCommander binary.

echo "Building ServerCommander..."

# Run Go build to create the binary
go build -o bin/server-commander .

# Check if the build was successful
if [ $? -eq 0 ]; then
  echo "Build successful. Output is in bin/server-commander"
else
  echo "Build failed."
  exit 1
fi
