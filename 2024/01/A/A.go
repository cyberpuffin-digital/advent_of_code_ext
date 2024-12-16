package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	listLeft, listRight := scanInput("input.txt")

	sort.Slice(listLeft, func(i, j int) bool {
		return listLeft[i] < listLeft[j]
	})
	sort.Slice(listRight, func(i, j int) bool {
		return listRight[i] < listRight[j]
	})

	var tally float64

	for index, valueLeft := range listLeft {
		tally += math.Abs(float64(valueLeft - listRight[index]))
	}
	fmt.Printf("List length - Left: %d; Right: %d\n", len(listLeft), len(listRight))
	fmt.Printf("Total distance: %d\n", int(tally))
}

// parseEntry converts the locationID string to an integer
func parseEntry(locationString string) int {
	if locationString == "" {
		return 0
	}

	locationID, err := strconv.Atoi(locationString)
	if err != nil {
		log.Fatalf("Unable to parse entry %q", locationString)
	}
	return locationID
}

// scanInput reads an input file line by line to create separate lists
func scanInput(inputFile string) (listLeft, listRight []int) {
	input, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		if len(entry) < 2 {
			log.Fatalf("Failed to parse input line.  Expected 2, received %d", len(entry))
			return
		}

		if entry[0] == "" || entry[1] == "" {
			continue
		}
		listLeft = append(listLeft, parseEntry(entry[0]))
		listRight = append(listRight, parseEntry(entry[1]))
	}

	return listLeft, listRight
}
