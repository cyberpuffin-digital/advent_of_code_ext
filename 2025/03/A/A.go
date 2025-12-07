package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
)

var debug, example, trace *bool

func main() {
	var exampleInput, challengeInput *batteryBanks
	var err error

	parseFlags()
	setupLog(*debug, *trace)

	exampleInput = new(batteryBanks)
	err = readIn("example.in", exampleInput)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("unable to open example input")

		os.Exit(1)
	}
	exampleInput.process()
	exampleInput.report()

	// Exit early if only testing the example
	if *example {
		os.Exit(0)
	}

	challengeInput = new(batteryBanks)
	err = readIn("input.txt", challengeInput)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("unable to open challenge input")

		os.Exit(1)
	}
	challengeInput.process()
	challengeInput.report()
}

func parseFlags() {
	debug = flag.Bool("debug", false, "set log level to debug")
	example = flag.Bool("example", false, "process the example only")
	trace = flag.Bool("trace", false, "set log level to trace")

	flag.Parse()
}
