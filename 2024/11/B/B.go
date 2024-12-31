package main

import (
	"log"
	"os"
	"runtime/pprof"
)

var cpuProfile string
var example bool
var verbose bool

func main() {
	parseCLI()
	var exampleInput, challengeInput fileInput

	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatalf("Failed to create CPU Profile file: %s", cpuProfile)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	exampleInput = readIn("example.in")
	exampleInput.report()
	exampleInput.cleanup()

	if !example {
		challengeInput = readIn("input.txt")
		challengeInput.report()
		challengeInput.cleanup()
	}
}
