package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func parseYAML(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid YAML: %s", err)
	}

	// normalize resulting map
	if result == nil {
		result = map[string]interface{}{}
	}

	fmt.Println(result)

	return result, nil
}
