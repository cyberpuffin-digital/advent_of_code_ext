package main

type LocationList struct {
	countCache map[int]int
	locations  []int
}

// countOccurrences scans through the slice of integers tallying occurrences
func (l *LocationList) countOccurrences(needle int) (count int) {
	// Return cached entry if it is set
	if _, ok := l.countCache[needle]; ok {
		return l.countCache[needle]
	}

	for _, value := range l.locations {
		if needle == value {
			count++
		}
	}

	l.countCache[needle] = count
	return count
}
