package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	supa "github.com/nedpals/supabase-go"
)

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
