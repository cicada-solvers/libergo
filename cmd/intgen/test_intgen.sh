#!/bin/bash

# Define the output file
output_file="intgen_output.txt"
filtered_output_file="filtered_intgen_output.txt"
number_only_output_file="number_only_output.txt"

# Remove the output files if they exist
if [ -f "$output_file" ]; then
  rm "$output_file"
fi

if [ -f "$filtered_output_file" ]; then
  rm "$filtered_output_file"
fi

if [ -f "$number_only_output_file" ]; then
  rm "$number_only_output_file"
fi

# Array of bit lengths
#bit_lengths=(8 16 32 64 128 256 512 1024 2048 4096)
bit_lengths=(8 16 32 64 128 256 512 1024)

# Loop through each bit length
for bits in "${bit_lengths[@]}"; do
  echo "Generating 100 ${bits}-bit numbers"
  for ((i=0; i<100; i++)); do
    echo "Generating number $i"
    ./intgen "$bits" >> "$output_file"
  done
done

# Filter the output file
grep "^Product of the two primes:" "$output_file" > "$filtered_output_file"

# Further filter to only leave number values
grep -o '[0-9]\+' "$filtered_output_file" > "$number_only_output_file"

sort -u "$number_only_output_file" | sort -n -o "$output_file"
rm "$filtered_output_file" "$number_only_output_file"

echo "Number generation complete. Filtered and sorted output written to $output_file"