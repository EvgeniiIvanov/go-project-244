package app

import (
	"code/internal/differ"
	"code/internal/formatter"
	"code/internal/parser"
	"fmt"
	"path/filepath"
)

func getFileFormat(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	default:
		return ""
	}
}

func Run(filePath1, filePath2 string, format string) (string, error) {
	format1 := getFileFormat(filePath1)
	format2 := getFileFormat(filePath2)

	if format1 != format2 {
		return "", fmt.Errorf("cannot compare files of different formats: %s vs %s", format1, format2)
	}

	data1, err := parser.Parse(filePath1)
	if err != nil {
		return "", fmt.Errorf("file1: %w", err)
	}

	data2, err := parser.Parse(filePath2)
	if err != nil {
		return "", fmt.Errorf("file2: %w", err)
	}

	raw := differ.Diff(data1, data2)

	formatted, err := formatter.Format(raw, format)
	if err != nil {
		return "", fmt.Errorf("format error: %w", err)
	}

	return formatted, nil
}
