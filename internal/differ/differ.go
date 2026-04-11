package differ

import (
	"fmt"
	"reflect"
)

type Change struct {
	Old interface{}
	New interface{}
}

type Differences struct {
	Same    map[string]interface{}
	Changed map[string]Change
	Removed map[string]interface{}
	Added   map[string]interface{}
}

func NewDifferences() Differences {
	return Differences{
		Same:    make(map[string]interface{}),
		Removed: make(map[string]interface{}),
		Added:   make(map[string]interface{}),
		Changed: make(map[string]Change),
	}
}

func Diff(data1, data2 map[string]interface{}) (Differences, error) {
	if data1 == nil || data2 == nil {
		return NewDifferences(), fmt.Errorf("data maps cannot be nil")
	}

	result := NewDifferences()
	for k, v1 := range data1 {
		v2, ok := data2[k]

		if !ok {
			result.Removed[k] = v1
			continue
		}

		if reflect.DeepEqual(v1, v2) {
			result.Same[k] = v1
		} else {
			result.Changed[k] = Change{v1, v2}
		}
	}

	for k, v2 := range data2 {
		if _, ok := data1[k]; !ok {
			result.Added[k] = v2
		}
	}

	return result, nil
}
