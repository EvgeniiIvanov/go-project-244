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
		check func(*testing.T, *DiffNode)
	}{
		{
			name:  "same values",
			data1: map[string]interface{}{"key": "value"},
			data2: map[string]interface{}{"key": "value"},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				assert.Len(t, node.Children, 1)
				assert.Equal(t, "unchanged", node.Children["key"].Status)
				assert.Equal(t, "value", node.Children["key"].OldValue)
				assert.Equal(t, "value", node.Children["key"].NewValue)
			},
		},
		{
			name:  "added key",
			data1: map[string]interface{}{"old": 1},
			data2: map[string]interface{}{"old": 1, "new": 2},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				assert.Len(t, node.Children, 2)
				assert.Equal(t, "unchanged", node.Children["old"].Status)
				assert.Equal(t, "added", node.Children["new"].Status)
				assert.Equal(t, 2, node.Children["new"].NewValue)
			},
		},
		{
			name:  "changed value",
			data1: map[string]interface{}{"key": 1},
			data2: map[string]interface{}{"key": 2},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				assert.Len(t, node.Children, 1)
				assert.Equal(t, "modified", node.Children["key"].Status)
				assert.Equal(t, 1, node.Children["key"].OldValue)
				assert.Equal(t, 2, node.Children["key"].NewValue)
			},
		},
		{
			name:  "removed key",
			data1: map[string]interface{}{"old": 1, "willBeRemoved": 2},
			data2: map[string]interface{}{"old": 1},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				assert.Len(t, node.Children, 2)
				assert.Equal(t, "unchanged", node.Children["old"].Status)
				assert.Equal(t, "removed", node.Children["willBeRemoved"].Status)
				assert.Equal(t, 2, node.Children["willBeRemoved"].OldValue)
			},
		},
		{
			name: "nested objects unchanged",
			data1: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			data2: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				assert.Len(t, node.Children, 1)
				userNode := node.Children["user"]
				assert.Equal(t, "modified", userNode.Status)
				assert.Len(t, userNode.Children, 2)
				assert.Equal(t, "unchanged", userNode.Children["name"].Status)
				assert.Equal(t, "unchanged", userNode.Children["age"].Status)
			},
		},
		{
			name: "nested objects with changes",
			data1: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			data2: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  31,
				},
			},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				assert.Len(t, node.Children, 1)
				userNode := node.Children["user"]
				assert.Equal(t, "modified", userNode.Status)
				assert.Len(t, userNode.Children, 2)
				assert.Equal(t, "unchanged", userNode.Children["name"].Status)
				assert.Equal(t, "modified", userNode.Children["age"].Status)
				assert.Equal(t, 30, userNode.Children["age"].OldValue)
				assert.Equal(t, 31, userNode.Children["age"].NewValue)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Diff(tc.data1, tc.data2)
			tc.check(t, got)
		})
	}
}
