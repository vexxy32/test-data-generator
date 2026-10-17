package tui

import (
	"path/filepath"
	"testing"

	"github.com/vexxy32/test-data-generator/internal/schema"
)

func TestFormConnectsGenerationAndExport(t *testing.T) {
	wantPath := filepath.Join(t.TempDir(), "output", "sample.json")
	generated := false
	exported := false
	m := NewModel(Services{
		Generate: func(fields []schema.FieldDefinition, count int, seed uint64) []map[string]any {
			generated = true
			if count != 3 || seed != 500 || len(fields) != 1 || fields[0].Name != "age" || fields[0].Type != "integer" {
				t.Fatalf("unexpected generation request: fields=%+v count=%d seed=%d", fields, count, seed)
			}
			return []map[string]any{{"age": 21}}
		},
		Export: func(records []map[string]any, name string) (string, error) {
			exported = true
			if name != "sample" || len(records) != 1 {
				t.Fatalf("unexpected export request: name=%q records=%v", name, records)
			}
			return wantPath, nil
		},
	}).(applicationModel)

	m.nameInput.SetValue("age")
	m.typeInput.SetValue("integer")
	m.minInput.SetValue("18")
	m.maxInput.SetValue("30")
	m.addField()
	if m.statusIsError || len(m.fields) != 1 {
		t.Fatalf("field was not added: status=%q fields=%v", m.status, m.fields)
	}

	m.countInput.SetValue("3")
	m.seedInput.SetValue("500")
	m.fileNameInput.SetValue("sample")
	m.preview()
	if !generated || exported || m.outputPath != "" {
		t.Fatalf("preview failed: generated=%v exported=%v path=%q", generated, exported, m.outputPath)
	}
	m.savePreview()
	if !generated || !exported || m.outputPath != wantPath || m.statusIsError {
		t.Fatalf("generation/export failed: generated=%v exported=%v path=%q status=%q", generated, exported, m.outputPath, m.status)
	}
}

func TestReservedTestIDCannotBeAdded(t *testing.T) {
	m := NewModel(Services{}).(applicationModel)
	m.nameInput.SetValue("_TEST_ID")
	m.typeInput.SetValue("integer")
	m.minInput.SetValue("1")
	m.maxInput.SetValue("10")
	m.addField()
	if !m.statusIsError || len(m.fields) != 0 {
		t.Fatalf("reserved field was accepted: status=%q fields=%v", m.status, m.fields)
	}
}
