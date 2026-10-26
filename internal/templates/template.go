package templates

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/vexxy32/test-data-generator/internal/schema"
)

const templatesDir = "templates"
const maxTemplateRecords = 10_000
const maxTemplateStringLength = 1_000

type FileStore struct{}

func (FileStore) Save(name string, template Template) (string, error) {
	return Save(name, template)
}

func (FileStore) Load(name string) (Template, error) {
	return Load(name)
}

func (FileStore) List() ([]string, error) {
	return List()
}

// Template хранит настройки генерации без сгенерированных данных.
type Template struct {
	Fields   []schema.FieldDefinition `json:"fields"`
	Count    int                      `json:"count"`
	FileName string                   `json:"fileName"`
}

// Validate проверяет поля шаблона по требованиям ТЗ.
func (t Template) Validate() error {
	if len(t.Fields) == 0 {
		return errors.New("шаблон не содержит полей")
	}
	if t.Count <= 0 || t.Count > maxTemplateRecords {
		return errors.New("количество записей должно быть от 1 до 10000")
	}
	fileName := strings.TrimSpace(t.FileName)
	if fileName == "" {
		return errors.New("имя выходного файла не должно быть пустым")
	}
	if strings.ContainsAny(fileName, `\/:*?"<>|`) || (filepath.Ext(fileName) != "" && !strings.EqualFold(filepath.Ext(fileName), ".json")) {
		return errors.New("некорректное имя выходного файла")
	}
	names := make(map[string]bool, len(t.Fields))
	for _, field := range t.Fields {
		if strings.TrimSpace(field.Name) == "" || strings.EqualFold(field.Name, "_test_id") || names[field.Name] {
			return fmt.Errorf("некорректное или повторяющееся имя поля %q", field.Name)
		}
		names[field.Name] = true
		if _, automatic := schema.AutomaticStringKind(field.Name); automatic && field.Type != "string" {
			return fmt.Errorf("поле %q должно иметь тип string", field.Name)
		}
		if math.IsNaN(field.Min) || math.IsNaN(field.Max) || math.IsInf(field.Min, 0) || math.IsInf(field.Max, 0) || field.Min > field.Max {
			return fmt.Errorf("некорректные границы поля %q", field.Name)
		}
		switch field.Type {
		case "integer":
			if math.Trunc(field.Min) != field.Min || math.Trunc(field.Max) != field.Max || field.Min < -1_000_000_000 || field.Max > 1_000_000_000 {
				return fmt.Errorf("некорректный диапазон integer для поля %q", field.Name)
			}
		case "float":
			if field.Precision < 0 || field.Precision > 6 {
				return fmt.Errorf("некорректная точность float для поля %q", field.Name)
			}
		case "string":
			if _, automatic := schema.AutomaticStringKind(field.Name); !automatic && (math.Trunc(field.Min) != field.Min || math.Trunc(field.Max) != field.Max || field.Min < 0 || field.Max > maxTemplateStringLength) {
				return fmt.Errorf("некорректная длина string для поля %q", field.Name)
			}
		case "boolean":
		default:
			return fmt.Errorf("неизвестный тип %q у поля %q", field.Type, field.Name)
		}
	}
	return nil
}

// normalizeTemplateName проверяет имя шаблона и добавляет .json при необходимости.
func normalizeTemplateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("имя шаблона не должно быть пустым")
	}

	if strings.ContainsAny(name, `\/:*?"<>|`) || name == "." || name == ".." {
		return "", errors.New("некорректное имя шаблона")
	}

	if !strings.HasSuffix(strings.ToLower(name), ".json") {
		name += ".json"
	}

	return name, nil
}
