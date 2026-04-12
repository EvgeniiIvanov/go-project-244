package formatter

import (
	"code/internal/differ"
	"fmt"
	"sort"
	"strings"
)

func collectSortedKeys(raw differ.Differences) []string {
	allKeys := make(map[string]struct{})

	for k := range raw.Same {
		allKeys[k] = struct{}{}
	}
	for k := range raw.Removed {
		allKeys[k] = struct{}{}
	}
	for k := range raw.Added {
		allKeys[k] = struct{}{}
	}
	for k := range raw.Changed {
		allKeys[k] = struct{}{}
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func Stylish(raw differ.Differences) string {
	keys := collectSortedKeys(raw)

	// make formatted string
	var b strings.Builder
	b.WriteString("{\n")
	for _, k := range keys {
		if v, ok := raw.Removed[k]; ok {
			fmt.Fprintf(&b, "  - %s: %v\n", k, v)
			continue
		}
		if v, ok := raw.Added[k]; ok {
			fmt.Fprintf(&b, "  + %s: %v\n", k, v)
			continue
		}
		if v, ok := raw.Changed[k]; ok {
			fmt.Fprintf(&b, "  - %s: %v\n", k, v.Old)
			fmt.Fprintf(&b, "  + %s: %v\n", k, v.New)
			continue
		}
		if v, ok := raw.Same[k]; ok {
			fmt.Fprintf(&b, "    %s: %v\n", k, v)
		}
	}
	b.WriteString("}\n")

	return b.String()
}
