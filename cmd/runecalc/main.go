package main

import (
	"characterrepo"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"runer"
	"sequences"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

func main() {
	a := app.NewWithID("com.cmbsolver.runecalc")
	a.Settings().SetTheme(&runeTheme{})
	w := a.NewWindow("Rune Calculator")

	// Create the new controls
	calculateTextLabel := widget.NewLabel("Text To Calculate:")
	entry := widget.NewEntry()
	entry.Resize(fyne.NewSize(300, entry.MinSize().Height)) // Set the width to 300

	repo := runelib.NewCharacterRepo()
	displayLabel := widget.NewLabel("Runes:")
	displayText := widget.NewLabel("")

	latinLabel := widget.NewLabel("Runeglish:")
	latinText := widget.NewLabel("")

	var gemValue int64
	gemLabel := widget.NewLabel("Gematria Sum:")
	gemText := widget.NewLabel("")
	gemPrimeCheckbox := widget.NewCheck("Is Prime", nil)
	gemSumEmirpCheckbox := widget.NewCheck("Is Emirp", nil)

	valuesLabel := widget.NewLabel("Rune Values:")
	wordValuesLabel := widget.NewLabel("Word Sums:")
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
			tmpText := displayText.Text + displayRune
			displayText.SetText(tmpText)

			latinRune := repo.GetCharFromRune(displayRune)
			latinTmpText := latinText.Text + latinRune
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
		entry.SetText("")
		values = nil
		wordValues = nil
	})

	// space Button
	spaceButtonLabel := "•"
	specialButtons[1] = widget.NewButton(spaceButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text + "•"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text + " "
		latinText.SetText(latinTmpText)
	})

	// tick Button
	tickButtonLabel := "'"
	specialButtons[2] = widget.NewButton(tickButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text + "'"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text + "'"
		latinText.SetText(latinTmpText)
	})

	// double tick Button
	doubleTickButtonLabel := "\""
	specialButtons[3] = widget.NewButton(doubleTickButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text + "\""
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text + "\""
		latinText.SetText(latinTmpText)
	})

	// period Button
	periodButtonLabel := "⊹"
	specialButtons[4] = widget.NewButton(periodButtonLabel, func() {
		mu.Lock()
		defer mu.Unlock()

		tmpText := displayText.Text + "⊹"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text + "."
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

	// Create copy buttons
	copyDisplayButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(displayText.Text)
	})
	copyDisplayButton.Importance = widget.LowImportance

	copyLatinButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(latinText.Text)
	})
	copyLatinButton.Importance = widget.LowImportance

	copyGemButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(gemText.Text)
	})
	copyGemButton.Importance = widget.LowImportance

	copyValuesButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(valuesText.Text)
	})
	copyValuesButton.Importance = widget.LowImportance

	copyWordValuesButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(wordValuesText.Text)
	})
	copyWordValuesButton.Importance = widget.LowImportance

	display := container.NewBorder(nil, nil, container.NewHBox(copyDisplayButton, displayLabel), backspaceButton, displayText)
	latin := container.NewHBox(copyLatinButton, latinLabel, latinText)
	gemSumBox := container.NewHBox(copyGemButton, gemLabel, gemText, gemPrimeCheckbox, gemSumEmirpCheckbox)
	valuesContainer := container.NewHBox(copyValuesButton, valuesLabel, valuesText)
	wordValuesContainer := container.NewHBox(copyWordValuesButton, wordValuesLabel, wordValuesText)

	options := []string{"From Latin", "From Runes"}
	combo := widget.NewSelect(options, nil)
	loadButton := widget.NewButton("Load", func() {
		selectedOption := combo.Selected
		if selectedOption == "From Latin" {
			latinText.SetText(runer.PrepLatinToRune(strings.ToUpper(entry.Text)))
			runes := runer.TransposeLatinToRune(latinText.Text)
			displayText.SetText(runes)
		} else if selectedOption == "From Runes" {
			displayText.SetText(entry.Text)
			latinText.SetText(runer.TransposeRuneToLatin(entry.Text))
		}

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

	controls := container.NewBorder(calculateTextLabel, nil, nil, container.NewHBox(combo, loadButton), entry)

	// Create the About button
	aboutButton := widget.NewButtonWithIcon("About", theme.HelpIcon(), func() {
		parsedURL, err := url.Parse("https://github.com/cmbsolver/libergo")
		if err == nil {
			_ = fyne.CurrentApp().OpenURL(parsedURL)
		}
	})
	aboutButton.Importance = widget.LowImportance

	content := container.NewVBox(
		container.NewBorder(nil, nil, aboutButton, nil),
		controls,
		display,
		latin,
		gemSumBox,
		wordValuesContainer,
		valuesContainer,
		specialButtonsGrid,
		container.NewGridWithColumns(4, buttonObjects...),
	)

	w.SetContent(content)
	w.ShowAndRun()
}

type ButtonInfo struct {
	ButtonLabel string
	LabelTwo    string
	Value       int
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
