package main

import (
	"flag"
	"fmt"
)

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
