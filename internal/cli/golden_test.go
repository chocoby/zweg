package cli

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "regenerate golden files under testdata/golden/")

func TestCLI_Run_Golden(t *testing.T) {
	cases := []struct {
		name       string
		inputFile  string
		goldenFile string
		trackName  string
	}{
		{
			name:       "single point",
			inputFile:  "single_point.json",
			goldenFile: "single_point.gpx",
			trackName:  "Single Point Test",
		},
		{
			// empty trackName exercises the tl → means → "Track" fallback chain.
			name:       "multi point with optional fields",
			inputFile:  "multi_point.json",
			goldenFile: "multi_point.gpx",
			trackName:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inputPath := filepath.Join("testdata", "input", tc.inputFile)
			goldenPath := filepath.Join("testdata", "golden", tc.goldenFile)

			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "out.gpx")

			if err := New(nil).Run(inputPath, outputPath, "", tc.trackName, 0); err != nil {
				t.Fatalf("Run: %v", err)
			}

			got, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("read output: %v", err)
			}

			if *updateGolden {
				if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
					t.Fatalf("update golden %s: %v", goldenPath, err)
				}
				t.Logf("updated golden: %s", goldenPath)
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden %s: %v", goldenPath, err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("output differs from %s\n\n--- got ---\n%s\n\n--- want ---\n%s", goldenPath, got, want)
			}
		})
	}
}
