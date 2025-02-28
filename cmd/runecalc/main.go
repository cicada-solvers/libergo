package main

import (
	"characterrepo"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"math/big"
	"os"
	"sequences"
)

func main() {
	a := app.New()
	w := a.NewWindow("Rune Calculator")

	repo := runelib.NewCharacterRepo()
	displayLabel := widget.NewLabel("Runes")
	displayText := widget.NewLabel("")

	latinLabel := widget.NewLabel("Latin")
	latinText := widget.NewLabel("")

	var gemValue = int64(0)
	gemLabel := widget.NewLabel("Gematria Sum")
	gemText := widget.NewLabel("")

	var gemProdValue = big.NewInt(0)
	gemProdLabel := widget.NewLabel("Gematria Product")
	gemProdText := widget.NewLabel("")

	buttons := make([]*widget.Button, 33)

	// clear Button
	clearButtonLabel := "CLR"
	buttons[0] = widget.NewButton(clearButtonLabel, func() {
		displayText.SetText("")
		latinText.SetText("")
		gemValue = int64(0)
		gemProdValue.SetInt64(int64(0))
		gemText.SetText("")
		gemProdText.SetText("")
	})

	primeSequence, _ := sequences.GetPrimeSequence(big.NewInt(int64(109)), false)
	btnCounter := 1
	for _, num := range primeSequence.Sequence {
		value := int(num.Int64())
		labelOne := repo.GetRuneFromValue(value)
		labelTwo := repo.GetCharFromRune(labelOne)

		buttonLabel := fmt.Sprintf("%s \\ %s", labelOne, labelTwo)
		buttons[btnCounter] = widget.NewButton(buttonLabel, func() {
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

			if gemProdValue.Int64() == 0 {
				gemProdValue.SetInt64(int64(value))
			} else {
				gemProdValue.Mul(gemProdValue, big.NewInt(int64(value)))
			}
			gemProdText.SetText(gemProdValue.String())
		})

		btnCounter++
	}

	// space Button
	spaceButtonLabel := "•"
	buttons[30] = widget.NewButton(spaceButtonLabel, func() {
		tmpText := displayText.Text
		tmpText = tmpText + "•"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text
		latinTmpText = latinTmpText + " "
		latinText.SetText(latinTmpText)
	})

	// tick Button
	tickButtonLabel := "'"
	buttons[31] = widget.NewButton(tickButtonLabel, func() {
		tmpText := displayText.Text
		tmpText = tmpText + "'"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text
		latinTmpText = latinTmpText + "'"
		latinText.SetText(latinTmpText)
	})

	// period Button
	periodButtonLabel := "⊹"
	buttons[32] = widget.NewButton(periodButtonLabel, func() {
		tmpText := displayText.Text
		tmpText = tmpText + "⊹"
		displayText.SetText(tmpText)

		latinTmpText := latinText.Text
		latinTmpText = latinTmpText + "."
		latinText.SetText(latinTmpText)
	})

	// Convert []*widget.Button to []fyne.CanvasObject
	buttonObjects := make([]fyne.CanvasObject, len(buttons))
	for i, btn := range buttons {
		buttonObjects[i] = btn
	}

	display := container.NewHBox(displayLabel, displayText)
	latin := container.NewHBox(latinLabel, latinText)
	gemSumBox := container.NewHBox(gemLabel, gemText)
	gemProdBox := container.NewHBox(gemProdLabel, gemProdText)
	grid := container.NewGridWithColumns(4, buttonObjects...)
	content := container.NewVBox(display, latin, gemSumBox, gemProdBox, grid)

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
					_, err := writer.Write([]byte(fmt.Sprintf("Runes: %s\nLatin: %s\nGematria Sum: %d\nGematria Product: %s",
						displayText.Text, latinText.Text, gemValue, gemProdText.Text)))
					if err != nil {
						return
					}
				}
			}, w)
		}),
		fyne.NewMenuItem("Exit", func() {
			a.Quit()
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
		fyne.NewMenuItem("Copy Gematria Product", func() {
			w.Clipboard().SetContent(gemProdText.Text)
		}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, copyMenu)
	w.SetMainMenu(mainMenu)

	w.SetContent(content)
	w.ShowAndRun()
}
