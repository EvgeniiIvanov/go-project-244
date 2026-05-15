package formatter

import (
	"code/internal/differ"
	"code/internal/models"
	"fmt"
	"sort"
)

func Format(raw *differ.DiffNode, format string) (string, error) {
	switch format {
	case models.OutputFormatStylish:
		return Stylish(raw, 0), nil
	case models.OutputFormatPlain:
		return Plain(raw, 0), nil
	case models.OutputFormatJSON:
		return Json(raw), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

// sortedKeys returns sorted keys from a map of DiffNodes
// Used by multiple formatters to ensure consistent output
func sortedKeys(children map[string]*differ.DiffNode) []string {
	keys := make([]string, 0, len(children))
	for k := range children {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// isMap checks if a value is a map and returns it
func isMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}

// isComplexValue checks if a value is a complex type (map or slice)
func isComplexValue(v interface{}) bool {
	if v == nil {
		return false
	}

	// Check for maps
	if _, ok := v.(map[string]interface{}); ok {
		return true
	}

	// Check for slices
	switch v.(type) {
	case []interface{}, []string, []int, []float64:
		return true
	}

	return false
}
