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

// scanInput reads an input file line by line to create separate lists
func (f *fileInput) countSafe() (safeReports int) {
	scanner := bufio.NewScanner(f.fileObj)

	for scanner.Scan() {
		var report Report
		report.generate(scanner.Text())

		if report.isSafe() {
			safeReports++
		}
	}

	return safeReports
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
