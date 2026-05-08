package parser

import (
	"fmt"
	"path/filepath"
)

func normalizeNumbers(v interface{}) interface{} {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int8, int16, int32, int64:
		return float64(val.(int64))
	case uint, uint8, uint16, uint32, uint64:
		return float64(val.(uint64))
	case float32:
		return float64(val)
	case map[string]interface{}:
		for k, v := range val {
			val[k] = normalizeNumbers(v)
		}
		return val
	case []interface{}:
		for i, v := range val {
			val[i] = normalizeNumbers(v)
		}
		return val
	default:
		return v
	}
}

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

	for k, v := range result {
		result[k] = normalizeNumbers(v)
	}

	return result, nil
}
