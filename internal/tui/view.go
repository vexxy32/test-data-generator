package tui

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func (m Model) inputRow(index int, label string, input textinput.Model) string {
	marker := "  "
	labelText := mutedStyle.Render(fmt.Sprintf("%-12s", label))
	if !m.schemaFocus && m.step == index {
		marker = accentStyle.Render("› ")
		labelText = accentStyle.Render(fmt.Sprintf("%-12s", label))
	}
	return marker + labelText + input.View()
}

// schemaView показывает добавленные поля и выделяет выбранное для правки.
func (m Model) schemaView(limit int) string {
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
			if field.Type == "float" {
				line += fmt.Sprintf("  точность %d", field.Precision)
			}
		} else if field.Type == "string" && !isAutomaticStringFieldName(field.Name) {
			line += fmt.Sprintf("  длина %g .. %g", field.Min, field.Max)
		}
		if m.schemaFocus && i == m.selectedField {
			line = accentStyle.Render(line)
		}
		lines = append(lines, line)
	}
	if len(m.fields) > limit {
		end := min(start+limit, len(m.fields))
		lines = append(lines, mutedStyle.Render(fmt.Sprintf("Поля %d–%d из %d · Tab на схему, ↑↓ листать", start+1, end, len(m.fields))))
	}
	return strings.Join(lines, "\n")
}

func (m Model) resultsView() string {
	if len(m.records) == 0 {
		return mutedStyle.Render("После сохранения здесь появится путь к файлу.")
	}
	if m.outputPath == "" {
		return mutedStyle.Render("Предпросмотр готов. Нажми Enter, чтобы создать JSON со всеми записями.")
	}
	return "Файл: " + m.outputPath
}

// previewView выводит первые три уже подготовленные записи.
func (m Model) previewView() string {
	if len(m.records) == 0 {
		return mutedStyle.Render("Нажми Ctrl+P, чтобы подготовить три записи для просмотра.")
	}

	previewCount := min(3, len(m.records))
	data, err := json.MarshalIndent(m.records[:previewCount], "", "  ")
	if err != nil {
		return errorStyle.Render("Не удалось подготовить предпросмотр: " + err.Error())
	}
	return string(data)
}

func (m Model) panel(heading, content string, width int) string {
	return panelStyle.Width(width).Render(accentStyle.Render(heading) + "\n" + content)
}

// templatePickerView показывает список файлов вместо формы, пока пользователь выбирает шаблон.
func (m Model) templatePickerView(width int) tea.View {
	limit := min(12, max(3, m.height-7))
	start := max(0, m.selectedTemplate-limit+1)
	var lines []string
	for i := start; i < len(m.templateNames) && i < start+limit; i++ {
		line := fmt.Sprintf("  %s", m.templateNames[i])
		if i == m.selectedTemplate {
			line = accentStyle.Render("› " + m.templateNames[i])
		}
		lines = append(lines, line)
	}
	content := strings.Join(lines, "\n")
	content += mutedStyle.Render(fmt.Sprintf("\n\n%d из %d", m.selectedTemplate+1, len(m.templateNames)))
	screen := titleStyle.Render("ВЫБОР ШАБЛОНА") + "\n" + m.panel("СОХРАНЁННЫЕ ШАБЛОНЫ", content, width)
	screen += "\n" + mutedStyle.Render("↑↓ выбрать · Enter загрузить · Esc назад")
	return tea.NewView(screen)
}

