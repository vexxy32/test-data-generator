package templates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Save сохраняет структуру Template в папку templates и возвращает полный путь к файлу.
func Save(templateName string, tmpl Template) (string, error) {
	fileName, err := normalizeTemplateName(templateName)
	if err != nil {
		return "", err
	}

	if err := tmpl.Validate(); err != nil {
		return "", err
	}

	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку шаблонов: %w", err)
	}

	relPath := filepath.Join(templatesDir, fileName)
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		absPath = relPath
	}

	data, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ошибка формирования JSON: %w", err)
	}

	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return "", fmt.Errorf("ошибка сохранения файла: %w", err)
	}

	return absPath, nil
}
