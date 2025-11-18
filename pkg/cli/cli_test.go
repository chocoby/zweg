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
			tmpDir := t.TempDir()

			// Create test JSON file
			inputPath := filepath.Join(tmpDir, tt.inputFile)
			if err := os.WriteFile(inputPath, []byte(tt.jsonContent), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create CLI instance
			cli := New(nil)

			// Test Run with empty outputFile (auto-generate)
			err := cli.Run(inputPath, "", "", "Test Track", 0)

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
	tmpDir := t.TempDir()

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
	err := cli.Run(inputPath, explicitOutput, "", "Test Track", 0)
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

	err := cli.Run(inputFile, "", "", "Test Track", 0)
	if err == nil {
		t.Errorf("Run() error = nil, want error for non-existent file")
	}
}

func TestCLI_Run_WithOutputDir(t *testing.T) {
	tmpDir := t.TempDir()

	jsonContent := `[
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
	]`
	inputPath := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(inputPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	outputDir := filepath.Join(tmpDir, "output")
	cli := New(nil)

	// Run with output directory specified
	err := cli.Run(inputPath, "", outputDir, "Test Track", 0)
	if err != nil {
		t.Errorf("Run() unexpected error = %v", err)
		return
	}

	// Check if the output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("Expected output directory %v was not created", outputDir)
	}

	// Check if a GPX file was created in the output directory
	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 file in output directory, got %d", len(files))
	}
	if len(files) > 0 && !strings.HasSuffix(files[0].Name(), ".gpx") {
		t.Errorf("Expected GPX file, got %s", files[0].Name())
	}
}

func TestCLI_Run_WithOutputDirAndOutputFile(t *testing.T) {
	tmpDir := t.TempDir()

	jsonContent := `[
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
	]`
	inputPath := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(inputPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	outputDir := filepath.Join(tmpDir, "ignored-dir")
	outputFile := filepath.Join(tmpDir, "explicit-output.gpx")
	cli := New(nil)

	// Run with both output directory and output file
	// The output file should take precedence
	err := cli.Run(inputPath, outputFile, outputDir, "Test Track", 0)
	if err != nil {
		t.Errorf("Run() unexpected error = %v", err)
		return
	}

	// Check that the output file was created at the explicit path
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %v was not created", outputFile)
	}

	// Check that the output directory was NOT created (ignored)
	if _, err := os.Stat(outputDir); err == nil {
		t.Errorf("Output directory %v should not have been created", outputDir)
	}
}

func TestCLI_Run_WithNestedOutputDir(t *testing.T) {
	tmpDir := t.TempDir()

	jsonContent := `[
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
	]`
	inputPath := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(inputPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	nestedOutputDir := filepath.Join(tmpDir, "level1", "level2", "level3")
	cli := New(nil)

	// Run with nested output directory
	err := cli.Run(inputPath, "", nestedOutputDir, "Test Track", 0)
	if err != nil {
		t.Errorf("Run() unexpected error = %v", err)
		return
	}

	// Check if the nested output directory was created
	if _, err := os.Stat(nestedOutputDir); os.IsNotExist(err) {
		t.Errorf("Expected nested output directory %v was not created", nestedOutputDir)
	}

	// Check if a GPX file was created in the nested directory
	files, err := os.ReadDir(nestedOutputDir)
	if err != nil {
		t.Fatalf("Failed to read nested output directory: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 file in nested output directory, got %d", len(files))
	}
}

func TestParseTimezoneOffset(t *testing.T) {
	tests := []struct {
		name    string
		offset  string
		want    int
		wantErr bool
	}{
		{
			name:    "UTC",
			offset:  "+00:00",
			want:    0,
			wantErr: false,
		},
		{
			name:    "JST with colon",
			offset:  "+09:00",
			want:    9 * 3600,
			wantErr: false,
		},
		{
			name:    "JST without colon",
			offset:  "+0900",
			want:    9 * 3600,
			wantErr: false,
		},
		{
			name:    "EST with colon",
			offset:  "-05:00",
			want:    -5 * 3600,
			wantErr: false,
		},
		{
			name:    "EST without colon",
			offset:  "-0500",
			want:    -5 * 3600,
			wantErr: false,
		},
		{
			name:    "India Standard Time",
			offset:  "+05:30",
			want:    5*3600 + 30*60,
			wantErr: false,
		},
		{
			name:    "Maximum positive offset",
			offset:  "+14:00",
			want:    14 * 3600,
			wantErr: false,
		},
		{
			name:    "Maximum negative offset",
			offset:  "-12:00",
			want:    -12 * 3600,
			wantErr: false,
		},
		{
			name:    "Invalid - empty string",
			offset:  "",
			wantErr: true,
		},
		{
			name:    "Invalid - no sign",
			offset:  "09:00",
			wantErr: true,
		},
		{
			name:    "Invalid - hours out of range",
			offset:  "+25:00",
			wantErr: true,
		},
		{
			name:    "Invalid - minutes out of range",
			offset:  "+09:70",
			wantErr: true,
		},
		{
			name:    "Invalid - over maximum positive",
			offset:  "+14:01",
			wantErr: true,
		},
		{
			name:    "Invalid - over maximum negative",
			offset:  "-13:00",
			wantErr: true,
		},
		{
			name:    "Invalid - wrong format",
			offset:  "+9:0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimezoneOffset(tt.offset)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseTimezoneOffset() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseTimezoneOffset() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTimezoneOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLI_Run_WithTimezoneOffset(t *testing.T) {
	tests := []struct {
		name           string
		timezoneOffset int
		wantPrefix     string
	}{
		{
			name:           "UTC timezone",
			timezoneOffset: 0,
			wantPrefix:     "20210101-000000", // 1609459200 in UTC
		},
		{
			name:           "JST timezone (+09:00)",
			timezoneOffset: 9 * 3600,
			wantPrefix:     "20210101-090000", // 1609459200 + 9 hours
		},
		{
			name:           "EST timezone (-05:00)",
			timezoneOffset: -5 * 3600,
			wantPrefix:     "20201231-190000", // 1609459200 - 5 hours
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

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

			cli := New(nil)

			err := cli.Run(inputPath, "", "", "Test Track", tt.timezoneOffset)
			if err != nil {
				t.Errorf("Run() unexpected error = %v", err)
				return
			}

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

			expectedFilename := tt.wantPrefix + ".gpx"
			if gpxFile != expectedFilename {
				t.Errorf("Generated filename = %v, want %v", gpxFile, expectedFilename)
			}
		})
	}
}

func TestValidateOutputPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid absolute path",
			path:    "/tmp/output",
			wantErr: false,
		},
		{
			name:    "valid relative path",
			path:    "output/gpx",
			wantErr: false,
		},
		{
			name:    "valid nested relative path",
			path:    "./output/2024/october",
			wantErr: false,
		},
		{
			name:    "path traversal with ..",
			path:    "../../../etc/passwd",
			wantErr: true,
			errMsg:  "invalid relative path components",
		},
		{
			name:    "path traversal in middle",
			path:    "output/../../etc",
			wantErr: true,
			errMsg:  "invalid relative path components",
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errMsg:  "output path is empty",
		},
		{
			name:    "current directory",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "parent directory only",
			path:    "..",
			wantErr: true,
			errMsg:  "invalid relative path components",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateOutputPath(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateOutputPath() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateOutputPath() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("validateOutputPath() unexpected error = %v", err)
				return
			}
			if !filepath.IsAbs(got) {
				t.Errorf("validateOutputPath() returned non-absolute path = %v", got)
			}
		})
	}
}

