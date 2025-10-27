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

var (
	// Version information - set via ldflags during build
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitFailure)
	}
}

func run() error {
	trackName := flag.String("track-name", "Track", "Name for the GPS track")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--track-name <name>] <input.json> [output.gpx]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Convert ZweiteGPS JSON format to GPX format.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  input.json    Input file in ZweiteGPS JSON format\n")
		fmt.Fprintf(os.Stderr, "  output.gpx    Output file in GPX format (optional, defaults to input.json.gpx)\n")
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("zweg version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
		return nil
	}

	nArgs := flag.NArg()
	if nArgs < 1 || nArgs > 2 {
		flag.Usage()
		return fmt.Errorf("1 or 2 arguments required (input file and optional output file)")
	}

	inputFile := flag.Arg(0)
	outputFile := ""
	if nArgs == 2 {
		outputFile = flag.Arg(1)
	} else {
		outputFile = inputFile + ".gpx"
	}

	c := cli.New(&cli.Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	return c.Run(inputFile, outputFile, *trackName)
}
