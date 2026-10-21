package tui

import (
	"path/filepath"
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/schema"
	"github.com/vexxy32/test-data-generator/internal/templates"
)

func TestFormConnectsGenerationAndExport(t *testing.T) {
	wantPath := filepath.Join(t.TempDir(), "output", "sample.json")
	generated := false
	exported := false
	m := NewModel(Services{
		GenerateStream: func(fields []schema.FieldDefinition, count int, seed uint64, onRecord func(map[string]any) error) error {
			generated = true
			if (count != 3 && count != 5) || seed != 500 || len(fields) != 1 || fields[0].Name != "age" || fields[0].Type != "integer" {
				t.Fatalf("unexpected generation request: fields=%+v count=%d seed=%d", fields, count, seed)
			}
			for i := 1; i <= count; i++ {
				if err := onRecord(map[string]any{"age": 20 + i, "_test_id": i}); err != nil {
					return err
				}
			}
			return nil
		},
		ExportStream: func(name string, produceRecords func(func(map[string]any) error) error) (string, error) {
			exported = true
			var records []map[string]any
			err := produceRecords(func(record map[string]any) error {
				records = append(records, record)
				return nil
			})
			if err != nil {
				return "", err
			}
			if name != "sample" || len(records) != 5 || records[0]["age"] != 21 {
				t.Fatalf("unexpected export request: name=%q records=%v", name, records)
			}
			return wantPath, nil
		},
	})

	m.nameInput.SetValue("age")
	m.typeInput.SetValue("integer")
	m.minInput.SetValue("18")
	m.maxInput.SetValue("30")
	m.addField()
	if m.statusIsError || len(m.fields) != 1 {
		t.Fatalf("field was not added: status=%q fields=%v", m.status, m.fields)
	}

	m.countInput.SetValue("5")
	m.seedInput.SetValue("500")
	m.fileNameInput.SetValue("sample")
	m.preview()
	if !generated || exported || m.outputPath != "" || len(m.records) != 3 {
		t.Fatalf("preview failed: generated=%v exported=%v path=%q records=%v", generated, exported, m.outputPath, m.records)
	}
	m.savePreview()
	if !generated || !exported || m.outputPath != wantPath || m.statusIsError {
		t.Fatalf("generation/export failed: generated=%v exported=%v path=%q status=%q", generated, exported, m.outputPath, m.status)
	}
}

func TestReservedTestIDCannotBeAdded(t *testing.T) {
	m := NewModel(Services{})
	m.nameInput.SetValue("_TEST_ID")
	m.typeInput.SetValue("integer")
	m.minInput.SetValue("1")
	m.maxInput.SetValue("10")
	m.addField()
	if !m.statusIsError || len(m.fields) != 0 {
		t.Fatalf("reserved field was accepted: status=%q fields=%v", m.status, m.fields)
	}
}

func TestChangingSeedInvalidatesPreview(t *testing.T) {
	m := NewModel(Services{})
	m.records = []map[string]any{{"_test_id": 1}}
	m.outputPath = "output/old.json"
	m.focus(seedStep)
	updated, _ := m.Update(tea.KeyPressMsg{Code: '4', Text: "4"})
	m = updated.(Model)
	if len(m.records) != 0 || m.outputPath != "" {
		t.Fatalf("old preview still active: records=%v path=%q", m.records, m.outputPath)
	}
}

func TestEmptyBoundsUseDefaults(t *testing.T) {
	m := NewModel(Services{})
	m.nameInput.SetValue("score")
	m.typeInput.SetValue("integer")
	m.addField()
	if m.statusIsError || len(m.fields) != 1 {
		t.Fatalf("field was not added: status=%q fields=%v", m.status, m.fields)
	}
	if m.fields[0].Min != defaultNumberMin || m.fields[0].Max != defaultNumberMax {
		t.Fatalf("unexpected number defaults: %+v", m.fields[0])
	}

	m.clearForm()
	m.nameInput.SetValue("comment")
	m.typeInput.SetValue("string")
	m.addField()
	if m.statusIsError || len(m.fields) != 2 {
		t.Fatalf("string field was not added: status=%q fields=%v", m.status, m.fields)
	}
	if m.fields[1].Min != defaultStringMin || m.fields[1].Max != defaultStringMax {
		t.Fatalf("unexpected string defaults: %+v", m.fields[1])
	}
}

