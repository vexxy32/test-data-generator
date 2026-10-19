package tui

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/schema"
	"github.com/vexxy32/test-data-generator/internal/templates"
)

const (
	defaultNumberMin = 0
	defaultNumberMax = 100
	defaultStringMin = 1
	defaultStringMax = 20
)

// clearForm очищает форму после добавления поля или отмены правки.
func (m *Model) clearForm() {
	m.nameInput.SetValue("")
	m.typeInput.SetValue("")
	m.minInput.SetValue("")
	m.maxInput.SetValue("")
	m.precisionInput.SetValue("2")
	m.editIndex = -1
}

// newTest очищает схему и результат, чтобы можно было создать другой набор данных.
func (m *Model) newTest() tea.Cmd {
	m.fields = nil
	m.records = nil
	m.outputPath = ""
	m.selectedField = 0
	m.clearForm()
	m.countInput.SetValue("5")
	m.seedInput.SetValue("")
	m.fileNameInput.SetValue("generated_tests")
	m.templateNameInput.SetValue("")
	m.setStatus("Новый тест: добавь поля для новой схемы.", false)
	return m.focus(nameStep)
}

// saveTemplate сохраняет текущую схему и настройки, без готовых записей и seed.
func (m *Model) saveTemplate() {
	if m.services.SaveTemplate == nil {
		m.setStatus("Сохранение шаблонов недоступно.", true)
		return
	}
	if strings.TrimSpace(m.templateNameInput.Value()) == "" {
		m.setStatus("Укажи имя шаблона в строке «Шаблон».", true)
		m.focus(templateNameStep)
		return
	}
	count, err := strconv.Atoi(strings.TrimSpace(m.countInput.Value()))
	if err != nil || count < 1 || count > maxRecordCount {
		m.setStatus("Количество должно быть целым числом от 1 до 10000.", true)
		m.focus(countStep)
		return
	}
	tmpl := templates.Template{
		Fields:   append([]schema.FieldDefinition(nil), m.fields...),
		Count:    count,
		FileName: strings.TrimSpace(m.fileNameInput.Value()),
	}
	path, err := m.services.SaveTemplate(m.templateNameInput.Value(), tmpl)
	if err != nil {
		m.setStatus("Не удалось сохранить шаблон: "+err.Error(), true)
		return
	}
	m.setStatus("Шаблон сохранён: "+path, false)
}

// loadTemplate подставляет сохранённую схему в форму.
func (m *Model) loadTemplate() tea.Cmd {
	if m.services.LoadTemplate == nil {
		m.setStatus("Загрузка шаблонов недоступна.", true)
		return nil
	}
	if strings.TrimSpace(m.templateNameInput.Value()) == "" {
		m.setStatus("Укажи имя шаблона в строке «Шаблон».", true)
		return m.focus(templateNameStep)
	}
	tmpl, err := m.services.LoadTemplate(m.templateNameInput.Value())
	if err != nil {
		m.setStatus("Не удалось загрузить шаблон: "+err.Error(), true)
		return nil
	}
	m.fields = append([]schema.FieldDefinition(nil), tmpl.Fields...)
	m.countInput.SetValue(strconv.Itoa(tmpl.Count))
	m.fileNameInput.SetValue(tmpl.FileName)
	m.records = nil
	m.outputPath = ""
	m.selectedField = 0
	m.clearForm()
	m.setStatus(fmt.Sprintf("Шаблон %q загружен: %d полей, %d записей.", strings.TrimSpace(m.templateNameInput.Value()), len(m.fields), tmpl.Count), false)
	return m.focus(nameStep)
}

// openTemplatePicker показывает сохранённые шаблоны для выбора стрелками.
func (m *Model) openTemplatePicker() {
	if m.services.ListTemplates == nil {
		m.setStatus("Список шаблонов недоступен.", true)
		return
	}
	names, err := m.services.ListTemplates()
	if err != nil {
		m.setStatus("Не удалось получить список шаблонов: "+err.Error(), true)
		return
	}
	if len(names) == 0 {
		m.setStatus("Шаблонов пока нет. Сначала сохрани схему через Ctrl+S.", true)
		return
	}
	m.templateNames = names
	m.selectedTemplate = 0
	m.templatePicker = true
	m.setStatus("Выбери шаблон стрелками и нажми Enter.", false)
}

