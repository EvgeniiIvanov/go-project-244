package parser

import (
	"code/internal/models"
	"fmt"
	"path/filepath"
)

func Parse(filePath string) (map[string]interface{}, error) {
	ext := filepath.Ext(filePath)

	switch ext {
	case models.ExtJSON:
		return parseJSON(filePath)
	case models.ExtYAML, models.ExtYML:
		return parseYAML(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}
