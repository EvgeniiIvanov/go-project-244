package differ

import (
	"reflect"
)

type DiffNode struct {
	Key      string
	Status   string
	OldValue interface{}
	NewValue interface{}
	Children map[string]*DiffNode
}

func NewDiffNode(key, status string) *DiffNode {
	return &DiffNode{
		Key:      key,
		Status:   status,
		Children: make(map[string]*DiffNode),
	}
}

// NewLeafNode creates a leaf node (no children)
func NewLeafNode(key, status string, old, new interface{}) *DiffNode {
	return &DiffNode{
		Key:      key,
		Status:   status,
		OldValue: old,
		NewValue: new,
		Children: make(map[string]*DiffNode),
	}
}

func Diff(data1, data2 map[string]interface{}) *DiffNode {
	root := NewDiffNode("", "unchanged")

	// Collect all keys
	allKeys := make(map[string]struct{})
	for k := range data1 {
		allKeys[k] = struct{}{}
	}
	for k := range data2 {
		allKeys[k] = struct{}{}
	}

	// Process each key
	for k := range allKeys {
		v1, ok1 := data1[k]
		v2, ok2 := data2[k]
		root.Children[k] = processKey(k, v1, ok1, v2, ok2)
	}

	return root
}

// processKey handles the logic for a single key
func processKey(key string, v1 interface{}, ok1 bool, v2 interface{}, ok2 bool) *DiffNode {
	switch {
	case !ok1 && ok2:
		return handleAdded(key, v2)
	case ok1 && !ok2:
		return handleRemoved(key, v1)
	default:
		return handleBothExist(key, v1, v2)
	}
}

// handleAdded processes a key that exists only in the second file
func handleAdded(key string, value interface{}) *DiffNode {
	if m, isMap := toMap(value); isMap {
		node := NewDiffNode(key, "added")
		node.Children = buildTreeForAdded(m)
		return node
	}
	return NewLeafNode(key, "added", nil, value)
}

// handleRemoved processes a key that exists only in the first file
func handleRemoved(key string, value interface{}) *DiffNode {
	if m, isMap := toMap(value); isMap {
		node := NewDiffNode(key, "removed")
		node.Children = buildTreeForRemoved(m)
		return node
	}
	return NewLeafNode(key, "removed", value, nil)
}

// handleBothExist processes a key that exists in both files
func handleBothExist(key string, v1, v2 interface{}) *DiffNode {
	// Both are maps - recursive diff
	if areMaps(v1, v2) {
		node := NewDiffNode(key, "modified")
		node.Children = Diff(
			v1.(map[string]interface{}),
			v2.(map[string]interface{}),
		).Children
		return node
	}

	// Values are equal
	if reflect.DeepEqual(v1, v2) {
		if m, isMap := toMap(v1); isMap {
			node := NewDiffNode(key, "unchanged")
			node.Children = buildTreeForUnchanged(m)
			return node
		}
		return NewLeafNode(key, "unchanged", v1, v2)
	}

	// Values are different
	return NewLeafNode(key, "modified", v1, v2)
}

// buildTreeForAdded builds a tree for added nested object
func buildTreeForAdded(m map[string]interface{}) map[string]*DiffNode {
	children := make(map[string]*DiffNode)
	for k, v := range m {
		if nestedMap, isMap := toMap(v); isMap {
			node := NewDiffNode(k, "added")
			node.Children = buildTreeForAdded(nestedMap)
			children[k] = node
		} else {
			children[k] = NewLeafNode(k, "added", nil, v)
		}
	}
	return children
}

// buildTreeForRemoved builds a tree for removed nested object
func buildTreeForRemoved(m map[string]interface{}) map[string]*DiffNode {
	children := make(map[string]*DiffNode)
	for k, v := range m {
		if nestedMap, isMap := toMap(v); isMap {
			node := NewDiffNode(k, "removed")
			node.Children = buildTreeForRemoved(nestedMap)
			children[k] = node
		} else {
			children[k] = NewLeafNode(k, "removed", v, nil)
		}
	}
	return children
}

// buildTreeForUnchanged builds a tree for unchanged nested object
func buildTreeForUnchanged(m map[string]interface{}) map[string]*DiffNode {
	children := make(map[string]*DiffNode)
	for k, v := range m {
		if nestedMap, isMap := toMap(v); isMap {
			node := NewDiffNode(k, "unchanged")
			node.Children = buildTreeForUnchanged(nestedMap)
			children[k] = node
		} else {
			children[k] = NewLeafNode(k, "unchanged", v, v)
		}
	}
	return children
}

// toMap safely converts interface{} to map[string]interface{}
func toMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}

// areMaps checks if both values are maps
func areMaps(a, b interface{}) bool {
	_, aIsMap := toMap(a)
	_, bIsMap := toMap(b)
	return aIsMap && bIsMap
}
