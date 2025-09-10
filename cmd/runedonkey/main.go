package main

import (
	"flag"
	"fmt"
	"liberdatabase"
	"runer"
	"strings"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Actions is an enumerated type representing various operations or instructions in the RuneDonkey application.
type Actions int

const (
	// GemSum represents an operation to calculate the gematria sum of a word.
	GemSum Actions = iota

	// WordLength represents an operation to calculate the length of a word in terms of the number of characters.
	WordLength

	// RuneLength represents an operation to calculate the number of runes in a string.
	RuneLength

	// RuneNoDoubletLength represents an operation to calculate the length of a string excluding consecutive duplicate runes.
	RuneNoDoubletLength

	// RuneglishLength represents an operation to determine the length of a word after converting it to Runeglish.
	RuneglishLength

	// RunePattern represents an operation for retrieving the pattern of runes in a string.
	RunePattern

	// RunePatternNoDoublet represents an operation for retrieving rune patterns excluding consecutive identical runes.
	RunePatternNoDoublet
)

// RuneDonkey is a struct that provides functionality for handling rune-based data and interacting with a database.
// DB is a pointer to a gorm.DB, which manages database operations related to rune processing.
type RuneDonkey struct {
	DB *gorm.DB
}

// GenerateExcelFromValues generates an Excel file based on input string data and operations, and saves it to the specified path.
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

// GetValuesFromString processes a given string and returns a list of values based on the specified text type and action.
func (rd *RuneDonkey) GetValuesFromString(value string, textType runer.TextType, whatToDo Actions) []string {
	var valuesToGetFromDB []string
	wordArray := rd.getWordsFromString(value)

	for _, word := range wordArray {
		word = strings.TrimSpace(word)
		if textType == runer.Runeglish || textType == runer.Latin {
			word = runer.PrepLatinToRune(word)
			word = runer.TransposeLatinToRune(word, false)
		}

		switch whatToDo {
		case GemSum:
			sum := runer.CalculateGemSum(word, textType, false)
			valuesToGetFromDB = append(valuesToGetFromDB, fmt.Sprintf("%d", sum))
			break
		case WordLength, RuneglishLength:
			length := len(word)
			valuesToGetFromDB = append(valuesToGetFromDB, fmt.Sprintf("%d", length))
			break
		case RuneLength, RuneNoDoubletLength:
			length := utf8.RuneCountInString(word)
			valuesToGetFromDB = append(valuesToGetFromDB, fmt.Sprintf("%d", length))
			break
		case RunePattern:
			pattern := liberdatabase.GetRunePattern(word)
			valuesToGetFromDB = append(valuesToGetFromDB, pattern)
			break
		case RunePatternNoDoublet:
			word = liberdatabase.RemoveDoublets(strings.Split(word, ""))
			pattern := liberdatabase.GetRunePattern(word)
			valuesToGetFromDB = append(valuesToGetFromDB, pattern)
			break
		}
	}

	return valuesToGetFromDB
}

// queryDatabase queries the database for words based on the specified operation (field) and value, returning a list of matching words.
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
	case RuneNoDoubletLength:
		rows = rd.DB.Table("dictionary_words").Where("dict_rune_no_doublet_length = ?", value).Select("dict_word")
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

// getWordsFromString splits the input string by the delimiter "•" and returns a slice of the resulting substrings.
func (rd *RuneDonkey) getWordsFromString(value string) []string {
	var words []string
	wordSplit := strings.Split(value, "•")
	for _, word := range wordSplit {
		words = append(words, word)
		word = ""
	}
	return words
}

// main initializes flags, connects to MySQL databases, performs operations on text, and generates Excel files with results.
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
			case RuneNoDoubletLength:
				actionName = "RuneNoDoubletLength"
				break
			case RunePattern:
				actionName = "RunePattern"
				break
			case RunePatternNoDoublet:
				actionName = "RunePatternNoDoublet"
				break
			}

			outfile := fmt.Sprintf("%s_%s_%s", corpus, actionName, *outputFile)
			fmt.Println("Generating Excel for:", corpus, actionName)

			genError := runeDonkey.GenerateExcelFromValues(*text, runer.Runes, action, outfile)
			if genError != nil {
				fmt.Println("Error:", genError)
				return
			}
		}

		// Get the underlying sql.DB object and defer its closure
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get database object")
		}
		cerr := sqlDB.Close()
		if cerr != nil {
			return
		}
	}
}
