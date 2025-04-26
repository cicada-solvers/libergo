#!/bin/bash

# Prompt for the input directory
read -p "Enter the input directory: " INPUT_DIRECTORY

# Remove trailing slash from the input directory if it exists
INPUT_DIRECTORY=$(echo "$INPUT_DIRECTORY" | sed 's:/*$::')

# Check if the provided input directory exists
if [ ! -d "$INPUT_DIRECTORY" ]; then
  echo "Error: $INPUT_DIRECTORY is not a directory"
  exit 1
fi

# Prompt for the output directory
read -p "Enter the output directory: " OUTPUT_DIRECTORY

# Remove trailing slash from the output directory if it exists
OUTPUT_DIRECTORY=$(echo "$OUTPUT_DIRECTORY" | sed 's:/*$::')

# Check if the provided output directory exists, if not create it
if [ ! -d "$OUTPUT_DIRECTORY" ]; then
  mkdir -p "$OUTPUT_DIRECTORY"
fi

# Iterate over each file in the input directory
for FILE in "$INPUT_DIRECTORY"/*; do
  if [ -f "$FILE" ]; then
    echo "Processing $FILE..."

    # Call the rdtext binary with the file name
    ./rdtext -create="true" -input="$FILE"
    sleep 300
  fi
done