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
zweg [--track-name <name>] <input.json> <output.gpx>
```

### Options

- `--track-name <name>`: Name for the GPS track (default: "Track")

### Arguments

- `input.json`: Path to the input ZweiteGPS JSON file
- `output.gpx`: Path to the output GPX file

### Examples

```bash
# With custom track name
zweg --track-name "My Morning Run" data.json output.gpx

# Using default track name
zweg data.json output.gpx

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
