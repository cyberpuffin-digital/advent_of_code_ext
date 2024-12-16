package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var listLeft, listRight LocationList
	listLeft.countCache = make(map[int]int)
	listRight.countCache = make(map[int]int)
	listLeft.locations, listRight.locations = scanInput("input.txt")
	// listLeft.locations, listRight.locations = scanInput("example.in")

	var tally int

	for _, valueLeft := range listLeft.locations {
		multiple := listRight.countOccurrences(valueLeft)
		tally += valueLeft * multiple
	}
	fmt.Printf("Similarity score: %d\n", tally)
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
