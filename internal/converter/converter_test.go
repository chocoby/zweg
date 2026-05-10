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

func TestGPXConverter_Convert_AccuracyFields(t *testing.T) {
	points := []models.Point{
		{Tm: 1609459200, Lo: 139.7671, La: 35.6812, Al: "10.5", Ha: 5.0, Va: 3.0},
		{Tm: 1609459260, Lo: 139.7672, La: 35.6813, Al: "11.0"}, // no accuracy
	}

	g, err := New(nil).Convert(points, "Acc Track")
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}

	trkPts := g.Trk[0].TrkSeg[0].TrkPt
	if got := trkPts[0].HDOP; got != 5.0 {
		t.Errorf("trkpt[0].HDOP = %v, want 5.0", got)
	}
	if got := trkPts[0].VDOP; got != 3.0 {
		t.Errorf("trkpt[0].VDOP = %v, want 3.0", got)
	}
	if got := trkPts[1].HDOP; got != 0 {
		t.Errorf("trkpt[1].HDOP = %v, want 0 (omitted)", got)
	}
	if got := trkPts[1].VDOP; got != 0 {
		t.Errorf("trkpt[1].VDOP = %v, want 0 (omitted)", got)
	}
}

func TestGPXConverter_Convert_DescriptionFromDp(t *testing.T) {
	points := []models.Point{
		{Tm: 1609459200, Lo: 139.7671, La: 35.6812, Al: "10.5", Dp: "start memo"},
		{Tm: 1609459260, Lo: 139.7672, La: 35.6813, Al: "11.0"},
		{Tm: 1609459320, Lo: 139.7673, La: 35.6814, Al: "11.5", Dp: "finish memo"},
	}

	g, err := New(nil).Convert(points, "Desc Track")
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}

	trkPts := g.Trk[0].TrkSeg[0].TrkPt
	if got, want := trkPts[0].Desc, "start memo"; got != want {
		t.Errorf("trkpt[0].Desc = %q, want %q", got, want)
	}
	if got := trkPts[1].Desc; got != "" {
		t.Errorf("trkpt[1].Desc = %q, want empty", got)
	}
	if got, want := trkPts[2].Desc, "finish memo"; got != want {
		t.Errorf("trkpt[2].Desc = %q, want %q", got, want)
	}

	if len(g.Wpt) != 2 {
		t.Fatalf("Wpt count = %d, want 2", len(g.Wpt))
	}
	if got, want := g.Wpt[0].Desc, "start memo"; got != want {
		t.Errorf("Start waypoint Desc = %q, want %q", got, want)
	}
	if got, want := g.Wpt[1].Desc, "finish memo"; got != want {
		t.Errorf("Goal waypoint Desc = %q, want %q", got, want)
	}
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
