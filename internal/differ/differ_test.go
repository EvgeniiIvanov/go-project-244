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
		{
			name: "added nested object",
			data1: map[string]interface{}{
				"key": "value",
			},
			data2: map[string]interface{}{
				"key": "value",
				"config": map[string]interface{}{
					"host": "localhost",
					"port": 8080,
				},
			},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				configNode := node.Children["config"]
				assert.Equal(t, "added", configNode.Status)
				assert.Len(t, configNode.Children, 2)
				assert.Equal(t, "added", configNode.Children["host"].Status)
				assert.Equal(t, "localhost", configNode.Children["host"].NewValue)
				assert.Equal(t, "added", configNode.Children["port"].Status)
				assert.Equal(t, 8080, configNode.Children["port"].NewValue)
			},
		},
		{
			name: "removed nested object",
			data1: map[string]interface{}{
				"cache": map[string]interface{}{
					"enabled": true,
					"ttl":     300,
				},
			},
			data2: map[string]interface{}{},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				cacheNode := node.Children["cache"]
				assert.Equal(t, "removed", cacheNode.Status)
				assert.Len(t, cacheNode.Children, 2)
				assert.Equal(t, "removed", cacheNode.Children["enabled"].Status)
				assert.Equal(t, true, cacheNode.Children["enabled"].OldValue)
				assert.Equal(t, "removed", cacheNode.Children["ttl"].Status)
				assert.Equal(t, 300, cacheNode.Children["ttl"].OldValue)
			},
		},
		{
			name: "unchanged nested object",
			data1: map[string]interface{}{
				"db": map[string]interface{}{
					"host": "localhost",
					"port": 5432,
				},
			},
			data2: map[string]interface{}{
				"db": map[string]interface{}{
					"host": "localhost",
					"port": 5432,
				},
			},
			check: func(t *testing.T, node *DiffNode) {
				assert.NotNil(t, node)
				dbNode := node.Children["db"]
				// When comparing nested maps, parent is marked as "modified" for processing
				// but all children are "unchanged"
				assert.Equal(t, "modified", dbNode.Status)
				assert.Len(t, dbNode.Children, 2)
				assert.Equal(t, "unchanged", dbNode.Children["host"].Status)
				assert.Equal(t, "localhost", dbNode.Children["host"].OldValue)
				assert.Equal(t, "unchanged", dbNode.Children["port"].Status)
				assert.Equal(t, 5432, dbNode.Children["port"].OldValue)
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
