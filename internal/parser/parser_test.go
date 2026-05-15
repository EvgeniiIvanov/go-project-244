package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, content string, ext string) string {
	t.Helper()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "test"+ext)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	return filePath
}

func createTempJSON(t *testing.T, content string) string {
	t.Helper()
	return createTempFile(t, content, ".json")
}

func createTempYAML(t *testing.T, content string) string {
	t.Helper()
	return createTempFile(t, content, ".yaml")
}

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]interface{}
	}{
		{
			name:    "simple object",
			content: `{"key": "value"}`,
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "with numbers",
			content: `{"age": 25, "score": 98.6}`,
			expected: map[string]interface{}{
				"age":   float64(25),
				"score": float64(98.6),
			},
		},
		{
			name:    "nested object",
			content: `{"user": {"name": "Alice", "age": 30}}`,
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  float64(30),
				},
			},
		},
		{
			name:     "empty object",
			content:  `{}`,
			expected: map[string]interface{}{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := createTempJSON(t, tc.content)

			result, err := Parse(tmpFile)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseYAML(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]interface{}
	}{
		{
			name:    "simple object",
			content: "key: value",
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "with numbers",
			content: "age: 25\nscore: 98.6",
			expected: map[string]interface{}{
				"age":   25,
				"score": 98.6,
			},
		},
		{
			name:    "nested object",
			content: "user:\n  name: Alice\n  age: 30",
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
		},
		{
			name:    "incomplete structure",
			content: "key:",
			expected: map[string]interface{}{
				"key": nil,
			},
		},
		{
			name:     "empty object",
			content:  "{}",
			expected: map[string]interface{}{},
		},
		{
			name:     "alternative empty syntax",
			content:  "",
			expected: map[string]interface{}{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := createTempYAML(t, tc.content)

			result, err := Parse(tmpFile)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseInvalidYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "invalid syntax",
			content: "key: [value",
		},
		{
			name:    "invalid indentation",
			content: "user:\n name: Alice\n   age: 30",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := createTempYAML(t, tc.content)

			_, err := Parse(tmpFile)
			assert.Error(t, err)
		})
	}
}

func TestParseUnsupportedFormat(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
	}{
		{
			name:     "txt file",
			filePath: "file.txt",
		},
		{
			name:     "xml file",
			filePath: "file.xml",
		},
		{
			name:     "no extension",
			filePath: "file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.filePath)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "unsupported")
		})
	}
}

