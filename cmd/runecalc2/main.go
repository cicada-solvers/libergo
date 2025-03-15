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
	"strings"
)

func main() {
	a := app.NewWithID("com.cmbsolver.runecalctwo")
	a.Settings().SetTheme(&runeTheme{})
	w := a.NewWindow("Lite Rune Calculator")

	// Create the About button
	aboutButton := widget.NewButtonWithIcon("Source Code", theme.HelpIcon(), func() {
		parsedURL, err := url.Parse("https://github.com/cmbsolver/libergo")
		if err == nil {
			_ = fyne.CurrentApp().OpenURL(parsedURL)
		}
	})
	aboutButton.Importance = widget.LowImportance

	redirectButton := widget.NewButtonWithIcon("cmbsolver.com", theme.HomeIcon(), func() {
		parsedURL, err := url.Parse("https://cmbsolver.com")
		if err == nil {
			_ = fyne.CurrentApp().OpenURL(parsedURL)
		}
	})
	redirectButton.Importance = widget.LowImportance

	// Create the new controls
	entry := widget.NewEntry()
	entry.Resize(fyne.NewSize(300, entry.MinSize().Height*3)) // Set the width to 300
	entry.MultiLine = true
	entry.SetPlaceHolder("Text to Calculate")

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

	webButtons := container.NewHBox(redirectButton, aboutButton)
	display := container.NewHBox(copyDisplayButton, displayLabel, displayText)
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

	controls := container.NewBorder(nil, nil, nil, container.NewHBox(combo, loadButton), entry)

	content := container.NewVBox(
		webButtons,
		controls,
		display,
		latin,
		gemSumBox,
		wordValuesContainer,
		valuesContainer,
	)

	w.SetContent(content)
	w.ShowAndRun()
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
