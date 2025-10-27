package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

// generateOutputFilename generates output filename based on GPS points timestamp.
// Returns YYYYMMDD-HHMMSS.gpx format.
// If outputDir is specified, the file is placed in that directory.
// Otherwise, it is placed in the same directory as the input file.
func (c *CLI) generateOutputFilename(inputFile string, outputDir string, points []models.Point) string {
	if len(points) == 0 {
		return inputFile + ".gpx"
	}

	firstPoint := points[0]
	timestamp := firstPoint.LocalTimestamp()
	baseName := timestamp.Format("20060102-150405") + ".gpx"

	dir := outputDir
	if dir == "" {
		dir = filepath.Dir(inputFile)
	}
	return filepath.Join(dir, baseName)
}

// Run executes the CLI command.
// If outputFile is empty, it will be auto-generated based on the track start time.
// outputDir is used only when outputFile is not specified.
func (c *CLI) Run(inputFile, outputFile, outputDir, trackName string) error {
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}

	points, err := c.reader.Read(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if outputFile == "" {
		outputFile = c.generateOutputFilename(inputFile, outputDir, points)
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
