package main

import (
	"encoding/json"
	"os"
)

type Acronym struct {
	Expanded   string `json:"expanded"`
	Definition string `json:"definition"`
	Index      int    `json:"index"`
}

type Dictionary map[string]Acronym

func saveDictionary(dict Dictionary, filename string) error {
	data, err := json.MarshalIndent(dict, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func loadDictionary(filename string) (Dictionary, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dict Dictionary
	err = json.Unmarshal(data, &dict)
	if err != nil {
		return nil, err
	}
	return dict, nil
}