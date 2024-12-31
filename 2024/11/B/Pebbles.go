package main

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type PebbleList struct {
	blinks       int
	elementCount int
	first        *Pebble
}

type Pebble struct {
	duplicates int
	next       *Pebble
	prev       *Pebble
	value      int
}

// blink runs the calculations N times
func (pl *PebbleList) blink(times int) {
	for blink := 1; blink <= times; blink++ {
		pl.calculateBlink()
		verbosef("Blink #%d pre-compression:\n", blink)
		pl.print()
		pl.compress()
		pl.blinks++
		verbosef("Blink #%d compressed:\n", blink)
		pl.print()
		printOut(fmt.Sprintf("Blink #%d has %d pebbles (compressed to %d).", blink, pl.count(), pl.elementCount))
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
			entryValue, err := strconv.Atoi(entry)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse input: %s", entry))
			}
			// Create a new pebble
			p := &Pebble{value: entryValue}

			// Handle the first pebble
			if pl.first == nil {
				pl.first = p
				currentPebble = p

				continue
			}

			// Handle remaining pebbles
			currentPebble.next = p
			p.prev = currentPebble
			currentPebble = p
		}
	}
}

// calculatePebble determines the action to take on pIn
func (p *Pebble) calculate() {
	if p.value == 0 {
		// If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
		p.value = 1
	} else if len(fmt.Sprint(p.value))%2 == 0 {
		// If the stone is engraved with a number that has an even number of digits, it is replaced by two stones. The
		// left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on
		// the new right stone. (The new numbers don't keep extra leading zeroes: 1000 would become stones 10 and 0.)
		p.splitPebble()
	} else {
		// If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by
		// 2024 is engraved on the new stone.
		p.value *= 2024
	}
}

// calculateBlink updates the pebble list based on challenge rules
func (pl *PebbleList) calculateBlink() {
	// var wg sync.WaitGroup
	p := pl.first

	for p != nil {
		next := p.next
		// 	wg.Add(1)
		// 	go func(pebble *Pebble) {
		// 		defer wg.Done()
		// 		pebble.calculate()
		// 	}(p)
		p.calculate()
		p = next
	}
	// wg.Wait()
}

// compress will cycle through the list of pebbles to look for and remove duplicates
func (pl *PebbleList) compress() {
	for p := pl.first; p != nil; p = p.next {
		for pCheck := p.next; pCheck != nil; pCheck = pCheck.next {
			if p.value == pCheck.value {
				p.duplicates++
				p.duplicates += pCheck.duplicates

				pCheck.prev.next = pCheck.next
				if pCheck.next != nil {
					pCheck.next.prev = pCheck.prev
				}
			}
		}
	}
}

// count cycles through the pebble list counting pebbles
func (pl *PebbleList) count() (pebbleCount int) {
	elementCount := 0

	for p := pl.first; p != nil; p = p.next {
		elementCount++
		pebbleCount += p.duplicates + 1
	}

	pl.elementCount = elementCount
	return pebbleCount
}

// intLength counts the number of digits in an integer
func intLength(i int) (count int) {
	if i == 0 {
		return 1
	}
	for i != 0 {
		i /= 10
		count++
	}

	return count
}

// print displays the board as a grid
func (pl *PebbleList) print() {
	if !verbose {
		return
	}

	var sb strings.Builder
	sb.WriteString("\t->")
	p := pl.first

	for {
		sb.WriteString(fmt.Sprintf(" %d", p.value))
		if p.next == nil {
			break
		}
		p = p.next
	}
	verbosef("%s\n", sb.String())
}

// readIn reads the content of the input file into the board
func (pl *PebbleList) readIn(f *fileInput) {
	pl.buildList(f)
	pl.print()
	verbosef("Blink #0 has %d pebbles\n", pl.count())
}

// splitPebble creates a new pebble from an existing one
func (p *Pebble) splitPebble() {
	divider := int(math.Pow10(intLength(p.value) / 2))
	// divider := int(math.Pow10(len(fmt.Sprint(p.value)) / 2))

	int1 := p.value / divider
	int2 := p.value % divider

	newPebble := &Pebble{
		duplicates: p.duplicates,
		next:       p.next,
		prev:       p,
		value:      int2,
	}

	if p.next != nil {
		p.next.prev = newPebble
	}
	p.next = newPebble
	p.value = int1
}