func TestParseJSONNested(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]interface{}
	}{
		{
			name: "deeply nested object",
			content: `{
				"level1": {
					"level2": {
						"level3": {
							"value": "deep"
						}
					}
				}
			}`,
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": map[string]interface{}{
							"value": "deep",
						},
					},
				},
			},
		},
		{
			name: "mixed types in nested structure",
			content: `{
				"user": {
					"name": "John",
					"age": 30,
					"address": {
						"city": "NYC",
						"zip": 10001
					},
					"contacts": ["email", "phone"]
				}
			}`,
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  float64(30),
					"address": map[string]interface{}{
						"city": "NYC",
						"zip":  float64(10001),
					},
					"contacts": []interface{}{"email", "phone"},
				},
			},
		},
		{
			name: "nested with null values",
			content: `{
				"outer": {
					"inner": null,
					"deep": {
						"value": null
					}
				}
			}`,
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": nil,
					"deep": map[string]interface{}{
						"value": nil,
					},
				},
			},
		},
		{
			name: "array of nested objects",
			content: `{
				"users": [
					{"name": "Alice", "age": 25},
					{"name": "Bob", "age": 30}
				]
			}`,
			expected: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "Alice",
						"age":  float64(25),
					},
					map[string]interface{}{
						"name": "Bob",
						"age":  float64(30),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := createTempJSON(t, tc.content)
			result, err := Parse(tmpFile)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseYAMLNested(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]interface{}
	}{
		{
			name: "deeply nested object",
			content: `
level1:
  level2:
    level3:
      value: deep
`,
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": map[string]interface{}{
							"value": "deep",
						},
					},
				},
			},
		},
		{
			name: "mixed types in nested structure",
			content: `
user:
  name: John
  age: 30
  address:
    city: NYC
    zip: 10001
  contacts:
    - email
    - phone
`,
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
					"address": map[string]interface{}{
						"city": "NYC",
						"zip":  10001,
					},
					"contacts": []interface{}{"email", "phone"},
				},
			},
		},
		{
			name: "nested with null values",
			content: `
outer:
  inner: null
  deep:
    value: null
`,
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": nil,
					"deep": map[string]interface{}{
						"value": nil,
					},
				},
			},
		},
		{
			name: "array of nested objects",
			content: `
users:
  - name: Alice
    age: 25
  - name: Bob
    age: 30
`,
			expected: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name": "Alice",
						"age":  25,
					},
					map[string]interface{}{
						"name": "Bob",
						"age":  30,
					},
				},
			},
		},
		{
			name: "inline nested objects",
			content: `
level1: {level2: {level3: {value: deep}}}
`,
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": map[string]interface{}{
							"value": "deep",
						},
					},
				},
			},
		},
		{
			name: "mixed indentation",
			content: `
root:
  child1:
    value1: test
    child2:
      value2: nested
  child3: simple
`,
			expected: map[string]interface{}{
				"root": map[string]interface{}{
					"child1": map[string]interface{}{
						"value1": "test",
						"child2": map[string]interface{}{
							"value2": "nested",
						},
					},
					"child3": "simple",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := createTempYAML(t, tc.content)
			result, err := Parse(tmpFile)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseDeepNesting(t *testing.T) {
	t.Run("JSON deep nesting", func(t *testing.T) {
		// 20 layers of nested
		content := `{"level1": `
		for i := 2; i <= 20; i++ {
			content += `{"level` + fmt.Sprintf("%d", i) + `": `
		}
		content += `{"value": "deep"}`
		for i := 1; i <= 20; i++ {
			content += `}`
		}

		tmpFile := createTempJSON(t, content)
		_, err := Parse(tmpFile)

		assert.NoError(t, err)
	})

	t.Run("YAML deep nesting", func(t *testing.T) {
		content := ""
		for i := 1; i <= 20; i++ {
			indent := ""
			for j := 0; j < i-1; j++ {
				indent += "  "
			}
			content += indent + "level" + fmt.Sprintf("%d", i) + ":\n"
		}
		indent := ""
		for j := 0; j < 20; j++ {
			indent += "  "
		}
		content += indent + "value: deep"

		tmpFile := createTempYAML(t, content)
		_, err := Parse(tmpFile)

		assert.NoError(t, err)
	})
}

func TestParseComplexNested(t *testing.T) {
	t.Run("JSON complex nested", func(t *testing.T) {
		content := `{
			"config": {
				"database": {
					"host": "localhost",
					"port": 5432,
					"credentials": {
						"username": "admin",
						"password": "secret"
					}
				},
				"features": {
					"logging": true,
					"cache": {
						"enabled": true,
						"ttl": 3600
					}
				}
			}
		}`

		expected := map[string]interface{}{
			"config": map[string]interface{}{
				"database": map[string]interface{}{
					"host": "localhost",
					"port": float64(5432),
					"credentials": map[string]interface{}{
						"username": "admin",
						"password": "secret",
					},
				},
				"features": map[string]interface{}{
					"logging": true,
					"cache": map[string]interface{}{
						"enabled": true,
						"ttl":     float64(3600),
					},
				},
			},
		}

		tmpFile := createTempJSON(t, content)
		result, err := Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("YAML complex nested", func(t *testing.T) {
		content := `
config:
  database:
    host: localhost
    port: 5432
    credentials:
      username: admin
      password: secret
  features:
    logging: true
    cache:
      enabled: true
      ttl: 3600
`
		expected := map[string]interface{}{
			"config": map[string]interface{}{
				"database": map[string]interface{}{
					"host": "localhost",
					"port": 5432,
					"credentials": map[string]interface{}{
						"username": "admin",
						"password": "secret",
					},
				},
				"features": map[string]interface{}{
					"logging": true,
					"cache": map[string]interface{}{
						"enabled": true,
						"ttl":     3600,
					},
				},
			},
		}

		tmpFile := createTempYAML(t, content)
		result, err := Parse(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestParseEmptyNested(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		content  string
		expected map[string]interface{}
	}{
		{
			name:   "JSON empty nested object",
			format: "json",
			content: `{
				"outer": {
					"inner": {}
				}
			}`,
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": map[string]interface{}{},
				},
			},
		},
		{
			name:   "YAML empty nested object",
			format: "yaml",
			content: `
outer:
  inner: {}
`,
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": map[string]interface{}{},
				},
			},
		},
		{
			name:   "YAML empty nested object alternative",
			format: "yaml",
			content: `
outer:
  inner:
`,
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": nil,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var tmpFile string
			if tc.format == "json" {
				tmpFile = createTempJSON(t, tc.content)
			} else {
				tmpFile = createTempYAML(t, tc.content)
			}

			result, err := Parse(tmpFile)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
