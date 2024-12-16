package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type fileInput struct {
	equations    []*Equation
	filename     string
	fileObj      *os.File
	highestTotal int64
	validTally   int64
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

// readIn creates a FileInput object
func readIn(filename string) (f fileInput) {
	f.filename = filename
	err := f.openInput()
	if err != nil {
		f.cleanup()
		log.Fatalf("Failed to open %q", filename)
	}

	// file init logic
	f.scanForEquations()

	return f
}

// report prints how many instances of string were found in the scan
func (f *fileInput) report() {
	fmt.Printf("%q valid operations tally up to %d.\n", f.filename, f.validTally)
}

// reset returns the cursor to the beginning of the file
func (f *fileInput) reset() {
	ret, err := f.fileObj.Seek(0, 0)
	if err != nil {
		log.Panicf("Failed to seek the beginning of %q", f.filename)
	}

	verbosef("%q reset, cursor: %v\n", f.filename, ret)
}

// scanForEquations will read the input file and build the file's equations object
func (f *fileInput) scanForEquations() {
	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	counter := 0

	for scanner.Scan() {
		inputLine := scanner.Text()
		verbosef("%s\n", inputLine)

		equationIn := new(Equation)
		equationIn.readIn(inputLine, counter)

		if equationIn.totalValue > f.highestTotal {
			f.highestTotal = equationIn.totalValue
		}

		f.equations = append(f.equations, equationIn)
		counter++
	}

	verbosef(
		"%d equations found in %q with the highest total of %d.\n",
		len(f.equations),
		f.filename,
		f.highestTotal,
	)
}

// scan validates the update file
func (f *fileInput) validateEquations() {
	for index, eq := range f.equations {
		eq.checkValidity()
		if eq.isValid {
			f.validTally += eq.totalValue
		}
		verbosef("Equation %d is valid: %t.\n", index, eq.isValid)
	}
	f.report()
}
