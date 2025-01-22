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

# Function to build distribution binaries
dist_binaries() {
  read -p "Enter version number: " VERSION
  DIST_DIR="dist/$VERSION"
  mkdir -p "$DIST_DIR/linux" "$DIST_DIR/mac" "$DIST_DIR/windows"

  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "Building $BINARY_NAME for Linux..."
      cd "$dir"
      BINARY_DIR="$DIST_DIR/linux/$BINARY_NAME"
      mkdir -p "$BINARY_DIR"
      GOOS=linux GOARCH=amd64 $GOCMD build -o "$BINARY_DIR/$BINARY_NAME"
      if [ $? -ne 0 ]; then
        echo "Failed to build $BINARY_NAME for Linux"
        exit 1
      fi
      cp "$dir\appsettings.json" "$BINARY_DIR/"

      echo "Building $BINARY_NAME for Mac..."
      BINARY_DIR="$DIST_DIR/mac/$BINARY_NAME"
      mkdir -p "$BINARY_DIR"
      GOOS=darwin GOARCH=amd64 $GOCMD build -o "$BINARY_DIR/$BINARY_NAME"
      if [ $? -ne 0 ]; then
        echo "Failed to build $BINARY_NAME for Mac"
        exit 1
      fi
      cp "$dir\appsettings.json" "$BINARY_DIR/"

      echo "Building $BINARY_NAME for Windows..."
      BINARY_DIR="$DIST_DIR/windows/$BINARY_NAME"
      mkdir -p "$BINARY_DIR"
      GOOS=windows GOARCH=amd64 $GOCMD build -o "$BINARY_DIR/$BINARY_NAME.exe"
      if [ $? -ne 0 ]; then
        echo "Failed to build $BINARY_NAME for Windows"
        exit 1
      fi
      cp "$dir\appsettings.json" "$BINARY_DIR/"
      cd - > /dev/null
    fi
  done
  echo "All distribution binaries built successfully."
}

# Ask the user for the action
echo "Choose an action: 1. clean, 2. build, 3. dist"
read action

case $action in
  1)
    clean_binaries
    ;;
  2)
    build_binaries
    ;;
  3)
    dist_binaries
    ;;
  *)
    echo "Invalid action. Please choose clean, build, or dist."
    exit 1
    ;;
esac