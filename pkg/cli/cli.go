package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chocoby/zweg/internal/converter"
	"github.com/chocoby/zweg/internal/fileio"
	"github.com/chocoby/zweg/internal/models"
)

// CLI represents the command-line interface.
type CLI struct {
	reader    fileio.Reader
	writer    fileio.Writer
	converter converter.Converter
	stdout    io.Writer
	stderr    io.Writer
}

// Config holds CLI configuration.
type Config struct {
	Reader    fileio.Reader
	Writer    fileio.Writer
	Converter converter.Converter
	Stdout    io.Writer
	Stderr    io.Writer
}

// New creates a new CLI instance.
func New(config *Config) *CLI {
	if config == nil {
		config = &Config{}
	}

	if config.Reader == nil {
		config.Reader = fileio.NewJSONReader()
	}

	if config.Writer == nil {
		config.Writer = fileio.NewGPXWriter("  ")
	}

	if config.Converter == nil {
		config.Converter = converter.New(nil)
	}

	return &CLI{
		reader:    config.Reader,
		writer:    config.Writer,
		converter: config.Converter,
		stdout:    config.Stdout,
		stderr:    config.Stderr,
	}
}

// validateOutputPath validates and sanitizes an output path to prevent path traversal attacks.
// It returns the cleaned absolute path and an error if the path is unsafe.
func validateOutputPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("output path is empty")
	}

	// Clean the path to remove any .., ., and redundant separators
	cleaned := filepath.Clean(path)

	// Convert to absolute path
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %q: %w", path, err)
	}

	// Check for path traversal attempts by looking for .. in the cleaned path
	// This catches relative paths that try to escape the current directory
	if strings.Contains(filepath.ToSlash(cleaned), "..") {
		return "", fmt.Errorf("path %q contains invalid relative path components", path)
	}

	return absPath, nil
}

// generateOutputFilename generates output filename based on GPS points timestamp.
// Returns YYYYMMDD-HHMMSS.gpx format.
// If outputDir is specified, the file is placed in that directory.
// Otherwise, it is placed in the same directory as the input file.
// The timezoneOffset parameter is used to adjust the timestamp (in seconds).
func (c *CLI) generateOutputFilename(inputFile string, outputDir string, points []models.Point, timezoneOffset int) (string, error) {
	if len(points) == 0 {
		return inputFile + ".gpx", nil
	}

	firstPoint := points[0]
	var timestamp time.Time
	if timezoneOffset == 0 {
		timestamp = firstPoint.Timestamp()
	} else {
		timestamp = firstPoint.TimestampWithOffset(timezoneOffset)
	}
	baseName := timestamp.Format("20060102-150405") + ".gpx"

	dir := outputDir
	if dir == "" {
		dir = filepath.Dir(inputFile)
	} else {
		// Validate output directory to prevent path traversal
		validatedDir, err := validateOutputPath(dir)
		if err != nil {
			return "", fmt.Errorf("invalid output directory: %w", err)
		}
		dir = validatedDir
	}

	return filepath.Join(dir, baseName), nil
}

// Run executes the CLI command.
// If outputFile is empty, it will be auto-generated based on the track start time.
// outputDir is used only when outputFile is not specified.
// timezoneOffset is the timezone offset in seconds for GPX timestamps and filename generation.
func (c *CLI) Run(inputFile, outputFile, outputDir, trackName string, timezoneOffset int) error {
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}

	points, err := c.reader.Read(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if outputFile == "" {
		outputFile, err = c.generateOutputFilename(inputFile, outputDir, points, timezoneOffset)
		if err != nil {
			return fmt.Errorf("failed to generate output filename: %w", err)
		}
	} else {
		// Validate explicitly specified output file path
		validatedOutput, err := validateOutputPath(outputFile)
		if err != nil {
			return fmt.Errorf("invalid output file path: %w", err)
		}
		outputFile = validatedOutput
	}

	if trackName == "" {
		trackName = "Track"
	}

	// Ensure output directory exists
	outputFileDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputFileDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	gpxData, err := c.converter.Convert(points, trackName)
	if err != nil {
		return fmt.Errorf("failed to convert data: %w", err)
	}

	if err := c.writer.Write(outputFile, gpxData); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	if c.stdout != nil {
		if _, err := fmt.Fprintf(c.stdout, "Successfully converted %d points to GPX: %s\n", len(points), outputFile); err != nil {
			return fmt.Errorf("failed to write output message: %w", err)
		}
	}

	return nil
}

// ParseTimezoneOffset parses a timezone offset string and returns the offset in seconds.
// Supported formats: ±HH:MM or ±HHMM (e.g., +09:00, -05:00, +0900, -0500)
// Valid range: -12:00 to +14:00
func ParseTimezoneOffset(offset string) (int, error) {
	if offset == "" {
		return 0, fmt.Errorf("timezone offset is empty")
	}

	if len(offset) < 3 {
		return 0, fmt.Errorf("invalid timezone offset format: %q (expected ±HH:MM or ±HHMM)", offset)
	}

	// Check sign
	var sign int
	switch offset[0] {
	case '+':
		sign = 1
		offset = offset[1:]
	case '-':
		sign = -1
		offset = offset[1:]
	default:
		return 0, fmt.Errorf("timezone offset must start with + or -: %q (expected ±HH:MM or ±HHMM format)", offset)
	}

	var hours, minutes int
	var err error

	// Parse HH:MM or HHMM format
	if strings.Contains(offset, ":") {
		parts := strings.Split(offset, ":")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid timezone offset format: %q (expected ±HH:MM format with exactly one colon)", offset)
		}
		if len(parts[0]) != 2 || len(parts[1]) != 2 {
			return 0, fmt.Errorf("invalid timezone offset format: %q (expected ±HH:MM format with 2-digit hours and minutes)", offset)
		}
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid hours in timezone offset %q: %w (expected numeric value 00-14)", offset, err)
		}
		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in timezone offset %q: %w (expected numeric value 00-59)", offset, err)
		}
	} else {
		// HHMM format
		if len(offset) != 4 {
			return 0, fmt.Errorf("invalid timezone offset format: %q (expected ±HHMM format with 4 digits)", offset)
		}
		hours, err = strconv.Atoi(offset[0:2])
		if err != nil {
			return 0, fmt.Errorf("invalid hours in timezone offset %q: %w (expected numeric value 00-14)", offset, err)
		}
		minutes, err = strconv.Atoi(offset[2:4])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in timezone offset %q: %w (expected numeric value 00-59)", offset, err)
		}
	}

	// Validate range
	if hours < 0 || hours > 14 {
		return 0, fmt.Errorf("hours out of valid range in %q: got %d (expected 0-14)", offset, hours)
	}
	if minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("minutes out of valid range in %q: got %d (expected 0-59)", offset, minutes)
	}
	if hours == 14 && minutes > 0 {
		return 0, fmt.Errorf("timezone offset out of valid range: %q (maximum offset is +14:00)", offset)
	}

	totalSeconds := sign * (hours*3600 + minutes*60)

	// Final range check: -12:00 to +14:00
	if totalSeconds < -12*3600 || totalSeconds > 14*3600 {
		return 0, fmt.Errorf("timezone offset out of valid range: %q (expected -12:00 to +14:00)", offset)
	}

	return totalSeconds, nil
}
