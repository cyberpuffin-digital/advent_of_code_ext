package main

import (
	"fmt"
	"log"
	"os"
)

type fileInput struct {
	filename   string
	fileObj    *os.File
	pebbleList *PebbleList
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
	f.pebbleList = new(PebbleList)
	f.pebbleList.readIn(&f)
	f.pebbleList.blink(25)

	return f
}

// report prints how many unique antinode locations are in the file.
func (f *fileInput) report() {
	printOut(fmt.Sprintf("%q has %d stones after %d blinks.\n", f.filename, f.pebbleList.count(), f.pebbleList.blinks))
}

// reset returns the cursor to the beginning of the file
func (f *fileInput) reset() {
	ret, err := f.fileObj.Seek(0, 0)
	if err != nil {
		log.Panicf("Failed to seek the beginning of %q", f.filename)
	}

	verbosef("%q reset, cursor: %v\n", f.filename, ret)
}
