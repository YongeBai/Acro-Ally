package main

import (
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

func createPopup(content fyne.CanvasObject) {
	newWindow := fyne.CurrentApp().NewWindow("Acronym Popup")
	popup := widget.NewPopUp(content, newWindow.Canvas())
	popup.Resize(fyne.NewSize(300, 200))

	x, y := robotgo.Location()

	pos := fyne.NewPos(float32(x), float32(y))
	popup.ShowAtPosition(pos)
}

