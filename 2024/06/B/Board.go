package main

import (
	"bufio"
	"fmt"
)

type Board struct {
	board                 [][]rune
	boardBackup           [][]rune
	cols                  int
	gameActive            bool
	playerPosition        []int
	playerPositionBackup  []int
	playerDirection       rune
	playerDirectionBackup rune
	rows                  int
	visitCount            int
}

const OBSTACLE = '#'
const PLAYERDOWN = 'v'
const PLAYERLEFT = '<'
const PLAYERRIGHT = '>'
const PLAYERUP = '^'
const UNVISITED = '.'
const VISITED = 'X'

// calculateCursor calculates the next or the last step based on player direction
func (b *Board) calculateCursor(row, col int, forward bool) (int, int) {
	switch b.playerDirection {
	case PLAYERDOWN:
		if forward {
			row++
		} else {
			row--
		}
	case PLAYERLEFT:
		if forward {
			col--
		} else {
			col++
		}
	case PLAYERRIGHT:
		if forward {
			col++
		} else {
			col--
		}
	case PLAYERUP:
		if forward {
			row--
		} else {
			row++
		}
	}
	return row, col
}

// checkEveryObstacle cycles through every possible position to check for loops
func (b *Board) checkEveryObstacle() (loopCount int) {
	// Cycle every element on the board to try to create a loop
	for row := range b.rows {
		for col := range b.cols {
			// Skip invalid obstacle positions
			switch b.board[row][col] {
			case OBSTACLE:
				continue
			case PLAYERDOWN, PLAYERLEFT, PLAYERRIGHT, PLAYERUP:
				continue
			}
			b.board[row][col] = OBSTACLE
			b.print()
			if b.checkForLoop() {
				loopCount++
			}
			verbosef("%d loops found.\n", loopCount)
			b.restore()
		}
	}

	return loopCount
}

// checkForLoop runs a game to see if the guard will get stuck in a loop
func (b *Board) checkForLoop() (looping bool) {
	for l := 0; l <= (b.cols * b.rows * 2); l++ {
		b.iterate()
		if !b.gameActive {
			break
		}
	}

	looping = b.gameActive
	if !b.gameActive {
		b.gameActive = true
	}

	return looping
}

// countDistinct will tally the number VISITED locations
func (b *Board) count() (distinct int) {
	for row := range b.rows {
		for col := range b.cols {
			switch b.board[row][col] {
			case PLAYERDOWN, PLAYERLEFT, PLAYERRIGHT, PLAYERUP, VISITED:
				distinct++
			}
		}
	}

	return distinct
}

// directionString gives the player's direction as a string
func (b *Board) directionString() (direction string) {
	switch b.playerDirection {
	case PLAYERDOWN:
		direction = "down"
	case PLAYERLEFT:
		direction = "left"
	case PLAYERRIGHT:
		direction = "right"
	case PLAYERUP:
		direction = "up"
	}

	return direction
}

// isSafe checks the bounds of the game board
func (b *Board) isSafe(row, col int) bool {
	// Bound check
	if row < 0 || row >= b.rows {
		// verbosef("Target Row %d is out of expected range (0:%d)\n", row, b.rows-1)

		return false
	} else if col < 0 || col >= b.cols {
		// verbosef("Target Col %d is out of expected range (0:%d)\n", col, b.cols-1)

		return false
	}

	return true
}

// iterate will calculate a board movement
func (b *Board) iterate() {
	calculating := true

	// set initial cursor position
	cursorCol := b.playerPosition[1]
	cursorRow := b.playerPosition[0]

	// loop until a valid move is found
	for i := 0; calculating; i++ {
		cursorRow, cursorCol = b.calculateCursor(cursorRow, cursorCol, true)

		if !b.isSafe(cursorRow, cursorCol) {
			calculating = false
			b.gameActive = false
			b.board[b.playerPosition[0]][b.playerPosition[1]] = VISITED

			return
		}

		switch b.board[cursorRow][cursorCol] {
		case OBSTACLE:
			cursorRow, cursorCol = b.calculateCursor(cursorRow, cursorCol, false)
			b.turnPlayer()
		case UNVISITED:
			if b.moveTo(cursorRow, cursorCol) {
				b.visitCount++
				calculating = false
			}
		case VISITED:
			if b.moveTo(cursorRow, cursorCol) {
				calculating = false
			}
		}
	}
	// b.print()
}

// moveTo moves a player to the requested coordinates
func (b *Board) moveTo(targetRow, targetCol int) bool {
	verbosef(
		"Moving player from %dx%d to %dx%d.\tvisitCount: %d\n",
		b.playerPosition[0],
		b.playerPosition[1],
		targetRow,
		targetCol,
		b.visitCount,
	)
	b.board[targetRow][targetCol] = b.board[b.playerPosition[0]][b.playerPosition[1]]
	b.board[b.playerPosition[0]][b.playerPosition[1]] = VISITED
	b.playerPosition[0] = targetRow
	b.playerPosition[1] = targetCol

	return true
}

