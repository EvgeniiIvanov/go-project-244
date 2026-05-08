package formatter

import (
	"code/internal/differ"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStylish(t *testing.T) {
	// Create a DiffNode tree representing the differences
	root := differ.NewDiffNode("", "unchanged")

	hostNode := differ.NewDiffNode("host", "unchanged")
	hostNode.OldValue = "hexlet.io"
	hostNode.NewValue = "hexlet.io"
	root.Children["host"] = hostNode

	portNode := differ.NewDiffNode("port", "modified")
	portNode.OldValue = 8080
	portNode.NewValue = 80
	root.Children["port"] = portNode

	timeoutNode := differ.NewDiffNode("timeout", "removed")
	timeoutNode.OldValue = 50
	root.Children["timeout"] = timeoutNode

	verboseNode := differ.NewDiffNode("verbose", "added")
	verboseNode.NewValue = true
	root.Children["verbose"] = verboseNode

	expected := `{
    host: hexlet.io
  - port: 8080
  + port: 80
  - timeout: 50
  + verbose: true
}
`
	result := Stylish(root, 0)
	assert.Equal(t, expected, result)
}

func TestStylishEmpty(t *testing.T) {
	root := differ.NewDiffNode("", "unchanged")
	result := Stylish(root, 0)
	assert.Equal(t, "{\n}\n", result)
}

func TestStylishNested(t *testing.T) {
	// Create a nested structure
	root := differ.NewDiffNode("", "unchanged")

	commonNode := differ.NewDiffNode("common", "modified")
	root.Children["common"] = commonNode

	setting1Node := differ.NewDiffNode("setting1", "unchanged")
	setting1Node.OldValue = "Value 1"
	setting1Node.NewValue = "Value 1"
	commonNode.Children["setting1"] = setting1Node

	setting2Node := differ.NewDiffNode("setting2", "removed")
	setting2Node.OldValue = 200
	commonNode.Children["setting2"] = setting2Node

	setting3Node := differ.NewDiffNode("setting3", "modified")
	setting3Node.OldValue = true
	setting3Node.NewValue = nil
	commonNode.Children["setting3"] = setting3Node

	result := Stylish(root, 0)

	// Check that nested structure is properly formatted
	assert.Contains(t, result, "    common: {")
	assert.Contains(t, result, "        setting1: Value 1")
	assert.Contains(t, result, "      - setting2: 200")
	assert.Contains(t, result, "      - setting3: true")
	assert.Contains(t, result, "      + setting3: null")
	assert.Contains(t, result, "    }")
}

func TestStylishMapToPrimitive(t *testing.T) {
	// Test case: value changes from map to primitive
	root := differ.NewDiffNode("", "unchanged")

	nestNode := differ.NewDiffNode("nest", "modified")
	nestNode.OldValue = map[string]interface{}{
		"key": "value",
	}
	nestNode.NewValue = "str"
	root.Children["nest"] = nestNode

	result := Stylish(root, 0)

	// Old value should be expanded as a map
	assert.Contains(t, result, "  - nest: {")
	assert.Contains(t, result, "        key: value")
	assert.Contains(t, result, "    }")
	// New value should be a simple string
	assert.Contains(t, result, "  + nest: str")
}

func TestStylishPrimitiveToMap(t *testing.T) {
	// Test case: value changes from primitive to map
	root := differ.NewDiffNode("", "unchanged")

	configNode := differ.NewDiffNode("config", "modified")
	configNode.OldValue = "simple"
	configNode.NewValue = map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	}
	root.Children["config"] = configNode

	result := Stylish(root, 0)

	// Old value should be simple
	assert.Contains(t, result, "  - config: simple")
	// New value should be expanded as a map
	assert.Contains(t, result, "  + config: {")
	assert.Contains(t, result, "        host: localhost")
	assert.Contains(t, result, "        port: 8080")
	assert.Contains(t, result, "    }")
}

func TestStylishNestedMapInModified(t *testing.T) {
	// Test case: deeply nested map in a modified value
	root := differ.NewDiffNode("", "unchanged")

	serverNode := differ.NewDiffNode("server", "modified")
	serverNode.OldValue = map[string]interface{}{
		"db": map[string]interface{}{
			"host": "old.example.com",
			"port": 5432,
		},
	}
	serverNode.NewValue = map[string]interface{}{
		"db": map[string]interface{}{
			"host": "new.example.com",
			"port": 5433,
		},
	}
	root.Children["server"] = serverNode

	result := Stylish(root, 0)

	// Both old and new should expand nested maps
	assert.Contains(t, result, "  - server: {")
	assert.Contains(t, result, "        db: {")
	assert.Contains(t, result, "            host: old.example.com")
	assert.Contains(t, result, "            port: 5432")

	assert.Contains(t, result, "  + server: {")
	assert.Contains(t, result, "        db: {")
	assert.Contains(t, result, "            host: new.example.com")
	assert.Contains(t, result, "            port: 5433")
}

func TestFormat(t *testing.T) {
	root := differ.NewDiffNode("", "unchanged")

	keyNode := differ.NewDiffNode("key", "modified")
	keyNode.OldValue = "old"
	keyNode.NewValue = "new"
	root.Children["key"] = keyNode

	t.Run("stylish format", func(t *testing.T) {
		result, err := Format(root, "stylish")
		assert.NoError(t, err)
		assert.Contains(t, result, "  - key: old")
		assert.Contains(t, result, "  + key: new")
	})

	t.Run("unknown format", func(t *testing.T) {
		_, err := Format(root, "unknown")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown format")
	})
}
