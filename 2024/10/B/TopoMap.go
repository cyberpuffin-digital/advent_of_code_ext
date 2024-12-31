package main

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"strconv"
)

type TopoMap struct {
	cols          int
	distinctPaths int
	knownPaths    []string
	rows          int
	topoMap       [][]int
	trailEnds     [][]int
	trailHeads    [][]int
}

// addPath creates a string path from start and end coordinates and adds it to known paths
func (t *TopoMap) addPath(headX, headY, endX, endY int) {
	trailPath := fmt.Sprintf("%dx%d -> %dx%d", headY, headX, endY, endX)
	t.distinctPaths++

	if !slices.Contains(t.knownPaths, trailPath) {
		t.knownPaths = append(t.knownPaths, trailPath)
	}
}

// calculatePath is a recursive function to check a point on the map for next steps
func (t *TopoMap) calculatePath(x int, y int, currentHeight int, trailHead []int) {
	checkX := x
	checkY := y

	// Check up
	if checkY > 0 && checkY < t.rows {
		checkY--
		if t.topoMap[checkY][checkX]-currentHeight == 1 {
			if t.topoMap[checkY][checkX] == 9 {
				t.addPath(trailHead[1], trailHead[0], checkX, checkY)
			}

			t.calculatePath(checkX, checkY, t.topoMap[checkY][checkX], trailHead)
		}
		checkY++
	}

	// Check down
	if checkY >= 0 && checkY < t.rows-1 {
		checkY++
		if t.topoMap[checkY][checkX]-currentHeight == 1 {
			if t.topoMap[checkY][checkX] == 9 {
				t.addPath(trailHead[1], trailHead[0], checkX, checkY)
			}

			t.calculatePath(checkX, checkY, t.topoMap[checkY][checkX], trailHead)
		}
		checkY--
	}

	// Check left
	if checkX > 0 && checkX < t.cols {
		checkX--
		if t.topoMap[checkY][checkX]-currentHeight == 1 {
			if t.topoMap[checkY][checkX] == 9 {
				t.addPath(trailHead[1], trailHead[0], checkX, checkY)
			}

			t.calculatePath(checkX, checkY, t.topoMap[checkY][checkX], trailHead)
		}
		checkX++
	}

	// Check right
	if checkX >= 0 && checkX < t.cols-1 {
		checkX++
		if t.topoMap[checkY][checkX]-currentHeight == 1 {
			if t.topoMap[checkY][checkX] == 9 {
				t.addPath(trailHead[1], trailHead[0], checkX, checkY)

				return
			}

			t.calculatePath(checkX, checkY, t.topoMap[checkY][checkX], trailHead)
		}
		checkX--
	}
}

// calculateTrailPaths cycles through the trailheads to find paths to trailends
func (t *TopoMap) calculateTrailPaths() {
	for _, trailHead := range t.trailHeads {
		verbosef("trailHead: %d.\n", trailHead)
		t.calculatePath(trailHead[1], trailHead[0], t.topoMap[trailHead[0]][trailHead[1]], trailHead)
	}
}

// initializeMaps initializes the Topographic maps
func (t *TopoMap) initializeMaps() {
	t.topoMap = make([][]int, t.rows, t.rows*t.cols)

	for row := range t.rows {
		topoRow := []int{}

		for range t.cols {
			topoRow = append(topoRow, -1)
		}

		t.topoMap[row] = topoRow
	}
}

// parseInputLine will scan an input line and store coordinates if necessary
func (t *TopoMap) parseInputLine(inputLine string, row int) (inputInts []int) {
	for index, char := range inputLine {
		var intIn int
		var err error

		if char == '.' {
			intIn = -1
		} else {
			intIn, err = strconv.Atoi(string(char))
			if err != nil {
				log.Panicf("Failed to parse character: %v.\n", char)
			}
		}
		inputInts = append(inputInts, intIn)

		// Scan for objects
		if intIn == 0 {
			coord := []int{row, index}
			t.trailHeads = append(t.trailHeads, coord)

			continue
		}

		if intIn == 9 {
			coord := []int{row, index}
			t.trailEnds = append(t.trailEnds, coord)

			continue
		}
	}

	return inputInts
}

// print displays the board as a grid
func (t *TopoMap) print() {
	for row := range t.rows {
		for col := range t.cols {
			verbosef("%02dx%02d: %02d | ", row, col, t.topoMap[row][col])
		}
		verbosef("\n")
	}
}

// readIn reads the content of the input file into the board
func (t *TopoMap) readIn(f *fileInput) {
	t.scanColsRows(f)
	t.initializeMaps()

	t.scanTopoMap(f)
	verbosef("Topographic map input:\n")
	t.print()

	t.calculateTrailPaths()
}

// scanTopoMap reads a file in to the board object
func (t *TopoMap) scanTopoMap(f *fileInput) {
	var row, width int

	// Read the board contents
	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	for scanner.Scan() {
		inputLine := scanner.Text()
		width = len(inputLine)
		if width <= 0 {
			continue
		}
		t.topoMap[row] = t.parseInputLine(inputLine, row)
		row++
	}

	verbosef("Found %d rows with a max of %d columns\n", t.rows, t.cols)
}

// scanColsRows will count the columns and rows of the board
func (t *TopoMap) scanColsRows(f *fileInput) {
	var height, width int

	f.reset()
	scanner := bufio.NewScanner(f.fileObj)

	// Count rows and columns
	for scanner.Scan() {
		width = 0
		inputLine := scanner.Text()

		width = len(inputLine)

		if t.cols < width {
			t.cols = width
		}

		if len(inputLine) > 0 {
			verbosef("Counting line %d: %s\n", height, inputLine)
			height++
		}
	}
	t.rows = height
}
