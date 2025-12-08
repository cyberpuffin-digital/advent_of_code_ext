package main

import (
	"bufio"
	"math/big"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// readIn creates a FileInput object
func readIn(filename string, inv *inventory) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Stack().
			Str("Filename", filename).
			Msg("Failed to open file")

		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputSlice := strings.Split(scanner.Text(), "-")
		switch len(inputSlice) {
		case 1:
			if strings.TrimSpace(inputSlice[0]) == "" {
				continue
			}

			ingredientID := new(big.Int)
			ingredientID.SetString(inputSlice[0], 10)
			inv.availableIDs = append(inv.availableIDs, ingredientID)
		case 2:
			inv.addRange(inputSlice[0], inputSlice[1])
		default:
			log.Error().Any("Input Slice", inputSlice).Msg("Input slice size is invalid")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).
			Str("Filename", filename).
			Msg("Scanning the input file experienced an error")
	}

	log.Info().
		Int("Ingredient Count", len(inv.availableIDs)).
		Int("Range Count", len(inv.ranges)).
		Msg("Finished scanning input file")

	return nil
}
