package main

import (
	"fmt"
	"github.com/jdkato/prose/v2"
	"liberdatabase"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

// Sentence represents a sentence with its content, output file name, and column index.
type Sentence struct {
	FileName   string
	Content    string
	Output     string
	PrimeValue int64
}

var fileMutex = sync.Mutex{}

func main() {
	outputFile := "output.txt"

	// Now we are going to remove the million records from the database.
	conn, connErr := liberdatabase.InitConnection()
	if connErr != nil {
		fmt.Printf("error initializing Postgres connection: %v", connErr)
	}

	// Gets all the file names from the database
	fmt.Println("Getting all file names from the database...")
	fileNames, _ := liberdatabase.GetAllFileNames(conn)

	// Presents a menu for the user to select a file
	fmt.Println("\nSelect a file to process:")
	fmt.Println("-------------------------------------")
	for i, fileName := range fileNames {
		fmt.Printf("%d. %s\n", i+1, fileName)
	}

	var selection int
	fmt.Print("\nEnter selection number: ")
	_, err := fmt.Scanln(&selection)
	if err != nil || selection < 1 || selection > len(fileNames) {
		fmt.Printf("Invalid selection: %v", err)
		return
	}

	fileName := fileNames[selection-1]
	fmt.Printf("Selected file: %s\n", fileName)

	totalRecords, _ := liberdatabase.GetRecordCountByFileName(conn, fileName)
	for totalRecords > int64(0) {
		// Get the top one million sentence records
		records, getErr := liberdatabase.GetTopMillionSentenceRecords(conn, fileName)
		if getErr != nil {
			fmt.Printf("error getting top million sentence records: %v", getErr)
		}

		var wg sync.WaitGroup
		sentenceChan := make(chan Sentence, 16384) // Increased buffer size

		go func() {
			for _, record := range records {
				// Create a new Sentence instance
				sentence := Sentence{
					FileName:   record.FileName,
					Content:    record.DictSentence,
					Output:     outputFile,
					PrimeValue: record.GemValue,
				}
				sentenceChan <- sentence
			}
			close(sentenceChan)
		}()

		numWorkers := runtime.NumCPU() // Adjusted number of workers
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go calculateProbabilityAndWriteToFile(sentenceChan, &wg)
		}

		wg.Wait()

		// Remove the million records from the database
		removeErr := liberdatabase.RemoveMillionSentenceRecords(conn, records)
		if removeErr != nil {
			fmt.Printf("error removing million sentence records: %v", removeErr)
		}

		totalRecords, _ = liberdatabase.GetRecordCountByFileName(conn, fileName)
	}
}

func calculateProbabilityAndWriteToFile(sentChan chan Sentence, wg *sync.WaitGroup) {
	defer wg.Done()

	for sentence := range sentChan {
		posCounts, totalWords := analyzeText(sentence.Content)
		probability := calculateSentenceProbability(posCounts, totalWords)

		if probability > 0 {
			fmt.Printf("Sentence Probability: %.2f%%\n", probability)

			outputText := fmt.Sprintf("%.2f%%|%d|%s\n\n", probability, sentence.PrimeValue, sentence.Content)

			// Write the content to the output file
			outputBytes := []byte(outputText)

			writeContentsToOutputFile(sentence.Output, outputBytes)
		}
	}
}

// analyzeText analyzes the given text and returns the part-of-speech counts and total word count.
func analyzeText(text string) (map[string]int, int) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	posCounts := map[string]int{
		"Noun":        0,
		"Verb":        0,
		"Adjective":   0,
		"Adverb":      0,
		"Determiner":  0,
		"Conjunction": 0,
		"Preposition": 0,
		"Pronoun":     0,
		"Punctuation": 0,
		"NamedEntity": 0,
	}
	totalWords := 0

	for _, tok := range doc.Tokens() {
		switch tok.Tag {
		case "NN", "NNS", "NNP", "NNPS":
			posCounts["Noun"]++
		case "VB", "VBD", "VBG", "VBN", "VBP", "VBZ":
			posCounts["Verb"]++
		case "JJ", "JJR", "JJS":
			posCounts["Adjective"]++
		case "RB", "RBR", "RBS":
			posCounts["Adverb"]++
		case "DT":
			posCounts["Determiner"]++
		case "CC":
			posCounts["Conjunction"]++
		case "IN":
			posCounts["Preposition"]++
		case "PRP", "PRP$", "WP", "WP$":
			posCounts["Pronoun"]++
		case ".", ",", ":", ";", "!", "?":
			posCounts["Punctuation"]++
		}
		totalWords++
	}

	posCounts["NamedEntity"] = len(doc.Entities())

	return posCounts, totalWords
}

// calculateSentenceProbability calculates the probability of a sentence being a valid English sentence.
func calculateSentenceProbability(posCounts map[string]int, totalWords int) float64 {
	if totalWords == 0 {
		return 0.0
	}

	probability := 0.0
	if posCounts["Noun"] > 0 && posCounts["Verb"] > 0 {
		probability = 50.0
		if posCounts["Adjective"] > 0 {
			probability += 10.0
		}
		if posCounts["Adverb"] > 0 {
			probability += 10.0
		}
		if posCounts["Determiner"] > 0 {
			probability += 5.0
		}
		if posCounts["Conjunction"] > 0 {
			probability += 5.0
		}
		if posCounts["Preposition"] > 0 {
			probability += 5.0
		}
		if posCounts["Pronoun"] > 0 {
			probability += 5.0
		}
		if posCounts["Punctuation"] > 0 {
			probability += 10.0
		}
		if posCounts["NamedEntity"] > 0 {
			probability += 5.0
		}
	}

	return probability
}

func writeContentsToOutputFile(outputFile string, outputBytes []byte) {
	for {
		fileMutex.Lock()
		file, openError := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if openError != nil {
			fmt.Printf("Failed to open file: %v\n", openError)
			fileMutex.Unlock()
			time.Sleep(100 * time.Millisecond) // Wait before retrying
			continue
		}

		_, writeErr := file.Write(outputBytes)
		if writeErr != nil {
			fmt.Printf("Failed to write to file: %v\n", writeErr)
			err := file.Close()
			if err != nil {
				fmt.Printf("Failed to close file: %v\n", err)
			}
			fileMutex.Unlock()
			time.Sleep(100 * time.Millisecond) // Wait before retrying
			continue
		}

		closeError := file.Close()
		if closeError != nil {
			fmt.Printf("Failed to close file: %v\n", closeError)
		}
		fileMutex.Unlock()
		break
	}
}
