package main

import (
	"fmt"
	"os"

	"github.com/chocoby/zweg/pkg/cli"
)

const (
	usageMessage = "Usage: zweg <input.json> <output.gpx> [track_name]"
	exitFailure  = 1
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitFailure)
	}
}

func run() error {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, usageMessage)
		return fmt.Errorf("insufficient arguments")
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	trackName := "Track"
	if len(os.Args) >= 4 {
		trackName = os.Args[3]
	}

	// Create CLI instance
	c := cli.New(&cli.Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	// Run conversion
	return c.Run(inputFile, outputFile, trackName)
}
