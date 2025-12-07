package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func setupLog(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Str("Puzzle", "B").
		Logger()

	log.Info().
		Any("Level", log.Logger.GetLevel()).
		Msg("Logging setup")
}
