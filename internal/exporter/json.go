package exporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportToFile writes records to the default output directory.
func ExportToFile(records []map[string]any, fileName string) (string, error) {
	return exportToDirectory(records, fileName, "output")
}

func exportToDirectory(records []map[string]any, fileName, outputDir string) (string, error) {
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

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку %s: %w", outputDir, err)
	}

	nameWithoutExtension := strings.TrimSuffix(fileName, extension)
	path := filepath.Join(outputDir, fileName)
	for copyNumber := 2; ; copyNumber++ {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("не удалось проверить файл: %w", err)
		}
		path = filepath.Join(outputDir, fmt.Sprintf("%s_%d%s", nameWithoutExtension, copyNumber, extension))
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
