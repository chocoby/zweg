# zweg

ZweiteGPS to GPX Converter - A command-line tool to convert ZweiteGPS JSON format to standard GPX format.

## Features

- Convert ZweiteGPS JSON data to GPX 1.1 format
- Clean, testable architecture with separation of concerns
- Comprehensive test coverage
- Easy to build and install

## Installation

### From source

```bash
git clone https://github.com/chocoby/zweg.git
cd zweg
make install
```

Or using Go directly:

```bash
go install github.com/chocoby/zweg/cmd/zweg@latest
```

## Usage

```bash
zweg [--track-name <name>] <input.json> [output.gpx]
```

### Options

- `--track-name <name>`: Name for the GPS track (default: "Track")

### Arguments

- `input.json`: Path to the input ZweiteGPS JSON file
- `output.gpx`: Path to the output GPX file (optional, defaults to `input.json.gpx`)

### Examples

```bash
# Auto-generate output filename (data.json -> data.json.gpx)
zweg data.json

# With custom output filename
zweg data.json output.gpx

# With custom track name and auto-generated output
zweg --track-name "My Morning Run" data.json

# With custom track name and custom output
zweg --track-name "My Morning Run" data.json output.gpx

# Show help
zweg --help
```

## Development

### Prerequisites

- Go 1.25.3 or later

### Building

```bash
make build
```

The binary will be created in `./bin/zweg`

### Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
make coverage
```

### Available Make Targets

- `make build` - Build the application
- `make test` - Run all tests
- `make coverage` - Run tests with coverage report
- `make bench` - Run benchmarks
- `make fmt` - Format code
- `make vet` - Run go vet
- `make lint` - Run golangci-lint (requires golangci-lint to be installed)
- `make install` - Install the binary
- `make clean` - Remove build artifacts
- `make check` - Run fmt, vet, and test
- `make help` - Show help message

## ZweiteGPS JSON Format Specification

The input JSON file should be an array of GPS trackpoints with the following fields:

| Field | Type   | Description                                      | Example      |
|-------|--------|--------------------------------------------------|--------------|
| `tm`  | number | Unix timestamp (seconds since epoch)             | 1609459200   |
| `la`  | number | Latitude in decimal degrees                      | 35.658581    |
| `lo`  | number | Longitude in decimal degrees                     | 139.745438   |
| `al`  | string | Altitude in meters                               | "150.0"      |
| `sp`  | string | Speed in meters per second                       | "1.0"        |
| `co`  | number | Course/bearing in degrees (0-360, -1 if unknown) | 90           |
| `th`  | number | True heading in degrees (0-360)                  | 85           |
| `he`  | number | Heading in degrees (0-360)                       | 88           |
| `ds`  | string | Distance in meters                               | "0.0"        |
| `ms`  | number | Milliseconds (optional)                          | 0            |
| `ow`  | string | Owner/device information (optional)              | "iPhone ..." |

### Example

```json
[
  {
    "tm": 1609459200,
    "lo": 139.745438,
    "la": 35.658581,
    "al": "150.0",
    "sp": "1.0",
    "co": 90,
    "th": 85,
    "he": 88,
    "ds": "0.0",
    "ms": 0,
    "ow": "iPhone [iPhone15,2 v18.0.1, ZweiteGPS, v48]"
  }
]
```

## Project Structure

```
zweg/
├── cmd/
│   └── zweg/           # Entry point
├── internal/
│   ├── converter/      # GPX conversion logic
│   ├── models/         # Data models
│   └── fileio/         # File I/O operations
├── pkg/
│   └── cli/            # CLI interface
├── Makefile            # Build tasks
└── README.md
```

## License

MIT
