package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// readIn creates a FileInput object
func readIn(filename string, wm *warehouseMap) (err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Error().Err(err).Stack().
			Str("Filename", filename).
			Msg("Failed to open file")

		return err
	}
	wm.data = strings.Split(string(data), "\n")

	return nil
}
