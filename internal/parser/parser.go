package parser

import (
	"fmt"
	"path/filepath"
)

func Parse(filePath string) (map[string]interface{}, error) {
	ext := filepath.Ext(filePath)

	var result map[string]interface{}
	var err error

	switch ext {
	case ".json":
		result, err = parseJSON(filePath)
	case ".yaml", ".yml":
		result, err = parseYAML(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
