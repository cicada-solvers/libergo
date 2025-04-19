#!/bin/bash

declare -a arr=("dwyl" "tfcom" "pgb")
declare -a ciphers=("caesar" "affine" "atbash" "trithemius" "autokey" "vigenere")

# Prompt the user for the depth
read -p "Enter the depth: " depth
read -p "Enter the input dir: " inputdir
read -p "Enter the output dir: " outputdir



for file in "$inputdir"/*
do
  basefile=$(basename "$file")
  echo "Base file name: $basefile"

  while IFS= read -r line
  do
    for ciphertype in "${ciphers[@]}"
    do
      mkdir "$outputdir/$ciphertype"

      for i in "${arr[@]}"
      do
          ./cipherval -passtwo="y" -ciphertype="$ciphertype" -maxdepth="$depth" -wordfile="$i.csv" -alphabet="rune" -text "$line" -output="$outputdir/$ciphertype/$basefile.$i.$ciphertype.$depth.p2.txt"
      done
    done
  done < "$file"
done