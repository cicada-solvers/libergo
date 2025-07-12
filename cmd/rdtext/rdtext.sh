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

# Create a temporary file to store file sizes and names
temp_file=$(mktemp)

# Get all the .xlsx files and their sizes
for FILE in "$INPUT_DIRECTORY"/*.xlsx; do
  if [ -f "$FILE" ]; then
    # Get file size in bytes
    size=$(stat -c %s "$FILE")
    echo "$size $FILE" >> "$temp_file"
  fi
done

# Sort files by size (ascending) and process them
while read size file; do
  echo "Processing $file... (Size: $size bytes)"

  # Call the rdtext binary with the file name
  ./rdtext -input="$file"
  ./rdtext -input="$file" -reverse="true"
done < <(sort -n "$temp_file")

# Remove the temporary file
rm "$temp_file"