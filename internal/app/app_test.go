package app

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunToStringWithFixtures(t *testing.T) {
	fixturesDir := filepath.Join("..", "..", "testdata", "fixtures")

	tests := []struct {
		name        string
		file1       string
		file2       string
		expected    []string
		notExpected []string
	}{
		{
			name:  "identical files",
			file1: "flat1.json",
			file2: "flat3.json",
			expected: []string{
				"    host: hexlet.io",
				"    timeout: 50",
				"    proxy: 123.234.53.22",
				"    follow: false",
			},
			notExpected: []string{
				"  -", "  +",
			},
		},
		{
			name:  "different structure",
			file1: "flat1.json",
			file2: "flat2.json",
			expected: []string{
				"  - follow: false",
				"    host: hexlet.io",
				"  - proxy: 123.234.53.22",
				"  - timeout: 50",
				"  + timeout: 20",
				"  + verbose: true",
			},
			notExpected: []string{},
		},
		{
			name:  "different values same keys",
			file1: "flat2.json",
			file2: "flat4.json",
			expected: []string{
				"  - host: hexlet.io",
				"  + host: coursera.org",
				"  - timeout: 20",
				"  + timeout: 30",
				"  - verbose: true",
				"  + verbose: false",
			},
		},
		{
			name:  "first file empty",
			file1: "empty.json",
			file2: "flat1.json",
			expected: []string{
				"  + follow: false",
				"  + host: hexlet.io",
				"  + proxy: 123.234.53.22",
				"  + timeout: 50",
			},
			notExpected: []string{"  - "},
		},
		{
			name:  "second file empty",
			file1: "flat1.json",
			file2: "empty.json",
			expected: []string{
				"  - follow: false",
				"  - host: hexlet.io",
				"  - proxy: 123.234.53.22",
				"  - timeout: 50",
			},
			notExpected: []string{"  + "},
		},
		{
			name:  "both empty",
			file1: "empty.json",
			file2: "empty.json",
			expected: []string{
				"{", "}",
			},
			notExpected: []string{"-", "+"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filePath1 := filepath.Join(fixturesDir, tc.file1)
			filePath2 := filepath.Join(fixturesDir, tc.file2)

			result, err := RunToString(filePath1, filePath2, "stylish")

			assert.NoError(t, err)

			for _, expected := range tc.expected {
				assert.Contains(t, result, expected,
					"Expected to contain: %s", expected)
			}

			for _, notExpected := range tc.notExpected {
				assert.NotContains(t, result, notExpected,
					"Should not contain: %s", notExpected)
			}
		})
	}
}
