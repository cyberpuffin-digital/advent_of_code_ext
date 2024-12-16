package main

import (
	"fmt"
	"log"
	"os"
)

type fileInput struct {
	filename     string
	fileObj      *os.File
	gameboard    GameBoard
	scan_results int
	scan_string  string
	scanned      bool
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

	return f
}

// report prints how many instances of string were found in the scan
func (f *fileInput) report() {
	if f.scanned {
		fmt.Printf("%q was found in %q %d times.\n", f.scan_string, f.filename, f.scan_results)
	} else {
		verbosef("Report called on %q before a scan was run.\n", f.filename)
	}
}

// reset returns the cursor to the beginning of the file
func (f *fileInput) reset() {
	ret, err := f.fileObj.Seek(0, 0)
	if err != nil {
		log.Panicf("Failed to seek the beginning of %q", f.filename)
	}

	verbosef("%q reset, cursor: %v\n", f.filename, ret)
}

// scan searches a file for a string
func (f *fileInput) scan(needle string) {
	f.scan_string = needle

	// Build the game board
	f.gameboard.scan(f)

	// Scan the gameboard
	f.scan_results = f.gameboard.scanForString(f.scan_string)

	// Apply and report the results
	f.scanned = true
	f.report()
}
