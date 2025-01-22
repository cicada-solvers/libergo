#!/bin/bash

# Find all files starting with "permutation" in the current and subdirectories
find . -type f -name 'permutation*' | while read -r file; do
  # Echo the file being processed
  echo "Processing file: $file"
  # Call processhashes with the file name
  ./processhashes "$file"
  # Delete the file after processing
  rm -vf "$file"
done