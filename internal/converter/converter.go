package converter

import (
	"fmt"
	"time"

	"github.com/chocoby/zweg/internal/models"
	"github.com/twpayne/go-gpx"
)

// Converter defines the interface for converting GPS data to GPX format.
type Converter interface {
	Convert(points []models.Point, trackName string) (*gpx.GPX, error)
}

// Config holds configuration for GPX conversion.
type Config struct {
	Version         string
	Creator         string
	IncludeWaypoint bool
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Version:         "1.1",
		Creator:         "zweg - ZweiteGPS to GPX Converter",
		IncludeWaypoint: true,
	}
}

// GPXConverter implements the Converter interface.
type GPXConverter struct {
	config *Config
}

// New creates a new GPXConverter with the given configuration.
func New(config *Config) *GPXConverter {
	if config == nil {
		config = DefaultConfig()
	}
	return &GPXConverter{
		config: config,
	}
}

// Convert converts ZweiteGPS points to GPX format.
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

	startTime := points[0].TimestampIn(time.UTC)
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

		timestamp := point.TimestampIn(time.UTC)

		segment.TrkPt = append(segment.TrkPt, &gpx.WptType{
			Lat:  point.La,
			Lon:  point.Lo,
			Ele:  alt,
			Time: timestamp,
			Desc: point.Dp,
			HDOP: point.Ha,
			VDOP: point.Va,
		})
	}

	track.TrkSeg = append(track.TrkSeg, segment)
	g.Trk = append(g.Trk, track)

	return g, nil
}

// addWaypoints adds start and end waypoints to the GPX document.
func (c *GPXConverter) addWaypoints(g *gpx.GPX, points []models.Point) error {
	start, err := waypointFrom(points[0], "Start")
	if err != nil {
		return err
	}
	goal, err := waypointFrom(points[len(points)-1], "Goal")
	if err != nil {
		return err
	}
	g.Wpt = append(g.Wpt, start, goal)
	return nil
}

func waypointFrom(p models.Point, name string) (*gpx.WptType, error) {
	alt, err := p.Altitude()
	if err != nil {
		return nil, fmt.Errorf("waypoint %q altitude: %w", name, err)
	}
	return &gpx.WptType{
		Lat:  p.La,
		Lon:  p.Lo,
		Ele:  alt,
		Time: p.TimestampIn(time.UTC),
		Name: name,
		Desc: p.Dp,
	}, nil
}
