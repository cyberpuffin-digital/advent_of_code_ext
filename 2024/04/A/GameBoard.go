package main

import (
	"bufio"
	"strings"
)

const FORWARD = "forward"
const REVERSE = "reverse"
const DOWN = "down"
const UP = "up"
const FORWARDDOWN = "forward down"
const FORWARDUP = "forward up"
const REVERSEDOWN = "reverse down"
const REVERSEUP = "reverse up"

type GameBoard struct {
	rows  int
	cols  int
	board [][]string
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

// getWords pulls the words in each of the 8 directions based on the coordianates
func (b *GameBoard) getWords(x, y, length int) (words map[string]string) {
	var word string
	words = make(map[string]string)

	// Check forward
	if (x + length) <= b.cols {
		for i := range length {
			word += b.board[y][x+i]
		}
		words[FORWARD] = word
		word = ""
	}

	// Check reverse
	if (x - length + 1) >= 0 {
		for i := range length {
			word += b.board[y][x-i]
		}
		words[REVERSE] = word
		word = ""
	}

	// Check down
	if (y + length) <= b.rows {
		for i := range length {
			word += b.board[y+i][x]
		}
		words[DOWN] = word
		word = ""
	}

	// Check up
	if (y - length + 1) >= 0 {
		for i := range length {
			word += b.board[y-i][x]
		}
		words[UP] = word
		word = ""
	}

	// Check forward-down
	if ((x + length) <= b.cols) && ((y + length) <= b.rows) {
		for i := range length {
			word += b.board[y+i][x+i]
		}
		words[FORWARDDOWN] = word
		word = ""
	}

	// Check forward-up
	if ((x + length) <= b.cols) && ((y - length + 1) >= 0) {
		for i := range length {
			word += b.board[y-i][x+i]
		}
		words[FORWARDUP] = word
		word = ""

	}

	// Check reverse-down
	if ((x - length + 1) >= 0) && ((y + length) <= b.rows) {
		for i := range length {
			word += b.board[y+i][x-i]
		}
		words[REVERSEDOWN] = word
		word = ""
	}

	// Check reverse-up
	if ((x - length + 1) >= 0) && ((y - length + 1) >= 0) {
		for i := range length {
			word += b.board[y-i][x-i]
		}
		words[REVERSEUP] = word
	}

	return words
}

// scanForString searches the gameboard for a string
func (b *GameBoard) scanForString(needle string) (tally int) {
	// loop over the game board
	for y, row := range b.board {
		rowCount := 0

		for x, char := range row {
			// Match on the first letter
			if strings.EqualFold(char, needle[0:1]) {
				// Get full word in each direction
				words := b.getWords(x, y, len(needle))

				for _, word := range words {
					if strings.EqualFold(word, needle) {
						rowCount++
						tally++
						verbosef("\tFound %q at row %d; col %d\n", needle, y, x)
					}
				}
			}
		}

		verbosef("Row %d has %d\n", y, rowCount)
	}

	return tally
}
