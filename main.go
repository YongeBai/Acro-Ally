package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/robotn/gohook"
	"github.com/go-vgo/robotgo"
)

var dict Dictionary

func main() {
	var err error
	dict, err = loadDictionary("acronyms.json")
	if err != nil {
		fmt.Println("No dictionary found, creating new one:", err)
		dict = make(Dictionary)
	}	
	fmt.Println(dict)
	
	myApp := app.New()
	myWindow := myApp.NewWindow("Acro-Ally")

	tree := createAcronymTree(dict)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for an acronym")
	searchEntry.OnSubmitted = func(text string) {
		lookUpOrDefine(myWindow, tree, dict, text)
	}


	content := container.NewBorder(
		container.NewVBox(
			searchEntry,
			widget.NewButton("Add Acronym", func() {
				addAcronymButton(myWindow, tree, dict)
			}),
		),
		widget.NewButton("Exit", func() {
			myApp.Quit()
		}),
		nil,
		nil,
		container.NewVScroll(tree),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 400))

	go setupGlobalHotkeys(myWindow, tree, dict)

	myWindow.ShowAndRun()
}


func simulateCopy() {
	robotgo.KeyTap("c", "Control")
}

func setupGlobalHotkeys(win fyne.Window, tree *widget.Tree, dict Dictionary) {
	hook.Register(hook.KeyDown, []string{"ctrl", "alt", "a"}, func(e hook.Event) {
		simulateCopy()
		text, err := clipboard.ReadAll()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if text != "" {
			lookUpOrDefine(win, tree, dict, text)
		}
	})

	s := hook.Start()
	<-hook.Process(s)
}

func lookUpOrDefine(win fyne.Window, tree *widget.Tree, dict Dictionary, acronym string) {
	if _, ok := dict[acronym]; !ok {
		addAcronym(win, tree, dict, acronym)
	} else {
		var definitions string
		for _, acro := range dict[acronym] {
			definitions += fmt.Sprintf("%s: %s\n", acro.Expanded, acro.Definition)
		}
		dialog.ShowInformation(
			fmt.Sprintf("Acronym %s found", acronym),
			definitions,
			win,
		)
	}
}