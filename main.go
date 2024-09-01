package main

import (
	"fmt"	
	
	"fyne.io/fyne/v2"
	// "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var dict Dictionary

// func main() {
// 	dict = make(Dictionary)
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Acro-Ally")

// 	acronymAccordion := widget.NewAccordion()

// 	searchEntry := widget.NewEntry()
// 	searchEntry.SetPlaceHolder("Enter an acronym")
// 	searchEntry.OnSubmitted = func(s string) {
// 		acronymAccordion.Open(dict[s].Index)
// 	}


// 	content := container.NewBorder(
// 		container.NewVBox(
// 			searchEntry,
// 			widget.NewButton("Add Acronym", func() {
// 				showAddAcronymDialog(myWindow, acronymAccordion)
// 			}),
// 		),
// 		widget.NewButton("Exit", func() {
// 			myApp.Quit()
// 		}),
// 		nil,
// 		nil,
// 		container.NewVScroll(acronymAccordion),
// 	)

// 	myWindow.SetContent(content)
// 	myWindow.Resize(fyne.NewSize(600, 400))
// 	myWindow.ShowAndRun()
// }

func createAcronymAccordionItem(acronym string, entry Acronym) *widget.AccordionItem {
	return widget.NewAccordionItem(
		acronym,
		container.NewVBox(
			widget.NewRichTextFromMarkdown(
				fmt.Sprintf("## %s\n%s", entry.Expanded, entry.Definition)),
		),
	)
}

func showAddAcronymDialog(win fyne.Window, accordion *widget.Accordion) {
	acronymEntry := widget.NewEntry()
	acronymEntry.SetPlaceHolder("Enter the acronym")
	fmt.Println(acronymEntry.Text)

	expandEntry := widget.NewEntry()
	expandEntry.SetPlaceHolder("Enter the expanded form")

	definitionEntry := widget.NewEntry()
	definitionEntry.SetPlaceHolder("Enter the definition")

	dialog.ShowForm("Add Acronym", "Add", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Acronym", acronymEntry),
		widget.NewFormItem("Expanded", expandEntry),
		widget.NewFormItem("Definition", definitionEntry),
	}, func(add bool) {
		if acronymEntry.Text != "" && expandEntry.Text != "" {
			dict[acronymEntry.Text] = Acronym {
				Expanded:     expandEntry.Text,
				Definition: definitionEntry.Text,
				Index: len(dict),
			}
			newItem := createAcronymAccordionItem(acronymEntry.Text, dict[acronymEntry.Text])
			accordion.Append(newItem)			
			saveDictionary(dict, "acronyms.json")
		}
	}, win)
}