// editField переносит выбранное поле из схемы обратно в форму.
func (m *Model) editField() tea.Cmd {
	if !m.schemaFocus || len(m.fields) == 0 {
		return nil
	}
	field := m.fields[m.selectedField]
	m.editIndex = m.selectedField
	m.nameInput.SetValue(field.Name)
	m.typeInput.SetValue(field.Type)
	if field.Type != "boolean" && !isAutomaticStringFieldName(field.Name) {
		m.minInput.SetValue(strconv.FormatFloat(field.Min, 'f', -1, 64))
		m.maxInput.SetValue(strconv.FormatFloat(field.Max, 'f', -1, 64))
	} else {
		m.minInput.SetValue("")
		m.maxInput.SetValue("")
	}
	if field.Type == "float" {
		m.precisionInput.SetValue(strconv.Itoa(field.Precision))
	} else {
		m.precisionInput.SetValue("2")
	}
	m.setStatus("Редактируешь поле. Измени значения и нажми + для сохранения.", false)
	return m.focus(nameStep)
}

// deleteField удаляет выбранное поле и сбрасывает старый результат генерации.
func (m *Model) deleteField() {
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

// addField проверяет введённые параметры и добавляет поле или сохраняет его правку.
func (m *Model) addField() tea.Cmd {
	name := strings.TrimSpace(m.nameInput.Value())
	fieldType := strings.ToLower(strings.TrimSpace(m.typeInput.Value()))
	if name == "" {
		m.setStatus("Укажи имя поля.", true)
		return m.focus(nameStep)
	}
	if strings.EqualFold(name, "_test_id") {
		m.setStatus("Имя _test_id зарезервировано программой.", true)
		return m.focus(nameStep)
	}
	if automaticType, automatic := schema.AutomaticStringType(name); automatic {
		fieldType = automaticType
		m.typeInput.SetValue(fieldType)
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
	if (fieldType == "integer" || fieldType == "float" || fieldType == "string") && !isAutomaticStringFieldName(name) {
		defaultMin, defaultMax := defaultBounds(fieldType)
		min, err := optionalFloat(m.minInput.Value(), defaultMin)
		if err != nil || math.IsNaN(min) || math.IsInf(min, 0) {
			m.setStatus("Минимум должен быть числом, а для string — длиной строки.", true)
			return m.focus(minStep)
		}
		max, err := optionalFloat(m.maxInput.Value(), defaultMax)
		if err != nil || math.IsNaN(max) || math.IsInf(max, 0) {
			m.setStatus("Максимум должен быть числом, а для string — длиной строки.", true)
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
	if fieldType == "float" {
		precision, err := strconv.Atoi(strings.TrimSpace(m.precisionInput.Value()))
		if err != nil || precision < 0 || precision > maxFloatPrecision {
			m.setStatus("Точность float должна быть целым числом от 0 до 6.", true)
			return m.focus(precisionStep)
		}
		newField.Precision = precision
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

// defaultBounds задаёт значение для каждой пустой границы.
func defaultBounds(fieldType string) (float64, float64) {
	if fieldType == "string" {
		return defaultStringMin, defaultStringMax
	}
	return defaultNumberMin, defaultNumberMax
}

func optionalFloat(value string, defaultValue float64) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseFloat(value, 64)
}

// preview создаёт записи один раз; эти же записи затем сохраняет savePreview.
func (m *Model) preview() tea.Cmd {
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

	seed, err := m.seed()
	if err != nil {
		m.setStatus("Seed должен быть целым неотрицательным числом.", true)
		m.focus(seedStep)
		return nil
	}

	m.records = m.services.Generate(m.fields, count, seed)
	m.outputPath = ""
	m.setStatus(fmt.Sprintf("Предпросмотр готов: %d записей. Seed: %d. Нажми Enter, чтобы сохранить именно их.", count, seed), false)
	return nil
}

// seed читает seed из формы или создаёт новый и оставляет его в форме для повторной генерации.
func (m *Model) seed() (uint64, error) {
	value := strings.TrimSpace(m.seedInput.Value())
	if value != "" {
		return strconv.ParseUint(value, 10, 64)
	}

	var bytes [8]byte
	if _, err := cryptorand.Read(bytes[:]); err != nil {
		return 0, err
	}
	seed := binary.LittleEndian.Uint64(bytes[:])
	m.seedInput.SetValue(strconv.FormatUint(seed, 10))
	return seed, nil
}

// savePreview записывает подготовленные данные в JSON и открывает папку результата.
// Если предпросмотра ещё нет, сначала подготавливает записи.
func (m *Model) savePreview() tea.Cmd {
	if len(m.records) == 0 || m.outputPath != "" {
		return m.preview()
	}

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
	m.setStatus(fmt.Sprintf("Готово. В файл сохранено записей: %d.", len(m.records)), false)
	return openOutputFolder(m.outputPath)
}
