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

	t.Run("plain format", func(t *testing.T) {
		result, err := Format(root, "plain")
		assert.NoError(t, err)
		assert.Contains(t, result, "Property 'key' was updated. From 'old' to 'new'")
	})

	t.Run("unknown format", func(t *testing.T) {
		_, err := Format(root, "unknown")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown format")
	})
}

func TestIsComplexValue(t *testing.T) {
	// Complex values
	assert.True(t, isComplexValue(map[string]interface{}{"key": "value"}))
	assert.True(t, isComplexValue([]interface{}{1, 2, 3}))
	assert.True(t, isComplexValue([]string{"a", "b"}))
	assert.True(t, isComplexValue([]int{1, 2}))
	assert.True(t, isComplexValue([]float64{1.0, 2.0}))

	// Simple values
	assert.False(t, isComplexValue(nil))
	assert.False(t, isComplexValue("string"))
	assert.False(t, isComplexValue(42))
	assert.False(t, isComplexValue(3.14))
	assert.False(t, isComplexValue(true))
}

func TestIsMap(t *testing.T) {
	// Is a map
	m := map[string]interface{}{"key": "value"}
	result, ok := isMap(m)
	assert.True(t, ok)
	assert.Equal(t, m, result)

	// Not a map
	_, ok = isMap("string")
	assert.False(t, ok)

	_, ok = isMap(42)
	assert.False(t, ok)

	_, ok = isMap(nil)
	assert.False(t, ok)
}

func TestPlain(t *testing.T) {
	t.Run("simple changes", func(t *testing.T) {
		root := differ.NewDiffNode("", "unchanged")

		addedNode := differ.NewDiffNode("added", "added")
		addedNode.NewValue = "value"
		root.Children["added"] = addedNode

		removedNode := differ.NewDiffNode("removed", "removed")
		removedNode.OldValue = 42
		root.Children["removed"] = removedNode

		modifiedNode := differ.NewDiffNode("modified", "modified")
		modifiedNode.OldValue = false
		modifiedNode.NewValue = true
		root.Children["modified"] = modifiedNode

		result := Plain(root, 0)

		assert.Contains(t, result, "Property 'added' was added with value: 'value'")
		assert.Contains(t, result, "Property 'removed' was removed")
		assert.Contains(t, result, "Property 'modified' was updated. From false to true")
	})

	t.Run("nested changes", func(t *testing.T) {
		root := differ.NewDiffNode("", "unchanged")

		parentNode := differ.NewDiffNode("parent", "modified")
		root.Children["parent"] = parentNode

		childNode := differ.NewDiffNode("child", "modified")
		childNode.OldValue = "old"
		childNode.NewValue = "new"
		parentNode.Children["child"] = childNode

		result := Plain(root, 0)

		assert.Contains(t, result, "Property 'parent.child' was updated. From 'old' to 'new'")
	})

	t.Run("complex values", func(t *testing.T) {
		root := differ.NewDiffNode("", "unchanged")

		complexNode := differ.NewDiffNode("config", "added")
		complexNode.NewValue = map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		}
		// Complex values have children
		complexNode.Children = map[string]*differ.DiffNode{
			"host": {Key: "host", Status: "added", NewValue: "localhost"},
		}
		root.Children["config"] = complexNode

		result := Plain(root, 0)

		assert.Contains(t, result, "Property 'config' was added with value: [complex value]")
	})

	t.Run("skips unchanged", func(t *testing.T) {
		root := differ.NewDiffNode("", "unchanged")

		unchangedNode := differ.NewDiffNode("unchanged", "unchanged")
		unchangedNode.OldValue = "value"
		unchangedNode.NewValue = "value"
		root.Children["unchanged"] = unchangedNode

		result := Plain(root, 0)

		// Unchanged values should not appear in plain format
		assert.NotContains(t, result, "unchanged")
	})
}
