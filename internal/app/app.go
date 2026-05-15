package app

import (
	"code/internal/differ"
	"code/internal/formatter"
	"code/internal/parser"
	"fmt"
)

func Run(filePath1, filePath2 string, format string) (string, error) {
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
