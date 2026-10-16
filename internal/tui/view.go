package tui

import (
	"encoding/json"
	"fmt"
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

func (m applicationModel) inputRow(index int, label string, input textinput.Model) string {
	marker := "  "
	labelText := mutedStyle.Render(fmt.Sprintf("%-12s", label))
	if !m.schemaFocus && m.step == index {
		marker = accentStyle.Render("› ")
		labelText = accentStyle.Render(fmt.Sprintf("%-12s", label))
	}
	return marker + labelText + input.View()
}

func (m applicationModel) schemaView(limit int) string {
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
	return strings.Join(lines, "\n")
}

func (m applicationModel) resultsView() string {
	if len(m.records) == 0 {
		return mutedStyle.Render("После сохранения здесь появится путь к файлу.")
	}
	if m.outputPath == "" {
		return mutedStyle.Render("Предпросмотр готов. Нажми Enter, чтобы сохранить эти записи.")
	}
	return "Файл: " + m.outputPath
}

func (m applicationModel) previewView() string {
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

func (m applicationModel) panel(heading, content string, width int) string {
	return panelStyle.Width(width).Render(accentStyle.Render(heading) + "\n" + content)
}

func (m applicationModel) View() tea.View {
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
	if fieldType == "boolean" || isAutomaticStringFieldName(m.nameInput.Value()) {
		form += "\n" + mutedStyle.Render("  Минимум       не используется")
		form += "\n" + mutedStyle.Render("  Максимум      не используется")
		if isAutomaticStringFieldName(m.nameInput.Value()) {
			form += "\n" + mutedStyle.Render("  Тип           string определяется автоматически")
		}
	} else if fieldType == "string" {
		form += "\n" + m.inputRow(minStep, "Мин. длина", m.minInput)
		form += "\n" + m.inputRow(maxStep, "Макс. длина", m.maxInput)
	} else {
		form += "\n" + m.inputRow(minStep, "Минимум", m.minInput)
		form += "\n" + m.inputRow(maxStep, "Максимум", m.maxInput)
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
	screen.WriteString("СОХРАНЕНИЕ " + m.inputRow(fileNameStep, "Имя файла", m.fileNameInput))
	screen.WriteString("\n")
	screen.WriteString(m.panel(fmt.Sprintf("СХЕМА · %d", len(m.fields)), m.schemaView(schemaRows), width))
	screen.WriteString("\n")
	previewTitle := "ПРЕДПРОСМОТР"
	if len(m.records) > 0 {
		previewTitle = fmt.Sprintf("ПРЕДПРОСМОТР · первые %d из %d записей", min(3, len(m.records)), len(m.records))
	}
	screen.WriteString(m.panel(previewTitle, m.previewView(), width))
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
	screen.WriteString(mutedStyle.Render("↑↓/Tab · + поле · P/Ctrl+P предпросмотр · Enter сохранить · E правка · Del удалить · O папка · Esc выход"))

	return tea.NewView(screen.String())
}
