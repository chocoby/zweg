package models

import (
	"fmt"
	"strconv"
	"time"
)

// Point represents a single GPS point from ZweiteGPS JSON format
type Point struct {
	Tm int64   `json:"tm"`           // Unix timestamp
	Lo float64 `json:"lo"`           // Longitude
	La float64 `json:"la"`           // Latitude
	Th int     `json:"th"`           // True heading
	Sp string  `json:"sp"`           // Speed
	Co int     `json:"co"`           // Course
	Al string  `json:"al"`           // Altitude
	He int     `json:"he"`           // Heading
	Ds string  `json:"ds"`           // Distance
	Ms int     `json:"ms,omitempty"` // Optional field
	Ow string  `json:"ow,omitempty"` // Owner/device info
}

// Timestamp returns the time.Time representation of the Unix timestamp in UTC.
func (p *Point) Timestamp() time.Time {
	return time.Unix(p.Tm, 0).UTC()
}

// LocalTimestamp returns the time.Time representation of the Unix timestamp in local timezone.
func (p *Point) LocalTimestamp() time.Time {
	return time.Unix(p.Tm, 0).Local()
}

// Altitude returns the altitude as a float64 value.
func (p *Point) Altitude() (float64, error) {
	if p.Al == "" {
		return 0, nil
	}
	alt, err := strconv.ParseFloat(p.Al, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse altitude %q: %w", p.Al, err)
	}
	return alt, nil
}

// Speed returns the speed as a float64 value.
func (p *Point) Speed() (float64, error) {
	if p.Sp == "" {
		return 0, nil
	}
	speed, err := strconv.ParseFloat(p.Sp, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse speed %q: %w", p.Sp, err)
	}
	return speed, nil
}

// Distance returns the distance as a float64 value.
func (p *Point) Distance() (float64, error) {
	if p.Ds == "" {
		return 0, nil
	}
	distance, err := strconv.ParseFloat(p.Ds, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse distance %q: %w", p.Ds, err)
	}
	return distance, nil
}
