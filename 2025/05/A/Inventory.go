package main

import (
	"math/big"

	"github.com/rs/zerolog/log"
)

type freshRange struct {
	end   *big.Int
	start *big.Int
}

type inventory struct {
	availableIDs []*big.Int
	freshIDs     []*big.Int
	ranges       []*freshRange
}

func newFreshRange(start, end string) (fr *freshRange) {
	fr = &freshRange{
		end:   new(big.Int),
		start: new(big.Int),
	}

	_, ok := fr.end.SetString(end, 10)
	if !ok {
		log.Error().Str("End range", end).Msg("Failed to set Fresh range end value")
	}

	_, ok = fr.start.SetString(start, 10)
	if !ok {
		log.Error().Str("Start range", start).Msg("Failed to set Fresh range start value")
	}

	return fr
}

func (i *inventory) checkInventory() {
	for index := range i.availableIDs {
		if i.checkItemIsFresh(i.availableIDs[index]) {
			i.freshIDs = append(i.freshIDs, i.availableIDs[index])
			log.Debug().
				Str("Ingredient ID", i.availableIDs[index].String()).
				Msg("Fresh Ingredient")
		} else {
			log.Debug().
				Str("Ingredient ID", i.availableIDs[index].String()).
				Msg("Item has spoiled")
		}
	}
}

func (i *inventory) checkItemIsFresh(item *big.Int) bool {
	for _, fr := range i.ranges {
		if fr.hasIngredient(item) {
			return true
		}
	}

	return false
}

func (fr *freshRange) hasIngredient(checkID *big.Int) bool {
	return checkID.Cmp(fr.start) >= 0 && checkID.Cmp(fr.end) <= 0
}

func (i *inventory) process() {
	log.Info().Msg("Process data object")
	i.checkInventory()
}

func (i *inventory) report() {
	log.Info().
		Int("Fresh Ingredient Count", len(i.freshIDs)).
		Msg("Report on processing")
}
