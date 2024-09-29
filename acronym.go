package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Acronym struct {
	Expanded   string `json:"expanded"`
	Definition string `json:"definition"`
}

type Dictionary map[string][]Acronym

func saveDictionary(dict Dictionary, filename string) error {
	dir := filepath.Dir(filename)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    data, err := json.MarshalIndent(dict, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal dictionary: %w", err)
    }

    if err := os.WriteFile(filename, data, 0644); err != nil {
        return fmt.Errorf("failed to write file: %w", err)
	}
	return nil	
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