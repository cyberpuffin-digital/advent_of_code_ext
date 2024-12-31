package main

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"strconv"
)

type DiskMap struct {
	checksum   int64
	diskMap    string
	fileSystem []string
}

// calculateChecksum calculates a checksum for filesystem
func (d *DiskMap) calculateChecksum() {
	tally := int64(0)

	for cursor := len(d.fileSystem) - 1; cursor > 0; cursor-- {
		if d.fileSystem[cursor] == "." {
			continue
		}
		fileId, err := strconv.Atoi(d.fileSystem[cursor])
		if err != nil {
			log.Panicf("Failed to parse %q as fileID.\n", d.fileSystem[cursor])
		}
		tally += int64(fileId * cursor)
	}

	d.checksum = tally
}

// compactFileSystem moves data from the end to free space
func (d *DiskMap) compactFileSystem() {
	var compareID int
	var currentFile []string
	var fileID, lastID string
	var moved bool
	var scanCursor int

	verbosef("Compact filesystem, size: %d.\n", len(d.fileSystem))
	for endCursor := len(d.fileSystem) - 1; endCursor >= 0; endCursor-- {
		if d.fileSystem[endCursor] == "." {
			verbosef("Skipping cursor %d.\n", endCursor)
			continue
		}

		// Check if the file ID at the cursor is smaller than the last ID
		cursorID, err := strconv.Atoi(d.fileSystem[endCursor])
		if err != nil {
			log.Panicf("Failed to parse %q.\n", d.fileSystem[endCursor])
		}

		if lastID == "" {
			compareID = cursorID + 1
		} else {
			compareID, err = strconv.Atoi(lastID)
			if err != nil {
				log.Panicf("Failed to parse %q.\n", lastID)
			}
		}

		if cursorID < compareID {
			fileID = d.fileSystem[endCursor]
		} else {
			continue
		}

		// Read in current fileID
		currentFile = []string{}
		moved = false
		lastID = fileID
		scanCursor = endCursor
		for keepGoing := true; keepGoing; keepGoing = d.fileSystem[scanCursor] == fileID {
			if d.fileSystem[scanCursor] == fileID {
				currentFile = append(currentFile, fileID)
			}
			scanCursor--
			if scanCursor < 0 {
				break
			}
		}
		verbosef("File ID %s, length of %d, read in. Cursor position: %d.\n", fileID, len(currentFile), endCursor)

		// Find first block of free space large enough to contain current fileID
		moveToIndex := d.getFreeSpaceIndex(len(currentFile), endCursor)

		if moveToIndex >= 0 && moveToIndex < endCursor {
			// Swap files
			for index, fID := range currentFile {
				d.fileSystem[moveToIndex+index] = fID
				d.fileSystem[endCursor-index] = "."
			}
			moved = true
		}

		if moved {
			d.print()
		} else {
			verbosef("No free space for file %s.\n", fileID)
		}
		endCursor -= (len(currentFile) - 1)
		// endCursor -= len(currentFile)
	}
}

// diskMapToFileSystem converts disk map to file system
func (d *DiskMap) diskMapToFileSystem() {
	fileId := 0
	d.fileSystem = []string{}
	freeSpace := false

	for _, char := range d.diskMap {
		charInt, err := strconv.Atoi(string(char))
		if err != nil {
			log.Panicf("Failed to parse character %q.", char)
		}
		if freeSpace {
			for range charInt {
				d.fileSystem = append(d.fileSystem, ".")
			}
		} else {
			for range charInt {
				d.fileSystem = append(d.fileSystem, fmt.Sprint(fileId))
			}
			fileId++
		}
		freeSpace = !freeSpace
	}
}

// getFreeSpaceIndex scans the file system for free space N large
func (d *DiskMap) getFreeSpaceIndex(size int, limit int) (startIndex int) {
	startIndex = -1

	for startCursor := 0; startCursor < len(d.fileSystem); startCursor++ {
		if d.fileSystem[startCursor] == "." {
			memoryChunk := slices.Clone(d.fileSystem[startCursor : startCursor+size])
			memoryChunk = slices.Compact(memoryChunk)

			if len(memoryChunk) == 1 {
				startIndex = startCursor
				break
				// } else {
				// startCursor += (size - 1)
			}

			if startCursor > limit {
				break
			}
		}
	}

	verbosef("Free space index: %d.\n", startIndex)
	return startIndex
}

// print diskmap information
func (d *DiskMap) print() {
	if len(d.diskMap) > 0 {
		verbosef("Disk map: %s\n", d.diskMap)
	}
	if len(d.fileSystem) > 0 {
		verbosef("File system: %s\n", d.fileSystem)
	}
}

// readDiskMap copies the content of the file to the DiskMap
func (d *DiskMap) readDiskMap(f *fileInput) {
	f.reset()
	scanner := bufio.NewScanner(f.fileObj)
	for scanner.Scan() {
		inputLine := scanner.Text()
		if len(inputLine) > 0 {
			d.diskMap = inputLine
		}
	}
}

// readIn reads the content of the input file into the board
func (d *DiskMap) readIn(f *fileInput) {
	// Read disk map
	d.readDiskMap(f)
	d.print()

	// Convert disk map to file system
	d.diskMapToFileSystem()
	d.print()

	// Compact file system
	d.compactFileSystem()
	d.print()

	// Calculate checksum
	d.calculateChecksum()
}
