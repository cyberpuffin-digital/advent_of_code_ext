package main

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	LEFT  = "left"
	RIGHT = "right"
)

type doorLock struct {
	data       []string
	pointingAt int // 0 - 99
	zeroes     int // number of times zero was found
}

// parseInputLine reads the line from the input and determines the direction
// and quantity of turns
func (d *doorLock) parseInputLine(dataIn string) (dir string, turns int) {
	if dataIn == "" {
		log.Debug().Msg("Skipping empty input")
		return dir, turns
	}

	var err error
	var numStr string

	log.Debug().
		Str("Data In", dataIn).
		Msg("Parsing input")

	if strings.HasPrefix(dataIn, "L") {
		dir = LEFT
		numStr = strings.TrimPrefix(dataIn, "L")
	} else if strings.HasPrefix(dataIn, "R") {
		dir = RIGHT
		numStr = strings.TrimPrefix(dataIn, "R")
	}

	if numStr == "" {
		log.Error().Msg("Failed to find prefix for input string")
		return dir, turns
	}

	turns, err = strconv.Atoi(numStr)
	if err != nil {
		log.Error().Err(err).Stack().
			Msg("failed to parse the number of turns")
	}

	log.Debug().
		Str("Direction", dir).
		Int("Turns", turns).
		Msg("processed input line")

	return dir, turns
}

// process reads through the data and turns the dial
func (d *doorLock) process() {
	d.setStart(50)
	for index := range d.data {
		dir, turns := d.parseInputLine(d.data[index])
		if turns > 0 {
			d.turnDial(dir, turns)
		}
	}
	d.report()
}

// report prints information on the data
func (d *doorLock) report() {
	log.Info().
		Int("Number of lines", len(d.data)).
		Int("Final position", d.pointingAt).
		Int("Zeroes", d.zeroes).
		Msg("Final report")
}

// setStart set's the door lock's starting position
func (d *doorLock) setStart(pos int) {
	d.pointingAt = pos
}

// turnDial turns the door lock's dial
func (d *doorLock) turnDial(dir string, times int) {
	log.Debug().
		Int("Starting Position", d.pointingAt).
		Str("direction", dir).
		Int("turns", times).
		Msg("Turning dial")

	for range times {
		if strings.EqualFold(dir, LEFT) {
			d.pointingAt--
			if d.pointingAt == 0 {
				d.zeroes++
				log.Info().
					Int("Zeroes", d.zeroes).
					Msg("Zero found!")
			} else if d.pointingAt < 0 {
				d.pointingAt = 100 + d.pointingAt
			}
		} else {
			d.pointingAt++
			if d.pointingAt == 100 {
				d.pointingAt = 0
				d.zeroes++
				log.Info().
					Int("Zeroes", d.zeroes).
					Msg("Zero found!")
			}
		}
	}

	log.Info().
		Int("Final position", d.pointingAt).
		Msg("Dial turned")
}
