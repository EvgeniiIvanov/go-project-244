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

	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.json")

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	return filePath
}

func TestParseJSON(t *testing.T) {
	tmpFile := createTempJSON(t, `{"key": "value"}`)

	result, err := Parse(tmpFile)

	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
}

func TestParseInvalidJSON(t *testing.T) {
	tmpFile := createTempFile(t, `{invalid json}`, ".json")

	_, err := Parse(tmpFile)
	assert.Error(t, err)
}

func TestParseUnsupportedFormat(t *testing.T) {
	_, err := Parse("file.yml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}