func TestNewTestClearsCurrentData(t *testing.T) {
	m := NewModel(Services{})
	m.fields = []schema.FieldDefinition{{Name: "age", Type: "integer"}}
	m.records = []map[string]any{{"age": 21}}
	m.outputPath = "output/result.json"
	m.seedInput.SetValue("500")
	m.fileNameInput.SetValue("old_result")

	m.newTest()
	if len(m.fields) != 0 || len(m.records) != 0 || m.outputPath != "" {
		t.Fatalf("new test kept old data: fields=%v records=%v path=%q", m.fields, m.records, m.outputPath)
	}
	if m.seedInput.Value() != "" || m.fileNameInput.Value() != "generated_tests" {
		t.Fatalf("new test kept old settings: seed=%q file=%q", m.seedInput.Value(), m.fileNameInput.Value())
	}
}

func TestTemplateSaveAndLoadUpdatesForm(t *testing.T) {
	want := templates.Template{
		Fields:   []schema.FieldDefinition{{Name: "age", Type: "integer", Min: 18, Max: 65}},
		Count:    12,
		FileName: "users",
	}
	saved := false
	m := NewModel(Services{
		SaveTemplate: func(name string, tmpl templates.Template) (string, error) {
			if name != "people" || !reflect.DeepEqual(tmpl, want) {
				t.Fatalf("unexpected save: name=%q template=%+v", name, tmpl)
			}
			saved = true
			return "templates/people.json", nil
		},
		LoadTemplate: func(name string) (templates.Template, error) {
			if name != "people" {
				t.Fatalf("unexpected load name: %q", name)
			}
			return want, nil
		},
	})
	m.fields = append([]schema.FieldDefinition(nil), want.Fields...)
	m.countInput.SetValue("12")
	m.fileNameInput.SetValue("users")
	m.templateNameInput.SetValue("people")
	m.saveTemplate()
	if !saved || m.statusIsError {
		t.Fatalf("save failed: %q", m.status)
	}
	m.fields = []schema.FieldDefinition{{Name: "old", Type: "boolean"}}
	m.records = []map[string]any{{"old": true}}
	m.outputPath = "output/old.json"
	m.loadTemplate()
	if !reflect.DeepEqual(m.fields, want.Fields) || m.countInput.Value() != "12" || m.fileNameInput.Value() != "users" || len(m.records) != 0 || m.outputPath != "" || m.statusIsError {
		t.Fatalf("load failed: fields=%+v count=%q file=%q records=%v path=%q status=%q", m.fields, m.countInput.Value(), m.fileNameInput.Value(), m.records, m.outputPath, m.status)
	}
}

func TestTemplatePickerLoadsSelectedTemplate(t *testing.T) {
	m := NewModel(Services{
		ListTemplates: func() ([]string, error) {
			return []string{"first.json", "second.json"}, nil
		},
		LoadTemplate: func(name string) (templates.Template, error) {
			if name != "second.json" {
				t.Fatalf("selected wrong template: %q", name)
			}
			return templates.Template{
				Fields: []schema.FieldDefinition{{Name: "city", Type: "string"}},
				Count:  10, FileName: "cities",
			}, nil
		},
	})
	updated, _ := m.Update(tea.KeyPressMsg{Code: 'l', Mod: tea.ModCtrl})
	m = updated.(Model)
	if !m.templatePicker || m.selectedTemplate != 0 {
		t.Fatalf("picker did not open: %+v", m)
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(Model)
	if m.selectedTemplate != 1 {
		t.Fatalf("arrow did not change selection: %d", m.selectedTemplate)
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)
	if m.templatePicker || m.templateNameInput.Value() != "second.json" || m.fields[0].Name != "city" || m.countInput.Value() != "10" {
		t.Fatalf("selected template was not loaded: %+v", m)
	}
}
