package fileio

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chocoby/zweg/internal/converter"
	"github.com/chocoby/zweg/internal/models"
)

func TestJSONReader_Decode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLen   int
		wantErr   bool
		errSubstr string
	}{
		{
			name: "valid single point",
			input: `[{
				"tm": 1609459200,
				"lo": 139.7671,
				"la": 35.6812,
				"th": 90,
				"sp": "5.0",
				"co": 90,
				"al": "10.5",
				"he": 90,
				"ds": "0"
			}]`,
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "valid multiple points",
			input: `[
				{"tm": 1609459200, "lo": 139.7671, "la": 35.6812, "al": "10.5", "sp": "5.0", "ds": "0", "th": 0, "co": 0, "he": 0},
				{"tm": 1609459260, "lo": 139.7672, "la": 35.6813, "al": "11.2", "sp": "5.5", "ds": "100", "th": 0, "co": 0, "he": 0}
			]`,
			wantLen: 2,
			wantErr: false,
		},
		{
			name:      "empty array",
			input:     `[]`,
			wantLen:   0,
			wantErr:   true,
			errSubstr: "no data points found",
		},
		{
			name:      "invalid json",
			input:     `{invalid json}`,
			wantLen:   0,
			wantErr:   true,
			errSubstr: "failed to parse JSON",
		},
		{
			name:      "empty input",
			input:     ``,
			wantLen:   0,
			wantErr:   true,
			errSubstr: "failed to parse JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewJSONReader()
			buf := strings.NewReader(tt.input)
			points, err := reader.Decode(buf)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Decode() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("Decode() error = %v, want substring %q", err, tt.errSubstr)
				}
				return
			}

			if err != nil {
				t.Errorf("Decode() unexpected error = %v", err)
				return
			}

			if len(points) != tt.wantLen {
				t.Errorf("Decode() points length = %d, want %d", len(points), tt.wantLen)
			}
		})
	}
}

func TestJSONReader_Read(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	t.Run("read valid file", func(t *testing.T) {
		filename := filepath.Join(tmpDir, "valid.json")
		content := `[{"tm": 1609459200, "lo": 139.7671, "la": 35.6812, "al": "10.5", "sp": "5.0", "ds": "0", "th": 0, "co": 0, "he": 0}]`
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		reader := NewJSONReader()
		points, err := reader.Read(filename)

		if err != nil {
			t.Errorf("Read() unexpected error = %v", err)
		}

		if len(points) != 1 {
			t.Errorf("Read() points length = %d, want 1", len(points))
		}
	})

	t.Run("read non-existent file", func(t *testing.T) {
		reader := NewJSONReader()
		_, err := reader.Read(filepath.Join(tmpDir, "nonexistent.json"))

		if err == nil {
			t.Error("Read() error = nil, want error for non-existent file")
		}
	})
}

func TestGPXWriter_Encode(t *testing.T) {
	// Create a simple GPX structure for testing
	points := []models.Point{
		{
			Tm: 1609459200,
			Lo: 139.7671,
			La: 35.6812,
			Al: "10.5",
		},
	}

	conv := converter.New(nil)
	gpxData, err := conv.Convert(points, "Test Track")
	if err != nil {
		t.Fatalf("Failed to create test GPX: %v", err)
	}

	t.Run("write to buffer", func(t *testing.T) {
		writer := NewGPXWriter("  ")
		var buf bytes.Buffer

		err := writer.Encode(&buf, gpxData)
		if err != nil {
			t.Errorf("Encode() unexpected error = %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "<gpx") {
			t.Error("Encode() output missing GPX root element")
		}
		if !strings.Contains(output, "Test Track") {
			t.Error("Encode() output missing track name")
		}
	})

	t.Run("write with custom indent", func(t *testing.T) {
		writer := NewGPXWriter("\t")
		var buf bytes.Buffer

		err := writer.Encode(&buf, gpxData)
		if err != nil {
			t.Errorf("Encode() unexpected error = %v", err)
		}

		output := buf.String()
		if len(output) == 0 {
			t.Error("Encode() output is empty")
		}
	})

	t.Run("write with empty indent uses default", func(t *testing.T) {
		writer := NewGPXWriter("")
		var buf bytes.Buffer

		err := writer.Encode(&buf, gpxData)
		if err != nil {
			t.Errorf("Encode() unexpected error = %v", err)
		}

		if buf.Len() == 0 {
			t.Error("Encode() output is empty")
		}
	})

	t.Run("output includes XML declaration", func(t *testing.T) {
		writer := NewGPXWriter("  ")
		var buf bytes.Buffer

		err := writer.Encode(&buf, gpxData)
		if err != nil {
			t.Errorf("Encode() unexpected error = %v", err)
		}

		output := buf.String()
		expectedDeclaration := "<?xml version=\"1.0\"?>"

		if !strings.HasPrefix(output, expectedDeclaration) {
			t.Errorf("Encode() output does not start with XML declaration, got: %s", output[:50])
		}

		// Verify the declaration is on its own line
		lines := strings.Split(output, "\n")
		if len(lines) < 2 {
			t.Error("Encode() output should have XML declaration on separate line")
		}
		if lines[0] != expectedDeclaration {
			t.Errorf("Encode() first line = %q, want %q", lines[0], expectedDeclaration)
		}
	})
}

func TestGPXWriter_Write(t *testing.T) {
	tmpDir := t.TempDir()

	points := []models.Point{
		{
			Tm: 1609459200,
			Lo: 139.7671,
			La: 35.6812,
			Al: "10.5",
		},
	}

	conv := converter.New(nil)
	gpxData, err := conv.Convert(points, "Test Track")
	if err != nil {
		t.Fatalf("Failed to create test GPX: %v", err)
	}

	t.Run("write to file", func(t *testing.T) {
		filename := filepath.Join(tmpDir, "output.gpx")
		writer := NewGPXWriter("  ")

		err := writer.Write(filename, gpxData)
		if err != nil {
			t.Errorf("Write() unexpected error = %v", err)
		}

		// Verify file exists and has content
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Errorf("Failed to read written file: %v", err)
		}

		if len(content) == 0 {
			t.Error("Write() created empty file")
		}

		if !strings.Contains(string(content), "<gpx") {
			t.Error("Write() output missing GPX root element")
		}

		// Verify XML declaration is present
		if !strings.HasPrefix(string(content), "<?xml version=\"1.0\"?>") {
			t.Error("Write() output missing XML declaration at start")
		}
	})

	t.Run("write to invalid path", func(t *testing.T) {
		writer := NewGPXWriter("  ")
		err := writer.Write("/invalid/path/output.gpx", gpxData)

		if err == nil {
			t.Error("Write() error = nil, want error for invalid path")
		}
	})
}
