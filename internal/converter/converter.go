package converter

import (
	"fmt"

	"github.com/chocoby/zweg/internal/models"
	"github.com/twpayne/go-gpx"
)

// Converter defines the interface for converting GPS data to GPX format
type Converter interface {
	Convert(points []models.Point, trackName string) (*gpx.GPX, error)
}

// Config holds configuration for GPX conversion
type Config struct {
	Version         string
	Creator         string
	IncludeWaypoint bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Version:         "1.1",
		Creator:         "zweg - ZweiteGPS to GPX Converter",
		IncludeWaypoint: true,
	}
}

// GPXConverter implements the Converter interface
type GPXConverter struct {
	config *Config
}

// New creates a new GPXConverter with the given configuration
func New(config *Config) *GPXConverter {
	if config == nil {
		config = DefaultConfig()
	}
	return &GPXConverter{
		config: config,
	}
}

// Convert converts ZweiteGPS points to GPX format
func (c *GPXConverter) Convert(points []models.Point, trackName string) (*gpx.GPX, error) {
	if len(points) == 0 {
		return nil, fmt.Errorf("no data points provided")
	}

	if trackName == "" {
		trackName = "Track"
	}

	g := &gpx.GPX{
		Version: c.config.Version,
		Creator: c.config.Creator,
	}

	startTime := points[0].Timestamp()
	g.Metadata = &gpx.MetadataType{
		Name: trackName,
		Time: startTime,
	}

	if c.config.IncludeWaypoint {
		if err := c.addWaypoints(g, points); err != nil {
			return nil, fmt.Errorf("failed to add waypoints: %w", err)
		}
	}

	track := &gpx.TrkType{
		Name: trackName,
	}

	segment := &gpx.TrkSegType{}

	for i, point := range points {
		alt, err := point.Altitude()
		if err != nil {
			return nil, fmt.Errorf("failed to parse altitude at point %d: %w", i, err)
		}

		timestamp := point.Timestamp()

		segment.TrkPt = append(segment.TrkPt, &gpx.WptType{
			Lat:  point.La,
			Lon:  point.Lo,
			Ele:  alt,
			Time: timestamp,
		})
	}

	track.TrkSeg = append(track.TrkSeg, segment)
	g.Trk = append(g.Trk, track)

	return g, nil
}

// addWaypoints adds start and end waypoints to the GPX document
func (c *GPXConverter) addWaypoints(g *gpx.GPX, points []models.Point) error {
	firstPoint := points[0]
	firstAlt, err := firstPoint.Altitude()
	if err != nil {
		return fmt.Errorf("failed to parse start altitude: %w", err)
	}

	g.Wpt = append(g.Wpt, &gpx.WptType{
		Lat:  firstPoint.La,
		Lon:  firstPoint.Lo,
		Ele:  firstAlt,
		Time: firstPoint.Timestamp(),
		Name: "Start",
	})

	lastPoint := points[len(points)-1]
	lastAlt, err := lastPoint.Altitude()
	if err != nil {
		return fmt.Errorf("failed to parse end altitude: %w", err)
	}

	g.Wpt = append(g.Wpt, &gpx.WptType{
		Lat:  lastPoint.La,
		Lon:  lastPoint.Lo,
		Ele:  lastAlt,
		Time: lastPoint.Timestamp(),
		Name: "Goal",
	})

	return nil
}
