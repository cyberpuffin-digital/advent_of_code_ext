package main

import (
	"log"
	"strconv"
	"strings"
)

type Report struct {
	text      string
	fields    []int
	direction int // 1 == ascending; -1 == descending
}

const ASCENDING = 1
const DESCENDING = -1
const STATIC = 0

// generate instantiates a new report
func (r *Report) generate(reportText string) {
	r.text = reportText

	err := r.parseText()
	if err != nil {
		log.Fatalf("Failed to parse text: %q", reportText)
	}
	r.setDirection()
}

// isSafe determines if the report is safe
func (r *Report) isSafe() bool {
	if r.direction == STATIC {
		return false
	}

	var lastValue int
	for index, value := range r.fields {
		if index == 0 {
			lastValue = value
			continue
		}

		// Neither increasing nor decreasing is unsafe
		if lastValue == value {
			return false
		}

		// Check each element against the last
		var delta int

		if r.direction == DESCENDING {
			delta = lastValue - value
		} else if r.direction == ASCENDING {
			delta = value - lastValue
		}

		if delta > 3 || delta < 1 {
			return false
		}
		lastValue = value
	}

	return true
}

// parseText scans for elements within the text field
func (r *Report) parseText() error {
	str_fields := strings.Fields(r.text)

	for _, strValue := range str_fields {
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			return err
		}
		r.fields = append(r.fields, intValue)
	}

	return nil
}

// setDirection determines the reports direction from the first two fields
func (r *Report) setDirection() {
	if len(r.fields) < 2 {
		log.Printf("Too few fields to set direction: %q", len(r.fields))
		return
	}

	if r.fields[0] > r.fields[1] {
		r.direction = DESCENDING
	} else if r.fields[0] < r.fields[1] {
		r.direction = ASCENDING
	} else {
		r.direction = STATIC
	}
}
