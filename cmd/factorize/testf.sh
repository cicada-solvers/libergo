#!/bin/bash

# Compile the factorize program
go build

# Check if the start and end numbers are provided as arguments
if [ $# -lt 2 ]; then
  echo "Please provide a start and end number as arguments."
  exit 1
fi

# Read input numbers
start=$1
end=$2

# Define the output file
output_file="factorize_output.txt"

# Remove the output file if it exists
if [ -f "$output_file" ]; then
  rm "$output_file"
fi

# Loop from start to end number
for ((i=start; i<=end; i++)); do
  echo "Factorizing $i"
  ./factorize "$i" >> "$output_file"
done