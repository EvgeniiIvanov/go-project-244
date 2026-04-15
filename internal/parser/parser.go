package parser

import (
	"fmt"
	"path/filepath"
)

func Parse(filePath string) (map[string]interface{}, error) {
	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		return parseJSON(filePath)
	case ".yaml", ".yml":
		return parseYAML(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}
