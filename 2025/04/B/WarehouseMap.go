package main

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type mapItem int

const (
	floor mapItem = 0
	paper mapItem = 1
)

type coord struct {
	col int
	row int
}

type warehouseMap struct {
	cols            int
	data            []string
	grid            [][]mapItem
	rollsAccessible int
	rollCoordinates []coord
	rows            int
	totalPaper      int
	totalRemoved    int
}

func (wm *warehouseMap) dataToMapGrid() {
	var colCount, rowCount int

	for index := range wm.data {
		if strings.TrimSpace(wm.data[index]) == "" {
			continue
		}
		var gridRow []mapItem

		rowCount++
		colCount = 0
		for rowIndex := range wm.data[index] {
			colCount++
			switch wm.data[index][rowIndex] {
			case '.':
				gridRow = append(gridRow, floor)
			case '@':
				gridRow = append(gridRow, paper)
			}
		}
		wm.grid = append(wm.grid, gridRow)
	}
	wm.cols = colCount
	wm.rows = rowCount
	wm.totalRemoved = 0

	log.Info().Any("Grid", wm.grid).Msg("Data converted to grid")
}

func (wm *warehouseMap) process() {
	log.Info().Msg("Process data object")
	wm.dataToMapGrid()
	for {
		wm.scanGrid()
		if wm.rollsAccessible == 0 {
			break
		}
		wm.removeAccessible()
	}
}

func (wm *warehouseMap) removeAccessible() {
	log.Info().
		Int("Accessible Count", wm.rollsAccessible).
		Msg("Clearing accessible paper rolls")

	for _, c := range wm.rollCoordinates {
		log.Debug().
			Any("Coord", c).
			Msg("Clearing paper roll")

		wm.grid[c.row][c.col] = floor
		wm.totalRemoved++
	}
}

func (wm *warehouseMap) report() {
	log.Info().Int("Paper Removed", wm.totalRemoved).Msg("Report on processing")
}

func (wm *warehouseMap) scanGrid() {
	wm.rollsAccessible = 0
	wm.rollCoordinates = []coord{}
	wm.totalPaper = 0

	for row := 0; row < wm.rows; row++ {
		for col := 0; col < wm.cols; col++ {
			switch wm.grid[row][col] {
			case floor:
				continue
			case paper:
				wm.totalPaper++
			}

			log.Debug().
				Int("Col", col).
				Int("Row", row).
				Msg("Scanning position")

			wm.scanPosition(col, row)
		}
	}
}

func (wm *warehouseMap) scanPosition(col, row int) {
	var paperCount int = 0

	for checkRow := -1; checkRow < 2; checkRow++ {
		for checkCol := -1; checkCol < 2; checkCol++ {
			// Skip self or out of bounds
			if (checkCol == 0 && checkRow == 0) ||
				(col+checkCol < 0) || (col+checkCol >= wm.cols) ||
				(row+checkRow < 0) || (row+checkRow >= wm.rows) {
				continue
			}

			if wm.grid[row+checkRow][col+checkCol] == paper {
				paperCount++
			}
		}
	}
	if paperCount < 4 {
		wm.rollsAccessible++
		wm.rollCoordinates = append(wm.rollCoordinates, coord{
			col: col,
			row: row,
		})

		log.Info().
			Int("Accessible Count", wm.rollsAccessible).
			Int("Adjacent Count", paperCount).
			Int("Col", col).
			Int("Row", row).
			Msg("Accessible paper found")
	}
}
