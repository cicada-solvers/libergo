#!/bin/bash

# Scan for .txt files in the current directory
for file in *.txt; do
  # Check if any .txt files exist
  if [ -f "$file" ]; then
      echo "Processing file: $file"
      ./checkmultihash "$file"
  else
    echo "No .txt files found in the current directory."
    break
  fi
done