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

	w.SetContent(content)
	w.ShowAndRun()
}
