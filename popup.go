package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	// "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

func createPopup(title string) fyne.Window {
	popup := fyne.CurrentApp().NewWindow(title)
	popup.Resize(fyne.NewSize(300, 200))

	x, y := robotgo.Location()
	// fmt.Printf("x:%d, y:%d\n", x, y)
	// y += 1000
	// fmt.Printf("x:%d, y:%d\n", x, y)

	pos := fyne.NewPos(float32(x), float32(y))
	popup.Canvas().Content().Move(pos)
	return popup
}

func showPopup(dict Dictionary, acronym string) {
	var content fyne.CanvasObject
	var popup fyne.Window
	if definitions, ok := dict[acronym]; !ok {
		popup = createPopup(fmt.Sprintf("Define: %s", acronym))
		content = createDefinitionPopup(popup, dict, acronym)
	} else {
		popup = createPopup(fmt.Sprintf("Lookup: %s", acronym))
		content = lookupPopup(popup, definitions, acronym)
	}
	popup.SetContent(content)
	popup.Show()
}

func lookupPopup(popup fyne.Window, definitions []Acronym, acronym string) fyne.CanvasObject {
	var definitionsText string
		for _, acro := range definitions {
			definitionsText += fmt.Sprintf("%s: %s\n", acro.Expanded, acro.Definition)
		}
	okButton := widget.NewButton("OK", func() {	
		popup.Close()
	})
	content := container.NewVBox(
		widget.NewLabel(definitionsText),
		okButton,
	)
	return content 
}

func createDefinitionPopup(popup fyne.Window, dict Dictionary, acronym string) fyne.CanvasObject {
	expandEntry := widget.NewEntry()
	expandEntry.SetPlaceHolder("Enter the expanded form")

	definitionEntry := widget.NewMultiLineEntry()
	definitionEntry.SetPlaceHolder("Enter the definition")

	addButton := widget.NewButton("Add", func() {
		if expandEntry.Text != "" && definitionEntry.Text != "" {
			newAcronym := Acronym{
				Expanded:   expandEntry.Text,
				Definition: definitionEntry.Text,
			}
			if _, ok := dict[acronym]; !ok {
				dict[acronym] = []Acronym{}
			}
			dict[acronym] = append(dict[acronym], newAcronym)
			tree.Refresh()
			err := saveDictionary(dict, "acronyms.json")
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Error",
					Content: "Failed to save dictionary: " + err.Error(),
				})
			}
		}
		popup.Close()
	})

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Add definition for: %s", acronym)),
		expandEntry,
		definitionEntry,
		addButton,
	)

	return content
}
