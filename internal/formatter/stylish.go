package formatter

import (
	"code/internal/differ"
	"fmt"
	"sort"
	"strings"
)

func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v)
}

// renderValue renders a value, either as a simple value or as a nested structure
// Used for modified nodes where one value is a map
func renderValue(v interface{}, indent string, prefix string) string {
	if m, ok := isMap(v); ok {
		// Value is a map, render as nested structure
		var result strings.Builder
		fmt.Fprintf(&result, "%s%s{\n", indent, prefix)

		// Get sorted keys
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Render each key-value pair with proper indentation
		childIndent := indent + "    "
		for _, k := range keys {
			result.WriteString(renderValue(m[k], childIndent, "    "+k+": "))
		}

		fmt.Fprintf(&result, "%s    }\n", indent)
		return result.String()
	}

	// Simple value
	return fmt.Sprintf("%s%s%s\n", indent, prefix, formatValue(v))
}

// formatLeaf formats a leaf node
func formatLeaf(node *differ.DiffNode, indent string, inheritFromParent bool) string {
	// If inheriting from parent, don't add status markers
	if inheritFromParent {
		value := node.OldValue
		if node.Status == "added" {
			value = node.NewValue
		}
		return fmt.Sprintf("%s    %s: %s\n", indent, node.Key, formatValue(value))
	}

	// Normal rendering with status markers
	switch node.Status {
	case "unchanged":
		return fmt.Sprintf("%s    %s: %s\n", indent, node.Key, formatValue(node.OldValue))
	case "removed":
		return fmt.Sprintf("%s  - %s: %s\n", indent, node.Key, formatValue(node.OldValue))
	case "added":
		return fmt.Sprintf("%s  + %s: %s\n", indent, node.Key, formatValue(node.NewValue))
	case "modified":
		// Modified leafs can contain map-to-primitive or primitive-to-map changes
		// Use renderValue for correct rendering of map values
		var result strings.Builder
		result.WriteString(renderValue(node.OldValue, indent, "  - "+node.Key+": "))
		result.WriteString(renderValue(node.NewValue, indent, "  + "+node.Key+": "))
		return result.String()
	}
	return ""
}

// formatNode formats a node (leaf or container)
func formatNode(node *differ.DiffNode, indent string, inheritFromParent bool) string {
	// For nodes with children
	if len(node.Children) > 0 {
		// Determine if children should inherit status (for added/removed/unchanged nodes)
		shouldInherit := node.Status == "added" || node.Status == "removed" || node.Status == "unchanged"

		// Determine prefix for opening brace
		openPrefix := "    "
		if !inheritFromParent {
			switch node.Status {
			case "removed":
				openPrefix = "  - "
			case "added":
				openPrefix = "  + "
			}
		}

		var result strings.Builder
		fmt.Fprintf(&result, "%s%s%s: {\n", indent, openPrefix, node.Key)
		result.WriteString(formatChildren(node, indent+"    ", shouldInherit))
		// Closing brace always uses "    " prefix (no status marker)
		fmt.Fprintf(&result, "%s    }\n", indent)
		return result.String()
	}

	// For leaf nodes
	return formatLeaf(node, indent, inheritFromParent)
}

// formatChildren formats all children of a node
func formatChildren(node *differ.DiffNode, indent string, inheritFromParent bool) string {
	var result strings.Builder
	for _, key := range sortedKeys(node.Children) {
		child := node.Children[key]
		result.WriteString(formatNode(child, indent, inheritFromParent))
	}
	return result.String()
}

func Stylish(node *differ.DiffNode, depth int) string {
	if node == nil || len(node.Children) == 0 {
		return "{\n}\n"
	}

	result := "{\n"
	result += formatChildren(node, "", false)
	result += "}\n"

	return result
}
