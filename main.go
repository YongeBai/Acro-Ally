package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
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
	mainWindow := myApp.NewWindow("Acro-Ally")
	mainWindow.SetMaster()

	tree := createAcronymTree(dict)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for an acronym")
	searchEntry.OnSubmitted = func(text string) {
		lookUpOrDefine(mainWindow, tree, dict, text)
	}


	content := container.NewBorder(
		container.NewVBox(
			searchEntry,
			widget.NewButton("Add Acronym", func() {
				addAcronymButton(mainWindow, tree, dict)
			}),
		),
		widget.NewButton("Exit", func() {
			myApp.Quit()
		}),
		nil,
		nil,
		container.NewVScroll(tree),
	)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(600, 400))

	go setupGlobalHotkeys(mainWindow, tree, dict)

	mainWindow.ShowAndRun()
}

// These are for when the user is in main window, its fine for now
func lookUpOrDefine(win fyne.Window, tree *widget.Tree, dict Dictionary, acronym string) {
	fmt.Println("Looking up or defining:", acronym)
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

func addAcronym(win fyne.Window, tree *widget.Tree, dict Dictionary, acronym string) {

	expandEntry := widget.NewEntry()
	expandEntry.SetPlaceHolder("Enter the expanded form")

	definitionEntry := widget.NewEntry()
	definitionEntry.SetPlaceHolder("Enter the definition")

	dialog.ShowForm(fmt.Sprintf("Add Acronym: %s", acronym), "Add", "Cancel", []*widget.FormItem{		
		widget.NewFormItem("Expanded", expandEntry),
		widget.NewFormItem("Definition", definitionEntry),
	}, func(add bool) {
		if add && expandEntry.Text != "" && definitionEntry.Text != "" {
			newAcronym := Acronym{
				Expanded:   expandEntry.Text,
				Definition: definitionEntry.Text,
			}
			if _, ok := dict[acronym]; !ok {
				dict[acronym] = []Acronym{}
			}
			dict[acronym] = append(dict[acronym], newAcronym)
			tree.Refresh()
			fmt.Printf("Dictionary after adding: %+v\n", dict)
			err := saveDictionary(dict, "acronyms.json")
			if err != nil {
				dialog.ShowError(err, win)
			}
		}
	}, win)
}


func simulateCopy() {
	robotgo.KeyTap("c", "Control")	
}

func setupGlobalHotkeys(win fyne.Window, tree *widget.Tree, dict Dictionary) {
	var lastPressed time.Time
	debounce := 300 * time.Millisecond
	hook.Register(hook.KeyDown, []string{"ctrl", "alt", "d"}, func(e hook.Event) {
		if time.Since(lastPressed) < debounce {
			return
		}
		lastPressed = time.Now()
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