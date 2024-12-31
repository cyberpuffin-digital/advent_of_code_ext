package main

import (
	"bufio"
	"fmt"
	"log"
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
	cursor := 0
	for endCursor := len(d.fileSystem) - 1; endCursor > cursor; endCursor-- {
		if d.fileSystem[endCursor] == "." {
			continue
		}
		for startCursor := 0; startCursor <= endCursor; startCursor++ {
			if d.fileSystem[startCursor] == "." {
				d.fileSystem[startCursor] = d.fileSystem[endCursor]
				d.fileSystem[endCursor] = "."
				break
			}
		}
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
