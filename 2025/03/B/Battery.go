package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const cursorSize int = 12

type batteryBanks struct {
	data         []string
	banks        []*bank
	totalJoltage uint64
}

type bank struct {
	batteries []int
	joltage   uint64
}

type cursor struct {
	index, value int
}

func (bb *batteryBanks) calculateJoltageForAllBanks() {
	for index := range bb.banks {
		bb.banks[index].calculateBankJoltage()
	}
}

func (b *bank) calculateBankJoltage() {
	var cursors [cursorSize]*cursor
	bankSize := len(b.batteries)

	// Initialize cursors for scanning a battery bank
	for index := cursorSize - 1; index >= 0; index-- {
		bankIndex := bankSize - index - 1

		cursors[cursorSize-1-index] = &cursor{
			index: bankIndex,
			value: b.batteries[bankIndex],
		}
	}
	log.Debug().Any("Bank", b.batteries).Msg("Cursors initialized")
	traceCursors(cursors[:])

	// Cycle over the cursors to find their highest values
	for index := range cursors {
		scanStart := 0
		if index > 0 {
			scanStart = cursors[index-1].index + 1
		}

		// cycle through the available numbers to find something larger
		for scanIndex := cursors[index].index; scanIndex >= scanStart; scanIndex-- {
			//		for ; scanIndex < cursors[index].index; scanIndex++ {
			if b.batteries[scanIndex] >= cursors[index].value {
				cursors[index].index = scanIndex
				cursors[index].value = b.batteries[scanIndex]
			}
		}
	}
	log.Debug().Msg("Cursors updated to highest values")
	traceCursors(cursors[:])

	var err error
	var joltageBuilder strings.Builder
	for cIndex := range cursors {
		fmt.Fprintf(&joltageBuilder, "%d", cursors[cIndex].value)
	}
	b.joltage, err = strconv.ParseUint(joltageBuilder.String(), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse joltage string")
		return
	}

	log.Info().Any("Batteries", b.batteries).Uint64("Joltage", b.joltage).Msg("Battery Joltage calculated.")
}

func (bb *batteryBanks) dataToBanks() {
	bb.banks = []*bank{}

	for index := range bb.data {
		newBank := new(bank)

		for _, joltageString := range bb.data[index] {
			log.Trace().Int("joltage", int(joltageString-'0')).Msg("Adding battery to bank")
			newBank.batteries = append(newBank.batteries, int(joltageString-'0'))
		}
		log.Debug().Any("Parsed bank", newBank.batteries).Msg("Data to battery bank")

		bb.banks = append(bb.banks, newBank)
	}

	log.Debug().Msg("Conversion complete")
}

func (bb *batteryBanks) process() {
	log.Info().Msg("Process data object")
	bb.dataToBanks()
	bb.calculateJoltageForAllBanks()
	bb.tallyBanks()
}

func (bb *batteryBanks) report() {
	log.Info().Int("Total Joltage", int(bb.totalJoltage)).Msg("Report on processing")
}

func (bb *batteryBanks) tallyBanks() {
	var tally uint64
	for index := range bb.banks {
		tally += bb.banks[index].joltage
	}
	bb.totalJoltage = tally
}

func traceCursors(cursors []*cursor) {
	for index := range cursors {
		log.Trace().
			Int("Index", cursors[index].index).
			Int("Value", cursors[index].value).
			Msg("Cursors")
	}
}
