package formatter

import (
	"code/internal/differ"
	"encoding/json"
)

type JSONDiffNode struct {
	Status   string                   `json:"status,omitempty"`
	OldValue interface{}              `json:"oldValue,omitempty"`
	NewValue interface{}              `json:"newValue,omitempty"`
	Children map[string]*JSONDiffNode `json:"children,omitempty"`
}

func Json(node *differ.DiffNode) string {
	result := convertToJSONFormat(node)
	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data)
}

func convertToJSONFormat(node *differ.DiffNode) map[string]*JSONDiffNode {
	result := make(map[string]*JSONDiffNode)

	for key, child := range node.Children {
		jsonNode := &JSONDiffNode{
			Status: child.Status,
		}

		// Add values for leaf nodes (no children)
		if len(child.Children) == 0 {
			if child.Status == "removed" || child.Status == "modified" || child.Status == "unchanged" {
				jsonNode.OldValue = child.OldValue
			}
			if child.Status == "added" || child.Status == "modified" {
				jsonNode.NewValue = child.NewValue
			}
		} else {
			// Recursively convert children
			jsonNode.Children = convertToJSONFormat(child)
		}

		result[key] = jsonNode
	}

	return result
}
