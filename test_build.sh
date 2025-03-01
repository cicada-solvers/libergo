#!/bin/bash

# Set the Go command
GOCMD=go

# Directories
CMD_DIR=cmd

# Function to build distribution binaries
dist_binaries() {
  DIST_DIR="dist"

  rm -rvf $DIST_DIR

  # Create the distribution directory
  mkdir -p $DIST_DIR

  echo "Cleaning binaries..."
  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "Removing $BINARY_NAME from $dir..."
      rm -vf "$dir/$BINARY_NAME"
    fi
  done

  # Create manifest.txt
  MANIFEST_FILE="manifest.txt"
  > "$MANIFEST_FILE"

  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "$BINARY_NAME" >> "$MANIFEST_FILE"

      echo "Building $BINARY_NAME for Linux..."
      cd "$dir"
      BINARY_DIR="$DIST_DIR"
      $GOCMD mod tidy
      $GOCMD mod download
      GOOS=linux GOARCH=amd64 $GOCMD build -o "../../$BINARY_DIR/$BINARY_NAME"
      if [ $? -ne 0 ]; then
        echo "Failed to build $BINARY_NAME for Linux"
        exit 1
      fi
      cp -v *.txt "../../$BINARY_DIR"
      cp -v *.sh "../../$BINARY_DIR"
      cd - > /dev/null
    fi
  done

  echo "Test binaries built successfully."
}

# Execute the dist_binaries function
dist_binaries