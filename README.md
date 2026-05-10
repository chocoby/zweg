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

- `--track-name <name>`: Name for the GPS track. When omitted, the fallback chain is used: `tl` (log title from JSON) → English `ms` name (Walking/Jogging/etc.) → `Track`.
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

> The JSON format specification can be viewed inside the ZweiteGPS app.

The input JSON file is an array of GPS trackpoints with the following fields.
Required fields must always be present; optional fields are populated only when the
device or session provides them.

### Core fields

| Field | Type   | Description                                      | Example    |
| ----- | ------ | ------------------------------------------------ | ---------- |
| `tm`  | number | Unix timestamp in seconds                        | 1609459200 |
| `lo`  | number | Longitude in decimal degrees                     | 139.745438 |
| `la`  | number | Latitude in decimal degrees                      | 35.658581  |
| `al`  | string | Altitude in meters                               | "150.0"    |
| `sp`  | string | Speed in meters per second                       | "1.0"      |
| `co`  | number | Course / bearing in degrees (0-360, -1 unknown)  | 90         |
| `th`  | number | True heading in degrees (0-360)                  | 85         |
| `he`  | number | Magnetic heading in degrees (0-360)              | 88         |
| `ds`  | string | Distance in meters                               | "0.0"      |

### Optional metadata

| Field | Type   | Description                                              | Example      |
| ----- | ------ | -------------------------------------------------------- | ------------ |
| `ms`  | number | Means of transportation (see table below) — used as track name fallback after `tl` | 1            |
| `ow`  | string | Owner / device information                               | "iPhone ..." |
| `tl`  | string | Log title — used as the default GPX track name           | "Tokyo Run"  |
| `dp`  | string | Description / memo — copied to GPX `desc`                | "checkpoint" |

#### Means (`ms`) values

| Value | Means       |
| ----- | ----------- |
| 0     | Walking     |
| 1     | Jogging     |
| 2     | Bicycle     |
| 3     | MotorCycle  |
| 4     | AutoMobile  |
| 5     | Train       |
| 6     | Misc        |

### Optional sensor data

| Field | Type   | Description                                              |
| ----- | ------ | -------------------------------------------------------- |
| `ap`  | number | Atmospheric pressure                                     |
| `ha`  | number | Horizontal accuracy in meters — copied to GPX `hdop`     |
| `va`  | number | Vertical accuracy in meters — copied to GPX `vdop`       |
| `xa`  | number | Heading accuracy                                         |
| `ra`  | number | Relative altitude                                        |
| `ws`  | number | Number of steps                                          |
| `gx`  | number | Gravity acceleration X                                   |
| `gy`  | number | Gravity acceleration Y                                   |
| `gz`  | number | Gravity acceleration Z                                   |
| `ax`  | number | User acceleration X                                      |
| `ay`  | number | User acceleration Y                                      |
| `az`  | number | User acceleration Z                                      |
| `ep`  | number | Pitch angle                                              |
| `er`  | number | Roll angle                                               |
| `ey`  | number | Yaw angle                                                |
| `pf`  | number | Peak frequency                                           |

### Field mapping to GPX

The following input fields are converted to GPX output:

| ZweiteGPS field | GPX element                                                                |
| --------------- | -------------------------------------------------------------------------- |
| `tm`            | `<time>` on track points / waypoints and `<metadata><time>`                |
| `la`            | `lat` attribute on track points / waypoints                                |
| `lo`            | `lon` attribute on track points / waypoints                                |
| `al`            | `<ele>`                                                                    |
| `ha`            | `<hdop>`                                                                   |
| `va`            | `<vdop>`                                                                   |
| `dp`            | `<desc>` on track points and start/goal waypoints                          |
| `tl`            | Track `<name>` (used when `--track-name` is not specified)                 |
| `ms`            | Track `<name>` fallback when both `--track-name` and `tl` are absent       |

### Unsupported fields

The following fields are recognized in the input but **not** written to the GPX output yet. They are read for forward compatibility and silently ignored during conversion:

| Field                | Description                          |
| -------------------- | ------------------------------------ |
| `sp`                 | Speed                                |
| `co`                 | Course / bearing                     |
| `th`                 | True heading                         |
| `he`                 | Magnetic heading                     |
| `ds`                 | Distance                             |
| `ow`                 | Owner / device information           |
| `ap`                 | Atmospheric pressure                 |
| `ra`                 | Relative altitude                    |
| `ws`                 | Number of steps                      |
| `xa`                 | Heading accuracy                     |
| `gx`, `gy`, `gz`     | Gravity acceleration (X / Y / Z)     |
| `ax`, `ay`, `az`     | User acceleration (X / Y / Z)        |
| `ep`, `er`, `ey`     | Pitch / Roll / Yaw angle             |
| `pf`                 | Peak frequency                       |

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
    "ow": "iPhone [iPhone15,2 v18.0.1, ZweiteGPS, v48]",
    "tl": "Morning Run",
    "ha": 5.0,
    "va": 3.0
  }
]
```

## License

MIT
