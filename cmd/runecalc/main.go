package main

import (
	"characterrepo"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"os"
	"runer"
	"sequences"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

type ButtonInfo struct {
	ButtonLabel string
	LabelTwo    string
	Value       int
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&runeTheme{})
	w := a.NewWindow("Rune Calculator")

	repo := runelib.NewCharacterRepo()
	displayLabel := widget.NewLabel("Runes")
	displayText := widget.NewLabel("")

	latinLabel := widget.NewLabel("Latin")
	latinText := widget.NewLabel("")

	var gemValue int64
	gemLabel := widget.NewLabel("Gematria Sum")
	gemText := widget.NewLabel("")
	gemPrimeCheckbox := widget.NewCheck("Is Prime", nil)
	gemSumEmirpCheckbox := widget.NewCheck("Is Emirp", nil)

	valuesLabel := widget.NewLabel("Values:")
	wordValuesLabel := widget.NewLabel("Word Values:")
	valuesText := widget.NewLabel("")
	wordValuesText := widget.NewLabel("")

	var values []int64
	var wordValues []int64

	buttons := make([]*widget.Button, 29)
	specialButtons := make([]*widget.Button, 5)

	primeSequence, _ := sequences.GetPrimeSequence64(int64(109), false)
	btnCounter := 0

	var buttonInfos []ButtonInfo

	for _, num := range primeSequence.Sequence {
		value := int(num)
		labelOne := repo.GetRuneFromValue(value)
		labelTwo := repo.GetCharFromRune(labelOne)

		buttonLabel := fmt.Sprintf("%s \\ %s", labelOne, labelTwo)
		buttonInfos = append(buttonInfos, ButtonInfo{
			ButtonLabel: buttonLabel,
			LabelTwo:    labelTwo,
			Value:       value,
		})
	}

	sort.Slice(buttonInfos, func(i, j int) bool {
		return buttonInfos[i].LabelTwo < buttonInfos[j].LabelTwo
	})

	var mu sync.Mutex

	for _, buttonInfo := range buttonInfos {
		value := buttonInfo.Value
		buttons[btnCounter] = widget.NewButton(buttonInfo.ButtonLabel, func() {
			mu.Lock()
			defer mu.Unlock()

			displayRune := repo.GetRuneFromValue(value)
			tmpText := displayText.Text
			tmpText = tmpText + displayRune
			displayText.SetText(tmpText)

			latinRune := repo.GetCharFromRune(displayRune)
			latinTmpText := latinText.Text
			latinTmpText = latinTmpText + latinRune
			latinText.SetText(latinTmpText)

			gemValue += int64(value)
			gemText.SetText(fmt.Sprintf("%d", gemValue))
			gemPrimeCheckbox.SetChecked(sequences.IsPrime64(gemValue))
			gemSumEmirpCheckbox.SetChecked(sequences.IsEmirp64(gemValue))

			if value > 0 {
				values = append(values, int64(value))
			}

			valuesText.SetText(fmt.Sprintf("%v", values))

			wordValues = calculateWordGemSums(displayText.Text, repo)
			wordValuesText.SetText(fmt.Sprintf("%v", wordValues))
		})

		btnCounter++
	}

	// clear Button
	clearButtonLabel := "CLR"
	specialButtons[0] = widget.NewButton(clearButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		displayText.SetText("")
		latinText.SetText("")
		gemValue = int64(0)
		gemText.SetText("")
		gemPrimeCheckbox.SetChecked(false)
		gemSumEmirpCheckbox.SetChecked(false)
		valuesText.SetText("")
		wordValuesText.SetText("")
		values = nil
		wordValues = nil
	})

	// space Button
	spaceButtonLabel := "•"
	specialButtons[1] = widget.NewButton(spaceButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text
		tmpText = tmpText + "•"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text
		latinTmpText = latinTmpText + " "
		latinText.SetText(latinTmpText)
	})

	// tick Button
	tickButtonLabel := "'"
	specialButtons[2] = widget.NewButton(tickButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text
		tmpText = tmpText + "'"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text
		latinTmpText = latinTmpText + "'"
		latinText.SetText(latinTmpText)
	})

	// double tick Button
	doubleTickButtonLabel := "\""
	specialButtons[3] = widget.NewButton(doubleTickButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text
		tmpText = tmpText + "\""
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text
		latinTmpText = latinTmpText + "\""
		latinText.SetText(latinTmpText)
	})

	// period Button
	periodButtonLabel := "⊹"
	specialButtons[4] = widget.NewButton(periodButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text
		tmpText = tmpText + "⊹"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text
		latinTmpText = latinTmpText + "."
		latinText.SetText(latinTmpText)
	})

	// Backspace Button
	backspaceButton := widget.NewButtonWithIcon("", theme.ContentUndoIcon(), func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text
		if utf8.RuneCountInString(tmpText) > 0 {
			_, size := utf8.DecodeLastRuneInString(tmpText)
			tmpText = tmpText[:len(tmpText)-size]
			displayText.SetText(tmpText)
		}

		latinText.SetText(runer.TransposeRuneToLatin(tmpText))

		gemValue = runer.CalculateGemSum(displayText.Text, runer.Runes)
		gemText.SetText(fmt.Sprintf("%d", gemValue))
		gemPrimeCheckbox.SetChecked(sequences.IsPrime64(gemValue))
		gemSumEmirpCheckbox.SetChecked(sequences.IsEmirp64(gemValue))

		values = nil
		for _, runeCharacter := range displayText.Text {
			runeValue := int64(repo.GetValueFromRune(string(runeCharacter)))
			if runeValue > 0 {
				values = append(values, runeValue)
			}
		}

		valuesText.SetText(fmt.Sprintf("%v", values))

		wordValues = calculateWordGemSums(displayText.Text, repo)
		wordValuesText.SetText(fmt.Sprintf("%v", wordValues))
	})
	backspaceButton.Importance = widget.LowImportance

	// Convert []*widget.Button to []fyne.CanvasObject
	buttonObjects := make([]fyne.CanvasObject, len(buttons))
	for i, btn := range buttons {
		buttonObjects[i] = btn
	}

	// Convert []*widget.Button to []fyne.CanvasObject
	specialButtonObjects := make([]fyne.CanvasObject, len(specialButtons))
	for i, btn := range specialButtons {
		specialButtonObjects[i] = btn
	}

	// Create a new grid for the specified buttons
	specialButtonsGrid := container.NewGridWithColumns(5, specialButtonObjects...)

	display := container.NewBorder(nil, nil, displayLabel, backspaceButton, displayText)
	latin := container.NewHBox(latinLabel, latinText)
	gemSumBox := container.NewBorder(nil, nil, nil, nil,
		container.NewHBox(gemLabel, gemText, gemPrimeCheckbox, gemSumEmirpCheckbox))
	valuesContainer := container.NewHBox(valuesLabel, valuesText)
	wordValuesContainer := container.NewHBox(wordValuesLabel, wordValuesText)

	content := container.NewVBox(
		display,
		latin,
		gemSumBox,
		valuesContainer,
		wordValuesContainer,
		specialButtonsGrid,
		container.NewGridWithColumns(4, buttonObjects...),
	)

	// Create the "About" window
	aboutWindow := func() {
		about := a.NewWindow("About")
		aboutLabel := widget.NewLabel("Written by cmbsolver")
		hyperlink := widget.NewHyperlink("https://github.com/cmbsolver/libergo", parseURL("https://github.com/cmbsolver/libergo"))
		aboutContent := container.NewVBox(aboutLabel, hyperlink)
		about.SetContent(aboutContent)
		about.Resize(fyne.NewSize(300, 100))
		about.Show()
	}

	// Create the menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Save", func() {
			dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err == nil && writer != nil {
					defer func(writer fyne.URIWriteCloser) {
						err := writer.Close()
						if err != nil {
							fmt.Println("Error closing file:", err)
							os.Exit(1)
						}
					}(writer)
					content := fmt.Sprintf("Runes: %s\nLatin: %s\nGematria Sum: %d\nIs Prime: %t\nIs Emirp: %t\nValues: %v\nWord Values: %v\n",
						displayText.Text, latinText.Text, gemValue, gemPrimeCheckbox.Checked, gemSumEmirpCheckbox.Checked, values, wordValues)
					_, err := writer.Write([]byte(content))
					if err != nil {
						return
					}
				}
			}, w)
		}),
	)

	copyMenu := fyne.NewMenu("Copy",
		fyne.NewMenuItem("Copy Runes", func() {
			w.Clipboard().SetContent(displayText.Text)
		}),
		fyne.NewMenuItem("Copy Latin", func() {
			w.Clipboard().SetContent(latinText.Text)
		}),
		fyne.NewMenuItem("Copy Gematria Sum", func() {
			w.Clipboard().SetContent(gemText.Text)
		}),
		fyne.NewMenuItem("Copy Values", func() {
			w.Clipboard().SetContent(valuesText.Text)
		}),
		fyne.NewMenuItem("Copy Word Values", func() {
			w.Clipboard().SetContent(wordValuesText.Text)
		}),
	)

	loadMenu := fyne.NewMenu("Load",
		fyne.NewMenuItem("Load from Latin", func() {
			entry := widget.NewEntry()
			dialog.ShowForm("Enter Latin Text", "Load", "Cancel", []*widget.FormItem{
				widget.NewFormItem("Latin Text", entry),
			}, func(b bool) {
				if b {
					entry.Text = strings.ToUpper(entry.Text)
					latinText.SetText(runer.PrepLatinToRune(entry.Text))
					runes := runer.TransposeLatinToRune(latinText.Text)
					displayText.SetText(runes)
					gemValue = runer.CalculateGemSum(runes, runer.Runes)
					gemText.SetText(fmt.Sprintf("%d", gemValue))
					gemPrimeCheckbox.SetChecked(sequences.IsPrime64(gemValue))
					gemSumEmirpCheckbox.SetChecked(sequences.IsEmirp64(gemValue))

					values = nil
					for _, runeCharacter := range runes {
						runeValue := int64(repo.GetValueFromRune(string(runeCharacter)))
						if runeValue > 0 {
							values = append(values, runeValue)
						}
					}

					valuesText.SetText(fmt.Sprintf("%v", values))

					wordValues = calculateWordGemSums(displayText.Text, repo)
					wordValuesText.SetText(fmt.Sprintf("%v", wordValues))
				}
			}, w)
		}),
		fyne.NewMenuItem("Load from Runes", func() {
			entry := widget.NewEntry()
			dialog.ShowForm("Enter Runes", "Load", "Cancel", []*widget.FormItem{
				widget.NewFormItem("Runes", entry),
			}, func(b bool) {
				if b {
					displayText.SetText(entry.Text)
					latinText.SetText(runer.TransposeRuneToLatin(entry.Text))

					gemValue = runer.CalculateGemSum(displayText.Text, runer.Runes)
					gemText.SetText(fmt.Sprintf("%d", gemValue))
					gemPrimeCheckbox.SetChecked(sequences.IsPrime64(gemValue))
					gemSumEmirpCheckbox.SetChecked(sequences.IsEmirp64(gemValue))

					values = nil
					for _, runeCharacter := range displayText.Text {
						runeValue := int64(repo.GetValueFromRune(string(runeCharacter)))
						if runeValue > 0 {
							values = append(values, runeValue)
						}
					}

					valuesText.SetText(fmt.Sprintf("%v", values))

					wordValues = calculateWordGemSums(displayText.Text, repo)
					wordValuesText.SetText(fmt.Sprintf("%v", wordValues))
				}
			}, w)
		}),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			aboutWindow()
		}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, copyMenu, loadMenu, helpMenu)
	w.SetMainMenu(mainMenu)

	w.SetContent(content)
	w.ShowAndRun()
}

// Helper function to parse URL
func parseURL(urlStr string) *url.URL {
	uri, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil
	}
	return uri
}

func calculateWordGemSums(sentence string, repo *runelib.CharacterRepo) []int64 {
	var values []int64
	var wordValues []int64

	for _, runeCharacter := range sentence {
		runeValue := int64(repo.GetValueFromRune(string(runeCharacter)))

		skip := false
		if string(runeCharacter) == "'" || string(runeCharacter) == "\"" {
			skip = true
		}

		if !skip {
			if runeValue > 0 {
				values = append(values, runeValue)
			} else {
				wordSum := int64(0)
				for _, value := range values {
					wordSum += value
				}

				if wordSum > 0 {
					wordValues = append(wordValues, wordSum)
				}

				values = nil
			}
		}
	}

	wordSum := int64(0)
	for _, value := range values {
		wordSum += value
	}

	if wordSum > 0 {
		wordValues = append(wordValues, wordSum)
	}

	return wordValues
}
