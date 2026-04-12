package formatter

import (
	"code/internal/differ"
	"testing"

	"github.com/stretchr/testify/assert"
)

// internal/formatter/stylish_test.go
func TestStylish(t *testing.T) {
	diff := differ.Differences{
		Same:    map[string]interface{}{"host": "hexlet.io"},
		Removed: map[string]interface{}{"timeout": 50},
		Added:   map[string]interface{}{"verbose": true},
		Changed: map[string]differ.Change{
			"port": {Old: 8080, New: 80},
		},
	}

	expected := `{
    host: hexlet.io
  - port: 8080
  + port: 80
  - timeout: 50
  + verbose: true
}
`
	result := Stylish(diff)
	assert.Equal(t, expected, result)
}

func TestStylishEmpty(t *testing.T) {
	diff := differ.NewDifferences()
	result := Stylish(diff)
	assert.Equal(t, "{\n}\n", result)
}
