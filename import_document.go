package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/dslipak/pdf"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

)

func readPdf(path string) (string, error) {
    r, err := pdf.Open(path)    
    if err != nil {
        return "", err
    }		

    var texts string
    totalPage := r.NumPage()
    for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
        p := r.Page(pageIndex)
        if p.V.IsNull() {
            continue
        }

		pageText := p.Content().Text
		for _, text := range pageText {
			texts += text.S
		}
    }

    return texts, nil
}

func importAcronyms(win fyne.Window, tree *widget.Tree, dict Dictionary) {
	dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()
		
		tmpFile, err := os.CreateTemp("", "tmp-*.pdf")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer os.Remove(tmpFile.Name())
		
		_, err = io.Copy(tmpFile, reader)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		content, err := readPdf(tmpFile.Name())
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		loadingDialog := showLoadingDialog(win)

		go func() {
			acronyms, err := extractAcronymsFromDocument(content)
			loadingDialog.Hide()

			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			
			addedAcronyms := make([]string, 0, len(acronyms))
			for _, acronym := range acronyms {
				// Logging each acronym as it's processed
				fmt.Printf("Processing acronym: %s\nExpanded: %s\nDefinition: %s\n\n", 
					acronym.Acronym, acronym.Expanded, acronym.Definition)

				if _, ok := dict[acronym.Acronym]; !ok {
					dict[acronym.Acronym] = []Acronym{}
				}
				dict[acronym.Acronym] = append(dict[acronym.Acronym], Acronym{
					Expanded: acronym.Expanded,
					Definition: acronym.Definition,
				})
				addedAcronyms = append(addedAcronyms, fmt.Sprintf("**%s** - %s", acronym.Acronym, acronym.Expanded))
			}
			err = saveDictionary(dict, dictPath)
			if err != nil {
				dialog.ShowError(err, win)
			}

			tree.Refresh()

			if len(addedAcronyms) > 0 {
				content := widget.NewRichTextFromMarkdown(strings.Join(addedAcronyms, "\n\n"))
				scroll := container.NewScroll(content)
				scroll.SetMinSize(fyne.NewSize(300, 200))
				dialog.ShowCustom("Added Acronyms", "OK", scroll, win)
				fmt.Printf("Added %d new acronyms\n", len(addedAcronyms))
			} else {
				dialog.ShowInformation("No New Acronyms", "No new acronyms were found in the document.", win)
				fmt.Println("No new acronyms found")
			}
		}()
		
	}, win)

	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	dialog.Show()
}

func showLoadingDialog(win fyne.Window) dialog.Dialog {
	progress := widget.NewProgressBarInfinite()
	loadingDialog := dialog.NewCustom("Processing", "Cancel", progress, win)
	loadingDialog.Show()
	return loadingDialog
}


type AcronymResult struct {
	Acronym string `json:"acronym"`
	Expanded string `json:"expanded"`
	Definition string `json:"definition"`
}

type AcronymResponse struct {
	Acronyms []AcronymResult `json:"acronyms"`
}

func extractAcronymsFromDocument(content string) ([]AcronymResult, error) {
	prompt := `
	You are an acronym extractor. Extract acronyms, their expanded form, and a brief definition from the given text. If there isn't enought context to create a sufficient definition, use the context to create a definition as if you were an expert in the subject matter.
	
	For example, if the context is a document about telecommunications, set 
	acronym: '2G'
	expanded: 'second-generation cellular network'
	definition: 'The second generation of wireless technology that transitioned from analog to digital signals, enhancing voice quality, enabling SMS services, and allowing for more efficient use of the radio frequency spectrum.'

	Return the result as a JSON object with an 'acronyms' array containing objects with 'acronym', 'expanded', and 'definition' fields.

	The 'expanded' field should contain the full phrase that the acronym stands for.
	The 'definition' field should contain a brief explanation or description of the concept.

	Example:
	[
		{
			"acronym": "API", 
			"expanded": "Application Programming Interface", 
			"definition": "A set of routines, protocols, and tools for building software applications. It specifies how software components should interact."
		},
		{
			"acronym": "HTML",
			"expanded": "Hypertext Markup Language",
			"definition": "A standard markup language used for creating web pages and web applications."
		},
		{
			"acronym": "CPU",
			"expanded": "Central Processing Unit",
			"definition": "The primary component of a computer that performs most of the processing inside a computer."
		}
	]

	Here is the text document to extract acronyms from:

	`
	
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not set")
	}
	client := openai.NewClient(apiKey)
	ctx := context.Background()
	
	var response AcronymResponse
	schema, err := jsonschema.GenerateSchemaForType(response)
	if err != nil {
		return nil, err
	}
	fmt.Println("Sending request to OpenAI")
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: prompt,
			},
			{
				Role: openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "Acronyms_Extraction",
				Schema: schema,
				Strict: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	err = schema.Unmarshal(resp.Choices[0].Message.Content, &response)
	if err != nil {
		return nil, err
	}
	
	return response.Acronyms, nil
}