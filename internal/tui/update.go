package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/schema"
)

type folderOpenedMsg struct{ err error }

// Update обрабатывает нажатия клавиш и сообщения от запущенных команд TUI.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.precisionInput.SetWidth(inputWidth)
		m.countInput.SetWidth(inputWidth)
		m.seedInput.SetWidth(inputWidth)
		m.fileNameInput.SetWidth(inputWidth)
		m.templateNameInput.SetWidth(inputWidth)
		return m, nil

	case tea.KeyPressMsg:
		if m.templatePicker {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.templatePicker = false
				m.setStatus("Выбор шаблона отменён.", false)
			case "up":
				if m.selectedTemplate > 0 {
					m.selectedTemplate--
				}
			case "down":
				if m.selectedTemplate+1 < len(m.templateNames) {
					m.selectedTemplate++
				}
			case "enter":
				m.templateNameInput.SetValue(m.templateNames[m.selectedTemplate])
				m.templatePicker = false
				return m, m.loadTemplate()
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			return m, m.newTest()
		case "ctrl+s":
			m.saveTemplate()
			return m, nil
		case "ctrl+l":
			m.openTemplatePicker()
			return m, nil
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
				return m, m.focus(templateNameStep)
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
			return m, m.savePreview()
		case "ctrl+p":
			return m, m.preview()
		}
		if m.schemaFocus {
			switch msg.String() {
			case "e":
				return m, m.editField()
			case "p":
				return m, m.preview()
			case "o":
				if m.outputPath == "" {
					m.setStatus("Файл ещё не сохранён. После экспорта путь появится здесь.", true)
					return m, nil
				}
				return m, openOutputFolder(m.outputPath)
			case "s":
				m.saveTemplate()
				return m, nil
			case "l":
				m.openTemplatePicker()
				return m, nil
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
	if m.step == nameStep {
		if automaticType, automatic := schema.AutomaticStringType(m.nameInput.Value()); automatic {
			m.typeInput.SetValue(automaticType)
		}
	}
	return m, cmd
}
