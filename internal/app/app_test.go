package app

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunToStringWithFixtures(t *testing.T) {
	formats := []struct {
		name string
		dir  string
		ext  string
	}{
		{name: "JSON", dir: "json", ext: ".json"},
		{name: "YAML", dir: "yaml", ext: ".yaml"},
	}

	testCases := []struct {
		name        string
		file1       string
		file2       string
		expected    []string
		notExpected []string
	}{
		{
			name:  "identical files",
			file1: "flat1",
			file2: "flat3",
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
			file1: "flat1",
			file2: "flat2",
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
			file1: "flat2",
			file2: "flat4",
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
			file1: "empty",
			file2: "flat1",
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
			file1: "flat1",
			file2: "empty",
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
			file1: "empty",
			file2: "empty",
			expected: []string{
				"{", "}",
			},
			notExpected: []string{"-", "+"},
		},
	}

	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			fixturesDir := filepath.Join("..", "..", "testdata", "fixtures", format.dir)

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					t.Helper()

					file1 := filepath.Join(fixturesDir, tc.file1+format.ext)
					file2 := filepath.Join(fixturesDir, tc.file2+format.ext)

					result, err := RunToString(file1, file2, "stylish")

					assert.NoError(t, err)

					for _, line := range tc.expected {
						assert.Contains(t, result, line, "Missing expected line: %s", line)
					}

					for _, line := range tc.notExpected {
						assert.NotContains(t, result, line, "Found unexpected line: %s", line)
					}
				})
			}
		})
	}
}
