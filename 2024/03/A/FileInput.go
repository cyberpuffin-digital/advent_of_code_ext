package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

type fileInput struct {
	filename string
	fileObj  *os.File
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
		f.fileObj.Close()
		log.Fatalf("Failed to open %q", filename)
	}

	return f
}

// tally adds multiples found in the corrupted memory
func (f *fileInput) tally() (instructionTally int) {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	scanner := bufio.NewScanner(f.fileObj)
	instructionTally = 0

	for scanner.Scan() {
		inputLine := scanner.Text()
		match := re.FindAllStringSubmatch(inputLine, -1)
		for _, instruction := range match {
			var x, y int
			var err error
			x, err = strconv.Atoi(instruction[1])
			if err != nil {
				log.Fatalf("Failed to convert digit %q", instruction[1])
			}
			y, err = strconv.Atoi(instruction[2])
			if err != nil {
				log.Fatalf("Failed to convert digit %q", instruction[2])
			}

			instructionTally += x * y
		}
	}

	return instructionTally
}
