package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type fileInput struct {
	filename               string
	fileObj                *os.File
	rules                  []*Rule
	updates                []*Update
	validUpdateMiddleCount int
}

// cleanup closes the open file handler
func (f *fileInput) cleanup() {
	f.fileObj.Close()
}

// openInput opens the input file for scanning
func (f *fileInput) openInput() error {
	var err error
	f.fileObj, err = os.Open(f.filename)
	if err != nil {
		return err
	}

	return nil
}

// parseFile scans the input file for rules and updates
func (f *fileInput) parseFile() {
	ruleMatch := regexp.MustCompile(`\d+\|\d+`)
	updateMatch := regexp.MustCompile(`\d+,\d+`)

	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	for scanner.Scan() {
		inputLine := scanner.Text()
		verbosef("Input line: %s\n", inputLine)

		if ruleMatch.MatchString(inputLine) {
			rule := new(Rule)
			rule.init(inputLine)
			f.rules = append(f.rules, rule)

			verbosef("\tNew Rule: %s\n", rule.string())
		} else if updateMatch.MatchString(inputLine) {
			update := new(Update)
			update.init(inputLine)
			f.updates = append(f.updates, update)

			verbosef("\tNew Update: %s\n", update.string())
		}
	}

	verbosef(
		"%q parseFile found %d rules and %d updates.\n",
		f.filename,
		len(f.rules),
		len(f.updates),
	)
}

// readIn creates a FileInput object
func readIn(filename string) (f fileInput) {
	f.filename = filename
	err := f.openInput()
	if err != nil {
		f.cleanup()
		log.Fatalf("Failed to open %q", filename)
	}

	return f
}

// report prints how many instances of string were found in the scan
func (f *fileInput) report() {
	fmt.Printf("The middle numbers of valid updates in %q tallies up to %d.\n", f.filename, f.validUpdateMiddleCount)
}

// reset returns the cursor to the beginning of the file
func (f *fileInput) reset() {
	ret, err := f.fileObj.Seek(0, 0)
	if err != nil {
		log.Panicf("Failed to seek the beginning of %q", f.filename)
	}

	verbosef("%q reset, cursor: %v\n", f.filename, ret)
}

// scan validates the update file
func (f *fileInput) scan() {
	f.parseFile()
	f.validateUpdates()
	f.sumMiddlePages()
	f.report()
}

// sumMiddlePages cycles through valid updates and tallies the middle page number
func (f *fileInput) sumMiddlePages() {
	var tally int
	for _, update := range f.updates {
		if update.wasInvalid {
			index := int(len(update.pages) / 2)
			verbosef("\tTally: %d; Add %d; New Tally: %d\n", tally, update.pages[index], tally+update.pages[index])
			tally += update.pages[index]
		}
	}
	f.validUpdateMiddleCount = tally
}

// validateUpdates cycles through updates to check compliance with rules
func (f *fileInput) validateUpdates() {
	for _, update := range f.updates {
		verbosef("Update pages: %s.\n", update.string())
		update.validate(update.pages, f.rules)
		verbosef("Update wasInvalid: %t.\n", update.wasInvalid)
	}
}
