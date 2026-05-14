### Quality & Testing

| Status | Badge |
|--------|-------|
| CI Pipeline | [![CI](https://github.com/EvgeniiIvanov/go-project-244/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/EvgeniiIvanov/go-project-244/actions/workflows/ci.yml) |
| Test Coverage | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=EvgeniiIvanov_go-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=EvgeniiIvanov_go-project-244) |

### Hexlet tests and linter status:
[![Actions Status](https://github.com/EvgeniiIvanov/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/EvgeniiIvanov/go-project-244/actions)

## Description

A command-line tool that compares two configuration files and shows the differences, similar to `diff`. Supports JSON and YAML formats, with multiple output formats (stylish, plain, and json) that clearly indicate added, removed, and modified values in nested structures.

## Usage

### How to install

```bash
make build
sudo cp bin/gendiff /usr/local/bin/
```

### How to use

```bash
gendiff [--format <format>] <file1> <file2>
```

**Options:**
- `--format, -f` - Output format (default: "stylish")
  - `stylish` - Tree-like format with symbols showing changes
  - `plain` - Text format with property paths
  - `json` - Machine-readable JSON format with status and values

**Supported file formats:**
- JSON (`.json`)
- YAML (`.yaml`, `.yml`)

**Examples:**
```bash
# Compare two JSON files (default stylish format)
gendiff file1.json file2.json

# Compare with plain text format
gendiff --format plain file1.json file2.json
gendiff -f plain file1.json file2.json

# Compare with JSON format
gendiff --format json file1.json file2.json
gendiff -f json file1.json file2.json

# Compare two YAML files
gendiff config1.yaml config2.yaml

# Compare nested configuration files
gendiff testdata/fixtures/json/nested1.json testdata/fixtures/json/nested2.json
```

### Output Formats

#### Stylish Format (default)

Tree-like format that shows structure and changes:
- `  key: value` - unchanged values
- `- key: value` - removed values
- `+ key: value` - added values
- Modified values are shown as both removed and added:
  ```
  - key: oldValue
  + key: newValue
  ```

**Example:**
```
{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
        setting6: {
            doge: {
              - wow:
              + wow: so much
            }
            key: value
          + ops: vops
        }
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
    }
}
```

#### Plain Format

Human-readable text format showing property paths:

**Example:**
```
Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting4' was added with value: 'blah blah'
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.doge.wow' was updated. From '' to 'so much'
Property 'common.setting6.ops' was added with value: 'vops'
Property 'group1.baz' was updated. From 'bas' to 'bars'
Property 'group1.nest' was updated. From [complex value] to 'str'
Property 'group2' was removed
Property 'group3' was added with value: [complex value]
```

**Notes:**
- Complex values (objects/arrays) are shown as `[complex value]`
- String values are quoted
- Unchanged properties are not shown

#### JSON Format

Machine-readable format that represents the diff as structured JSON data:

**Example:**
```json
{
  "common": {
    "status": "modified",
    "children": {
      "follow": {
        "status": "added",
        "newValue": false
      },
      "setting1": {
        "status": "unchanged",
        "oldValue": "Value 1"
      },
      "setting2": {
        "status": "removed",
        "oldValue": 200
      },
      "setting3": {
        "status": "modified",
        "oldValue": true,
        "newValue": null
      },
      "setting6": {
        "status": "modified",
        "children": {
          "doge": {
            "status": "modified",
            "children": {
              "wow": {
                "status": "modified",
                "oldValue": "",
                "newValue": "so much"
              }
            }
          }
        }
      }
    }
  }
}
```

**Notes:**
- Each node has a `status` field: `added`, `removed`, `modified`, or `unchanged`
- Leaf nodes have `oldValue` and/or `newValue` depending on the status
- Container nodes have `children` with nested structure
- Uses camelCase naming convention (JavaScript standard)
- Perfect for programmatic processing or API responses

### How to uninstall

```bash
sudo rm /usr/local/bin/gendiff
```

## Asciinema demo

[![asciicast](https://asciinema.org/a/1051136.svg)](https://asciinema.org/a/1051136)

## Development part

### Project structure

```bash
.
├── Makefile
├── cmd
│   └── gendiff
│       └── main.go          # Entry point
├── internal
│   ├── app
│   │   ├── app.go           # Main application logic
│   │   └── app_test.go      # Integration tests
│   ├── differ
│   │   ├── differ.go        # Diff algorithm (tree-based)
│   │   └── differ_test.go   # Differ unit tests
│   ├── formatter
│   │   ├── formatter.go     # Formatter dispatcher and shared utilities
│   │   ├── formatter_test.go
│   │   ├── json.go          # JSON output formatter
│   │   ├── plain.go         # Plain text output formatter
│   │   └── stylish.go       # Stylish tree output formatter
│   └── parser
│       ├── json.go          # JSON parser
│       ├── parser.go        # Parser factory
│       ├── parser_test.go   # Parser tests
│       └── yaml.go          # YAML parser
├── testdata
│   └── fixtures
│       ├── json             # JSON test fixtures
│       │   ├── nested1.json
│       │   ├── nested2.json
│       │   └── ...
│       └── yaml             # YAML test fixtures
│           ├── nested1.yaml
│           ├── nested2.yaml
│           └── ...
├── go.mod
├── go.sum
└── README.md
```

### Key features

- **Recursive nested structure support**: Handles deeply nested objects
- **Multiple input format support**: JSON and YAML
- **Multiple output formats**:
  - Stylish: Tree-like format with visual indicators
  - Plain: Human-readable text with property paths
  - JSON: Machine-readable structured format
- **Type change handling**: Correctly displays when values change from primitives to objects or vice versa
- **Status inheritance**: Children of added/removed nodes are displayed without redundant markers
- **Comprehensive testing**: Unit tests, integration tests, and fixture-based testing
- **High test coverage**: ~85% average coverage across all packages

### How to lint

```bash
make lint
```

This runs:
- `go fmt` - code formatting
- `go vet` - static analysis
- `golangci-lint` - comprehensive linting

### How to run tests

```bash
make test
```

All tests with verbose output:
```bash
go test ./... -v
```

### How to build

```bash
make build
```

This creates the binary at `bin/gendiff`.