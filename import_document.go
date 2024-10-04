package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/ledongthuc/pdf"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
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
		
		tmpFile, err := os.CreateTemp("", "tmp.pdf")
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
		
		// acronyms, err := extractAcronymsFromDocument(content)
		// if err != nil {
		// 	dialog.ShowError(err, win)
		// 	return
		// }
		
		tree.Refresh()
	}, win)

	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	dialog.Show()
}



type AcronymResult struct {
	Acronym string `json:"acronym"`
	Expanded string `json:"expanded"`
	Definition string `json:"definition"`
}

func extractAcronymsFromDocument(content string) ([]AcronymResult, error) {
	prompt := "You are an acronym extractor. Extract acronyms, their expanded forms, and definitions from the given text. " +
		"Return the result as a JSON object with an 'acronyms' array containing objects with 'acronym', 'expanded', and 'definition' fields."

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	ctx := context.Background()
	
	var results []AcronymResult
	schema, err := jsonschema.GenerateSchemaForType(results)
	if err != nil {
		return nil, err
	}

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
	err = schema.Unmarshal(resp.Choices[0].Message.Content, &results)
	if err != nil {
		return nil, err
	}
	fmt.Println(results)
	return results, nil
}