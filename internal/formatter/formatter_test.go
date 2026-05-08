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
