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
echo "2) Generate arrays"
echo "3) Process a single number"
read -p "Enter your choice (1, 2, or 3): " choice

if [ "$choice" -eq 1 ]; then
  echo "Cleaning data..."
  find . -type d -name '[0-9]*' | while read -r dir; do
    echo "Deleting directory: $dir"
    rm -rf "$dir"
  done
  echo "All subdirectories deleted."
elif [ "$choice" -eq 2 ]; then
  read -p "Enter the starting number: " start
  read -p "Enter the ending number: " end
  read -p "Enter the number of files per package: " files_per_package
  for i in $(seq "$start" "$end"); do
    echo "Creating byte arrays for $i"
    ./generatebytearrays "$i"
  done
  for i in $(seq "$start" "$end"); do
    folder=$(printf "%010d" "$i")
    cd "$folder" || exit
    file_count=0
    package_count=1
    find . -type f -name 'permutations_*.txt' | while read -r file; do
      if (( file_count % files_per_package == 0 )); then
        zip_file="package_$package_count.zip"
        ((package_count++))
      fi
      zip -m "$zip_file" "$file"
      ((file_count++))
    done
    cd ..
  done
elif [ "$choice" -eq 3 ]; then
  read -p "Enter the number to process: " number
  read -p "Enter the number of files per package: " files_per_package
  echo "Creating byte arrays for $number"
  ./generatebytearrays "$number"
  folder=$(printf "%010d" "$number")
  cd "$folder" || exit
  file_count=0
  package_count=1
  find . -type f -name 'permutations_*.txt' | while read -r file; do
    if (( file_count % files_per_package == 0 )); then
      zip_file="package_$package_count.zip"
      ((package_count++))
    fi
    zip -m "$zip_file" "$file"
    ((file_count++))
  done
  cd ..
else
  echo "Invalid choice. Exiting."
fi