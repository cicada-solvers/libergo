#!/bin/bash

echo "  ░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓████████▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░   ░▒▓█▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░   ░▒▓█▓▒░"
echo "  ░▒▓███████▓▒░ ░▒▓██████▓▒░   ░▒▓█▓▒░   ░▒▓██████▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░"
echo "  ░▒▓███████▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓████████▓▒░"
echo ""
echo ""
echo "   ░▒▓██████▓▒░░▒▓███████▓▒░░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░"
echo "  ░▒▓████████▓▒░▒▓███████▓▒░░▒▓███████▓▒░░▒▓████████▓▒░░▒▓██████▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░"
echo "  ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░"
echo ""
echo ""
echo "   ░▒▓███████▓▒░░▒▓██████▓▒░░▒▓███████▓▒░░▒▓█▓▒░▒▓███████▓▒░▒▓████████▓▒░"
echo "  ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░"
echo "  ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░"
echo "   ░▒▓██████▓▒░░▒▓█▓▒░      ░▒▓███████▓▒░░▒▓█▓▒░▒▓███████▓▒░  ░▒▓█▓▒░"
echo "         ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░        ░▒▓█▓▒░"
echo "         ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░        ░▒▓█▓▒░"
echo "  ░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░        ░▒▓█▓▒░"

echo "Choose an option:"
echo "1) Clean data"
echo "2) Generate array for a number"
echo "3) Generate arrays for a range of numbers"
read -p "Enter your choice (1, 2, or 3): " choice

if [ "$choice" -eq 1 ]; then
  echo "Cleaning data..."
  find . -type d -name '[0-9]*' | while read -r dir; do
    echo "Deleting directory: $dir"
    rm -rf "$dir"
  done
  echo "All subdirectories deleted."
elif [ "$choice" -eq 2 ]; then
  read -p "Enter the number to create arrays: " number
  echo "Creating byte arrays for $number"
  ./generatebytearrays "$number"
elif [ "$choice" -eq 3 ]; then
  read -p "Enter the start number: " start
  read -p "Enter the end number: " end
  for ((i=start; i<=end; i++)); do
    echo "Creating byte arrays for $i"
    ./generatebytearrays "$i"
  done
else
  echo "Invalid choice. Exiting."
fi