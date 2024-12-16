package main

import (
	"fmt"
	"log"
	"os"
)

type fileInput struct {
	filename  string
	fileObj   *os.File
	board     *Board
	loopCount int
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
	f.board = new(Board)
	f.board.scan(&f)
	f.board.save()
	f.board.gameActive = true

	return f
}

// report prints how many instances of string were found in the scan
func (f *fileInput) report() {
	fmt.Printf("%d locations will cause the guard to loop in %q.\n", f.loopCount, f.filename)
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
func (f *fileInput) predictGuardPath() {
	f.loopCount += f.board.checkEveryObstacle()

	f.report()
}
