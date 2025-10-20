package converter

import (
	"testing"

	"github.com/chocoby/zweg/internal/models"
)

func TestGPXConverter_Convert(t *testing.T) {
	tests := []struct {
		name      string
		points    []models.Point
		trackName string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid single point",
			points: []models.Point{
				{
					Tm: 1609459200,
					Lo: 139.7671,
					La: 35.6812,
					Al: "10.5",
					Sp: "5.0",
					Ds: "0",
				},
			},
			trackName: "Test Track",
			wantErr:   false,
		},
		{
			name: "valid multiple points",
			points: []models.Point{
				{
					Tm: 1609459200,
					Lo: 139.7671,
					La: 35.6812,
					Al: "10.5",
					Sp: "5.0",
					Ds: "0",
				},
				{
					Tm: 1609459260,
					Lo: 139.7672,
					La: 35.6813,
					Al: "11.2",
					Sp: "5.5",
					Ds: "100",
				},
			},
			trackName: "Multi Point Track",
			wantErr:   false,
		},
		{
			name:      "empty points",
			points:    []models.Point{},
			trackName: "Empty Track",
			wantErr:   true,
			errMsg:    "no data points provided",
		},
		{
			name: "invalid altitude",
			points: []models.Point{
				{
					Tm: 1609459200,
					Lo: 139.7671,
					La: 35.6812,
					Al: "invalid",
					Sp: "5.0",
					Ds: "0",
				},
			},
			trackName: "Invalid Alt Track",
			wantErr:   true,
			errMsg:    "failed to parse altitude",
		},
		{
			name: "empty track name defaults to Track",
			points: []models.Point{
				{
					Tm: 1609459200,
					Lo: 139.7671,
					La: 35.6812,
					Al: "10.5",
				},
			},
			trackName: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(nil)
			gpx, err := c.Convert(tt.points, tt.trackName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Convert() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() == "" {
					t.Errorf("Convert() error message is empty, want %q", tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Convert() unexpected error = %v", err)
				return
			}

			if gpx == nil {
				t.Error("Convert() returned nil GPX")
				return
			}

			// Verify GPX structure
			if gpx.Version != "1.1" {
				t.Errorf("GPX version = %q, want %q", gpx.Version, "1.1")
			}

			if gpx.Metadata == nil {
				t.Error("GPX metadata is nil")
				return
			}

			expectedName := tt.trackName
			if expectedName == "" {
				expectedName = "Track"
			}
			if gpx.Metadata.Name != expectedName {
				t.Errorf("GPX metadata name = %q, want %q", gpx.Metadata.Name, expectedName)
			}

			if len(gpx.Trk) != 1 {
				t.Errorf("GPX tracks count = %d, want 1", len(gpx.Trk))
				return
			}

			if len(gpx.Trk[0].TrkSeg) != 1 {
				t.Errorf("GPX track segments count = %d, want 1", len(gpx.Trk[0].TrkSeg))
				return
			}

			if len(gpx.Trk[0].TrkSeg[0].TrkPt) != len(tt.points) {
				t.Errorf("GPX track points count = %d, want %d", len(gpx.Trk[0].TrkSeg[0].TrkPt), len(tt.points))
			}
		})
	}
}

func TestGPXConverter_Convert_WithConfig(t *testing.T) {
	points := []models.Point{
		{
			Tm: 1609459200,
			Lo: 139.7671,
			La: 35.6812,
			Al: "10.5",
		},
	}

	t.Run("with waypoints enabled", func(t *testing.T) {
		config := &Config{
			IncludeWaypoint: true,
		}
		c := New(config)
		gpx, err := c.Convert(points, "Test")

		if err != nil {
			t.Fatalf("Convert() unexpected error = %v", err)
		}

		if len(gpx.Wpt) != 2 {
			t.Errorf("Waypoints count = %d, want 2 (start and goal)", len(gpx.Wpt))
		}
	})

	t.Run("with waypoints disabled", func(t *testing.T) {
		config := &Config{
			IncludeWaypoint: false,
		}
		c := New(config)
		gpx, err := c.Convert(points, "Test")

		if err != nil {
			t.Fatalf("Convert() unexpected error = %v", err)
		}

		if len(gpx.Wpt) != 0 {
			t.Errorf("Waypoints count = %d, want 0", len(gpx.Wpt))
		}
	})

	t.Run("with custom creator", func(t *testing.T) {
		config := &Config{
			Creator: "Custom Creator",
		}
		c := New(config)
		gpx, err := c.Convert(points, "Test")

		if err != nil {
			t.Fatalf("Convert() unexpected error = %v", err)
		}

		if gpx.Creator != "Custom Creator" {
			t.Errorf("Creator = %q, want %q", gpx.Creator, "Custom Creator")
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Version != "1.1" {
		t.Errorf("Default version = %q, want %q", config.Version, "1.1")
	}

	if config.Creator == "" {
		t.Error("Default creator is empty")
	}

	if !config.IncludeWaypoint {
		t.Error("Default IncludeWaypoint = false, want true")
	}
}
