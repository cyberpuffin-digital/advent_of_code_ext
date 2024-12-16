package main

import (
	"bufio"
)

type GameBoard struct {
	rows  int
	cols  int
	board [][]string
}

// checkCross checks for MAS crossings with the A at the specified coordinates
func (b *GameBoard) checkCross(x, y int) bool {
	// Can't cross with the middle on an edge
	if x == 0 || x == (b.cols-1) || y == 0 || y == (b.rows-1) {
		return false
	}

	var check1, check2 string

	// check M above
	check1 = b.board[y-1][x-1]
	check2 = b.board[y-1][x+1]
	if check1 == "M" && check2 == "M" {
		check1 = b.board[y+1][x-1]
		check2 = b.board[y+1][x+1]
		if check1 == "S" && check2 == "S" {
			return true
		}
	}

	// check M in front
	check1 = b.board[y-1][x+1]
	check2 = b.board[y+1][x+1]
	if check1 == "M" && check2 == "M" {
		check1 = b.board[y-1][x-1]
		check2 = b.board[y+1][x-1]
		if check1 == "S" && check2 == "S" {
			return true
		}
	}

	// check M below
	check1 = b.board[y+1][x-1]
	check2 = b.board[y+1][x+1]
	if check1 == "M" && check2 == "M" {
		check1 = b.board[y-1][x-1]
		check2 = b.board[y-1][x+1]
		if check1 == "S" && check2 == "S" {
			return true
		}
	}

	// check M behind
	check1 = b.board[y-1][x-1]
	check2 = b.board[y+1][x-1]
	if check1 == "M" && check2 == "M" {
		check1 = b.board[y-1][x+1]
		check2 = b.board[y+1][x+1]
		if check1 == "S" && check2 == "S" {
			return true
		}
	}

	return false
}

// readIn reads the content of the input file into the gameboard
func (b *GameBoard) readIn(f *fileInput) {
	var height int
	b.board = make([][]string, b.rows, b.rows*b.cols)

	// Read the gameboard contents
	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	for scanner.Scan() {
		inputLine := scanner.Text()
		inputString := make([]string, b.cols)

		for index, char := range inputLine {
			inputString[index] = string(char)
		}

		b.board[height] = inputString
		height++
	}
	verbosef("Found %d rows with a max of %d columns\n", b.rows, b.cols)
}

// scan will read the file input and build the board map from the content
func (b *GameBoard) scan(f *fileInput) {
	b.scanColsRows(f)
	b.readIn(f)
}

func (b *GameBoard) scanColsRows(f *fileInput) {
	var height, width int

	f.reset()
	scanner := bufio.NewScanner(f.fileObj)

	// Count rows and columns
	for scanner.Scan() {
		width = 0
		inputLine := scanner.Text()
		verbosef("Input line: %s\n", inputLine)

		for index, char := range inputLine {
			if example {
				verbosef("\tR%dC%d, Index %d == %c\n", height, width, index, char)
			}
			width++
		}
		if b.cols < width {
			b.cols = width
		}

		if len(inputLine) > 0 {
			height++
		}
	}
	b.rows = height
}

// scanForCrossMAS searches for crossing MAS strings
func (b *GameBoard) scanForCrossMAS() (tally int) {
	for y, row := range b.board {
		rowCount := 0

		for x, char := range row {
			// Find middle A
			if char == "A" && b.checkCross(x, y) {
				rowCount++
				tally++
				verbosef("\tFound at row %d; col %d\n", y, x)
			}
		}
	}
	return tally
}

// // scanForString searches the gameboard for a string
// func (b *GameBoard) scanForString(needle string) (tally int) {
// 	// height := 0
// 	// size := len(needle)
// 	// width := 0

// 	// loop over the game board
// 	for y, row := range b.board {
// 		rowCount := 0

// 		for x, char := range row {
// 			// Match on the first letter
// 			if strings.EqualFold(char, needle[0:1]) {
// 				// Get full word in each direction
// 				words := b.getWords(x, y, len(needle))

// 				for _, word := range words {
// 					if strings.EqualFold(word, needle) {
// 						rowCount++
// 						tally++
// 						verbosef("\tFound %q at row %d; col %d\n", needle, y, x)
// 					}
// 				}
// 			}
// 		}

// 		verbosef("Row %d has %d\n", y, rowCount)
// 	}

// 	return tally
// }
