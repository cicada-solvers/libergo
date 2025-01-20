#!/bin/bash

# Set the Go command
GOCMD=go

# Directories
CMD_DIR=cmd

# Loop through each subdirectory in the cmd directory
for dir in $CMD_DIR/*; do
  if [ -d "$dir" ]; then
    # Get the base name of the directory
    BINARY_NAME=$(basename "$dir")
    # Change to the directory
    cd "$dir"
    # Build the binary
    echo "Building $BINARY_NAME in $(pwd)..."
    CMD="$GOCMD build ."
    $CMD
    if [ $? -ne 0 ]; then
      echo "Failed to build $BINARY_NAME"
      exit 1
    fi
    # Change back to the original directory
    cd - > /dev/null
  fi
done

echo "All binaries built successfully."