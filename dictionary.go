package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type acronymLayout struct{}

func (al *acronymLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 2 {
		return
	}
	expanded := objects[0]
	definition := objects[1]

	// Set a maximum width for the expanded label (e.g., 30% of the total width)
	maxExpandedWidth := float32(size.Width) * 0.3
	expandedSize := expanded.MinSize()
	if expandedSize.Width > maxExpandedWidth {
		expandedSize.Width = maxExpandedWidth
	}

	// Position and resize the expanded label
	expanded.Resize(expandedSize)
	expanded.Move(fyne.NewPos(0, 0))

	// Calculate remaining space for the definition
	defX := expandedSize.Width + theme.Padding()
	defWidth := size.Width - defX
	defHeight := size.Height

	// Resize and position the definition label
	definition.Resize(fyne.NewSize(defWidth, defHeight))
	definition.Move(fyne.NewPos(defX, 0))
}

func (al *acronymLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) != 2 {
		return fyne.NewSize(0, 0)
	}
	expanded := objects[0]
	definition := objects[1]

	// Calculate minimum width based on both labels
	minWidth := expanded.MinSize().Width + theme.Padding() + definition.MinSize().Width
	
	// Use the height of the taller label
	minHeight := fyne.Max(expanded.MinSize().Height, definition.MinSize().Height)
	
	return fyne.NewSize(minWidth, minHeight)
}

func createAcronymTree(dict Dictionary) *widget.Tree {
		tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			if id == "" {
				return getSortedAcronyms(dict)
			}
			if entry, ok := dict[id]; ok {
				children := make([]widget.TreeNodeID, len(entry))
				for i, _ := range entry {
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
			expanded.TextStyle = fyne.TextStyle{Bold: true}

			definition := widget.NewLabel("")
			definition.Wrapping = fyne.TextWrapWord

			return container.New(&acronymLayout{}, expanded, definition)
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			expanded := container.Objects[0].(*widget.Label)
			definition := container.Objects[1].(*widget.Label)

			if branch {
				expanded.SetText(id)
				definition.SetText("")
			} else {
				parts := strings.SplitN(id, ":", 2)
				if len(parts) == 2 {
					acronym := parts[0]
					index, _ := strconv.Atoi(parts[1])
					entry := dict[acronym][index]
					expanded.SetText(entry.Expanded)
					definition.SetText(entry.Definition)
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
			fmt.Printf("Dictionary after adding: %+v\n", dict)
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