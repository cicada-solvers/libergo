package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/blake2b"
	"hashinglib"
	"os"
)

func main() {
	existingHash := "36367763ab73783c7af284446c59466b4cd653239a311cb7116d4618dee09a8425893dc7500b464fdaf1672d7bef5e891c6e2274568926a49fb4f45132c2a8b4"

	// Check if the user provided an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: checkhashes <input-string>")
		return
	}

	// Get the input string from the arguments
	inputString := os.Args[1]

	// Convert the string to a byte array
	byteArray := []byte(inputString)

	// Print the byte array in hexadecimal format
	fmt.Printf("Input string: %s\n", inputString)
	fmt.Printf("Byte array (hex): %s\n", hex.EncodeToString(byteArray))

	hashes := generateHashes(byteArray)

	for hashName, hash := range hashes {
		fmt.Printf("%s: %s\n", hashName, hash)

		if hash == existingHash {
			fmt.Printf("Found matching hash: %s\n", hash)
		} else {
			fmt.Printf("No match for %s\n", hashName)
		}
	}

	fmt.Print("Done\n\n\n")
}

func generateHashes(data []byte) map[string]string {
	hashes := make(map[string]string)

	sha512Hash := sha512.Sum512(data)
	hashes["SHA-512"] = hex.EncodeToString(sha512Hash[:])

	whirlpoolHash := whirlpool.New()
	whirlpoolHash.Write(data)
	whirlHash := whirlpoolHash.Sum(nil)
	hashes["Whirlpool-512"] = hex.EncodeToString(whirlHash[:])

	blake2bHash := blake2b.Sum512(data)
	hashes["Blake2b-512"] = hex.EncodeToString(blake2bHash[:])

	blake512Hash := hashinglib.Blake512Hash(data)
	hashes["Blake-512"] = hex.EncodeToString(blake512Hash)

	return hashes
}
