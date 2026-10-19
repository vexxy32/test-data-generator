package templates

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// List возвращает имена JSON-шаблонов в алфавитном порядке.
func List() ([]string, error) {
	entries, err := os.ReadDir(templatesDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать папку шаблонов: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
