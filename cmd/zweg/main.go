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
	outputDir := flag.String("d", "", "Output directory (ignored if output file is specified)")
	flag.StringVar(outputDir, "output-dir", "", "Output directory (ignored if output file is specified)")
	timezoneOffsetStr := flag.String("timezone-offset", "+00:00", "Timezone offset for GPX timestamps (e.g., +09:00, -05:00)")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <input.json> [output.gpx]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Convert ZweiteGPS JSON format to GPX format.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  input.json    Input file in ZweiteGPS JSON format\n")
		fmt.Fprintf(os.Stderr, "  output.gpx    Output file in GPX format (optional, defaults to YYYYMMDD-HHMMSS.gpx based on track start time)\n")
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
	}

	// Parse timezone offset for filename generation
	timezoneOffset, err := cli.ParseTimezoneOffset(*timezoneOffsetStr)
	if err != nil {
		return fmt.Errorf("invalid timezone offset: %w", err)
	}

	c := cli.New(&cli.Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	return c.Run(inputFile, outputFile, *outputDir, *trackName, timezoneOffset)
}
