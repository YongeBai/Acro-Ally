package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)



func importDictionary(currentDict *Dictionary, filename string, merge bool) (Dictionary, error) {
	importDict, err := loadDictionary(filename)
	if err != nil {
		return nil, err
	}
	if merge {
		for key, value := range importDict {
			if _, ok := (*currentDict)[key]; !ok {
				(*currentDict)[key] = []Acronym{}
			}
			(*currentDict)[key] = append((*currentDict)[key], value...)
		}
		return *currentDict, nil
	}
	return importDict, nil
		
}

func exportDictionary(dict Dictionary, filename string) error {
	return saveDictionary(dict, filename)
}

func exportDictionaryDialog(win fyne.Window, dict *Dictionary) {
    dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
        if err != nil {
            dialog.ShowError(err, win)
            return
        }
        if writer == nil {
            return
        }
        defer writer.Close()

        err = exportDictionary(*dict, writer.URI().Path())
        if err != nil {
            dialog.ShowError(err, win)
        } else {
            dialog.ShowInformation("Success", "Dictionary exported successfully", win)
        }
    }, win)
}

func importDictionaryDialog(win fyne.Window, tree *widget.Tree, dict *Dictionary) {
    dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
        if err != nil {
            dialog.ShowError(err, win)
            return
        }
        if reader == nil {
            return
        }
        defer reader.Close()

        dialog.ShowConfirm("Import Options", "Do you want to merge with the current dictionary or replace it?", func(merge bool) {
            newDict, err := importDictionary(dict, reader.URI().Path(), merge)
            if err != nil {
                dialog.ShowError(err, win)
            } else {
                *dict = newDict
                tree.Refresh()
                dialog.ShowInformation("Success", "Dictionary imported successfully", win)
            }
        }, win)
    }, win)
}