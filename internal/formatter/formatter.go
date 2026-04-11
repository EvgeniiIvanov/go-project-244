package formatter

import (
	"code/internal/differ"
	"fmt"
)

func Format(raw differ.Differences, format string) (string, error) {
	switch format {
	case "stylish":
		return Stylish(raw), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
