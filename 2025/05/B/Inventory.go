package main

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
)

type freshRange struct {
	end   *big.Int
	start *big.Int
}

type inventory struct {
	availableIDs []*big.Int
	ranges       []*freshRange
	uniqueRanges []*freshRange
}

// newFreshRange instantiates the big ints and returns a a freshRange struct
func newFreshRange(start, end string) (fr *freshRange) {
	fr = &freshRange{
		end:   new(big.Int),
		start: new(big.Int),
	}

	_, ok := fr.end.SetString(end, 10)
	if !ok {
		log.Error().Str("End range", end).Msg("Failed to set Fresh range end value")
	}

	_, ok = fr.start.SetString(start, 10)
	if !ok {
		log.Error().Str("Start range", start).Msg("Failed to set Fresh range start value")
	}

	return fr
}

// addRange adds a fresh ingredients range of product IDs to the inventory
func (i *inventory) addRange(rangeStart, rangeEnd string) {
	fr := newFreshRange(rangeStart, rangeEnd)
	i.ranges = append(i.ranges, fr)
}

// calculateUniqueRanges scans i.ranges for overlapping ranges and consolidates to a unique list of ranges
func (i *inventory) calculateUniqueRanges() {
	var changed bool
	var startCount int
	var uRanges []*freshRange = i.ranges
	startCount = len(uRanges)

	for {
		uRanges, changed = consolidateRanges(uRanges)
		if !changed {
			break
		}
	}

	i.uniqueRanges = uRanges

	log.Debug().
		Int("Starting Count", startCount).
		Int("Ending Count", len(uRanges)).
		Msg("Consolidated ranges")
}

// consolidateRanges takes a list of ranges and returns a consolidated list where overlapping ranges are combined,
// and a boolean indicator of changing the list
func consolidateRanges(rangeList []*freshRange) (uRanges []*freshRange, changed bool) {
	// Read in the full list of ranges with first level overlap checking
ranges:
	for rangeIndex := range rangeList {
		// Start unique range with first entry
		if uRanges == nil {
			newRange := &freshRange{
				end:   rangeList[rangeIndex].end,
				start: rangeList[rangeIndex].start,
			}
			uRanges = append(uRanges, newRange)
			continue ranges
		}

		for uniqueIndex := range uRanges {
			if rangesOverlap(rangeList[rangeIndex], uRanges[uniqueIndex]) {
				changed = true
				uRanges[uniqueIndex] = mergeRanges(rangeList[rangeIndex], uRanges[uniqueIndex])
				continue ranges
			}
		}

		uRanges = append(uRanges, rangeList[rangeIndex])
	}

	return uRanges, changed
}

// total calculates the total available IDs within a range
func (fr *freshRange) total() *big.Int {
	var bigOne *big.Int = big.NewInt(1)
	rangeTotal := new(big.Int)

	// Get the difference in the range
	rangeTotal.Sub(fr.end, fr.start)

	// add one back, because the ranges are inclusive
	rangeTotal.Add(rangeTotal, bigOne)
	if rangeTotal.Cmp(bigOne) >= 0 {
	}

	return rangeTotal
}

// countFreshRanges tallies the possible number of fresh ingredient IDs available
func (i *inventory) countFreshRanges() string {
	fidAvailable := new(big.Int)

	for index := range i.uniqueRanges {
		fidAvailable.Add(fidAvailable, i.uniqueRanges[index].total())
	}

	return fidAvailable.String()
}

// mergeRanges makes the the largest range possible from the ranges passed in
func mergeRanges(fr1, fr2 *freshRange) *freshRange {
	fullRange := &freshRange{
		end:   new(big.Int),
		start: new(big.Int),
	}

	if fr1.end.Cmp(fr2.end) >= 0 {
		fullRange.end.Set(fr1.end)
	} else {
		fullRange.end.Set(fr2.end)
	}

	if fr1.start.Cmp(fr2.start) <= 0 {
		fullRange.start.Set(fr1.start)
	} else {
		fullRange.start.Set(fr2.start)
	}

	log.Debug().
		Str("range1", fmt.Sprintf("%s-%s", fr1.start, fr1.end)).
		Str("range2", fmt.Sprintf("%s-%s", fr2.start, fr2.end)).
		Str("merged", fmt.Sprintf("%s-%s", fullRange.start, fullRange.end)).
		Msg("Ranges merged")

	return fullRange
}

// process process the read in data
func (i *inventory) process() {
	log.Info().Msg("Processing data file")
	i.calculateUniqueRanges()
}

func rangesOverlap(fr1, fr2 *freshRange) bool {
	return (fr1.start.Cmp(fr2.start) >= 0 && fr1.start.Cmp(fr2.end) <= 0) ||
		(fr1.end.Cmp(fr2.start) >= 0 && fr1.end.Cmp(fr2.end) <= 0) ||
		(fr2.start.Cmp(fr1.start) >= 0 && fr2.start.Cmp(fr1.end) <= 0) ||
		(fr2.end.Cmp(fr1.start) >= 0 && fr2.end.Cmp(fr1.end) <= 0)
}

// report on the inventory
func (i *inventory) report() {
	log.Info().
		Str("Fresh Ingredients Possible", i.countFreshRanges()).
		Msg("Report on processing")
}
