package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func exportToFile(records []map[string]any, fileName string) (string, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return "", errors.New("укажи имя файла")
	}
	if strings.ContainsAny(fileName, `\\/:*?"<>|`) {
		return "", errors.New("имя файла не должно содержать \\ / : * ? \" < > или |")
	}

	extension := filepath.Ext(fileName)
	if extension == "" {
		extension = ".json"
		fileName += extension
	} else if !strings.EqualFold(extension, ".json") {
		return "", errors.New("используй расширение .json или не указывай расширение")
	}

	if err := os.MkdirAll("output", 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку output: %w", err)
	}

	nameWithoutExtension := strings.TrimSuffix(fileName, extension)
	path := filepath.Join("output", fileName)
	for copyNumber := 2; ; copyNumber++ {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("не удалось проверить файл: %w", err)
		}
		path = filepath.Join("output", fmt.Sprintf("%s_%d%s", nameWithoutExtension, copyNumber, extension))
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}
