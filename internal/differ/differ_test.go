package differ

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name  string
		data1 map[string]interface{}
		data2 map[string]interface{}
		want  Differences
	}{
		{
			name:  "same values",
			data1: map[string]interface{}{"key": "value"},
			data2: map[string]interface{}{"key": "value"},
			want: Differences{
				Same:    map[string]interface{}{"key": "value"},
				Removed: map[string]interface{}{},
				Added:   map[string]interface{}{},
				Changed: map[string]Change{},
			},
		},
		{
			name:  "added key",
			data1: map[string]interface{}{"old": 1},
			data2: map[string]interface{}{"old": 1, "new": 2},
			want: Differences{
				Same:    map[string]interface{}{"old": 1},
				Removed: map[string]interface{}{},
				Added:   map[string]interface{}{"new": 2},
				Changed: map[string]Change{},
			},
		},
		{
			name:  "changed value",
			data1: map[string]interface{}{"key": 1},
			data2: map[string]interface{}{"key": 2},
			want: Differences{
				Same:    map[string]interface{}{},
				Removed: map[string]interface{}{},
				Added:   map[string]interface{}{},
				Changed: map[string]Change{"key": {Old: 1, New: 2}},
			},
		},
		{
			name:  "removed key",
			data1: map[string]interface{}{"old": 1, "willBeRemoved": 2},
			data2: map[string]interface{}{"old": 1},
			want: Differences{
				Same:    map[string]interface{}{"old": 1},
				Removed: map[string]interface{}{"willBeRemoved": 2},
				Added:   map[string]interface{}{},
				Changed: map[string]Change{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := Diff(tc.data1, tc.data2)
			assert.Equal(t, tc.want, got)
		})
	}
}
