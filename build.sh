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
      rm -vf "$dir/*.txt"
      rm -vf "$dir/*.bin"
    fi
  done
  rm -rvf dist
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
  mkdir -p "$DIST_DIR/linux" "$DIST_DIR/mac_amd64" "$DIST_DIR/mac_arm64" "$DIST_DIR/windows"

  echo "Cleaning binaries..."
  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "Removing $BINARY_NAME from $dir..."
      rm -vf "$dir/$BINARY_NAME"
      rm -vf "$dir/*.txt"
      rm -vf "$dir/*.bin"
    fi
  done

  # Create manifest.txt
  MANIFEST_FILE="manifest.txt"
  echo "" > "$MANIFEST_FILE"

  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "$BINARY_NAME" >> "$MANIFEST_FILE"

      echo "Building $BINARY_NAME for Linux..."
      pushd "$dir" > /dev/null
      BINARY_DIR="$DIST_DIR/linux"
      GOOS=linux GOARCH=amd64 $GOCMD build -o "../../$BINARY_DIR/$BINARY_NAME"
      if [ $? -ne 0 ]; then
        echo "Failed to build $BINARY_NAME for Linux"
        exit 1
      fi

      if [ "$BINARY_NAME" != "runecalc" ]; then
        echo "Building $BINARY_NAME for Mac (amd64)..."
        BINARY_DIR="$DIST_DIR/mac_amd64"
        GOOS=darwin GOARCH=amd64 $GOCMD build -o "../../$BINARY_DIR/$BINARY_NAME"
        if [ $? -ne 0 ]; then
          echo "Failed to build $BINARY_NAME for Mac (amd64)"
          exit 1
        fi

        echo "Building $BINARY_NAME for Mac (arm64)..."
        BINARY_DIR="$DIST_DIR/mac_arm64"
        GOOS=darwin GOARCH=arm64 $GOCMD build -o "../../$BINARY_DIR/$BINARY_NAME"
        if [ $? -ne 0 ]; then
          echo "Failed to build $BINARY_NAME for Mac (arm64)"
          exit 1
        fi

        echo "Building $BINARY_NAME for Windows..."
        BINARY_DIR="$DIST_DIR/windows"
        GOOS=windows GOARCH=amd64 $GOCMD build -o "../../$BINARY_DIR/$BINARY_NAME.exe"
        if [ $? -ne 0 ]; then
          echo "Failed to build $BINARY_NAME for Windows"
          exit 1
        fi
      else
        echo "Skipping $BINARY_NAME for Mac and Windows..."
        rm -fv runecalc.tar.xz
        fyne package -os linux
        cp -f runecalc.tar.xz "../../$BINARY_DIR/runecalc.tar.xz"
      fi
      popd > /dev/null
    fi
  done

  echo "Copying additional files..."
  cp install.sh "$DIST_DIR/linux"
  cp install.sh "$DIST_DIR/mac_amd64"
  cp install.sh "$DIST_DIR/mac_arm64"
  cp uninstall.sh "$DIST_DIR/linux"
  cp uninstall.sh "$DIST_DIR/mac_amd64"
  cp uninstall.sh "$DIST_DIR/mac_arm64"
  cp create_podman_db.sh "$DIST_DIR/linux"
  cp create_podman_db.sh "$DIST_DIR/mac_amd64"
  cp create_podman_db.sh "$DIST_DIR/mac_arm64"
  cp create_podman_db.ps1 "$DIST_DIR/windows"
  cp docker-compose.yml "$DIST_DIR/windows"
  cp "$MANIFEST_FILE" "$DIST_DIR/linux"
  cp "$MANIFEST_FILE" "$DIST_DIR/mac_amd64"
  cp "$MANIFEST_FILE" "$DIST_DIR/mac_arm64"
  cp appsettings.json "$DIST_DIR/linux"
  cp appsettings.json "$DIST_DIR/mac_amd64"
  cp appsettings.json "$DIST_DIR/mac_arm64"
  cp appsettings.json "$DIST_DIR/windows"
  cp words.txt "$DIST_DIR/linux"
  cp words.txt "$DIST_DIR/mac_amd64"
  cp words.txt "$DIST_DIR/mac_arm64"
  cp words.txt "$DIST_DIR/windows"
  cp definitions_flat "$DIST_DIR/linux"
  cp definitions_flat "$DIST_DIR/mac_amd64"
  cp definitions_flat "$DIST_DIR/mac_arm64"
  cp definitions_flat "$DIST_DIR/windows"

  echo "Compressing directories..."
  tar -czvf "$DIST_DIR/linux_$VERSION.tar.gz" -C "$DIST_DIR" linux
  tar -czvf "$DIST_DIR/mac_amd64_$VERSION.tar.gz" -C "$DIST_DIR" mac_amd64
  tar -czvf "$DIST_DIR/mac_arm64_$VERSION.tar.gz" -C "$DIST_DIR" mac_arm64
  zip -r "$DIST_DIR/windows_$VERSION.zip" "$DIST_DIR/windows"

  echo "Removing uncompressed directories..."
  rm -rf "$DIST_DIR/mac_amd64" "$DIST_DIR/mac_arm64" "$DIST_DIR/windows"

  echo "Cleaning binaries..."
  for dir in $CMD_DIR/*; do
    if [ -d "$dir" ]; then
      BINARY_NAME=$(basename "$dir")
      echo "Removing $BINARY_NAME from $dir..."
      rm -vf "$dir/$BINARY_NAME"
      rm -vf "$dir/*.txt"
      rm -vf "$dir/*.bin"
    fi
  done

  echo "All distribution binaries built, compressed, and cleaned up successfully."
}

# Ask the user for the action
read -p "Choose an action: 1) Clean, 2) Build, 3) Dist: " action

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