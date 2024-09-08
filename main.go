package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	// "github.com/atotto/clipboard"
	// "github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
)

var dict Dictionary
var lastPressedTime time.Time
var debounceTime = 300 * time.Millisecond
var tree *widget.Tree
var dictPath = "dict/acronyms.json"

func main() {
	var err error
	dict, err = loadDictionary(dictPath)
	if err != nil {
		fmt.Println("No dictionary found, creating new one:", err)
		dict = make(Dictionary)
	}	
	fmt.Println(dict)
	
	myApp := app.New()
	mainWindow := myApp.NewWindow("Acro-Ally")	
	mainWindow.SetMaster()

	tree = createAcronymTree(dict)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for an acronym")
	searchEntry.OnSubmitted = func(text string) {
		lookUpOrDefineSearch(mainWindow, tree, dict, text)
	}


	content := container.NewBorder(
		container.NewVBox(
			searchEntry,
			widget.NewButton("Add Acronym", func() {
				addAcronymButton(mainWindow, tree, dict)
			}),
		),
		container.NewVBox(
			widget.NewButton("Exit", func() {
				myApp.Quit()
			}),
        	layout.NewSpacer(),
        	container.NewHBox(
            	layout.NewSpacer(),
				widget.NewButton("Import Acronyms", func() {
					importDictionaryDialog(mainWindow, tree, &dict)
				}),
				widget.NewButton("Export Acronyms", func() {
					exportDictionaryDialog(mainWindow, &dict)
				}),				
            layout.NewSpacer(),
        ),
		),
		nil,
		nil,
		container.NewVScroll(tree),
	)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(600, 400))

	go setupGlobalHotkeys(mainWindow, dict)

	mainWindow.ShowAndRun()
}

// These are for when the user is in main window, its fine for now
func lookUpOrDefineSearch(win fyne.Window, tree *widget.Tree, dict Dictionary, acronym string) {
	fmt.Println("Looking up or defining:", acronym)
	if _, ok := dict[acronym]; !ok {
		addAcronymSearch(win, tree, dict, acronym)
		return
	} 
	var definitions string
	for _, acro := range dict[acronym] {
		definitions += fmt.Sprintf("%s: %s\n", acro.Expanded, acro.Definition)
	}
	d := dialog.NewInformation(
		fmt.Sprintf("Acronym %s found", acronym),
		definitions,
		win,
	)
	d.Show()
	win.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyReturn || ke.Name == fyne.KeyEnter {
			d.Hide()
		}
	})
}

func addAcronymSearch(win fyne.Window, tree *widget.Tree, dict Dictionary, acronym string) {	
	expandEntry := widget.NewEntry()
	expandEntry.SetPlaceHolder("Enter the expanded form")

	definitionEntry := widget.NewEntry()
	definitionEntry.SetPlaceHolder("Enter the definition")

	dialog.NewForm(
		fmt.Sprintf("Add Acronym: %s", acronym), 
		"Add", 
		"Cancel", 
		[]*widget.FormItem{		
			widget.NewFormItem("Expanded", expandEntry),
			widget.NewFormItem("Definition", definitionEntry),
		}, 
		func(add bool) {
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
			err := saveDictionary(dict, dictPath)
			if err != nil {
				dialog.ShowError(err, win)
			}
		}
	}, win)
}

func getHighlightedText() (string, error) {
    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "darwin":
        // Use AppleScript to get the selected text on macOS
        script := `
            tell application "System Events"
                keystroke "c" using {command down}
            end tell
            delay 0.1
            return the clipboard
        `
        cmd = exec.Command("osascript", "-e", script)
    case "linux":
        // Use xsel to get the primary selection on Linux
        cmd = exec.Command("xsel", "-p")
    case "windows":
        // Use PowerShell to simulate Ctrl+C and get clipboard content
        script := `
            Add-Type -AssemblyName System.Windows.Forms
            [System.Windows.Forms.SendKeys]::SendWait("^c")
            Start-Sleep -Milliseconds 100
            Get-Clipboard
        `
        cmd = exec.Command("powershell", "-command", script)
    default:
        return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
    }

    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get highlighted text: %w", err)
    }

    return string(output), nil
}

func setupGlobalHotkeys(win fyne.Window, dict Dictionary) {	
	hook.Register(hook.KeyDown, []string{"ctrl", "alt", "d"}, func(e hook.Event) {		
		if time.Since(lastPressedTime) < debounceTime {						
			return			
		}
		fmt.Println("Hotkey pressed")
		
		text, err := getHighlightedText()
		lastPressedTime = time.Now()

		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if text != "" {
			showPopup(dict, text)
		}		
	})

	s := hook.Start()
	<-hook.Process(s)
}