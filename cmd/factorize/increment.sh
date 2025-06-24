#!/bin/bash

rm -fv factorize_output.txt

# Start from 2
start=2
# 32-bit maximum (2^31 - 1)
end=2147483647

for ((i=start; i<=end; i++)); do
    # Call the factorize command with the current number
    ./factorize "$i"
done