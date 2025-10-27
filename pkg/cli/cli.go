package cli

import (
	"fmt"
	"io"

	"github.com/chocoby/zweg/internal/converter"
	"github.com/chocoby/zweg/internal/fileio"
)

// CLI represents the command-line interface
type CLI struct {
	reader    fileio.Reader
	writer    fileio.Writer
	converter converter.Converter
	stdout    io.Writer
	stderr    io.Writer
}

// Config holds CLI configuration
type Config struct {
	Reader    fileio.Reader
	Writer    fileio.Writer
	Converter converter.Converter
	Stdout    io.Writer
	Stderr    io.Writer
}

// New creates a new CLI instance
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

// Run executes the CLI command
func (c *CLI) Run(inputFile, outputFile, trackName string) error {
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}
	if outputFile == "" {
		return fmt.Errorf("output file is required")
	}
	if trackName == "" {
		trackName = "Track"
	}

	points, err := c.reader.Read(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
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
