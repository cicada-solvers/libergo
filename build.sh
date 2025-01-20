#!/bin/bash

# Set the Go command
GOCMD=go

# Directories
CMD_DIR=cmd
BIN_DIR=bin

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
  mkdir -p $BIN_DIR
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

# Function to install binaries
install_binaries() {
  if [ "$EUID" -ne 0 ]; then
    echo "Please run as root to install binaries."
    exit 1
  fi
  for binary in $BIN_DIR/*; do
    if [ -f "$binary" ]; then
      echo "Installing $(basename "$binary") to /usr/bin..."
      cp "$binary" /usr/bin/
      if [ $? -ne 0 ]; then
        echo "Failed to install $(basename "$binary")"
        exit 1
      fi
    fi
  done
  echo "All binaries installed successfully."
}

# Ask the user for the action
echo "Choose an action: 1. clean, 2. build, or 3. install"
read action

case $action in
  1)
    clean_binaries
    ;;
  2)
    build_binaries
    ;;
  3)
    install_binaries
    ;;
  *)
    echo "Invalid action. Please choose clean, build, or install."
    exit 1
    ;;
esac