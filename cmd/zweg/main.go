package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chocoby/zweg/pkg/cli"
)

const (
	exitFailure = 1
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitFailure)
	}
}

func run() error {
	// Define flags
	trackName := flag.String("track-name", "Track", "Name for the GPS track")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--track-name <name>] <input.json> <output.gpx>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Convert ZweiteGPS JSON format to GPX format.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  input.json   Input file in ZweiteGPS JSON format\n")
		fmt.Fprintf(os.Stderr, "  output.gpx   Output file in GPX format\n")
	}

	flag.Parse()

	// Check for required positional arguments
	if flag.NArg() != 2 {
		flag.Usage()
		return fmt.Errorf("exactly 2 arguments required (input and output files)")
	}

	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)

	// Create CLI instance
	c := cli.New(&cli.Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	// Run conversion
	return c.Run(inputFile, outputFile, *trackName)
}
