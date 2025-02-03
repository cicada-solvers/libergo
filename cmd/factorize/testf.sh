#!/bin/bash

# Compile the factorize program
go build

# Check if the number is provided as an argument
if [ $# -lt 1 ]; then
  echo "Please provide a number to be factorized as an argument."
  exit 1
fi

# Read input number
number=$1

# Define the output file
output_file="factorize_output.txt"

# Remove the output file if it exists
if [ -f "$output_file" ]; then
  rm "$output_file"
fi

# Loop from 0 to the provided number
for ((i=0; i<=number; i++)); do
  echo "Factorizing $i"
  ./factorize "$i" >> "$output_file"
done