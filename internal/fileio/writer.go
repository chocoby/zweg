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
	defer file.Close()

	return w.WriteTo(file, g)
}

// WriteTo writes GPX data to an io.Writer
func (w *GPXWriter) WriteTo(writer io.Writer, g *gpx.GPX) error {
	if err := g.WriteIndent(writer, "", w.indent); err != nil {
		return fmt.Errorf("failed to write GPX: %w", err)
	}
	return nil
}
