package main

import (
	"bufio"
	"log"
	"os"
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
	enabled := true
	inputLine := ""
	instructionTally = 0
	scanner := bufio.NewScanner(f.fileObj)

	for scanner.Scan() {
		lineTally := 0
		inputLine = scanner.Text()
		lineTally, enabled = tallyLine(inputLine, enabled)
		instructionTally += lineTally
	}

	return instructionTally
}