// View собирает экран из формы, схемы, предпросмотра и результата.
func (m Model) View() tea.View {
	width := m.width - 4
	if width < 28 {
		width = 28
	}
	if width > 94 {
		width = 94
	}
	if m.templatePicker {
		return m.templatePickerView(width)
	}
	// На терминале 80x24 оставляем место для формы, схемы, результата и подсказок.
	schemaRows := m.visibleSchemaRows()

	var screen strings.Builder
	recordCount := len(m.records)
	if m.outputPath != "" {
		if count, err := strconv.Atoi(strings.TrimSpace(m.countInput.Value())); err == nil {
			recordCount = count
		}
	}
	screen.WriteString(titleStyle.Render("TEST DATA LAB"))
	screen.WriteString(mutedStyle.Render(fmt.Sprintf("   /   %d полей   /   %d записей", len(m.fields), recordCount)))
	screen.WriteString("\n")

	form := strings.Join([]string{
		m.inputRow(nameStep, "Имя", m.nameInput),
		m.inputRow(typeStep, "Тип", m.typeInput),
	}, "\n")
	fieldType := strings.ToLower(strings.TrimSpace(m.typeInput.Value()))
	if fieldType == "boolean" || isAutomaticStringFieldName(m.nameInput.Value()) {
		form += "\n" + mutedStyle.Render("  Минимум       не используется")
		form += "\n" + mutedStyle.Render("  Максимум      не используется")
		if isAutomaticStringFieldName(m.nameInput.Value()) {
			form += "\n" + mutedStyle.Render("  Тип           string определяется автоматически")
		}
	} else if fieldType == "string" {
		form += "\n" + m.inputRow(minStep, "Мин. длина", m.minInput)
		form += "\n" + m.inputRow(maxStep, "Макс. длина", m.maxInput)
		form += "\n" + mutedStyle.Render("  Пустая длина: 1 .. 20")
	} else {
		form += "\n" + m.inputRow(minStep, "Минимум", m.minInput)
		form += "\n" + m.inputRow(maxStep, "Максимум", m.maxInput)
		form += "\n" + mutedStyle.Render("  Пустые границы: 0 .. 100")
		if fieldType == "float" {
			form += "\n" + m.inputRow(precisionStep, "Точность", m.precisionInput)
		}
	}
	formTitle := "НОВОЕ ПОЛЕ"
	if m.editIndex >= 0 {
		formTitle = fmt.Sprintf("РЕДАКТИРОВАНИЕ ПОЛЯ %02d", m.editIndex+1)
	}
	screen.WriteString(m.panel(formTitle, form, width))
	screen.WriteString("\n")
	screen.WriteString("ГЕНЕРАЦИЯ   " + m.inputRow(countStep, "Количество", m.countInput))
	screen.WriteString("\n")
	screen.WriteString("ГЕНЕРАЦИЯ   " + m.inputRow(seedStep, "Seed", m.seedInput))
	screen.WriteString("\n")
	screen.WriteString("СОХРАНЕНИЕ " + m.inputRow(fileNameStep, "Имя файла", m.fileNameInput))
	screen.WriteString("\n")
	screen.WriteString("ШАБЛОН     " + m.inputRow(templateNameStep, "Имя шаблона", m.templateNameInput))
	screen.WriteString("\n")
	screen.WriteString(m.panel(fmt.Sprintf("СХЕМА · %d", len(m.fields)), m.schemaView(schemaRows), width))
	screen.WriteString("\n")
	previewTitle := "ПРЕДПРОСМОТР"
	if len(m.records) > 0 {
		previewTitle = fmt.Sprintf("ПРЕДПРОСМОТР · первые %d из %s записей", len(m.records), strings.TrimSpace(m.countInput.Value()))
	}
	screen.WriteString(m.panel(previewTitle, m.previewView(), width))
	screen.WriteString("\n")
	resultTitle := "РЕЗУЛЬТАТ"
	if m.outputPath != "" {
		resultTitle = fmt.Sprintf("РЕЗУЛЬТАТ · %d записей по %d полям", recordCount, len(m.fields))
	}
	screen.WriteString(m.panel(resultTitle, m.resultsView(), width))
	screen.WriteString("\n")

	if m.statusIsError {
		screen.WriteString(errorStyle.Render("! " + m.status))
	} else {
		screen.WriteString(successStyle.Render("● " + m.status))
	}
	screen.WriteString("\n")
	screen.WriteString(mutedStyle.Render("Tab/↑↓ выбор · + поле · Ctrl+P просмотр · Enter создать JSON · Ctrl+N заново"))
	screen.WriteString("\n")
	screen.WriteString(mutedStyle.Render("Ctrl+S шаблон сохранить · Ctrl+L выбрать · E/Del поле · O папка · Esc выход"))

	return tea.NewView(screen.String())
}
