package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func parseYAML(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	return result, nil
}
