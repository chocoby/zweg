package models

import (
	"fmt"
	"strconv"
	"time"
)

// Means represents a means of transportation recorded in the ZweiteGPS data.
type Means int

const (
	MeansWalking Means = iota
	MeansJogging
	MeansBicycle
	MeansMotorCycle
	MeansAutoMobile
	MeansTrain
	MeansMisc
)

var meansNames = [...]string{"Walking", "Jogging", "Bicycle", "MotorCycle", "AutoMobile", "Train", "Misc"}

// String returns the English name of the means of transportation.
// Unknown values return an empty string.
func (m Means) String() string {
	if m < 0 || int(m) >= len(meansNames) {
		return ""
	}
	return meansNames[m]
}

// Point represents a single GPS point from ZweiteGPS JSON format.
// Field semantics follow the official ZweiteGPS JSON specification.
type Point struct {
	Tm int64   `json:"tm"` // Unix timestamp
	Lo float64 `json:"lo"` // Longitude
	La float64 `json:"la"` // Latitude
	Al string  `json:"al"` // Altitude
	Sp string  `json:"sp"` // Speed
	Co int     `json:"co"` // Course (true bearing of motion)
	Th int     `json:"th"` // True heading
	He int     `json:"he"` // Magnetic heading
	Ds string  `json:"ds"` // Distance

	Ap float64 `json:"ap,omitempty"` // Atmospheric pressure
	Dp string  `json:"dp,omitempty"` // Description / memo
	Ha float64 `json:"ha,omitempty"` // Horizontal accuracy
	Ms *Means  `json:"ms,omitempty"` // Means of transportation (nil when absent in JSON)
	Ow string  `json:"ow,omitempty"` // Owner / device info
	Ra float64 `json:"ra,omitempty"` // Relative altitude
	Tl string  `json:"tl,omitempty"` // Title (log name)
	Va float64 `json:"va,omitempty"` // Vertical accuracy
	Ws int     `json:"ws,omitempty"` // Number of steps
	Xa float64 `json:"xa,omitempty"` // Heading accuracy

	Gx float64 `json:"gx,omitempty"` // Gravity acceleration X
	Gy float64 `json:"gy,omitempty"` // Gravity acceleration Y
	Gz float64 `json:"gz,omitempty"` // Gravity acceleration Z
	Ax float64 `json:"ax,omitempty"` // User acceleration X
	Ay float64 `json:"ay,omitempty"` // User acceleration Y
	Az float64 `json:"az,omitempty"` // User acceleration Z
	Ep float64 `json:"ep,omitempty"` // Pitch angle
	Er float64 `json:"er,omitempty"` // Roll angle
	Ey float64 `json:"ey,omitempty"` // Yaw angle
	Pf float64 `json:"pf,omitempty"` // Peak frequency
}

// TimestampIn returns the time.Time representation of the Unix timestamp in the given location.
// Callers should pass time.UTC for GPX-spec output, or a *time.Location built via
// time.FixedZone for filename generation in a specific offset.
func (p *Point) TimestampIn(loc *time.Location) time.Time {
	return time.Unix(p.Tm, 0).In(loc)
}

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

// FirstTitle returns the first non-empty Tl (log title) found in the slice.
// Returns an empty string if no point has a title.
func FirstTitle(points []Point) string {
	for _, p := range points {
		if p.Tl != "" {
			return p.Tl
		}
	}
	return ""
}

// FirstMeans returns the first recorded means of transportation in the slice.
// The bool is false when no point has Ms set.
func FirstMeans(points []Point) (Means, bool) {
	for _, p := range points {
		if p.Ms != nil {
			return *p.Ms, true
		}
	}
	return 0, false
}
