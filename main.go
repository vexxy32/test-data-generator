package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/vexxy32/test-data-generator/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel(tui.Services{
		Generate: generateRecords,
		Export:   exportToFile,
	}))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка запуска TUI: %v\n", err)
		os.Exit(1)
	}
}
