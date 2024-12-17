package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Equation struct {
	id         int
	isValid    bool
	mutex      *sync.Mutex
	testValues []int64
	totalValue int64
	validMask  string
}

// calculateMasks will calculate all the possible operator masks to be applied to the equation
func (e *Equation) calculateMasks() (masks []string) {
	var mask, validOperators []rune
	var maskString string

	validOperators = []rune{'+', '*', '|'}

	// create a "zero" mask
	mask = []rune(strings.Repeat(fmt.Sprintf("%c", validOperators[0]), len(e.testValues)-1))

	// Calculate the number of possible combinations
	combinations := math.Pow(float64(len(validOperators)), float64(len(mask)))

	// for all possible combinations
	for combo := int(combinations); combo > 0; combo-- {
		// Get the binary representation of this combination integer
		// length of 50 set arbitrarily.  Not sure how to set it by variable.
		// maskSet := fmt.Sprintf("%050b", combo)
		maskSet := fmt.Sprintf("%050s", strconv.FormatInt(int64(combo), len(validOperators)))

		// For the length of the mask
		for maskCursor := len(mask) - 1; maskCursor >= 0; maskCursor-- {
			e.mutex.Lock()
			verbosef(
				"E%d: %d | %d: mask: %c; combinations: %v; combo: %d; maskSet: %s; maskCursor: %d\n",
				e.id,
				e.totalValue,
				e.testValues,
				mask,
				combinations,
				combo,
				maskSet,
				maskCursor,
			)
			e.mutex.Unlock()
			var err error
			var distance, maskIndex, operatorIndex int
			var operator rune

			distance = len(mask) - maskCursor - 1

			maskIndex = len(maskSet) - 1 - distance
			operatorIndex, err = strconv.Atoi(fmt.Sprintf("%c", maskSet[maskIndex]))
			if err != nil {
				panic("Failed to get operator index in Equation.checkValidity()")
			}

			operator = validOperators[operatorIndex]

			e.mutex.Lock()
			verbosef(
				"Cursor distance: %d; Mask Set value: %c; Operator[%d]: %c\n",
				distance,
				maskSet[maskIndex],
				operatorIndex,
				operator,
			)
			e.mutex.Unlock()
			mask[maskCursor] = operator
		}
		maskString = ""
		for _, o := range mask {
			maskString += fmt.Sprintf("%c", o)
		}
		masks = append(masks, maskString)
	}

	e.mutex.Lock()
	printOut(fmt.Sprintf("E%d has %d possible masks.", e.id, len(masks)))
	verbosef("\t%s\n", masks)
	e.mutex.Unlock()

	return masks
}

// checkMasks cycles through available masks to see if equation is valid
func (e *Equation) checkMasks(masks []string) {
	var err error
	var tally int64

	for _, mask := range masks {
		tally = e.testValues[0]

		for operIndex := 0; operIndex < len(mask); operIndex++ {
			if tally > e.totalValue {
				continue
			}
			switch mask[operIndex] {
			case '+':
				tally += e.testValues[operIndex+1]
			case '*':
				tally *= e.testValues[operIndex+1]
			case '|':
				tally, err = strconv.ParseInt(fmt.Sprintf("%d%d", tally, e.testValues[operIndex+1]), 10, 64)
				if err != nil {
					panic("Failed to parse integer with concat (||) instruction.")
				}
			}
		}

		if tally == e.totalValue {
			e.validMask = mask
			e.isValid = true
			e.mutex.Lock()
			printOut(fmt.Sprintf("\tE%d (%d) is valid with mask %q.", e.id, e.testValues, e.validMask))
			e.mutex.Unlock()

			return
		}

		e.mutex.Lock()
		verbosef("\tE%d with mask %s totals %d (want: %d).\n", e.id, mask, tally, e.totalValue)
		e.mutex.Unlock()
	}
}

// checkValidity will scan the testValues to see if any combination of operands
// can create the totalValue
func (e *Equation) checkValidity(wg *sync.WaitGroup, mutex *sync.Mutex) {
	funcTime := time.Now()

	e.mutex = mutex
	masks := e.calculateMasks()
	maskCalcTime := time.Since(funcTime)
	funcTime = time.Now()
	e.checkMasks(masks)
	maskCheckTime := time.Since(funcTime)
	verbosef("Calculate Masks took %s.\n", maskCalcTime)
	verbosef("Check Masks took %s.\n", maskCheckTime)
	wg.Done()
}

// readIn will scan the input line for total and test values
// Does not use go routines
func (e *Equation) readIn(inputLine string, id int) {
	var err error
	inputs := strings.Split(inputLine, ":")
	e.id = id
	e.isValid = false
	e.totalValue, err = strconv.ParseInt(inputs[0], 10, 64)
	if err != nil {
		verbosef("Failed to readIn equation total: %q.\n", inputs[0])
		panic(err)
	}

	testValuesIn := strings.Split(inputs[1], " ")
	for _, val := range testValuesIn {
		if val == "" {
			continue
		}

		var testIn int64
		testIn, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			verbosef("Failed to readIn test value: %q.\n", val)
			continue
		}
		e.testValues = append(e.testValues, testIn)
	}

	printOut(fmt.Sprintf("E%d read in, total: %d, test values: %v.", e.id, e.totalValue, e.testValues))
}
