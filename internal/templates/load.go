package templates

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Load читает JSON-файл шаблона из папки templates, проверяет структуру и валидирует данные.
func Load(templateName string) (Template, error) {
	var tmpl Template

	fileName, err := normalizeTemplateName(templateName)
	if err != nil {
		return tmpl, err
	}

	filePath := filepath.Join(templatesDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return tmpl, errors.New("файл не найден")
		}
		return tmpl, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&tmpl); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF), errors.Is(err, io.EOF):
			return tmpl, errors.New("JSON в файле повреждён")
		case errors.As(err, &unmarshalTypeErr), strings.HasPrefix(err.Error(), "json: unknown field"):
			return tmpl, errors.New("JSON не соответствует структуре шаблона")
		default:
			return tmpl, errors.New("JSON в файле повреждён")
		}
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return tmpl, errors.New("после шаблона обнаружены лишние данные")
	}

	if err := tmpl.Validate(); err != nil {
		return tmpl, err
	}

	return tmpl, nil
}
