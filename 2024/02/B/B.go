package main

import (
	"fmt"
)

func main() {
	var exampleInput, challengeInput fileInput

	exampleInput = readIn("example.in")
	fmt.Printf("%q has %d safe reports.\n", exampleInput.filename, exampleInput.countSafe())
	exampleInput.cleanup()

	challengeInput = readIn("input.txt")
	fmt.Printf("%q has %d safe reports.\n", challengeInput.filename, challengeInput.countSafe())
	challengeInput.cleanup()
}
