package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/twpayne/go-gpx"
)

// ZweiteGPSPoint represents a single GPS point from ZweiteGPS JSON
type ZweiteGPSPoint struct {
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

func convertJSONtoGPX(jsonData []ZweiteGPSPoint, trackName string) (*gpx.GPX, error) {
	if len(jsonData) == 0 {
		return nil, fmt.Errorf("no data points provided")
	}

	// Create GPX document
	g := &gpx.GPX{
		Version: "1.1",
		Creator: "zweg - ZweiteGPS to GPX Converter",
	}

	// Set metadata
	startTime := time.Unix(jsonData[0].Tm, 0)
	g.Metadata = &gpx.MetadataType{
		Name: trackName,
		Time: startTime,
	}

	// Add start waypoint
	firstAlt := parseFloat(jsonData[0].Al)
	g.Wpt = append(g.Wpt, &gpx.WptType{
		Lat:  jsonData[0].La,
		Lon:  jsonData[0].Lo,
		Ele:  firstAlt,
		Time: startTime,
		Name: "Start",
	})

	// Add goal waypoint
	lastAlt := parseFloat(jsonData[len(jsonData)-1].Al)
	endTime := time.Unix(jsonData[len(jsonData)-1].Tm, 0)
	g.Wpt = append(g.Wpt, &gpx.WptType{
		Lat:  jsonData[len(jsonData)-1].La,
		Lon:  jsonData[len(jsonData)-1].Lo,
		Ele:  lastAlt,
		Time: endTime,
		Name: "Goal",
	})

	// Create track
	track := &gpx.TrkType{
		Name: trackName,
	}

	// Create track segment
	segment := &gpx.TrkSegType{}

	// Add track points
	for _, point := range jsonData {
		alt := parseFloat(point.Al)
		timestamp := time.Unix(point.Tm, 0)

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

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gpx-converter <input.json> <output.gpx> [track_name]")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	trackName := "Track"
	if len(os.Args) >= 4 {
		trackName = os.Args[3]
	}

	// Read JSON file
	jsonBytes, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON
	var points []ZweiteGPSPoint
	if err := json.Unmarshal(jsonBytes, &points); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Convert to GPX
	g, err := convertJSONtoGPX(points, trackName)
	if err != nil {
		fmt.Printf("Error converting to GPX: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating GPX file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := g.WriteIndent(file, "", "  "); err != nil {
		fmt.Printf("Error writing GPX XML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %d points to GPX: %s\n", len(points), outputFile)
}
