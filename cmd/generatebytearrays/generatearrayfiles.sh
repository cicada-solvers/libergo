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
echo "3) Generate array for a single number"
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
  files_per_package=10000

  for i in $(seq "$start" "$end"); do
    echo "Creating byte arrays for $i"
    ./generatebytearrays "$i"

    # Compress the files in the directory after generatebytearrays runs
    package_count=1
    folder=$(printf "%010d" "$i")
    cd "$folder" || exit
    file_count=0
    find . -type f -name 'permutations_*.txt' | while read -r file; do
      if (( file_count % files_per_package == 0 )); then
        zip_file="package_$package_count.zip"
        ((package_count++))
      fi
      zip -q "$zip_file" "$file"
      rm "$file"
      ((file_count++))
    done
    cd ..
  done
elif [ "$choice" -eq 3 ]; then
  read -p "Enter the number to create arrays: " number
  files_per_package=10000
  echo "Creating byte arrays for $number"
  ./generatebytearrays "$number"

  # Compress the files in the directory after generatebytearrays runs
  folder=$(printf "%010d" "$number")
  package_count=1
  cd "$folder" || exit
  file_count=0
  find . -type f -name 'permutations_*.txt' | while read -r file; do
    if (( file_count % files_per_package == 0 )); then
      zip_file="package_$package_count.zip"
      ((package_count++))
    fi
    zip -q "$zip_file" "$file"
    rm "$file"
    ((file_count++))
  done
  cd ..
else
  echo "Invalid choice. Exiting."
fi