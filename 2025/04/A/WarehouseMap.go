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

type warehouseMap struct {
	cols            int
	data            []string
	grid            [][]mapItem
	rollsAccessible int
	rows            int
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

	log.Info().Any("Grid", wm.grid).Msg("Data converted to grid")
}

func (wm *warehouseMap) process() {
	log.Info().Msg("Process data object")
	wm.dataToMapGrid()
	wm.scanGrid()
}

func (wm *warehouseMap) report() {
	log.Info().Int("Accessible Rolls", wm.rollsAccessible).Msg("Report on processing")
}

func (wm *warehouseMap) scanGrid() {
	for row := 0; row < wm.rows; row++ {
		for col := 0; col < wm.cols; col++ {
			if wm.grid[row][col] == floor {
				continue
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

				log.Debug().
					Int("Col", checkCol).
					Int("Count", paperCount).
					Int("Row", checkRow).
					Msg("Paper found")
			}
		}
	}
	if paperCount < 4 {
		wm.rollsAccessible++
		log.Info().
			Int("Accessible Count", wm.rollsAccessible).
			Int("Adjacent Count", paperCount).
			Int("Col", col).
			Int("Row", row).
			Msg("Accessible paper found")
	}
}
