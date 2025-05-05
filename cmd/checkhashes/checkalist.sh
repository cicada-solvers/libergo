#!/bin/bash

# Scan for .txt files in the current directory
for file in *.txt; do
  # Check if any .txt files exist
  if [ -f "$file" ]; then
    echo "Processing file: $file"
    # Get the total number of lines in the file
    total_lines=$(wc -l < "$file")
    current_line=0

    # Read and process each line in the .txt file
    while IFS= read -r line; do
      current_line=$((current_line + 1))
      percentage=$((current_line * 100 / total_lines))
      echo -ne "Count: $current_line/$total_lines - Process: $percentage%\r"
      ./checkhashes "$line" >> "$file.output.txt"
    done < "$file"
    echo -e "\nFinished processing $file."
  else
    echo "No .txt files found in the current directory."
    break
  fi
done