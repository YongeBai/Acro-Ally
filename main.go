package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
	hook "github.com/robotn/gohook"
)

var dict Dictionary
var lastPressedTime time.Time
var debounceTime = 300 * time.Millisecond
var tree *widget.Tree
var dictPath = "dict/acronyms.json"

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	myApp := app.New()
	mainWindow := myApp.NewWindow("Acro-Ally")
	mainWindow.SetMaster()	
	if _, err := os.Stat("license_key.txt"); err == nil {
		// License key file exists, read it
		licenseKey, err := os.ReadFile("license_key.txt")
		if err == nil && validateLicenseKey(string(licenseKey)) {
			// License key is valid, proceed with app
		} else {
			// License key is invalid or not found, prompt for a new one
			checkLicenseKey(mainWindow)
		}
	} else {
		// No license key file, prompt for a new one
		checkLicenseKey(mainWindow)
	}
	var err error
	dict, err = loadDictionary(dictPath)
	if err != nil {
		fmt.Println("No dictionary found, creating new one")
		fmt.Println(err)
		dict = make(Dictionary)
	}

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

	definitions := container.NewVBox()
	for _, acro := range dict[acronym] {
		acroLabel := widget.NewLabel(acro.Expanded)
		acroLabel.TextStyle = fyne.TextStyle{Bold: true}
		acroLabel.Wrapping = fyne.TextWrapWord

		definitionLabel := widget.NewLabel(acro.Definition)
		definitionLabel.Wrapping = fyne.TextWrapWord

		definitionBox := container.NewVBox(
			acroLabel,
			definitionLabel,
			widget.NewSeparator(),
		)
		definitions.Add(definitionBox)
	}

	scrollContainer := container.NewScroll(definitions)
	scrollContainer.SetMinSize(fyne.NewSize(400, 300))

	d := dialog.NewCustom(
		fmt.Sprintf("Acronym %s found", acronym),
		"Close",
		scrollContainer,
		win,
	)
	d.Resize(fyne.NewSize(350, 250)) // Adjust dialog size as needed
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

	formDialog := dialog.NewForm(
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
				// fmt.Printf("Dictionary after adding: %+v\n", dict)
				err := saveDictionary(dict, dictPath)
				if err != nil {
					dialog.ShowError(err, win)
				}
			}
		},
		win,
	)
	formDialog.Show()
}
func setupGlobalHotkeys(win fyne.Window, dict Dictionary) {
	hook.Register(hook.KeyDown, []string{"ctrl", "alt", "d"}, func(e hook.Event) {
		if time.Since(lastPressedTime) < debounceTime {
			return
		}
		fmt.Println("Hotkey pressed")

		// read clipboard
		text, err := clipboard.ReadAll()
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


func checkLicenseKey(win fyne.Window) bool {
    licenseEntry := widget.NewEntry()
    licenseEntry.SetPlaceHolder("Enter your license key")

    var isValid bool

    formDialog := dialog.NewForm(
        "License Key Required",
        "Submit",
        "Cancel",
        []*widget.FormItem{
            widget.NewFormItem("License Key", licenseEntry),
        },
        func(submit bool) {
            if submit {
                licenseKey := licenseEntry.Text
                isValid = validateLicenseKey(licenseKey)
                if !isValid {
                    os.Exit(0)
                } else {
                    saveLicenseKey(licenseKey)
                }
            } else {
                // If canceled, exit the application
                os.Exit(0)
            }
        },
        win,
    )

    formDialog.Show()

    return isValid
}

func validateLicenseKey(licenseKey string) bool {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	
	// Check if the environment variables are set
	if supabaseURL == "" || supabaseKey == "" {
		fmt.Println("Supabase URL or Key is not set")
		return false
	}

	supabaseClient := supa.CreateClient(supabaseURL, supabaseKey)
	
	var results []map[string]interface{}

	err := supabaseClient.DB.From("licenses").
		Select("*").
		Eq("license_key", licenseKey).
		Execute(&results)
	
	if err != nil {
		fmt.Println("Error checking license key:", err)
		return false
	}

	return len(results) > 0
}

func saveLicenseKey(licenseKey string) {
	// Save the license key to a file or use a more secure method
	err := os.WriteFile("license_key.txt", []byte(licenseKey), 0644)
	if err != nil {
		fmt.Println("Error saving license key:", err)
	}
}
