### Quality & Testing

| Status | Badge |
|--------|-------|
| CI Pipeline | [![CI](https://github.com/EvgeniiIvanov/go-project-244/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/EvgeniiIvanov/go-project-244/actions/workflows/ci.yml) |
| Test Coverage | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=EvgeniiIvanov_go-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=EvgeniiIvanov_go-project-244) |

### Hexlet tests and linter status:
[![Actions Status](https://github.com/EvgeniiIvanov/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/EvgeniiIvanov/go-project-244/actions)

## Description

A command-line tool that compares two configuration files and shows the differences, similar to `diff`. Supports JSON and YAML formats, with a stylish output format that clearly indicates added, removed, and modified values in nested structures.

## Usage

### How to install

```bash
make build
sudo cp bin/gendiff /usr/local/bin/
```

### How to use

```bash
gendiff <file1> <file2>
```

Supported formats:
- JSON (`.json`)
- YAML (`.yaml`, `.yml`)

Examples:
```bash
# Compare two JSON files
gendiff file1.json file2.json

# Compare two YAML files
gendiff config1.yaml config2.yaml

# Compare nested configuration files
gendiff testdata/fixtures/json/nested1.json testdata/fixtures/json/nested2.json
```

### Output format

The tool uses a "stylish" format that shows:
- `  key: value` - unchanged values
- `- key: value` - removed values
- `+ key: value` - added values
- Modified values are shown as both removed and added:
  ```
  - key: oldValue
  + key: newValue
  ```

Example output:
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

### How to uninstall

```bash
sudo rm /usr/local/bin/gendiff
```

## Asciinema demo

[![asciicast](https://asciinema.org/a/1025741.svg)](https://asciinema.org/a/1025741)

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
│   │   ├── formatter.go     # Formatter interface
│   │   ├── formatter_test.go
│   │   └── stylish.go       # Stylish output formatter
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
- **Multiple format support**: JSON and YAML
- **Type change handling**: Correctly displays when values change from primitives to objects or vice versa
- **Status inheritance**: Children of added/removed nodes are displayed without redundant markers
- **Comprehensive testing**: Unit tests, integration tests, and fixture-based testing

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