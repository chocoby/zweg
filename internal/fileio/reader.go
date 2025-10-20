package fileio

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chocoby/zweg/internal/models"
)

// Reader defines the interface for reading GPS data
type Reader interface {
	Read(filename string) ([]models.Point, error)
}

// JSONReader implements Reader for JSON files
type JSONReader struct{}

// NewJSONReader creates a new JSONReader
func NewJSONReader() *JSONReader {
	return &JSONReader{}
}

// Read reads and parses ZweiteGPS JSON data from a file
func (r *JSONReader) Read(filename string) ([]models.Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filename, err)
	}
	defer file.Close()

	return r.ReadFrom(file)
}

// ReadFrom reads and parses ZweiteGPS JSON data from an io.Reader
func (r *JSONReader) ReadFrom(reader io.Reader) ([]models.Point, error) {
	var points []models.Point
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&points); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(points) == 0 {
		return nil, fmt.Errorf("no data points found in JSON")
	}

	return points, nil
}
