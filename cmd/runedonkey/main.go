package main

import (
	"flag"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lgstructs"
	"runer"
	"strings"
)

type Actions int

const (
	GemSum Actions = iota
	WordLength
	RuneLength
	RuneglishLength
	RunePattern
	RunePatternNoDoublet
)

type RuneDonkey struct {
	DB *gorm.DB
}

func (rd *RuneDonkey) GenerateExcelFromValues(value string, textType runer.TextType, whatToDo Actions, outputFile string) error {
	values := rd.GetValuesFromString(value, textType, whatToDo)
	f := excelize.NewFile()
	_, sheetError := f.NewSheet("Sheet1")
	if sheetError != nil {
		fmt.Printf(sheetError.Error())
		return sheetError
	}

	wordArray := rd.getWordsFromString(value)

	topval := fmt.Sprintf("Original:%s - Delimiters:%s - Words:%s", value, "•", strings.Join(wordArray, ","))
	f.SetCellValue("Sheet1", "A1", topval)
	f.MergeCell("Sheet1", "A1", "Z1")

	for colIndex, val := range values {
		words := rd.queryDatabase(whatToDo, val)
		cell := fmt.Sprintf("%c2", 'A'+colIndex)
		f.SetCellValue("Sheet1", cell, val)
		for rowIndex, word := range words {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, rowIndex+3)
			f.SetCellValue("Sheet1", cell, word)
		}
	}

	if err := f.SaveAs(outputFile); err != nil {
		return err
	}

	return nil
}

func (rd *RuneDonkey) GetValuesFromString(value string, textType runer.TextType, whatToDo Actions) []string {
	var valuesToGetFromDB []string
	wordArray := rd.getWordsFromString(value)

	for _, word := range wordArray {
		word = strings.TrimSpace(word)
		if textType == runer.Runeglish || textType == runer.Latin {
			word = runer.PrepLatinToRune(word)
			word = runer.TransposeLatinToRune(word)
		}

		switch whatToDo {
		case GemSum:
			sum := runer.CalculateGemSum(word, textType)
			valuesToGetFromDB = append(valuesToGetFromDB, fmt.Sprintf("%d", sum))
			break
		case WordLength, RuneLength, RuneglishLength:
			length := len(word)
			valuesToGetFromDB = append(valuesToGetFromDB, fmt.Sprintf("%d", length))
			break
		case RunePattern:
			pattern := lgstructs.GetRunePattern(word)
			valuesToGetFromDB = append(valuesToGetFromDB, pattern)
			break
		case RunePatternNoDoublet:
			word = lgstructs.RemoveDoublets(word)
			pattern := lgstructs.GetRunePattern(word)
			valuesToGetFromDB = append(valuesToGetFromDB, pattern)
			break
		}
	}

	return valuesToGetFromDB
}

func (rd *RuneDonkey) queryDatabase(field Actions, value string) []string {
	var results []string
	var rows *gorm.DB

	switch field {
	case WordLength:
		rows = rd.DB.Table("dictionary_words").Where("dict_word_length = ?", value).Select("dict_word")
		break
	case RuneLength:
		rows = rd.DB.Table("dictionary_words").Where("dict_rune_length = ?", value).Select("dict_word")
		break
	case RuneglishLength:
		rows = rd.DB.Table("dictionary_words").Where("dict_runeglish_length = ?", value).Select("dict_word")
		break
	case RunePattern:
		rows = rd.DB.Table("dictionary_words").Where("rune_pattern = ?", value).Select("dict_word")
		break
	case RunePatternNoDoublet:
		rows = rd.DB.Table("dictionary_words").Where("rune_pattern_no_doublet = ?", value).Select("dict_word")
		break
	default:
		rows = rd.DB.Table("dictionary_words").Where("gem_sum = ?", value).Select("dict_word")
		break
	}

	rows.Scan(&results)
	return results
}

func (rd *RuneDonkey) getWordsFromString(value string) []string {
	var words []string
	wordSplit := strings.Split(value, "•")
	for _, word := range wordSplit {
		words = append(words, word)
		word = ""
	}
	return words
}

func main() {
	// Define flags
	text := flag.String("text", "", "Text to process")
	outputFile := flag.String("output", "generated_excel.xlsx", "Output Excel file")
	flag.Parse()

	for port := 3306; port <= 3308; port++ {
		// Update the connection string with your MySQL credentials
		dsn := fmt.Sprintf("%s%d%s", "runedonkey:dpasswd@tcp(localhost:", port, ")/wordsdb?charset=utf8mb4&parseTime=True&loc=Local")
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		for action := GemSum; action <= RunePatternNoDoublet; action++ {
			runeDonkey := RuneDonkey{DB: db}

			corpus := ""
			switch port {
			case 3306:
				corpus = "dwyl"
				break
			case 3307:
				corpus = "tfcom"
				break
			case 3308:
				corpus = "pgb"
				break
			}

			actionName := ""
			switch action {
			case GemSum:
				actionName = "GemSum"
				break
			case WordLength:
				actionName = "WordLength"
				break
			case RuneLength:
				actionName = "RuneLength"
				break
			case RuneglishLength:
				actionName = "RuneglishLength"
				break
			case RunePattern:
				actionName = "RunePattern"
				break
			case RunePatternNoDoublet:
				actionName = "RunePatternNoDoublet"
				break
			}

			outfile := fmt.Sprintf("%s_%s_%s.xlsx", corpus, actionName, *outputFile)
			fmt.Println("Generating Excel for:", corpus, actionName)

			genError := runeDonkey.GenerateExcelFromValues(*text, runer.Runes, action, outfile)
			if genError != nil {
				fmt.Println("Error:", genError)
				return
			}
		}
	}

}
