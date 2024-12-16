package main

import (
	"fmt"
)

func main() {
	var exampleInput, challengeInput fileInput

	exampleInput = readIn("example.in")
	fmt.Printf("%q adds up to %d.\n", exampleInput.filename, exampleInput.tally())
	exampleInput.cleanup()

	challengeInput = readIn("input.txt")
	fmt.Printf("%q adds up to %d.\n", challengeInput.filename, challengeInput.tally())
	challengeInput.cleanup()
}
