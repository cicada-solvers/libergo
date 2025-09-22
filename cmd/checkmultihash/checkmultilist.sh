#!/bin/bash

# Prompt the user for how to read the file content
read -r -p "Read files as (1) string byte array or (2) plain string? [1/2]: " choice

case "$choice" in
  1) bytefile_flag="true" ;;
  2) bytefile_flag="false" ;;
  *) echo "Invalid choice. Defaulting to plain string."; bytefile_flag="false" ;;
esac

# Scan for .txt files in the current directory
for file in *.txt; do
  # Check if any .txt files exist
  if [ -f "$file" ]; then
      echo "Processing file: $file"
      # Call checkmultihash with filename and bytefile flags
      ./checkmultihash --filename "$file" --bytefile "$bytefile_flag" >> output.sh
  else
    echo "No .txt files found in the current directory."
    break
  fi
done