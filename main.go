package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/spf13/pflag"
)

// Function to handle quiet mode
func handleQuietMode(quietMode bool, line string) {
	if !quietMode {
		fmt.Println(line)
	}
}

// Function to handle dry-run mode (print to stdout, no file writing)
func handleDryRun(dryRun bool, fn string, f io.WriteCloser, line string) {
	if !dryRun && fn != "" {
		fmt.Fprintf(f, "%s\n", line)
	}
}

// Function to handle trimming of lines
func handleTrim(trim bool, line string) string {
	if trim {
		return strings.TrimSpace(line)
	}
	return line
}

// Function to handle wet-run (printing repeated lines)
func handleWetRun(wetRun bool, lines map[string]bool, line string) {
	if wetRun && lines[line] {
		fmt.Println(line)
	}
}

// Function to handle repeated lines
func handleRepeat(repeat bool, lines map[string]bool, line string) {
	if repeat && lines[line] {
		fmt.Println(line)
	}
}

func main() {
	// Define flags with pflag
	var quietMode bool
	var dryRun bool
	var wetRun bool
	var trim bool
	var repeat bool

	pflag.BoolVarP(&quietMode, "quiet", "q", false, "quiet mode (no output at all)")
	pflag.BoolVarP(&dryRun, "dry-run", "d", false, "don't append anything to the file, just print the new lines to stdout")
	pflag.BoolVarP(&wetRun, "wet-run", "s", false, "don't append anything to the file, just print the repeated lines to stdout")
	pflag.BoolVarP(&trim, "trim", "t", false, "trim leading and trailing whitespace before comparison")
	pflag.BoolVarP(&repeat, "repeat", "r", false, "show the repeated lines")

	// Set a custom Usage function to replace the default flag usage
	pflag.Usage = func() {
		fmt.Println("Usage: ain [OPTIONS] <FILE>")
		fmt.Println("Options:")
		pflag.PrintDefaults() // This prints the flags as usual
	}

	// Parse the flags using pflag
	pflag.Parse()

	// Retrieve the first non-flag argument (likely the filename)
	fn := pflag.Arg(0)

	// Initialize a map to track lines (for uniqueness)
	lines := make(map[string]bool)

	var f io.WriteCloser

	if fn != "" {
		// Read the whole file into a map if it exists
		r, err := os.Open(fn)
		if err == nil {
			sc := bufio.NewScanner(r)

			for sc.Scan() {
				line := sc.Text()
				line = handleTrim(trim, line)
				lines[line] = true
			}
			r.Close()
		} else {
			fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
			return
		}

		if !dryRun {
			// Re-open the file for appending new content
			f, err = os.OpenFile(fn, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file for writing: %s\n", err)
				return
			}
			defer f.Close()
		}
	} else {
		fmt.Println("you need to specify a file")
		return
	}

	// Read the lines from stdin and process them
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		line = handleTrim(trim, line)

		// Skip duplicate lines
		if lines[line] {
			handleRepeat(repeat, lines, line)
			handleWetRun(wetRun, lines, line)
			continue
		}

		// Add the line to the map to prevent duplicates
		lines[line] = true

		// Handle quiet mode (only print if not quiet)
		handleQuietMode(quietMode, line)

		// Handle dry run (print to stdout but don't write to the file)
		handleDryRun(dryRun, fn, f, line)

		// Handle wet-run and repeat logic
		handleWetRun(wetRun, lines, line)
		handleRepeat(repeat, lines, line)
	}
}
