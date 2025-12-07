package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// readIn creates a FileInput object
func readIn(filename string, dObj *batteryBanks) (err error) {
	log.Logger = log.With().Str("Filename", filename).Stack().Logger()
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file")

		return err
	}
	// dObj.data = strings.Split(string(data), "\n")
	dObj.data = strings.FieldsFunc(string(data), func(c rune) bool {
		return c == '\n'
	})

	log.Trace().Any("Data", dObj.data).Msg("Data read in")

	return nil
}
