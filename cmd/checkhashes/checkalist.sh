#!/bin/bash

# Scan for .txt files in the current directory
for file in *.txt; do
  # Check if any .txt files exist
  if [ -f "$file" ]; then
    echo "Processing file: $file"
    # Read and output each line in the .txt file
    while IFS= read -r line; do
      ./checkhashes "$line" >> "$file.output.txt"
    done < "$file"
  else
    echo "No .txt files found in the current directory."
    break
  fi
done