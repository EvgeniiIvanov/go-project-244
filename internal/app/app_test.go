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
			name:  "nested structure",
			file1: "nested1",
			file2: "nested2",
			expected: []string{
				"    common: {",
				"      + follow: false",
				"        setting1: Value 1",
				"      - setting2: 200",
				"      - setting3: true",
				"      + setting3: null",
				"      + setting4: blah blah",
				"      + setting5: {",
				"            key5: value5",
				"        }",
				"        setting6: {",
				"            doge: {",
				"              - wow: ",
				"              + wow: so much",
				"            }",
				"            key: value",
				"          + ops: vops",
				"        }",
				"    }",
				"    group1: {",
				"      - baz: bas",
				"      + baz: bars",
				"        foo: bar",
				"      - nest: {",
				"            key: value",
				"        }",
				"      + nest: str",
				"    }",
				"  - group2: {",
				"        abc: 12345",
				"        deep: {",
				"            id: 45",
				"        }",
				"    }",
				"  + group3: {",
				"        deep: {",
				"            id: {",
				"                number: 45",
				"            }",
				"        }",
				"        fee: 100500",
				"    }",
			},
			notExpected: []string{},
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

func TestRunToStringErrors(t *testing.T) {
	t.Run("file1 does not exist", func(t *testing.T) {
		_, err := RunToString("nonexistent.json", "testdata/fixtures/json/flat1.json", "stylish")
		assert.Error(t, err)
	})

	t.Run("file2 does not exist", func(t *testing.T) {
		_, err := RunToString("testdata/fixtures/json/flat1.json", "nonexistent.json", "stylish")
		assert.Error(t, err)
	})

	t.Run("unsupported file format", func(t *testing.T) {
		_, err := RunToString("testdata/fixtures/json/flat1.txt", "testdata/fixtures/json/flat2.txt", "stylish")
		assert.Error(t, err)
	})
}

func TestRun(t *testing.T) {
	t.Run("successful run", func(t *testing.T) {
		file1 := filepath.Join("..", "..", "testdata", "fixtures", "json", "flat1.json")
		file2 := filepath.Join("..", "..", "testdata", "fixtures", "json", "flat2.json")

		// Run should not return an error for valid files
		err := Run(file1, file2, "stylish")
		assert.NoError(t, err)
	})

	t.Run("error with nonexistent file", func(t *testing.T) {
		err := Run("nonexistent.json", filepath.Join("..", "..", "testdata", "fixtures", "json", "flat1.json"), "stylish")
		assert.Error(t, err)
	})
}

func TestPlainFormat(t *testing.T) {
	formats := []struct {
		name string
		dir  string
		ext  string
	}{
		{name: "JSON", dir: "json", ext: ".json"},
		{name: "YAML", dir: "yaml", ext: ".yaml"},
	}

	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			file1 := filepath.Join("..", "..", "testdata", "fixtures", format.dir, "nested1"+format.ext)
			file2 := filepath.Join("..", "..", "testdata", "fixtures", format.dir, "nested2"+format.ext)

			result, err := RunToString(file1, file2, "plain")
			assert.NoError(t, err)

			expectedLines := []string{
				"Property 'common.follow' was added with value: false",
				"Property 'common.setting2' was removed",
				"Property 'common.setting3' was updated. From true to null",
				"Property 'common.setting4' was added with value: 'blah blah'",
				"Property 'common.setting5' was added with value: [complex value]",
				"Property 'common.setting6.doge.wow' was updated. From '' to 'so much'",
				"Property 'common.setting6.ops' was added with value: 'vops'",
				"Property 'group1.baz' was updated. From 'bas' to 'bars'",
				"Property 'group1.nest' was updated. From [complex value] to 'str'",
				"Property 'group2' was removed",
				"Property 'group3' was added with value: [complex value]",
			}

			// Check that all expected lines are present
			for _, line := range expectedLines {
				assert.Contains(t, result, line, "Missing expected line: %s", line)
			}
		})
	}
}
