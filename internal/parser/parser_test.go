package parser

import (
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
		{
			name:    "incomplete structure",
			content: "key:",
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