func TestCLI_Run_PathTraversalPrevention(t *testing.T) {
	tests := []struct {
		name       string
		outputDir  string
		outputFile string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "safe output directory",
			outputDir:  "safe-output",
			outputFile: "",
			wantErr:    false,
		},
		{
			name:       "path traversal in output directory",
			outputDir:  "../../../etc",
			outputFile: "",
			wantErr:    true,
			errMsg:     "invalid output directory",
		},
		{
			name:       "safe absolute output file",
			outputDir:  "",
			outputFile: "/tmp/safe-output.gpx",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			jsonContent := `[
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
			]`
			inputPath := filepath.Join(tmpDir, "test.json")
			if err := os.WriteFile(inputPath, []byte(jsonContent), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			cli := New(nil)

			outputFile := tt.outputFile
			if tt.outputFile != "" && !filepath.IsAbs(tt.outputFile) {
				outputFile = filepath.Join(tmpDir, tt.outputFile)
			}

			err := cli.Run(inputPath, outputFile, tt.outputDir, "Test Track", 0)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Run() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Run() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Run() unexpected error = %v", err)
			}
		})
	}
}

func TestParseTimezoneOffset_ErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		offset        string
		wantErrSubstr string
	}{
		{
			name:          "error message includes format for no sign",
			offset:        "09:00",
			wantErrSubstr: "±HH:MM or ±HHMM",
		},
		{
			name:          "error message includes format for wrong length",
			offset:        "+9:0",
			wantErrSubstr: "2-digit hours and minutes",
		},
		{
			name:          "error message includes range for invalid hours",
			offset:        "+25:00",
			wantErrSubstr: "0-14",
		},
		{
			name:          "error message includes range for invalid minutes",
			offset:        "+09:70",
			wantErrSubstr: "0-59",
		},
		{
			name:          "error message mentions maximum for over limit",
			offset:        "+14:01",
			wantErrSubstr: "+14:00",
		},
		{
			name:          "error message for short input",
			offset:        "+9",
			wantErrSubstr: "±HH:MM or ±HHMM",
		},
		{
			name:          "error message for invalid format without colon",
			offset:        "+090",
			wantErrSubstr: "4 digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTimezoneOffset(tt.offset)
			if err == nil {
				t.Errorf("ParseTimezoneOffset() error = nil, want error")
				return
			}
			if !strings.Contains(err.Error(), tt.wantErrSubstr) {
				t.Errorf("ParseTimezoneOffset() error = %v, want error containing %q", err, tt.wantErrSubstr)
			}
		})
	}
}
