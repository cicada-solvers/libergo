package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hashinglib"
	"liberdatabase"
	"sync"
	"time"

	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/blake2b"
)

// processTasks processes the tasks
func processTasks(tasks chan liberdatabase.ReadPermutation, wg *sync.WaitGroup, existingHash string, done chan struct{}, once *sync.Once, rowCount *int) {
	defer wg.Done()

	hashCount := 0
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	go func() {
		colors := []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m", "\033[90m", "\033[91m", "\033[92m"}
		colorIndex := 0
		for range ticker.C {
			aps := hashCount / 4
			*rowCount = *rowCount - aps
			fmt.Printf("%sArrays per minute: %d - %d Remaining\033[0m\n", colors[colorIndex], aps, *rowCount)
			hashCount = 0
			colorIndex = (colorIndex + 1) % len(colors)
		}
	}()

	var idsToRemove []string

	for {
		select {
		case perm, ok := <-tasks:
			if !ok {
				return
			}

			idsToRemove = append(idsToRemove, perm.ID)
			if (len(idsToRemove) % 60000) == 0 {
				db, _ := liberdatabase.InitConnection()
				liberdatabase.RemoveItems(db, idsToRemove)
				idsToRemove = nil
				closeError := liberdatabase.CloseConnection(db)
				if closeError != nil {
					return
				}
			}

			hashes := generateHashes(perm.StartArray)
			for hashName, hash := range hashes {
				hashCount++
				if hash == existingHash {
					taskStr := fmt.Sprintf("%v", perm)
					output := fmt.Sprintf("Match found: %s, Hash Name: %s, Byte Array: %s\n", taskStr, hashName, perm.StartArray)
					fmt.Print(output)
					err := liberdatabase.InsertFoundHash(taskStr, hashName)
					if err != nil {
						fmt.Printf("Error inserting found hash: %v\n", err)
					}
					once.Do(func() { close(done) })
				}
			}
		case <-done:
			return
		}
	}
}

// generateHashes generates hashes for a given byte array
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