// print displays the board as a grid
func (b *Board) print() {
	for row := range b.rows {
		verbosef("Row %d | ", row)
		for col := range b.cols {
			verbosef("%d %c | ", col, b.board[row][col])

		}
		verbosef("\n")
	}
}

// readIn reads the content of the input file into the board
func (b *Board) readIn(f *fileInput) {
	var height int
	var msg string
	b.board = make([][]rune, b.rows, b.rows*b.cols)

	// Read the board contents
	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	for scanner.Scan() {
		inputLine := scanner.Text()
		inputString := make([]rune, b.cols)
		msg += fmt.Sprintf("Row %d | ", height)

		for index, char := range inputLine {
			inputString[index] = rune(char)
			msg += fmt.Sprintf("%d %q ", index, char)

			// Scan for objects
			switch char {
			case OBSTACLE:
			case PLAYERDOWN:
				b.playerDirection = PLAYERDOWN
				b.playerPosition = []int{height, index}
			case PLAYERLEFT:
				b.playerDirection = PLAYERLEFT
				b.playerPosition = []int{height, index}
			case PLAYERRIGHT:
				b.playerDirection = PLAYERRIGHT
				b.playerPosition = []int{height, index}
			case PLAYERUP:
				b.playerDirection = PLAYERUP
				b.playerPosition = []int{height, index}
			case UNVISITED:
			case VISITED:
			}
		}

		b.board[height] = inputString
		height++
		msg += "|\n"
	}
	if example {
		verbosef(msg)
	}
	verbosef("Found %d rows with a max of %d columns\n", b.rows, b.cols)
}

// restore will copy the boardBackup back to the game board
func (b *Board) restore() {
	b.board = make([][]rune, b.rows, b.rows*b.cols)
	for row := range b.rows {
		b.board[row] = make([]rune, b.cols)
		copy(b.board[row], b.boardBackup[row])
	}
	b.playerDirection = b.playerDirectionBackup
	b.playerPosition = make([]int, 2)
	copy(b.playerPosition, b.playerPositionBackup)
	b.visitCount = 0

	verbosef("Game board restored from backup.\n")
}

// save will copy the game board to boardBackup
func (b *Board) save() {
	b.boardBackup = make([][]rune, b.rows, b.rows*b.cols)
	for row := range b.rows {
		b.boardBackup[row] = make([]rune, b.cols)
		copy(b.boardBackup[row], b.board[row])
	}
	b.playerDirectionBackup = b.playerDirection
	b.playerPositionBackup = make([]int, 2)
	copy(b.playerPositionBackup, b.playerPosition)

	verbosef("Game board saved to backup.\n")
}

// scan will read the file input and build the board map from the content
func (b *Board) scan(f *fileInput) {
	b.scanColsRows(f)
	b.readIn(f)

	fmt.Printf(
		"Player is at %dx%d facing %s.\n",
		b.playerPosition[0],
		b.playerPosition[1],
		b.directionString(),
	)
}

// scanColsRows will count the columns and rows of the game board
func (b *Board) scanColsRows(f *fileInput) {
	var height, width int

	f.reset()
	scanner := bufio.NewScanner(f.fileObj)

	// Count rows and columns
	for scanner.Scan() {
		width = 0
		inputLine := scanner.Text()
		verbosef("Input line: %s\n", inputLine)

		width = len(inputLine)

		if b.cols < width {
			b.cols = width
		}

		if len(inputLine) > 0 {
			height++
		}
	}
	b.rows = height
}

// turnPlayer rotates player right by 90 degrees
func (b *Board) turnPlayer() {
	verbosef("Player starts facing %c ", b.playerDirection)
	switch b.playerDirection {
	case PLAYERDOWN:
		b.playerDirection = PLAYERLEFT
		b.board[b.playerPosition[0]][b.playerPosition[1]] = PLAYERLEFT
	case PLAYERLEFT:
		b.playerDirection = PLAYERUP
		b.board[b.playerPosition[0]][b.playerPosition[1]] = PLAYERUP
	case PLAYERRIGHT:
		b.playerDirection = PLAYERDOWN
		b.board[b.playerPosition[0]][b.playerPosition[1]] = PLAYERDOWN
	case PLAYERUP:
		b.playerDirection = PLAYERRIGHT
		b.board[b.playerPosition[0]][b.playerPosition[1]] = PLAYERRIGHT
	}
	verbosef("and finished facing %c\n", b.playerDirection)
}
