package exporter

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestExportToFileStreamKeepsPreviousResult(t *testing.T) {
	outputDir := t.TempDir()
	produce := func(onRecord func(map[string]any) error) error {
		for _, age := range []int{25, 30} {
			if err := onRecord(map[string]any{"age": age}); err != nil {
				return err
			}
		}
		return nil
	}
	first, err := exportToDirectoryStream("people", outputDir, produce)
	if err != nil {
		t.Fatal(err)
	}
	second, err := exportToDirectoryStream("people", outputDir, produce)
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
	if len(decoded) != 2 || decoded[0]["age"] != 25 || decoded[1]["age"] != 30 {
		t.Fatalf("unexpected JSON: %s", data)
	}
}

func TestExportToFileStreamRemovesIncompleteFile(t *testing.T) {
	outputDir := t.TempDir()
	wantErr := errors.New("generation failed")
	path, err := exportToDirectoryStream("people", outputDir, func(onRecord func(map[string]any) error) error {
		if err := onRecord(map[string]any{"age": 25}); err != nil {
			return err
		}
		return wantErr
	})
	if !errors.Is(err, wantErr) || path != "" {
		t.Fatalf("path=%q, err=%v", path, err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "people.json")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("incomplete file remains: %v", err)
	}
}
