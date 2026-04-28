package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const (
	indent string = "  "
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := make([]string, len(pages))
	pageData := make([]PageData, len(pages))

	i := 0
	for k, _ := range pages {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	for i := range keys {
		pageData[i] = pages[keys[i]]
	}

	jsonData, err := json.MarshalIndent(pageData, "", indent)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	os.WriteFile(filename, jsonData, 0444)

	return nil
}
