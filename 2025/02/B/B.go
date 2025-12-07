package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
)

var debug, example *bool

func main() {
	var exampleInput, challengeInput *products
	var err error

	parseFlags()
	setupLog(*debug)

	log.Info().Msg("Reading example.in")
	exampleInput = new(products)
	err = readIn("example.in", exampleInput)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("unable to open example input")

		cleanup(1)
	}
	exampleInput.process()
	exampleInput.report()

	// Exit early if only testing the example
	if *example {
		cleanup(0)
	}

	log.Info().Msg("Reading input.txt")
	challengeInput = new(products)
	err = readIn("input.txt", challengeInput)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("unable to open challenge input")

		cleanup(1)
	}
	challengeInput.process()
	challengeInput.report()
}

func cleanup(exitCode int) {
	log.Info().Msg("Finished")
	os.Exit(exitCode)
}

func parseFlags() {
	debug = flag.Bool("debug", false, "set log level to debug")
	example = flag.Bool("example", false, "process the example only")

	flag.Parse()
}
