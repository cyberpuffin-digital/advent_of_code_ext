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
	checkStr := fmt.Sprintf("%d", pID)
	strLength := len(checkStr)
	if strLength%2 != 0 {
		log.Debug().
			Uint64("Product ID", pID).
			Msg("Odd length product ID, not evenly divisible.")

		return
	}

	checkA, checkB := checkStr[0:(strLength/2)], checkStr[(strLength/2):]
	log.Debug().
		Uint64("Product ID", pID).
		Str("Check A", checkA).
		Str("Check B", checkB).
		Msg("Checking pattern")

	if checkA == checkB {
		log.Info().Uint64("ProductID", pID).Msg("Invalid pid found")
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
