package main

import (
	"regexp"
	"strconv"
)

type Mult struct {
	coord int
	sum   int
	x     int
	y     int
}

// generate will check subString for multiple instructions and create the struct
func (m *Mult) generate(subString string, coord int) {
	multiples := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	regResults := multiples.FindAllStringSubmatch(subString, -1)
	if len(regResults) > 0 {
		i, err := strconv.Atoi(regResults[0][1])
		if err == nil {
			m.x = i
		}
		i, err = strconv.Atoi(regResults[0][2])
		if err == nil {
			m.y = i
		}
		m.sum = m.x * m.y
	}
	m.coord = coord
}

// getEarlierCoord compares two sets of coordinates to choose the earlier one
// useTwo denotes coord2 is the bigger
func getEarlierCoord(coord1, coord2 []int) (earlierIndex int, useTwo bool) {
	if coord1 == nil && coord2 == nil {
		earlierIndex = -1
	} else if coord1 == nil {
		earlierIndex = coord2[0]
		useTwo = true
	} else if coord2 == nil {
		earlierIndex = coord1[0]
		useTwo = false
	} else if coord1[0] < coord2[0] {
		earlierIndex = coord1[0]
		useTwo = false
	} else {
		earlierIndex = coord2[0]
		useTwo = true
	}

	return earlierIndex, useTwo
}

// tallyLine parses a line of instructions and returns the total
func tallyLine(inputLine string, enabled bool) (int, bool) {
	tally := 0

	multiples := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	setDisable := regexp.MustCompile(`don\'t\(\)`)
	setEnable := regexp.MustCompile(`do\(\)`)

	// calculate next change
	nextDisable := setDisable.FindStringIndex(inputLine)
	nextEnable := setEnable.FindStringIndex(inputLine)
	nextChange, toSet := getEarlierCoord(nextDisable, nextEnable)

	// find next multiple instruction
	nextMultiple := multiples.FindStringIndex(inputLine)

	for keepParsing := true; keepParsing; keepParsing = (nextChange >= 0 || len(nextMultiple) > 0) {

		verbosef("Input line: %s\n\tTally: %d\n\tEnabled: %t\n", inputLine, tally, enabled)

		// If multiples found and the first coordinate is before the next change
		if nextMultiple != nil && (nextMultiple[0] < nextChange || nextChange < 0) {
			var m Mult
			m.generate(inputLine[nextMultiple[0]:nextMultiple[1]], nextMultiple[0])
			if enabled {
				verbosef("\t%d * %d + %d = %d\n", m.x, m.y, tally, tally+m.sum)
				tally += m.sum
			} else {
				verbosef("\tskip %d * %d\n", m.x, m.y)
			}
			inputLine = inputLine[nextMultiple[1]:]
		} else if nextDisable != nil || nextEnable != nil {
			enabled = toSet
			if toSet {
				inputLine = inputLine[nextEnable[1]:]
				verbosef("\tExecute: do() -> Enabled: %t\n", enabled)
			} else {
				inputLine = inputLine[nextDisable[1]:]
				verbosef("\tExecute: don't() -> Enabled: %t\n", enabled)
			}
		}

		// calculate next change
		nextDisable = setDisable.FindStringIndex(inputLine)
		nextEnable = setEnable.FindStringIndex(inputLine)
		nextChange, toSet = getEarlierCoord(nextDisable, nextEnable)

		// find next multiple instruction
		nextMultiple = multiples.FindStringIndex(inputLine)

	}

	return tally, enabled
}
