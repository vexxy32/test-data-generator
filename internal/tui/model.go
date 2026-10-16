package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/schema"
)

const (
	nameStep = iota
	typeStep
	minStep
	maxStep
	countStep
	fileNameStep
	inputCount
)

const (
	maxRecordCount  = 10_000
	maxStringLength = 1_000
)

// Services connects the interface to application operations without importing their implementations.
type Services struct {
	Generate func([]schema.FieldDefinition, int) []map[string]any
	Export   func([]map[string]any, string) (string, error)
}

type applicationModel struct {
	services Services
	fields   []schema.FieldDefinition

	nameInput     textinput.Model
	typeInput     textinput.Model
	minInput      textinput.Model
	maxInput      textinput.Model
	countInput    textinput.Model
	fileNameInput textinput.Model
	step          int
	schemaFocus   bool
	selectedField int
	editIndex     int // -1 означает добавление нового поля; иначе сохраняем изменения в fields[editIndex].

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

func NewModel(services Services) tea.Model {
	m := applicationModel{
		services:      services,
		nameInput:     newInput("например, age"),
		typeInput:     newInput("integer / float / string / boolean"),
		minInput:      newInput("например, 0"),
		maxInput:      newInput("например, 120"),
		countInput:    newInput("сколько записей создать"),
		fileNameInput: newInput("например, purchases"),
		status:        "Заполни поле и нажми +, затем укажи количество и имя файла.",
		editIndex:     -1,
		width:         80,
		height:        24,
	}
	m.countInput.SetValue("5")
	m.fileNameInput.SetValue("generated_tests")
	m.nameInput.Focus()
	return m
}

func (m applicationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *applicationModel) focusedInput() *textinput.Model {
	switch m.step {
	case nameStep:
		return &m.nameInput
	case typeStep:
		return &m.typeInput
	case minStep:
		return &m.minInput
	case maxStep:
		return &m.maxInput
	case countStep:
		return &m.countInput
	default:
		return &m.fileNameInput
	}
}

func (m *applicationModel) focus(step int) tea.Cmd {
	m.nameInput.Blur()
	m.typeInput.Blur()
	m.minInput.Blur()
	m.maxInput.Blur()
	m.countInput.Blur()
	m.fileNameInput.Blur()
	m.schemaFocus = false
	m.step = step
	return m.focusedInput().Focus()
}

func (m *applicationModel) focusSchema() {
	m.nameInput.Blur()
	m.typeInput.Blur()
	m.minInput.Blur()
	m.maxInput.Blur()
	m.countInput.Blur()
	m.fileNameInput.Blur()
	m.schemaFocus = true
	if m.selectedField >= len(m.fields) {
		m.selectedField = len(m.fields) - 1
	}
}

func (m applicationModel) usesBounds() bool {
	return strings.ToLower(strings.TrimSpace(m.typeInput.Value())) != "boolean"
}

func (m *applicationModel) moveForm(direction int) tea.Cmd {
	next := m.step + direction
	if !m.usesBounds() && (next == minStep || next == maxStep) {
		if direction > 0 {
			next = countStep
		} else {
			next = typeStep
		}
	}
	if next >= inputCount || next < 0 {
		if len(m.fields) > 0 {
			m.focusSchema()
			return nil
		}
		next = (next + inputCount) % inputCount
	}
	return m.focus(next)
}

func (m *applicationModel) setStatus(message string, isError bool) {
	m.status = message
	m.statusIsError = isError
}

func (m applicationModel) visibleSchemaRows() int {
	rows := m.height - 24
	if rows < 1 {
		rows = 1
	}
	if rows > 5 {
		rows = 5
	}
	return rows
}
