# LiberGo
This is a Liber Primus Analysis Toolkit. I have been using it for investigating the Liber Primus. It is **very early** in its development and I am working on it as I have the time to do it. I have also used it on some other puzzles so its application is not limited to the Liber Primus.

## Program Information

LiberGo includes many command-line utilities for cryptographic analysis, text processing, and mathematical operations:

### Text Processing Tools
- **rdtext**: Processes Excel files to generate sentence permutations and checks if their gemmatria sums are prime numbers. Can process in normal or reverse word order.
- **rdcheck**: Validates and checks text data against various criteria.
- **detangle**: Untangles encoded or interleaved text.
- **rencode**: Encodes text using various methods.
- **winchafftext**: Processes text using chaff techniques.
- **railit**: Applies rail fence cipher to text.
- **corpuslist**: Analyzes text corpora.
- **wordlecheat**: Helper for Wordle puzzles.
- **calrunestats**: Calculates statistics for runes/characters.
- **runedonkey**: Performs operations on rune-encoded text.

### Mathematical Tools
- **gemsum**: Calculates gemmatria sums for text.
- **gemproduct**: Calculates gemmatria products.
- **isprime**: Checks if a number is prime.
- **primedb**: Database operations for prime numbers.
- **factorize**: Factorizes numbers into prime factors.
- **intgen**: Generates integers based on specified patterns.
- **genseq**: Generates numerical sequences.
- **invgoldbach**: Inventory Goldbach operations.
- **mob**: Mobius operations tool.
- **fanalysis**: Factorization verification tool.

### Conversion Tools
- **base60**: Converts to/from base 60.
- **binfile**: Binary file operations.
- **binstring**: Binary string operations.
- **binvert**: Binary inversion operations.
- **byte2bin**: Converts bytes to binary.
- **byte2binstring**: Converts bytes to binary strings.
- **byte2int**: Converts bytes to integers.
- **hex2bytearray**: Converts hex to byte arrays.
- **hex2int**: Converts hex to integers.
- **decode64**: Base64 decoder.

### Network & Specialized Tools
- **ipaddressgen**: Generates IP addresses.
- **primeaddressgen**: Generates addresses based on prime numbers.
- **checkhashes**: Validates hashes.
- **checkmultihash**: Validates multiple hashes.
- **cipherval**: Calculates cipher values.
- **permute**: Generates permutations.
- **libergo**: Main application that provides a framework for other tools.

## Building release binaries
You will need to have go installed on your system. You can get it from https://golang.org.

This script is for Linux and not Windows or Mac.
**This will not destroy the database!**

## Installation
If you're on Linux, then you can use install.sh. It will copy the binaries to /opt. Then it will create the symlinks to the bin directory.

There are no installers for Windows at this time. You will need to copy the files to a directory of your choosing.

After you get it installed, run the following command. It will set everything up for you.

**Note: If you are on Windows, you will need to copy in the text files into the .libergo directory in your user directory.**

I don't use Windows so if anyone wants to help, it would be much appreciated.

## Database Required!!!
Some tools require a Postgres database to be present.

I prefer to use podman to host the database. You can use docker if you want. You can use the following command to start the database for podman.

### Bash
`chmod 777 create_podman_db.sh ./create_podman_db.sh`

### libergo
If you use your own server, you will need to set the variables in the appsettings.json file in the .libergo directory in you user directory.
`libergo -initdbserver`

## Upcoming Tools available in libergo
- circular shift
- least significant bits in message.
- least significant bit in number strings.
- skip and take
- dictionary checker
- clock angle calculator
- letter frequency analysis stuff
- Scytale
- Text spiraler