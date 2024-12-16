package main

var example bool
var verbose bool

func main() {
	parseCLI()
	var exampleInput, challengeInput fileInput

	exampleInput = readIn("example.in")
	exampleInput.scan()
	exampleInput.cleanup()

	if !example {
		challengeInput = readIn("input.txt")
		challengeInput.scan()
		challengeInput.cleanup()
	}
}
