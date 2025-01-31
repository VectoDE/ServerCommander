#!/bin/bash

# This script cleans up build artifacts.

echo "Cleaning up..."

# Remove the build directory and binary
rm -rf bin

# Remove any Go module cache that is no longer needed
go clean -modcache

echo "Cleanup complete."
