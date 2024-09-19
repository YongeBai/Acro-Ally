package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createPopup(title string) fyne.Window {
	popup := fyne.CurrentApp().NewWindow(title)		
	popup.SetPadded(false)
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
		content = createLookupPopup(popup, definitions)
	}
	popup.SetContent(content)
	popup.Resize(fyne.NewSize(400, 200))
	
	popup.Show()
}

func createLookupPopup(popup fyne.Window, definitions []Acronym) fyne.CanvasObject {
	var definitionsText string
	for _, acro := range definitions {
		definitionsText += fmt.Sprintf("**%s**: %s\n\n", acro.Expanded, acro.Definition)
	}

	richText := widget.NewRichTextFromMarkdown(definitionsText)
	richText.Wrapping = fyne.TextWrapWord
	
	okButton := widget.NewButton("OK", func() {	
		popup.Close()
	})

	content := container.NewVBox(
		richText,
		okButton,
	)	
	scrollContainer := container.NewVScroll(content)

	popup.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyReturn || ke.Name == fyne.KeyEnter {
			popup.Close()
		}
	})
	
	return scrollContainer
}

func createDefinitionPopup(popup fyne.Window, dict Dictionary, acronym string) fyne.CanvasObject {
	expandEntry := widget.NewEntry()
	expandEntry.SetPlaceHolder("Enter the expanded form")

	definitionEntry := widget.NewMultiLineEntry()
	definitionEntry.SetPlaceHolder("Enter the definition")
	definitionEntry.Wrapping = fyne.TextWrapWord

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

	popup.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyReturn || ke.Name == fyne.KeyEnter {
			popup.Close()
		}
	})

	return content
}
