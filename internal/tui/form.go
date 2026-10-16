package tui

import (
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/schema"
)

func (m *applicationModel) clearForm() {
	m.nameInput.SetValue("")
	m.typeInput.SetValue("")
	m.minInput.SetValue("")
	m.maxInput.SetValue("")
	m.editIndex = -1
}

func (m *applicationModel) editField() tea.Cmd {
	if !m.schemaFocus || len(m.fields) == 0 {
		return nil
	}
	field := m.fields[m.selectedField]
	m.editIndex = m.selectedField
	m.nameInput.SetValue(field.Name)
	m.typeInput.SetValue(field.Type)
	if field.Type != "boolean" {
		m.minInput.SetValue(strconv.FormatFloat(field.Min, 'f', -1, 64))
		m.maxInput.SetValue(strconv.FormatFloat(field.Max, 'f', -1, 64))
	} else {
		m.minInput.SetValue("")
		m.maxInput.SetValue("")
	}
	m.setStatus("Редактируешь поле. Измени значения и нажми + для сохранения.", false)
	return m.focus(nameStep)
}

func (m *applicationModel) deleteField() {
	if !m.schemaFocus || len(m.fields) == 0 {
		m.setStatus("Сначала выбери поле в схеме клавишей Tab.", true)
		return
	}
	removed := m.fields[m.selectedField].Name
	m.fields = append(m.fields[:m.selectedField], m.fields[m.selectedField+1:]...)
	m.records = nil
	m.outputPath = ""
	if len(m.fields) == 0 {
		m.focus(nameStep)
	} else if m.selectedField >= len(m.fields) {
		m.selectedField = len(m.fields) - 1
	}
	m.setStatus(fmt.Sprintf("Поле %q удалено.", removed), false)
}

func (m *applicationModel) addField() tea.Cmd {
	name := strings.TrimSpace(m.nameInput.Value())
	fieldType := strings.ToLower(strings.TrimSpace(m.typeInput.Value()))
	if name == "" {
		m.setStatus("Укажи имя поля.", true)
		return m.focus(nameStep)
	}

	switch fieldType {
	case "integer", "float", "string", "boolean":
	default:
		m.setStatus("Тип должен быть integer, float, string или boolean.", true)
		return m.focus(typeStep)
	}

	for i, field := range m.fields {
		if i != m.editIndex && field.Name == name {
			m.setStatus("Поле с таким именем уже добавлено.", true)
			return m.focus(nameStep)
		}
	}

	newField := schema.FieldDefinition{Name: name, Type: fieldType}
	if fieldType == "integer" || fieldType == "float" || fieldType == "string" {
		min, err := strconv.ParseFloat(strings.TrimSpace(m.minInput.Value()), 64)
		if err != nil || math.IsNaN(min) || math.IsInf(min, 0) {
			m.setStatus("Введи минимум: число, а для string — длину строки.", true)
			return m.focus(minStep)
		}
		max, err := strconv.ParseFloat(strings.TrimSpace(m.maxInput.Value()), 64)
		if err != nil || math.IsNaN(max) || math.IsInf(max, 0) {
			m.setStatus("Введи максимум: число, а для string — длину строки.", true)
			return m.focus(maxStep)
		}
		if max < min {
			m.setStatus("Максимум должен быть не меньше минимума.", true)
			return m.focus(maxStep)
		}
		if fieldType == "integer" {
			if math.Trunc(min) != min || math.Trunc(max) != max {
				m.setStatus("Для integer обе границы должны быть целыми.", true)
				return m.focus(minStep)
			}
			// Пока ограничиваем диапазон, чтобы IntN получил безопасный положительный аргумент.
			if min < -1_000_000_000 || max > 1_000_000_000 {
				m.setStatus("Для integer используй границы от -1 млрд до 1 млрд.", true)
				return m.focus(minStep)
			}
		}
		if fieldType == "string" {
			if math.Trunc(min) != min || math.Trunc(max) != max || min < 0 || max > maxStringLength {
				m.setStatus("Длина string должна быть целым числом от 0 до 1000.", true)
				return m.focus(minStep)
			}
		}
		newField.Min, newField.Max = min, max
	}

	wasEditing := m.editIndex >= 0
	if wasEditing {
		m.fields[m.editIndex] = newField
		m.selectedField = m.editIndex
	} else {
		m.fields = append(m.fields, newField)
		m.selectedField = len(m.fields) - 1
	}
	m.records = nil // Старый результат построен по прежней схеме.
	m.outputPath = ""
	m.clearForm()
	if wasEditing {
		m.setStatus(fmt.Sprintf("Поле %q обновлено.", name), false)
	} else {
		m.setStatus(fmt.Sprintf("Поле %q добавлено. В схеме полей: %d.", name, len(m.fields)), false)
	}
	return m.focus(nameStep)
}

func (m *applicationModel) generate() tea.Cmd {
	if len(m.fields) == 0 {
		m.setStatus("Сначала добавь хотя бы одно поле клавишей +.", true)
		return nil
	}

	count, err := strconv.Atoi(strings.TrimSpace(m.countInput.Value()))
	if err != nil || count < 1 || count > maxRecordCount {
		m.setStatus("Количество должно быть целым числом от 1 до 10000.", true)
		m.focus(countStep)
		return nil
	}

	m.records = m.services.Generate(m.fields, count)
	m.outputPath = ""
	path, err := m.services.Export(m.records, m.fileNameInput.Value())
	if err != nil {
		m.setStatus("Не удалось сохранить JSON: "+err.Error(), true)
		return nil
	}
	m.outputPath, err = filepath.Abs(path)
	if err != nil {
		m.outputPath = path
		m.setStatus("Файл сохранён, но не удалось определить полный путь: "+err.Error(), true)
		return nil
	}
	m.setStatus(fmt.Sprintf("Готово. Полей в схеме: %d. Записей: %d.", len(m.fields), count), false)
	return openOutputFolder(m.outputPath)
}
