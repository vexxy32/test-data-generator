package main

import (
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

var (
	accent = lipgloss.Color("#78DCCA")
	muted  = lipgloss.Color("#8993A5")
	light  = lipgloss.Color("#E8EDF7")

	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(light)
	accentStyle  = lipgloss.NewStyle().Bold(true).Foreground(accent)
	mutedStyle   = lipgloss.NewStyle().Foreground(muted)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8792"))
	successStyle = lipgloss.NewStyle().Foreground(accent)
	panelStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#506079")).Padding(0, 1)
)

// model хранит введённые значения, схему и результат. Генератор получает только fields и count.
type model struct {
	fields []FieldDefinition

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

func initialModel() model {
	m := model{
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

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) focusedInput() *textinput.Model {
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

func (m *model) focus(step int) tea.Cmd {
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

func (m *model) focusSchema() {
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

func (m model) usesBounds() bool {
	return strings.ToLower(strings.TrimSpace(m.typeInput.Value())) != "boolean"
}

func (m *model) moveForm(direction int) tea.Cmd {
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

func (m *model) clearForm() {
	m.nameInput.SetValue("")
	m.typeInput.SetValue("")
	m.minInput.SetValue("")
	m.maxInput.SetValue("")
	m.editIndex = -1
}

func (m *model) editField() tea.Cmd {
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

func (m *model) deleteField() {
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

func (m *model) setStatus(message string, isError bool) {
	m.status = message
	m.statusIsError = isError
}

func (m *model) addField() tea.Cmd {
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

	newField := FieldDefinition{Name: name, Type: fieldType}
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

func (m *model) generate() tea.Cmd {
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

	m.records = generateRecords(m.fields, count)
	m.outputPath = ""
	path, err := exportToFile(m.records, m.fileNameInput.Value())
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

func (m model) visibleSchemaRows() int {
	rows := m.height - 24
	if rows < 1 {
		rows = 1
	}
	if rows > 5 {
		rows = 5
	}
	return rows
}

type folderOpenedMsg struct{ err error }

func openOutputFolder(path string) tea.Cmd {
	return func() tea.Msg {
		return folderOpenedMsg{err: exec.Command("explorer.exe", filepath.Dir(path)).Start()}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case folderOpenedMsg:
		if msg.err != nil {
			m.setStatus("Не удалось открыть папку: "+msg.err.Error(), true)
		} else {
			m.setStatus("Папка с результатом открыта.", false)
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputWidth := msg.Width - 25
		if inputWidth < 12 {
			inputWidth = 12
		}
		if inputWidth > 42 {
			inputWidth = 42
		}
		m.nameInput.SetWidth(inputWidth)
		m.typeInput.SetWidth(inputWidth)
		m.minInput.SetWidth(inputWidth)
		m.maxInput.SetWidth(inputWidth)
		m.countInput.SetWidth(inputWidth)
		m.fileNameInput.SetWidth(inputWidth)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.editIndex >= 0 {
				m.clearForm()
				m.setStatus("Редактирование отменено.", false)
				return m, m.focus(nameStep)
			}
			return m, tea.Quit
		case "tab":
			if m.schemaFocus {
				return m, m.focus(nameStep)
			}
			return m, m.moveForm(1)
		case "shift+tab":
			if m.schemaFocus {
				return m, m.focus(countStep)
			}
			return m, m.moveForm(-1)
		case "up":
			if m.schemaFocus {
				if m.selectedField > 0 {
					m.selectedField--
				}
				return m, nil
			}
			return m, m.moveForm(-1)
		case "down":
			if m.schemaFocus {
				if m.selectedField+1 < len(m.fields) {
					m.selectedField++
				}
				return m, nil
			}
			return m, m.moveForm(1)
		case "+":
			if m.schemaFocus {
				m.clearForm()
				m.setStatus("Введи данные нового поля и нажми + ещё раз.", false)
				return m, m.focus(nameStep)
			}
			return m, m.addField()
		case "ctrl+x", "delete":
			if m.schemaFocus {
				m.deleteField()
				return m, nil
			}
		case "enter":
			return m, m.generate()
		}
		if m.schemaFocus {
			switch msg.String() {
			case "e":
				return m, m.editField()
			case "o":
				if m.outputPath == "" {
					m.setStatus("Файл ещё не сохранён. После экспорта путь появится здесь.", true)
					return m, nil
				}
				return m, openOutputFolder(m.outputPath)
			}
			return m, nil
		}
		if msg.String() == "ctrl+o" {
			if m.outputPath == "" {
				m.setStatus("Файл ещё не сохранён. После экспорта путь появится здесь.", true)
				return m, nil
			}
			return m, openOutputFolder(m.outputPath)
		}
	}

	input := m.focusedInput()
	var cmd tea.Cmd
	*input, cmd = input.Update(msg)
	return m, cmd
}

func (m model) inputRow(index int, label string, input textinput.Model) string {
	marker := "  "
	labelText := mutedStyle.Render(fmt.Sprintf("%-12s", label))
	if !m.schemaFocus && m.step == index {
		marker = accentStyle.Render("› ")
		labelText = accentStyle.Render(fmt.Sprintf("%-12s", label))
	}
	return marker + labelText + input.View()
}

func (m model) schemaView(limit int) string {
	if len(m.fields) == 0 {
		return mutedStyle.Render("Пока пусто. Заполни форму и нажми +.")
	}

	var lines []string
	start := m.selectedField - limit + 1
	if start < 0 {
		start = 0
	}
	for i := start; i < len(m.fields) && i < start+limit; i++ {
		field := m.fields[i]
		marker := "  "
		if m.schemaFocus && i == m.selectedField {
			marker = "› "
		}
		line := fmt.Sprintf("%s%02d  %-16s %-8s", marker, i+1, field.Name, field.Type)
		if field.Type == "integer" || field.Type == "float" {
			line += fmt.Sprintf("  %g .. %g", field.Min, field.Max)
		} else if field.Type == "string" {
			line += fmt.Sprintf("  длина %g .. %g", field.Min, field.Max)
		}
		if m.schemaFocus && i == m.selectedField {
			line = accentStyle.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m model) resultsView() string {
	if len(m.records) == 0 {
		return mutedStyle.Render("После генерации здесь появится путь к файлу.")
	}
	if m.outputPath == "" {
		return errorStyle.Render("Не удалось сохранить файл. Причина указана ниже.")
	}
	return "Файл: " + m.outputPath
}

func (m model) panel(heading, content string, width int) string {
	return panelStyle.Width(width).Render(accentStyle.Render(heading) + "\n" + content)
}

func (m model) View() tea.View {
	width := m.width - 4
	if width < 28 {
		width = 28
	}
	if width > 94 {
		width = 94
	}
	// На терминале 80x24 оставляем место для формы, схемы, результата и подсказок.
	schemaRows := m.visibleSchemaRows()

	var screen strings.Builder
	screen.WriteString(titleStyle.Render("TEST DATA LAB"))
	screen.WriteString(mutedStyle.Render(fmt.Sprintf("   /   %d полей   /   %d записей", len(m.fields), len(m.records))))
	screen.WriteString("\n")

	form := strings.Join([]string{
		m.inputRow(nameStep, "Имя", m.nameInput),
		m.inputRow(typeStep, "Тип", m.typeInput),
	}, "\n")
	fieldType := strings.ToLower(strings.TrimSpace(m.typeInput.Value()))
	if fieldType == "boolean" {
		form += "\n" + mutedStyle.Render("  Минимум       не используется")
		form += "\n" + mutedStyle.Render("  Максимум      не используется")
	} else if fieldType == "string" {
		form += "\n" + m.inputRow(minStep, "Мин. длина", m.minInput)
		form += "\n" + m.inputRow(maxStep, "Макс. длина", m.maxInput)
	} else {
		form += "\n" + m.inputRow(minStep, "Минимум", m.minInput)
		form += "\n" + m.inputRow(maxStep, "Максимум", m.maxInput)
	}
	formTitle := "НОВОЕ ПОЛЕ"
	if m.editIndex >= 0 {
		formTitle = fmt.Sprintf("РЕДАКТИРОВАНИЕ ПОЛЯ %02d", m.editIndex+1)
	}
	screen.WriteString(m.panel(formTitle, form, width))
	screen.WriteString("\n")
	screen.WriteString("ГЕНЕРАЦИЯ   " + m.inputRow(countStep, "Количество", m.countInput))
	screen.WriteString("\n")
	screen.WriteString("СОХРАНЕНИЕ " + m.inputRow(fileNameStep, "Имя файла", m.fileNameInput))
	screen.WriteString("\n")
	screen.WriteString(m.panel(fmt.Sprintf("СХЕМА · %d", len(m.fields)), m.schemaView(schemaRows), width))
	screen.WriteString("\n")
	resultTitle := "РЕЗУЛЬТАТ"
	if len(m.records) > 0 {
		resultTitle = fmt.Sprintf("РЕЗУЛЬТАТ · %d записей по %d полям", len(m.records), len(m.fields))
	}
	screen.WriteString(m.panel(resultTitle, m.resultsView(), width))
	screen.WriteString("\n")

	if m.statusIsError {
		screen.WriteString(errorStyle.Render("! " + m.status))
	} else {
		screen.WriteString(successStyle.Render("● " + m.status))
	}
	screen.WriteString("\n")
	screen.WriteString(mutedStyle.Render("↑↓/Tab · + поле · E правка · Del удалить · Enter создать · O папка · Esc выход"))

	return tea.NewView(screen.String())
}
