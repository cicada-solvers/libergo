#!/bin/bash

# Ensure the script is running as sudo
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exec sudo "$0" "$@"
  exit
fi

# Define the target directories and files
TARGET_DIR="$HOME/.libergo"
BIN_DIR="/opt/libergo"
MANIFEST_FILE="manifest.txt"

# Read the binary names from the manifest file
if [ ! -f "$MANIFEST_FILE" ]; then
  echo "Manifest file $MANIFEST_FILE not found!"
  exit 1
fi

BINARIES=()
while IFS= read -r line; do
  BINARIES+=("$line")
done < "$MANIFEST_FILE"

# Remove symbolic links in /usr/bin
for BINARY in "${BINARIES[@]}"; do
  if [ -L "/usr/bin/$BINARY" ]; then
    rm "/usr/bin/$BINARY"
    echo "Removed symbolic link: /usr/bin/$BINARY"
  fi
done

# Remove the binary files from the binary directory
if [ -d "$BIN_DIR" ]; then
  rm -rvf "$BIN_DIR"
  echo "Removed directory: $BIN_DIR"
fi

# Remove the target directory
if [ -d "$TARGET_DIR" ]; then
  rm -rvf "$TARGET_DIR"
  echo "Removed directory: $TARGET_DIR"
fi

echo "Uninstallation completed."