package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func createAcronymTree(dict Dictionary) *widget.Tree {
	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			if id == "" {
				return getSortedAcronyms(dict)
			}
			if entry, ok := dict[id]; ok {
				children := make([]widget.TreeNodeID, len(entry))
				for i := range entry {
					children[i] = fmt.Sprintf("%s:%d", id, i)
				}
				return children
			}
			return nil
		},
		func(id widget.TreeNodeID) bool {
			return !strings.Contains(id, ":")
		},
		func(branch bool) fyne.CanvasObject {
			expanded := widget.NewLabel("")
			expanded.Wrapping = fyne.TextWrapWord

			return container.NewVBox(expanded)
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {			
			container := o.(*fyne.Container)			
			expanded := container.Objects[0].(*widget.Label)

			if branch {
				expanded.SetText(id)
			} else {				
				parts := strings.SplitN(id, ":", 2)
				
				if len(parts) == 2 {
					acronym := parts[0]
					index, _ := strconv.Atoi(parts[1])
					entry := dict[acronym][index]
					expanded.SetText(entry.Expanded)
				}
			}
		},
	)

	return tree
}


func addAcronymButton(win fyne.Window, tree *widget.Tree, dict Dictionary) {
	acronymEntry := widget.NewEntry()
	acronymEntry.SetPlaceHolder("Enter the acronym")
	
	expandEntry := widget.NewEntry()
	expandEntry.SetPlaceHolder("Enter the expanded form")

	definitionEntry := widget.NewEntry()
	definitionEntry.SetPlaceHolder("Enter the definition")

	dialog.ShowForm(fmt.Sprintf("Add New Acronym %s", acronymEntry.Text), "Add", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Acronym", acronymEntry),
		widget.NewFormItem("Expanded", expandEntry),
		widget.NewFormItem("Definition", definitionEntry),
	}, func(add bool) {
		if add && expandEntry.Text != "" && definitionEntry.Text != "" {
			newAcronym := Acronym{
				Expanded:   expandEntry.Text,
				Definition: definitionEntry.Text,
			}
			if _, ok := dict[acronymEntry.Text]; !ok {
				dict[acronymEntry.Text] = []Acronym{}
			}
			dict[acronymEntry.Text] = append(dict[acronymEntry.Text], newAcronym)
			tree.Refresh()
			err := saveDictionary(dict, dictPath)
			if err != nil {
				dialog.ShowError(err, win)
			}
		}
	}, win)
}

func getSortedAcronyms(dict Dictionary) []string {
	acronyms := make([]string, 0, len(dict))
	for acronym := range dict {
		acronyms = append(acronyms, acronym)
	}
	sort.Strings(acronyms)
	return acronyms
}