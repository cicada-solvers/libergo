package main

import (
	runelib "characterrepo"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	charRepo := runelib.NewCharacterRepo()
	text := flag.String("text", "", "The text to sub out")
	outputFile := flag.String("output", "", "The output file to write the results")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *text == "" || *outputFile == "" {
		flag.Usage()
		return
	}

	alphabetArray := charRepo.GetGematriaRunes()
	alphabetArray = append(alphabetArray, "")
	alphabetArray = append(alphabetArray, "ᚠᚢ")
	alphabetArray = append(alphabetArray, "ᚩᚱᚢᛉ")
	alphabetArray = append(alphabetArray, "ᛄᛋᛈᛉ")
	alphabetArray = append(alphabetArray, "ᛟᛞᚪᚫ")

	for _, c := range alphabetArray {
		runeCounter := make(map[string]int)
		tmpText := ""

		if len(strings.Split(c, "")) > 0 {
			for _, r := range strings.Split(c, "") {
				tmpText = strings.ReplaceAll(*text, r, "")
			}
		} else {
			tmpText = strings.ReplaceAll(*text, c, "")
		}

		tmpArray := strings.Split(tmpText, "")

		for _, r := range tmpArray {
			if _, ok := runeCounter[r]; ok {
				runeCounter[r]++
			} else {
				runeCounter[r] = 1
			}
		}

		writeMapToExcel(runeCounter, *outputFile, c)
	}
}

func writeMapToExcel(mapOfRunes map[string]int, outputFileName string, runeVal string) {
	var (
		f        *excelize.File
		err      error
		existing bool
	)

	if _, err = os.Stat(outputFileName); err == nil {
		existing = true
		f, err = excelize.OpenFile(outputFileName)
		if err != nil {
			log.Printf("failed to open existing Excel file %q: %v", outputFileName, err)
			return
		}
	} else {
		f = excelize.NewFile()
		// Keep default sheet; we'll just add a new one.
	}

	sheetName := fmt.Sprintf("Del Rune %s", runeVal)
	idx, err := f.NewSheet(sheetName)
	if err != nil {
		log.Printf("failed to create new sheet %q: %v", sheetName, err)
		return
	}

	// Headers
	if err := f.SetCellValue(sheetName, "A1", "Rune"); err != nil {
		log.Printf("failed to set header: %v", err)
		return
	}
	if err := f.SetCellValue(sheetName, "B1", "Count"); err != nil {
		log.Printf("failed to set header: %v", err)
		return
	}

	// Write rows sorted by rune for stable output
	keys := make([]string, 0, len(mapOfRunes))
	for k := range mapOfRunes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	row := 2
	for _, k := range keys {
		if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), k); err != nil {
			log.Printf("failed to write rune at row %d: %v", row, err)
			return
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), mapOfRunes[k]); err != nil {
			log.Printf("failed to write count at row %d: %v", row, err)
			return
		}
		row++
	}

	f.SetActiveSheet(idx)

	if existing {
		if err := f.Save(); err != nil {
			log.Printf("failed to save existing Excel file: %v", err)
		}
	} else {
		if err := f.SaveAs(outputFileName); err != nil {
			log.Printf("failed to save new Excel file %q: %v", outputFileName, err)
		}
	}
}
