package charindex

import (
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"os"
	"path/filepath"
	"runer"
	"strings"
)

func IndexCharactersFromDirectory(directory string) error {
	db, connError := liberdatabase.InitConnection()
	if connError != nil {
		fmt.Println("Error connecting to database: ", connError)
	}

	return readDirectoryContents(db, directory)
}

func readDirectoryContents(db *gorm.DB, directory string) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".bmp" && ext != ".tiff" && ext != ".webp" && ext != ".svg" && ext != ".pdf" && ext != ".zip" && ext != ".rar" && ext != ".7z" && ext != ".tar" && ext != ".gz" && ext != ".bz2" && ext != ".xz" && ext != ".mov" && ext != ".mp4" && ext != ".avi" && ext != ".mkv" && ext != ".mp3" && ext != ".wav" && ext != ".flac" && ext != ".ogg" && ext != ".wma" && ext != ".aac" && ext != ".m4a" && ext != ".opus" && ext != ".webm" && ext != ".flv" && ext != ".wmv" {
			return readAndIndexFileContents(db, path)
		}
		return nil
	})
	return err
}

func readAndIndexFileContents(db *gorm.DB, file string) error {
	lines, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	var textDocument = liberdatabase.TextDocument{
		FileName: filepath.Base(file),
	}
	docId, insertError := liberdatabase.InsertTextDocument(db, &textDocument)
	if insertError != nil {
		fmt.Println("Error inserting text document: ", insertError)
	}

	// textDocumentCharacters
	var textDocumentCharacters []liberdatabase.TextDocumentCharacter
	for _, line := range strings.Split(string(lines), "\n") {
		for _, char := range strings.ToUpper(line) {
			charStr := string(char)
			found := false
			for i, tdc := range textDocumentCharacters {
				if tdc.Character == charStr {
					textDocumentCharacters[i].Count++
					found = true
					break
				}
			}
			if !found {
				textDocumentCharacters = append(textDocumentCharacters, liberdatabase.TextDocumentCharacter{
					Character:      charStr,
					Count:          1,
					TextDocumentId: int64(docId),
				})
			}
		}
	}

	for _, tdc := range textDocumentCharacters {
		_, err := liberdatabase.InsertTextDocumentCharacter(db, &tdc)
		if err != nil {
			fmt.Println("Error inserting text document character: ", err)
		}
	}

	// runeglish characters
	var liberTextDocumentCharacters []liberdatabase.LiberTextDocumentCharacter
	for _, line := range strings.Split(string(lines), "\n") {
		runeglishLine := runer.PrepLatinToRune(line)
		for _, char := range runeglishLine {
			charStr := string(char)
			found := false
			for i, tdc := range liberTextDocumentCharacters {
				if tdc.Character == charStr {
					liberTextDocumentCharacters[i].Count++
					found = true
					break
				}
			}
			if !found {
				liberTextDocumentCharacters = append(liberTextDocumentCharacters,
					liberdatabase.LiberTextDocumentCharacter{
						Character:      charStr,
						Count:          1,
						TextDocumentId: int64(docId),
					})
			}
		}
	}

	for _, tdc := range liberTextDocumentCharacters {
		_, err := liberdatabase.InsertLiberTextDocumentCharacter(db, &tdc)
		if err != nil {
			fmt.Println("Error inserting text document character: ", err)
		}
	}

	// rune characters
	var runeTextDocumentCharacters []liberdatabase.RuneTextDocumentCharacter
	for _, line := range strings.Split(string(lines), "\n") {
		runeglishLine := runer.PrepLatinToRune(line)
		runeLine := runer.TransposeLatinToRune(runeglishLine)
		for _, char := range runeLine {
			charStr := string(char)
			found := false
			for i, tdc := range runeTextDocumentCharacters {
				if tdc.Character == charStr {
					runeTextDocumentCharacters[i].Count++
					found = true
					break
				}
			}
			if !found {
				runeTextDocumentCharacters = append(runeTextDocumentCharacters,
					liberdatabase.RuneTextDocumentCharacter{
						Character:      charStr,
						Count:          1,
						TextDocumentId: int64(docId),
					})
			}
		}
	}

	for _, tdc := range runeTextDocumentCharacters {
		_, err := liberdatabase.InsertRuneTextDocumentCharacter(db, &tdc)
		if err != nil {
			fmt.Println("Error inserting text document character: ", err)
		}
	}

	return nil
}
