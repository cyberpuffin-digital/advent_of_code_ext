package main

import (
	"flag"
	"fmt"
)

var example bool
var verbose bool

func main() {
	parseCLI()
	var exampleInput, challengeInput fileInput

	exampleInput = readIn("example.in")
	fmt.Printf("%q adds up to %d.\n", exampleInput.filename, exampleInput.tally())
	exampleInput.cleanup()

	if !example {
		challengeInput = readIn("input.txt")
		fmt.Printf("%q adds up to %d.\n", challengeInput.filename, challengeInput.tally())
		challengeInput.cleanup()
	}
}

func parseCLI() {
	flag.BoolVar(&example, "example", false, "Only run on example input.")
	flag.BoolVar(&verbose, "verbose", false, "Increase program output.")
	flag.Parse()
}

func verbosef(msg string, args ...interface{}) {
	if verbose {
		fmt.Printf(msg, args...)
	}
}
