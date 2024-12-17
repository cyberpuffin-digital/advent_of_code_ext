package main

import (
	"flag"
	"fmt"
)

func parseCLI() {
	flag.StringVar(&cpuProfile, "cpuprofile", "", "write cpu profile to file.")
	flag.BoolVar(&example, "example", false, "Only run on example input.")
	flag.BoolVar(&verbose, "verbose", false, "Increase program output.")
	flag.Parse()
}

func printOut(msg string) {
	fmt.Println(msg)
}

func verbosef(msg string, args ...interface{}) {
	if verbose {
		fmt.Printf(msg, args...)
	}
}
