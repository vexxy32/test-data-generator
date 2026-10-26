package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/schema"
	"github.com/vexxy32/test-data-generator/internal/templates"
)

const (
	nameStep = iota
	typeStep
	minStep
	maxStep
	precisionStep
	allowEmptyStep
	countStep
	seedStep
	fileNameStep
	templateNameStep
	inputCount
)

const (
	maxRecordCount    = 100_000_000
	maxStringLength   = 1_000
	maxFloatPrecision = 6
)

// Exporter сохраняет поток записей и возвращает путь к результату.
type Exporter interface {
	Export(string, func(func(map[string]any) error) error) (string, error)
}

// TemplateStore сохраняет, загружает и перечисляет шаблоны.
type TemplateStore interface {
	Save(string, templates.Template) (string, error)
	Load(string) (templates.Template, error)
	List() ([]string, error)
}

// Services хранит операции, переданные в TUI из main.
type Services struct {
	GenerateStream func([]schema.FieldDefinition, int, uint64, func(map[string]any) error) error
	ExportStream   Exporter
	Templates      TemplateStore
}

// Model хранит состояние формы, схемы и подготовленного результата.
type Model struct {
	services Services
	fields   []schema.FieldDefinition

	nameInput         textinput.Model
	typeInput         textinput.Model
	minInput          textinput.Model
	maxInput          textinput.Model
	precisionInput    textinput.Model
	allowEmpty        bool
	countInput        textinput.Model
	seedInput         textinput.Model
	fileNameInput     textinput.Model
	templateNameInput textinput.Model
	step              int
	schemaFocus       bool
	selectedField     int
	editIndex         int // -1 означает добавление нового поля; иначе сохраняем изменения в fields[editIndex].
	templatePicker    bool
	templateNames     []string
	selectedTemplate  int

	records       []map[string]any
	outputPath    string
	status        string
	statusIsError bool
	width         int
	height        int
}

func newInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = 80
	input.SetWidth(28)
	return input
}

// NewModel создаёт начальное состояние интерфейса с настройками по умолчанию.
func NewModel(services Services) Model {
	m := Model{
		services:          services,
		nameInput:         newInput("например, age"),
		typeInput:         newInput("integer / float / string / boolean"),
		minInput:          newInput("например, 0"),
		maxInput:          newInput("например, 120"),
		precisionInput:    newInput("от 0 до 6"),
		countInput:        newInput("сколько записей создать"),
		seedInput:         newInput("пусто — новый случайный seed"),
		fileNameInput:     newInput("например, purchases"),
		templateNameInput: newInput("например, users"),
		status:            "Заполни поле и нажми +, затем укажи количество и имя файла.",
		editIndex:         -1,
		width:             80,
		height:            24,
	}
	m.countInput.SetValue("5")
	m.precisionInput.SetValue("2")
	m.fileNameInput.SetValue("generated_tests")
	m.nameInput.Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) focusedInput() *textinput.Model {
	switch m.step {
	case nameStep:
		return &m.nameInput
	case typeStep:
		return &m.typeInput
	case minStep:
		return &m.minInput
	case maxStep:
		return &m.maxInput
	case precisionStep:
		return &m.precisionInput
	case allowEmptyStep:
		return nil
	case countStep:
		return &m.countInput
	case seedStep:
		return &m.seedInput
	case fileNameStep:
		return &m.fileNameInput
	default:
		return &m.templateNameInput
	}
}

// focus переводит ввод на указанный элемент формы.
func (m *Model) focus(step int) tea.Cmd {
	m.nameInput.Blur()
	m.typeInput.Blur()
	m.minInput.Blur()
	m.maxInput.Blur()
	m.precisionInput.Blur()
	m.countInput.Blur()
	m.seedInput.Blur()
	m.fileNameInput.Blur()
	m.templateNameInput.Blur()
	m.schemaFocus = false
	m.step = step
	if input := m.focusedInput(); input != nil {
		return input.Focus()
	}
	return nil
}

func (m *Model) focusSchema() {
	m.nameInput.Blur()
	m.typeInput.Blur()
	m.minInput.Blur()
	m.maxInput.Blur()
	m.precisionInput.Blur()
	m.countInput.Blur()
	m.seedInput.Blur()
	m.fileNameInput.Blur()
	m.templateNameInput.Blur()
	m.schemaFocus = true
	if m.selectedField >= len(m.fields) {
		m.selectedField = len(m.fields) - 1
	}
}

// isInputAvailable определяет, можно ли перейти к параметру для выбранного типа.
func (m Model) isInputAvailable(step int) bool {
	fieldType := strings.ToLower(strings.TrimSpace(m.typeInput.Value()))
	switch step {
	case minStep, maxStep:
		return fieldType != "boolean" && !isAutomaticStringFieldName(m.nameInput.Value())
	case precisionStep:
		return fieldType == "float"
	default:
		return true
	}
}

func isAutomaticStringFieldName(name string) bool {
	_, automatic := schema.AutomaticStringType(name)
	return automatic
}

func (m *Model) moveForm(direction int) tea.Cmd {
	next := m.step
	for {
		next += direction
		if next >= inputCount || next < 0 {
			if len(m.fields) > 0 {
				m.focusSchema()
				return nil
			}
			next = (next + inputCount) % inputCount
		}
		if m.isInputAvailable(next) {
			return m.focus(next)
		}
	}
}

func (m *Model) setStatus(message string, isError bool) {
	m.status = message
	m.statusIsError = isError
}

func (m Model) visibleSchemaRows() int {
	rows := m.height - 24
	if rows < 1 {
		rows = 1
	}
	if rows > 5 {
		rows = 5
	}
	return rows
}
