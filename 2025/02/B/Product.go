package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type products struct {
	data   []string
	badIDs []uint64
}

// accumulateBadIDs tallies the bad product IDs and returns the total
func (p *products) accumulateBadIDs() (tally uint64) {
	for index := range p.badIDs {
		tally += p.badIDs[index]
	}

	return tally
}

// check_for_pattern examines a single product ID for a pattern
func (p *products) check_for_pattern(pID uint64) {
	var pidIsValid bool = true
	checkStr := fmt.Sprintf("%d", pID)
	pidLength := len(checkStr)
	sublog := log.With().Uint64("Product ID", pID).Int("pid length", pidLength).Logger()

	// Patterns can be as long as half the length of the pid
outer:
	for patternLength := 1; patternLength <= pidLength/2; patternLength++ {
		if pidLength%patternLength > 0 {
			sublog.Debug().
				Int("pattern length", patternLength).
				Msg("Test pattern won't fit")

			continue
		}
		sublog.Debug().Int("Test length", patternLength).Msg("Pattern length to test")

		// 2121212124 - 10 spaces
		// 21 - test pattern - pattern length == 2
		// 5x copies of the pattern
		// 4x test cases

		testCases := pidLength / patternLength
		// Split the product ID into test pieces
		pieces := make([]string, testCases)
		for index := range testCases {
			pieces[index] = checkStr[patternLength*index : (patternLength*index)+patternLength]
		}
		sublog.Debug().Any("Test Pieces", pieces).Msg("Split product ID into test pieces")

		for index := range testCases {
			if pieces[index] != pieces[0] {
				sublog.Debug().
					Str("Zero", pieces[0]).
					Int("Index", index).
					Str("Index String", pieces[index]).
					Msg("Doesn't match")

				continue outer
			}
		}

		pidIsValid = false
		break outer
	}

	if pidIsValid {
		sublog.Debug().Msg("Product ID is valid")
	} else {
		sublog.Info().Msg("Invalid Product ID found")
		p.badIDs = append(p.badIDs, pID)
	}
}

func (p *products) check_range(index int) {
	if p.data[index] == "" {
		log.Debug().Msg("Empty data")

		return
	}

	var rangeStart, rangeEnd uint64
	var err error

	rangeString := strings.Split(p.data[index], "-")
	if len(rangeString) != 2 {
		log.Error().Err(err).
			Int("Range size", len(rangeString)).
			Msg("Unexpecte range syntax, should have two numbers")
	}

	rangeStart, err = strconv.ParseUint(strings.TrimSpace(rangeString[0]), 10, 64)
	if err != nil {
		log.Error().Err(err).
			Str("Start #", rangeString[0]).
			Msg("Failed to parse start range number")
	}

	rangeEnd, err = strconv.ParseUint(strings.TrimSpace(rangeString[1]), 10, 64)
	if err != nil {
		log.Error().Err(err).
			Str("End #", rangeString[1]).
			Msg("Failed to parse end range number")
	}

	for pID := rangeStart; pID <= rangeEnd; pID++ {
		log.Debug().
			Uint64("Current", pID).
			Msg("Checking product ID for patterns")

		p.check_for_pattern(pID)
	}
}

func (p *products) process() {
	log.Debug().
		Int("Range Quantity", len(p.data)).
		Msg("Processing product ID ranges.")

	for index := range p.data {
		p.check_range(index)
	}

	log.Info().Msg("Processing complete")
}

func (p *products) report() {
	log.Info().
		Int("Invalid count", len(p.badIDs)).
		Uint64("Tally", p.accumulateBadIDs()).
		Msg("")
}
