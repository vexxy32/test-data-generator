package exporter

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type JSONExporter struct{}

func (e JSONExporter) Export(
	fileName string,
	produceRecords func(func(map[string]any) error) error,
) (string, error) {
	return ExportToFileStream(fileName, produceRecords)
}

// ExportToFileStream открывает файл и записывает в него записи по мере их получения от produceRecords.
func ExportToFileStream(
	fileName string,
	produceRecords func(onRecord func(record map[string]any) error) error,
) (string, error) {
	return exportToDirectoryStream(fileName, "output", produceRecords)
}

func exportToDirectoryStream(
	fileName, outputDir string,
	produceRecords func(onRecord func(record map[string]any) error) error,
) (path string, err error) {
	if produceRecords == nil {
		return "", errors.New("функция-источник записей не задана")
	}

	filePath, err := outputFilePath(fileName, outputDir)
	if err != nil {
		return "", err
	}

	newFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return "", err
	}

	defer func() {
		if closeErr := newFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
		if err != nil {
			os.Remove(filePath)
			path = ""
		}
	}()
	writer := bufio.NewWriter(newFile)
	if _, err := writer.WriteString("[\n"); err != nil {
		return "", err
	}

	firstRecord := true

	onRecord := func(record map[string]any) error {
		data, err := json.MarshalIndent(record, "  ", "  ")
		if err != nil {
			return err
		}

		if !firstRecord {
			if _, err := writer.WriteString(",\n"); err != nil {
				return err
			}
		}

		if _, err := writer.WriteString("  "); err != nil {
			return err
		}
		if _, err := writer.Write(data); err != nil {
			return err
		}

		firstRecord = false
		return nil
	}

	if err := produceRecords(onRecord); err != nil {
		return "", err
	}

	if _, err := writer.WriteString("\n]\n"); err != nil {
		return "", err
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}

	return filePath, nil
}

// outputFilePath проверяет имя и подбирает новый путь, не затирая старый файл.
func outputFilePath(fileName, outputDir string) (string, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		fileName = "generated_tests"
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

	return path, nil
}
