package main

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func createAcronymTree(dict Dictionary) *widget.Tree {
	tree := widget.NewTree(
		func(acro widget.TreeNodeID) []widget.TreeNodeID {
			if acro == "" {
				return getSortedAcronyms(dict)
			}
			if entry, ok := dict[acro]; ok {
				children := make([]widget.TreeNodeID, 0, len(entry))
				for _, definition := range entry {
					children = append(children, fmt.Sprintf("%s: %s", definition.Expanded, definition.Definition))
				}
				return children
			}
			return []widget.TreeNodeID{}
		},
		func(acro widget.TreeNodeID) bool {
			return acro == "" || len(dict[acro]) > 0
		},
		func(branch bool) fyne.CanvasObject {
			return container.NewStack(container.NewVScroll(widget.NewRichText()))
		},
		func(acro widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			scrollContainer := obj.(*fyne.Container)
			richText := scrollContainer.Objects[0].(*container.Scroll).Content.(*widget.RichText)
			richText.Wrapping = fyne.TextWrapWord

			if acro == "" {
				richText.ParseMarkdown("**Acronyms**")
			} else if branch {
				richText.ParseMarkdown(fmt.Sprintf("**%s**", acro))
			} else {
				parts := splitAcronymDefinition(acro)
				if len(parts) == 2 {
					richText.ParseMarkdown(fmt.Sprintf("**%s**: %s", parts[0], parts[1]))
				} else {
					richText.ParseMarkdown(acro)
				}
			}
		},
	)
	
	return tree
}

// Helper function to split the acronym definition string
func splitAcronymDefinition(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
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
			fmt.Printf("Dictionary after adding: %+v\n", dict)
			err := saveDictionary(dict, "acronyms.json")
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