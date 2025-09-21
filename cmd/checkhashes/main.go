package main

import (
	"blake"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"fnv5120"
	"fnv5121"
	"groestl"
	cube2 "hashinglib/cube"
	jh2 "jh"
	"keccak3"
	"lsh"
	"md6"
	skein2 "skein"
	"streebog"
	whirlpool2 "whirlpool"

	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"

	"os"
)

// main is the entry point of the program; it processes input strings and checks for matching cryptographic hash values.
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

// generateHashes computes cryptographic hashes (SHA-512, Whirlpool-512, Blake2b-512, and Blake-512) for the given input data.
// It returns a map where the keys are the hash algorithm names and the values are the corresponding hexadecimal hash strings.
func generateHashes(data []byte) map[string]string {
	hashes := make(map[string]string)

	sha512Hash := sha512.Sum512(data)
	hashes["SHA-512"] = hex.EncodeToString(sha512Hash[:])

	sha3512Hash := sha3.New512()
	_, err := sha3512Hash.Write(data)
	if err != nil {
		return nil
	}
	sha3Hash := sha3512Hash.Sum(nil)
	hashes["SHA3-512"] = hex.EncodeToString(sha3Hash)

	whirlpoolHash := whirlpool.New()
	whirlpoolHash.Write(data)
	whirlHash := whirlpoolHash.Sum(nil)
	hashes["Whirlpool-512"] = hex.EncodeToString(whirlHash[:])

	blake2bHash := blake2b.Sum512(data)
	hashes["Blake2b-512"] = hex.EncodeToString(blake2bHash[:])

	blake512Hash := blake.Blake512Hash(data)
	hashes["Blake-512"] = hex.EncodeToString(blake512Hash)

	jh := jh2.NewJH512()
	jh.Write(data)
	jhHash := jh.Sum(nil)
	hashes["JH-512"] = hex.EncodeToString(jhHash)

	skein := skein2.NewSkein512()
	skein.Write(data)
	skeinHash := skein.Sum(nil)
	hashes["Skein-512"] = hex.EncodeToString(skeinHash)

	grostl := groestl.NewGroestl512()
	grostl.Write(data)
	grostlHash := grostl.Sum(nil)
	hashes["Groestl-512"] = hex.EncodeToString(grostlHash)

	keccak3512 := keccak3.Keccak3_512(data)
	hashes["Keccak3-512"] = hex.EncodeToString(keccak3512)

	whirlpool0 := whirlpool2.Whirlpool0(data)
	hashes["Whirlpool-0"] = hex.EncodeToString(whirlpool0)

	whirlpoolT := whirlpool2.WhirlpoolT(data)
	hashes["Whirlpool-T"] = hex.EncodeToString(whirlpoolT)

	cube := cube2.CubeHash512(data)
	hashes["Cube-512"] = hex.EncodeToString(cube)

	streebogHash := streebog.Hash512(data)
	hashes["Streebog-512"] = hex.EncodeToString(streebogHash)

	fnv5120hash := fnv5120.Hash512(data)
	hashes["FNV5120-512"] = hex.EncodeToString(fnv5120hash)

	fnv5120ahash := fnv5120.Hash512a(data)
	hashes["FNV5120A-512"] = hex.EncodeToString(fnv5120ahash)

	fnv5121hash := fnv5121.Hash(data)
	hashes["FNV5121-512"] = hex.EncodeToString(fnv5121hash)

	fnv5121ahash := fnv5121.HashA(data)
	hashes["FNV5121A-512"] = hex.EncodeToString(fnv5121ahash)

	md6hash := md6.Sum512(data)
	hashes["MD6-512"] = hex.EncodeToString(md6hash)

	lshHash := lsh.Sum512(data)
	hashes["LSHH-512"] = hex.EncodeToString(lshHash)

	return hashes
}
