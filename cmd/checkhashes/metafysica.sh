#!/bin/bash

# Check if metafysica.txt exists
if [ ! -f "metafysica.txt" ]; then
  echo "Error: metafysica.txt not found."
  exit 1
fi

# Clear or create metafysica.output.txt
> metafysica.output.txt

./checkhashes "http://www.metafysica.nl/" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/turing/" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/holism/" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/ontology/" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/nature/" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/wings/" >> metafysica.output.txt

./checkhashes "http://www.metafysica.nl" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/turing" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/holism" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/ontology" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/nature" >> metafysica.output.txt
./checkhashes "http://www.metafysica.nl/wings" >> metafysica.output.txt

# Read each line from metafysica.txt
while IFS= read -r line; do
  # Call the checkhashes program and append the output to metafysica.output.txt
  ./checkhashes "http://www.metafysica.nl/$line" >> metafysica.output.txt
done < metafysica.txt

echo "Processing complete. Output written to metafysica.output.txt."