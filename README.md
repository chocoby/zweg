# zweg

A command-line tool to convert ZweiteGPS JSON format to standard GPX format.

## Features

- Convert ZweiteGPS JSON data to GPX 1.1 format
- Auto-generate output filenames based on track start time

## Installation

### Download pre-built binaries

Download the latest release for your platform from the [Releases](https://github.com/chocoby/zweg/releases) page.

Extract the archive and move the binary to a directory in your PATH.

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
zweg [options] <input.json> [output.gpx]
```

### Options

- `--track-name <name>`: Name for the GPS track (default: "Track")
- `-d, --output-dir <directory>`: Output directory for the GPX file (ignored if output file is specified)
- `--timezone-offset <offset>`: Timezone offset for auto-generated filename in ±HH:MM or ±HHMM format (default: "+00:00" UTC). **Note: This only affects the filename; GPX timestamps are always in UTC per GPX 1.1 specification.**
- `--version`: Show version information

### Arguments

- `input.json`: Path to the input ZweiteGPS JSON file
- `output.gpx`: Path to the output GPX file (optional, defaults to YYYYMMDD-HHMMSS.gpx based on track start time in the specified timezone)

### Examples

```bash
# Auto-generate output filename based on track start time
# Output: YYYYMMDD-HHMMSS.gpx (same directory as input file)
zweg data.json

# With custom output filename
zweg data.json output.gpx

# With output directory
zweg -d ./gpx-output data.json

# With timezone offset
# Only affects the filename, GPX content stays in UTC
zweg --timezone-offset +09:00 data.json

# With custom track name
zweg --track-name "My Morning Run" data.json

# Show help
zweg --help
```

### Batch Conversion Examples

Convert all JSON files in the current directory

```bash
# Basic batch conversion
for f in *.json; do zweg "$f"; done

# With output directory
for f in *.json; do zweg -d ./gpx "$f"; done

# With timezone offset
for f in *.json; do zweg --timezone-offset +09:00 "$f"; done
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

```bash
make test
```

### Linting

```bash
make lint
```

Requires golangci-lint to be installed.

## Timezone Offset Support

### GPX Timestamps

**All timestamps in GPX files are always in UTC**, as required by the GPX 1.1 specification. This ensures maximum compatibility with GPS software like Garmin Connect, Strava, and other applications that expect standard-compliant GPX files.

### Filename Generation (Timezone Offset Support)

The `--timezone-offset` option allows you to specify a timezone for **auto-generated filenames only**. This makes it easier to identify files based on your local time without compromising GPX compatibility.

### Supported Formats

- `±HH:MM` format (e.g., `+09:00`, `-05:00`)
- `±HHMM` format (e.g., `+0900`, `-0500`)

### Examples with Different Timezones

Given an input with Unix timestamp `1609459200` (2021-01-01 00:00:00 UTC):

```bash
# UTC (default)
zweg data.json
# Output filename: 20210101-000000.gpx
# GPX timestamps: 2021-01-01T00:00:00Z (all timestamps in UTC)

# Japan Standard Time (JST: UTC+9)
zweg --timezone-offset +09:00 data.json
# Output filename: 20210101-090000.gpx (JST time)
# GPX timestamps: 2021-01-01T00:00:00Z (still UTC)

# Eastern Standard Time (EST: UTC-5)
zweg --timezone-offset -05:00 data.json
# Output filename: 20201231-190000.gpx (EST time)
# GPX timestamps: 2021-01-01T00:00:00Z (still UTC)
```

## ZweiteGPS JSON Format Specification

> **Note**: This format specification is based on reverse engineering and observation, not official documentation from ZweiteGPS.

The input JSON file should be an array of GPS trackpoints with the following fields:

| Field | Type   | Description                                      | Example      |
| ----- | ------ | ------------------------------------------------ | ------------ |
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

## License

MIT
