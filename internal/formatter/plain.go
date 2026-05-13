package formatter

import (
	"code/internal/differ"
	"fmt"
	"sort"
	"strings"
)

// Plain formats the diff in plain text format
func Plain(node *differ.DiffNode, depth int) string {
	if node == nil {
		return ""
	}

	var lines []string
	collectPlainLines(node, "", &lines)

	// Sort lines for consistent output
	sort.Strings(lines)

	return strings.Join(lines, "\n") + "\n"
}

// collectPlainLines recursively collects plain format lines
func collectPlainLines(node *differ.DiffNode, path string, lines *[]string) {
	keys := sortedKeys(node.Children)

	for _, key := range keys {
		child := node.Children[key]
		currentPath := buildPath(path, child.Key)

		switch child.Status {
		case "added":
			if len(child.Children) > 0 {
				// Added complex value (nested object)
				*lines = append(*lines, fmt.Sprintf("Property '%s' was added with value: [complex value]", currentPath))
			} else {
				// Added simple value
				*lines = append(*lines, fmt.Sprintf("Property '%s' was added with value: %s", currentPath, formatPlainValue(child.NewValue)))
			}

		case "removed":
			*lines = append(*lines, fmt.Sprintf("Property '%s' was removed", currentPath))

		case "modified":
			if len(child.Children) > 0 {
				// Modified nested object - recurse into children
				collectPlainLines(child, currentPath, lines)
			} else {
				// Modified leaf value
				*lines = append(*lines, fmt.Sprintf("Property '%s' was updated. From %s to %s",
					currentPath,
					formatPlainValue(child.OldValue),
					formatPlainValue(child.NewValue)))
			}

		case "unchanged":
			// Skip unchanged values in plain format
			if len(child.Children) > 0 {
				// But recurse into children to find changes
				collectPlainLines(child, currentPath, lines)
			}
		}
	}
}

// buildPath constructs the property path
func buildPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

// formatPlainValue formats a value for plain output
func formatPlainValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	// Check if it's a complex value (map or slice)
	if isComplexValue(v) {
		return "[complex value]"
	}

	// String values should be quoted
	if str, ok := v.(string); ok {
		return fmt.Sprintf("'%s'", str)
	}

	// Numbers, booleans, etc.
	return fmt.Sprintf("%v", v)
}
