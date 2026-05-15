package differ

import (
	"code/internal/models"
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
	root := NewDiffNode("", models.StatusUnchanged)

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
		node := NewDiffNode(key, models.StatusAdded)
		node.Children = buildFullTree(m, models.StatusAdded)
		return node
	}
	return NewLeafNode(key, models.StatusAdded, nil, value)
}

// handleRemoved processes a key that exists only in the first file
func handleRemoved(key string, value interface{}) *DiffNode {
	if m, isMap := toMap(value); isMap {
		node := NewDiffNode(key, models.StatusRemoved)
		node.Children = buildFullTree(m, models.StatusRemoved)
		return node
	}
	return NewLeafNode(key, models.StatusRemoved, value, nil)
}

// handleBothExist processes a key that exists in both files
func handleBothExist(key string, v1, v2 interface{}) *DiffNode {
	// Both are maps - recursive diff
	if areMaps(v1, v2) {
		node := NewDiffNode(key, models.StatusModified)
		node.Children = Diff(
			v1.(map[string]interface{}),
			v2.(map[string]interface{}),
		).Children
		return node
	}

	// Values are equal
	if reflect.DeepEqual(v1, v2) {
		if m, isMap := toMap(v1); isMap {
			node := NewDiffNode(key, models.StatusUnchanged)
			node.Children = buildFullTree(m, models.StatusUnchanged)
			return node
		}
		return NewLeafNode(key, models.StatusUnchanged, v1, v2)
	}

	// Values are different
	return NewLeafNode(key, models.StatusModified, v1, v2)
}

// buildFullTree builds a tree for a nested object with the given status
func buildFullTree(m map[string]interface{}, status string) map[string]*DiffNode {
	children := make(map[string]*DiffNode)
	for k, v := range m {
		if nestedMap, isMap := toMap(v); isMap {
			node := NewDiffNode(k, status)
			node.Children = buildFullTree(nestedMap, status)
			children[k] = node
		} else {
			// Determine old/new values based on status
			var old, new interface{}
			switch status {
			case models.StatusAdded:
				old, new = nil, v
			case models.StatusRemoved:
				old, new = v, nil
			case models.StatusUnchanged:
				old, new = v, v
			}
			children[k] = NewLeafNode(k, status, old, new)
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
