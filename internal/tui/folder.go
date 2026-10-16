package tui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "charm.land/bubbletea/v2"
)

func openOutputFolder(path string) tea.Cmd {
	return func() tea.Msg {
		directory := filepath.Dir(path)
		var command *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			command = exec.Command("explorer.exe", directory)
		case "darwin":
			command = exec.Command("open", directory)
		case "linux":
			command = exec.Command("xdg-open", directory)
		default:
			return folderOpenedMsg{err: fmt.Errorf("открытие папки не поддерживается в %s", runtime.GOOS)}
		}
		return folderOpenedMsg{err: command.Run()}
	}
}
