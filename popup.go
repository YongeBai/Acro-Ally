package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

func createPopup(title string, content fyne.CanvasObject) fyne.Window {
	popup := fyne.CurrentApp().NewWindow(title)
	popup.SetContent(content)
	popup.Resize(fyne.NewSize(300, 200))

	x, y := robotgo.Location()
	y -= 200
	pos := fyne.NewPos(float32(x), float32(y))
	popup.Canvas().Content().Move(pos)
	popup.Show()
	return popup
}

func showPopup(dict Dictionary, acronym string) {
	var content fyne.CanvasObject
	if definitions, ok := dict[acronym]; !ok {
		// TODO: show set definition popup
		fmt.Println("Acronym not found in dictionary")
	} else {
		var definitionsText string
		for _, acro := range definitions {
			definitionsText += fmt.Sprintf("%s: %s\n", acro.Expanded, acro.Definition)
		}
		content = widget.NewLabel(definitionsText)
	}
	createPopup(fmt.Sprintf("Lookup: %s", acronym), content)
}