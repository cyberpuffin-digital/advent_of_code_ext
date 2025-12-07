package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
)

var debug, example *bool

func main() {
	var exampleInput, challengeInput *doorLock
	var err error

	parseFlags()
	setupLog(*debug)

	exampleInput = new(doorLock)
	err = readIn("example.in", exampleInput)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("unable to open example input")

		os.Exit(1)
	}
	exampleInput.process()

	// Exit early if only testing the example
	if *example {
		os.Exit(0)
	}

	challengeInput = new(doorLock)
	err = readIn("input.txt", challengeInput)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("unable to open challenge input")

		os.Exit(1)
	}
	challengeInput.process()
}

func parseFlags() {
	debug = flag.Bool("debug", false, "set log level to debug")
	example = flag.Bool("example", false, "process the example only")

	flag.Parse()
}
