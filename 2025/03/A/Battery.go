package main

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

type batteryBanks struct {
	data         []string
	banks        []*bank
	totalJoltage int
}

type bank struct {
	batteries []int
	joltage   int
}

func (bb *batteryBanks) calculateJoltageForAllBanks() {
	for index := range bb.banks {
		bb.banks[index].calculateBankJoltage()
	}
}

func (b *bank) calculateBankJoltage() {
	var err error
	var tensIndex, tensValue, onesValue int
	bankSize := len(b.batteries)

	for index, value := range b.batteries {
		if value > tensValue && index != bankSize-1 {
			tensIndex = index
			tensValue = value
		}
	}

	for index := tensIndex + 1; index < bankSize; index++ {
		if b.batteries[index] > onesValue {
			onesValue = b.batteries[index]
		}
	}

	log.Debug().
		Any("Batteries", b.batteries).
		Int("tensValue", tensValue).
		Int("tensIndex", tensIndex).
		Int("onesValue", onesValue).
		Msg("High values found")

	joltageString := fmt.Sprintf("%d%d", tensValue, onesValue)
	b.joltage, err = strconv.Atoi(joltageString)
	if err != nil {
		log.Error().Err(err).Str("Joltage String", joltageString).Msg("Failed to parse joltage string")

		return
	}

	log.Info().Int("Joltage", b.joltage).Msg("Joltage calculated")
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
}

func (bb *batteryBanks) report() {
	log.Info().Int("Total Joltage", bb.tallyBanks()).Msg("Report on processing")
}

func (bb *batteryBanks) tallyBanks() (tally int) {
	for index := range bb.banks {
		tally += bb.banks[index].joltage
	}

	return tally
}
