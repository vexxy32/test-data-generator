package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/exporter"
	"github.com/vexxy32/test-data-generator/internal/generator"
	"github.com/vexxy32/test-data-generator/internal/templates"
	"github.com/vexxy32/test-data-generator/internal/tui"
)

func main() {
	// Здесь соединяются части программы: TUI получает функции генерации и экспорта.
	services := tui.Services{
		GenerateStream: generator.GenerateRecordsStream,
		ExportStream:   exporter.ExportToFileStream,
		SaveTemplate:   templates.Save,
		LoadTemplate:   templates.Load,
		ListTemplates:  templates.List,
	}

	model := tui.NewModel(services)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		fmt.Printf(" Ошибка запуска UI: %+v\n", err)
		os.Exit(1)
	}
}
