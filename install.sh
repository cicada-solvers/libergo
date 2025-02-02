#!/bin/bash

# Ensure the script is running as sudo
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exec sudo "$0" "$@"
  exit
fi

# Define the target directories and files
TARGET_DIR="$HOME/.libergo"
SOURCE_FILE="create_podman_db.sh"
TARGET_FILE="$TARGET_DIR/$SOURCE_FILE"
BIN_DIR="/opt/libergo"
MANIFEST_FILE="manifest.txt"
APPSETTINGS_FILE="appsettings.json"
TARGET_APPSETTINGS_FILE="$TARGET_DIR/$APPSETTINGS_FILE"

# Remove the binary directory if it exists
if [ -d "$BIN_DIR" ]; then
  rm -rf "$BIN_DIR"
  echo "Removed existing directory: $BIN_DIR"
fi

# Create the target directory if it does not exist
if [ ! -d "$TARGET_DIR" ]; then
  mkdir -p "$TARGET_DIR"
  echo "Created directory: $TARGET_DIR"
fi

# Copy the source file to the target directory
cp "$SOURCE_FILE" "$TARGET_FILE"
echo "Copied $SOURCE_FILE to $TARGET_FILE"

# Mark the target file as executable
chmod +x "$TARGET_FILE"
echo "Marked $TARGET_FILE as executable"

# Copy the appsettings.json file to the target directory
cp "$APPSETTINGS_FILE" "$TARGET_APPSETTINGS_FILE"
echo "Copied $APPSETTINGS_FILE to $TARGET_APPSETTINGS_FILE"

# Create the binary directory if it does not exist
if [ ! -d "$BIN_DIR" ]; then
  mkdir -p "$BIN_DIR"
  echo "Created directory: $BIN_DIR"
fi

# Read the binary names from the manifest file
if [ ! -f "$MANIFEST_FILE" ]; then
  echo "Manifest file $MANIFEST_FILE not found!"
  exit 1
fi

BINARIES=()
while IFS= read -r line; do
  BINARIES+=("$line")
done < "$MANIFEST_FILE"

# Remove existing symbolic links
for BINARY in "${BINARIES[@]}"; do
  if [ -L "/usr/bin/$BINARY" ]; then
    rm "/usr/bin/$BINARY"
    echo "Removed existing symbolic link: /usr/bin/$BINARY"
  fi
done

# Copy the binary files to the binary directory and mark them as executable by all users
for BINARY in "${BINARIES[@]}"; do
  cp "$BINARY" "$BIN_DIR"
  chmod 755 "$BIN_DIR/$BINARY"
  echo "Copied and marked $BINARY as executable in $BIN_DIR"
done

# Create symbolic links in /usr/bin
for BINARY in "${BINARIES[@]}"; do
  ln -sf "$BIN_DIR/$BINARY" "/usr/bin/$BINARY"
  echo "Created symbolic link for $BINARY in /usr/bin"
done

# Ensure the .libergo directory and all files in it are owned by the user
chown -R "$SUDO_USER:$SUDO_USER" "$TARGET_DIR"
chmod 755 "$TARGET_FILE"
echo "Changed ownership of $TARGET_DIR to $SUDO_USER"