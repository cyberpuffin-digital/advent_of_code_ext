package main

import (
	"bufio"
)

type AntennaMap struct {
	antennas    [][]rune
	antinodes   [][]bool
	cols        int
	frequencies map[rune][][]int
	rows        int
}

// calculateAntinodeLocation cycles through the frequencies to determine antinode locations
func (a *AntennaMap) calculateAntinodeLocation() {
	for freq, coords := range a.frequencies {
		verbosef("Freq %c: %d\n", freq, coords)
		for index1, coord1 := range coords {
			for index2, coord2 := range coords {
				// Skip self check
				if index1 == index2 {
					continue
				}

				// Distance math
				xdiff := coord2[1] - coord1[1]
				ydiff := coord2[0] - coord1[0]

				xAnti := coord1[1] - xdiff
				yAnti := coord1[0] - ydiff

				// Bounds checking
				if xAnti < a.cols && xAnti >= 0 {
					if yAnti < a.rows && yAnti >= 0 {
						a.antinodes[yAnti][xAnti] = true
					}
				}
			}
			verbosef("\t%2d: %d\n", index1, coord1)
		}
	}
}

// countAntinodes will tally the number antinodes in the map
func (a *AntennaMap) countAntinodes() (distinct int) {
	for row := range a.rows {
		for col := range a.cols {
			if a.antinodes[row][col] {
				distinct++
			}
		}
	}

	return distinct
}

// initializeMaps initializes the antenna maps
func (a *AntennaMap) initializeMaps() {
	a.antennas = make([][]rune, a.rows, a.rows*a.cols)
	a.antinodes = make([][]bool, a.rows, a.rows*a.cols)
	a.frequencies = make(map[rune][][]int)

	for row := range a.rows {
		antennaRow := []rune{}
		antinodeRow := []bool{}

		for range a.cols {
			antennaRow = append(antennaRow, '.')
			antinodeRow = append(antinodeRow, false)
		}

		a.antennas[row] = antennaRow
		a.antinodes[row] = antinodeRow
	}
}

// parseInputLine will scan an input line and store coordinates if necessary
func (a *AntennaMap) parseInputLine(inputLine string, row int) (inputRunes []rune) {
	for index, char := range inputLine {
		runeIn := rune(char)
		inputRunes = append(inputRunes, runeIn)

		// Scan for objects
		if runeIn == '.' {
			continue
		}

		freqCoords, ok := a.frequencies[runeIn]
		if !ok {
			freqCoords = [][]int{}
		}
		freqCoords = append(freqCoords, []int{row, index})
		a.frequencies[runeIn] = freqCoords
	}

	return inputRunes
}

// print displays the board as a grid
func (a *AntennaMap) print(withAntinodes bool) {
	for row := range a.rows {
		verbosef("Row %d | ", row)
		for col := range a.cols {
			if withAntinodes && a.antinodes[row][col] {
				verbosef("%d X | ", col)
			} else {
				verbosef("%d %c | ", col, a.antennas[row][col])
			}
		}
		verbosef("\n")
	}
}

// readIn reads the content of the input file into the board
func (a *AntennaMap) readIn(f *fileInput) {
	a.scanColsRows(f)
	a.initializeMaps()

	a.scanAntennaMap(f)
	verbosef("Antenna map input:\n")
	a.print(false)

	a.calculateAntinodeLocation()
	verbosef("Antinodes calculated:\n")
	a.print(true)
}

// scanAntennaMap reads a file in to the board object
func (a *AntennaMap) scanAntennaMap(f *fileInput) {
	var row, width int
	var msg string

	// Read the board contents
	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	for scanner.Scan() {
		inputLine := scanner.Text()
		width = len(inputLine)
		if width <= 0 {
			continue
		}
		a.antennas[row] = a.parseInputLine(inputLine, row)
		row++
	}
	if example {
		verbosef(msg)
	}

	verbosef("Found %d rows with a max of %d columns\n", a.rows, a.cols)
}

// scanColsRows will count the columns and rows of the board
func (a *AntennaMap) scanColsRows(f *fileInput) {
	var height, width int

	f.reset()
	scanner := bufio.NewScanner(f.fileObj)

	// Count rows and columns
	for scanner.Scan() {
		width = 0
		inputLine := scanner.Text()

		width = len(inputLine)

		if a.cols < width {
			a.cols = width
		}

		if len(inputLine) > 0 {
			verbosef("Counting line in: %s\n", inputLine)
			height++
		}
	}
	a.rows = height
}
