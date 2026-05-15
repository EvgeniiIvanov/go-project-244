package parser

import (
	"code/internal/models"
	"fmt"
	"path/filepath"
)

func Parse(filePath string) (map[string]interface{}, error) {
	ext := filepath.Ext(filePath)

	var result map[string]interface{}
	var err error

	switch ext {
	case models.ExtJSON:
		result, err = parseJSON(filePath)
	case models.ExtYAML, models.ExtYML:
		result, err = parseYAML(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
