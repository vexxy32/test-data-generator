package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExportToFileKeepsPreviousResult(t *testing.T) {
	outputDir := t.TempDir()
	records := []map[string]any{{"age": 25}}
	first, err := exportToDirectory(records, "people", outputDir)
	if err != nil {
		t.Fatal(err)
	}
	second, err := exportToDirectory(records, "people", outputDir)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(first) != "people.json" || filepath.Base(second) != "people_2.json" {
		t.Fatalf("unexpected names: %q, %q", first, second)
	}

	data, err := os.ReadFile(first)
	if err != nil {
		t.Fatal(err)
	}
	var decoded []map[string]int
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 1 || decoded[0]["age"] != 25 {
		t.Fatalf("unexpected JSON: %s", data)
	}
}
