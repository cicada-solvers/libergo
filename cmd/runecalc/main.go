package main

import (
	"characterrepo"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"math/big"
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

	buttons := make([]*widget.Button, 29)

	primeSequence, _ := sequences.GetPrimeSequence(big.NewInt(int64(109)), false)
	btnCounter := 0
	for _, num := range primeSequence.Sequence {
		value := int(num.Int64())
		buttonLabel := repo.GetRuneFromValue(value)
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
		})

		btnCounter++
	}

	// Convert []*widget.Button to []fyne.CanvasObject
	buttonObjects := make([]fyne.CanvasObject, len(buttons))
	for i, btn := range buttons {
		buttonObjects[i] = btn
	}

	display := container.NewHBox(displayLabel, displayText)
	latin := container.NewHBox(latinLabel, latinText)
	gemSumBox := container.NewHBox(gemLabel, gemText)
	grid := container.NewGridWithColumns(4, buttonObjects...)
	content := container.NewVBox(display, latin, gemSumBox, grid)

	w.SetContent(content)
	w.ShowAndRun()
}
