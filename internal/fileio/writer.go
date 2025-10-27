package fileio

import (
	"fmt"
	"io"
	"os"

	"github.com/twpayne/go-gpx"
)

// Writer defines the interface for writing GPX data
type Writer interface {
	Write(filename string, g *gpx.GPX) error
}

// GPXWriter implements Writer for GPX files
type GPXWriter struct {
	indent string
}

// NewGPXWriter creates a new GPXWriter
func NewGPXWriter(indent string) *GPXWriter {
	if indent == "" {
		indent = "  "
	}
	return &GPXWriter{
		indent: indent,
	}
}

// Write writes GPX data to a file
func (w *GPXWriter) Write(filename string, g *gpx.GPX) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", filename, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	return w.Encode(file, g)
}

// Encode writes GPX data to an io.Writer
func (w *GPXWriter) Encode(writer io.Writer, g *gpx.GPX) error {
	// Write XML declaration manually since go-gpx's WriteIndent does not include it.
	// This ensures better compatibility with XML parsers and GPX readers.
	if _, err := writer.Write([]byte("<?xml version=\"1.0\"?>\n")); err != nil {
		return fmt.Errorf("failed to write XML declaration: %w", err)
	}

	if err := g.WriteIndent(writer, "", w.indent); err != nil {
		return fmt.Errorf("failed to write GPX: %w", err)
	}
	return nil
}
