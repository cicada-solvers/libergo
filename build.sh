#!/bin/bash

# Set the Go command
GOCMD=go

# Directories
CMD_DIR=cmd

# Function to clean binaries
clean_binaries() {
  echo "Cleaning binaries..."
  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "Removing $BINARY_NAME from $dir..."
      rm -vf "$dir/$BINARY_NAME"
    fi
  done
  echo "Binaries cleaned."
}

# Function to build binaries
build_binaries() {
  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "Building $BINARY_NAME in $dir..."
      cd "$dir"
      $GOCMD build .
      if [ $? -ne 0 ]; then
        echo "Failed to build $BINARY_NAME"
        exit 1
      fi
      cd - > /dev/null
    fi
  done
  echo "All binaries built successfully."
}

# Ask the user for the action
echo "Choose an action: 1. clean, 2. build"
read action

case $action in
  1)
    clean_binaries
    ;;
  2)
    build_binaries
    ;;
  *)
    echo "Invalid action. Please choose clean or build."
    exit 1
    ;;
esac