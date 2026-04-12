package app

import (
	"code/internal/differ"
	"code/internal/formatter"
	"code/internal/parser"
	"fmt"
)

func RunToString(filePath1, filePath2 string, format string) (string, error) {
	data1, err := parser.Parse(filePath1)
	if err != nil {
		return "", fmt.Errorf("file1: %w", err)
	}

	data2, err := parser.Parse(filePath2)
	if err != nil {
		return "", fmt.Errorf("file2: %w", err)
	}

	raw, err := differ.Diff(data1, data2)
	if err != nil {
		return "", fmt.Errorf("diff error: %w", err)
	}

	formatted, err := formatter.Format(raw, format)
	if err != nil {
		return "", fmt.Errorf("format error: %w", err)
	}

	return formatted, nil
}

func Run(filePath1, filePath2 string, format string) error {
	result, err := RunToString(filePath1, filePath2, format)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
