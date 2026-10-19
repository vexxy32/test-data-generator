package templates

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/vexxy32/test-data-generator/internal/schema"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Chdir(t.TempDir())
	want := Template{
		Fields: []schema.FieldDefinition{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "integer", Min: 18, Max: 65},
		},
		Count:    25,
		FileName: "users",
	}
	path, err := Save("users", want)
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(path) || filepath.Base(path) != "users.json" {
		t.Fatalf("unexpected saved path: %q", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"name": "age"`) || strings.Contains(string(data), `"seed"`) {
		t.Fatalf("unexpected JSON contents: %s", data)
	}
	got, err := Load("users")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("loaded template = %+v; want %+v", got, want)
	}
}

func TestLoadRejectsInvalidTemplate(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.Mkdir("templates", 0755); err != nil {
		t.Fatal(err)
	}
	tests := []struct{ name, contents string }{
		{"trailing", `{"fields":[{"Name":"age","Type":"integer","Min":0,"Max":10}],"count":1,"fileName":"data"} {}`},
		{"bad-range", `{"fields":[{"Name":"age","Type":"integer","Min":10,"Max":0}],"count":1,"fileName":"data"}`},
		{"unknown-type", `{"fields":[{"Name":"age","Type":"mystery"}],"count":1,"fileName":"data"}`},
	}
	for _, tc := range tests {
		if err := os.WriteFile(filepath.Join("templates", tc.name+".json"), []byte(tc.contents), 0644); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(tc.name); err == nil {
			t.Errorf("Load(%q) accepted invalid template", tc.name)
		}
	}
	if _, err := Load("../outside"); err == nil || !strings.Contains(err.Error(), "имя шаблона") {
		t.Fatalf("path traversal was accepted: %v", err)
	}
}

func TestListTemplates(t *testing.T) {
	t.Chdir(t.TempDir())
	if names, err := List(); err != nil || len(names) != 0 {
		t.Fatalf("missing folder: names=%v err=%v", names, err)
	}
	if err := os.Mkdir("templates", 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"z.json", "a.JSON", "notes.txt"} {
		if err := os.WriteFile(filepath.Join("templates", name), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	names, err := List()
	if err != nil || !reflect.DeepEqual(names, []string{"a.JSON", "z.json"}) {
		t.Fatalf("List() = %v, %v", names, err)
	}
}
