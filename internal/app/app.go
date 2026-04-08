package app

import (
	"code/internal/parser"
	"fmt"
)

func Run(filePath1, filePath2 string) error {
	data1, err := parser.Parse(filePath1)
	if err != nil {
		return fmt.Errorf("file1: %w", err)
	}

	data2, err := parser.Parse(filePath2)
	if err != nil {
		return fmt.Errorf("file2: %w", err)
	}

	fmt.Println(data1, data2)
	return nil
}
