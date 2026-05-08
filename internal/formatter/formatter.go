package formatter

import (
	"code/internal/differ"
	"fmt"
)

func Format(raw *differ.DiffNode, format string) (string, error) {
	switch format {
	case "stylish":
		return Stylish(raw, 0), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
