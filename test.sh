#!/bin/bash

# Find all subdirectories under pkg
directories=$(find ./pkg -type d)

# Loop through each directory and run go test
for dir in $directories; do
    cd $dir
    go test
    cd - > /dev/null
done