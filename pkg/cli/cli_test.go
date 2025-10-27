package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_Run_AutoGenerateOutputFilename(t *testing.T) {
	tests := []struct {
		name        string
		jsonContent string
		inputFile   string
		wantPrefix  string // Expected prefix for timestamp-based filename
		wantErr     bool
	}{
		{
			name: "valid JSON with single point - auto-generate filename",
			jsonContent: `[
				{
					"tm": 1729411200,
					"lo": 139.7454,
					"la": 35.6812,
					"th": 0,
					"sp": "0",
					"co": 0,
					"al": "0",
					"he": 0,
					"ds": "0"
				}
			]`,
			inputFile:  "test.json",
			wantPrefix: "20241020-",
			wantErr:    false,
		},
		{
			name: "valid JSON with multiple points - auto-generate filename",
			jsonContent: `[
				{
					"tm": 1609459200,
					"lo": 139.7454,
					"la": 35.6812,
					"th": 0,
					"sp": "0",
					"co": 0,
					"al": "0",
					"he": 0,
					"ds": "0"
				},
				{
					"tm": 1609459300,
					"lo": 139.7455,
					"la": 35.6813,
					"th": 0,
					"sp": "10",
					"co": 0,
					"al": "5",
					"he": 0,
					"ds": "100"
				}
			]`,
			inputFile:  "test.json",
			wantPrefix: "20210101-",
			wantErr:    false,
		},
		{
			name:        "invalid JSON - should return error",
			jsonContent: `{invalid json}`,
			inputFile:   "test.json",
			wantErr:     true,
		},
		{
			name:        "empty JSON array - should return error",
			jsonContent: `[]`,
			inputFile:   "test.json",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "zweg-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(tmpDir); err != nil {
					t.Errorf("Failed to remove temp dir: %v", err)
				}
			}()

			// Create test JSON file
			inputPath := filepath.Join(tmpDir, tt.inputFile)
			if err := os.WriteFile(inputPath, []byte(tt.jsonContent), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create CLI instance
			cli := New(nil)

			// Test Run with empty outputFile (auto-generate)
			err = cli.Run(inputPath, "", "Test Track")

			if tt.wantErr {
				if err == nil {
					t.Errorf("Run() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Run() unexpected error = %v", err)
				return
			}

			// Check if output file was created with correct name format
			files, err := os.ReadDir(tmpDir)
			if err != nil {
				t.Fatalf("Failed to read temp dir: %v", err)
			}

			var gpxFile string
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".gpx") {
					gpxFile = file.Name()
					break
				}
			}

			if gpxFile == "" {
				t.Errorf("No GPX file was created")
				return
			}

			// Check if it has the timestamp format
			if !strings.HasPrefix(gpxFile, tt.wantPrefix) {
				t.Errorf("Generated filename = %v, want prefix %v", gpxFile, tt.wantPrefix)
			}

			// Check format: YYYYMMDD-HHMMSS.gpx (should be 19 characters)
			if len(gpxFile) != 19 {
				t.Errorf("Generated filename length = %d, want 19 (YYYYMMDD-HHMMSS.gpx)", len(gpxFile))
			}
		})
	}
}

func TestCLI_Run_ExplicitOutputFilename(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "zweg-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test JSON file
	jsonContent := `[
		{
			"tm": 1609459200,
			"lo": 139.7454,
			"la": 35.6812,
			"th": 0,
			"sp": "0",
			"co": 0,
			"al": "0",
			"he": 0,
			"ds": "0"
		}
	]`
	inputPath := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(inputPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create CLI instance
	cli := New(nil)

	// Test Run with explicit output filename
	explicitOutput := filepath.Join(tmpDir, "custom-output.gpx")
	err = cli.Run(inputPath, explicitOutput, "Test Track")
	if err != nil {
		t.Errorf("Run() unexpected error = %v", err)
		return
	}

	// Check if the explicit output file was created
	if _, err := os.Stat(explicitOutput); os.IsNotExist(err) {
		t.Errorf("Expected output file %v was not created", explicitOutput)
	}
}

func TestCLI_Run_NonExistentFile(t *testing.T) {
	cli := New(nil)
	inputFile := "/nonexistent/file/that/does/not/exist.json"

	err := cli.Run(inputFile, "", "Test Track")
	if err == nil {
		t.Errorf("Run() error = nil, want error for non-existent file")
	}
}
