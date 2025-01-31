#!/bin/bash

# This script runs the tests for ServerCommander.

echo "Running tests..."

# Run Go tests
go test ./...

# Check if tests passed
if [ $? -eq 0 ]; then
  echo "Tests passed successfully."
else
  echo "Tests failed."
  exit 1
fi
