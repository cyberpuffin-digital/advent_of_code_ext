package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type PebbleList struct {
	blinks int
	first  *Pebble
	last   *Pebble
}

type Pebble struct {
	next  *Pebble
	prev  *Pebble
	value int
}

// blink runs the calculations N times
func (pl *PebbleList) blink(times int) {
	for b := range times {
		verbosef("Blink #%d:\n", b+1)
		pl.calculateBlink()
		pl.print()
		pl.blinks++
	}
}

// calculateBlink updates the pebble list based on challenge rules
func (pl *PebbleList) calculateBlink() {
	p := pl.first

	for {
		if p == nil {
			break
		}

		// If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
		if p.value == 0 {
			p.value = 1
			p = p.next
			continue
		}

		// If the stone is engraved with a number that has an even number of digits, it is replaced by two stones. The
		// left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on
		// the new right stone. (The new numbers don't keep extra leading zeroes: 1000 would become stones 10 and 0.)
		strValue := fmt.Sprint(p.value)
		if len(strValue)%2 == 0 {
			// Split the string evenly
			strSplit := strings.Split(strValue, "")
			value1 := strSplit[0 : len(strSplit)/2]
			value2 := strSplit[len(strSplit)/2:]

			// Parse the new pebble numbers
			int1, err := strconv.Atoi(strings.Join(value1, ""))
			if err != nil {
				panic("Failed to parse int1")
			}
			int2, err := strconv.Atoi(strings.Join(value2, ""))
			if err != nil {
				panic("Failed to parse int2")
			}

			// Create a new pebble and add to the list
			newPebble := new(Pebble)

			// Setup new pebble
			newPebble.value = int2
			newPebble.next = p.next
			newPebble.prev = p
			if pl.last == p {
				pl.last = newPebble
			}

			// Update old pebble
			p.value = int1
			p.next = newPebble

			// Continue to next pebble
			p = newPebble.next
			continue
		}

		// If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by
		// 2024 is engraved on the new stone.
		p.value = p.value * 2024

		if p.next == nil {
			break
		} else {
			p = p.next
		}
	}
}

// buildList reads the file input and builds the list of pebbles
func (pl *PebbleList) buildList(f *fileInput) {
	var currentPebble *Pebble

	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	for scanner.Scan() {
		inputLine := scanner.Text()

		for _, entry := range strings.Split(inputLine, " ") {
			var err error

			// Create a new pebble
			p := new(Pebble)
			p.value, err = strconv.Atoi(entry)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse input: %s", entry))
			}

			// Handle the first pebble
			if pl.first == nil {
				pl.first = p
				pl.last = p
				currentPebble = p

				continue
			}

			// Handle remaining pebbles
			currentPebble.next = p
			p.prev = currentPebble
			pl.last = p
			currentPebble = p
		}
	}
}

// count cycles through the pebble list counting pebbles
func (pl *PebbleList) count() (pebbleCount int) {
	p := pl.first

	if p == nil {
		verbosef("Empty pebble list.  PL: %v\n", pl)
		return pebbleCount
	}

	for {
		pebbleCount++
		if p.next == nil {
			break
		}
		p = p.next
	}

	return pebbleCount
}

// print displays the board as a grid
func (pl *PebbleList) print() {
	p := pl.first
	returnStr := "\t-> "

	for {
		returnStr = fmt.Sprintf("%s %d", returnStr, p.value)
		if p.next == nil {
			break
		}
		p = p.next
	}
	verbosef("%s\n", returnStr)
}

// readIn reads the content of the input file into the board
func (pl *PebbleList) readIn(f *fileInput) {
	pl.buildList(f)
	verbosef("Starting list (%d pebbles)\n", pl.count())
	pl.print()
}
