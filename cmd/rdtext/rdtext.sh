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

# Iterate over each file in the input directory
for FILE in "$INPUT_DIRECTORY"/*; do
  if [ -f "$FILE" ]; then
    # If the file does not end with .xlsx, we want to skip it.
    if [[ ! "$FILE" =~ \.xlsx$ ]]; then
        continue
    fi
    
    echo "Processing $FILE..."

    # Call the rdtext binary with the file name
    ./rdtext -input="$FILE"
    ./rdtext -input="$FILE" -reverse="true"
  fi
done