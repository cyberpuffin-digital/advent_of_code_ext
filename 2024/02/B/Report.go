package main

import (
	"log"
	"strconv"
	"strings"
)

type Report struct {
	dampened  bool
	direction int // 1 == ascending; -1 == descending
	fields    []int
	text      string
}

const ASCENDING = 1
const DESCENDING = -1
const STATIC = 0

// areFieldsSafe checks a slice of integers to see if they are safe
func (r *Report) areFieldsSafe(fieldsIn []int) bool {
	var lastValue int
	for index, value := range fieldsIn {
		if index == 0 {
			lastValue = value
			continue
		}

		if lastValue == value {
			return false
		}

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

// generate instantiates a new report
func (r *Report) generate(reportText string) {
	r.text = reportText
	r.dampened = false

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

	if r.areFieldsSafe(r.fields) {
		return true
	}

	for index := range r.fields {
		var reducedFields []int
		if index == 0 {
			reducedFields = r.fields[1:]
		} else if index == len(r.fields) {
			reducedFields = append(reducedFields, r.fields[0:index]...)
		} else {
			reducedFields = append(reducedFields, r.fields[0:index]...)
			reducedFields = append(reducedFields, r.fields[index+1:]...)
		}
		if r.areFieldsSafe(reducedFields) {
			return true
		}
	}

	return false
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

// setDirection determines the reports direction from the sum of the changes
func (r *Report) setDirection() {
	var directionality, lastValue int
	for index, value := range r.fields {
		if index == 0 {
			lastValue = value
			continue
		}
		directionality += lastValue - value
		lastValue = value
	}

	if directionality > 0 {
		r.direction = DESCENDING
	} else if directionality < 0 {
		r.direction = ASCENDING
	} else {
		r.direction = STATIC
	}
}
