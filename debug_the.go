package main

import (
	"fmt"
	"strings"
)

func testTheLogic() {
	questionDesc := "last response status code"
	fmt.Printf("Original: %q\n", questionDesc)

	// Clean up "last response" to make it more readable
	if strings.HasPrefix(questionDesc, "last response ") {
		questionDesc = "the " + questionDesc
	}

	fmt.Printf("After 'the' addition: %q\n", questionDesc)

	// Check what happens if we apply capitalization logic
	if len(questionDesc) > 0 && questionDesc[0] >= 'a' && questionDesc[0] <= 'z' {
		capitalized := strings.ToUpper(questionDesc[:1]) + questionDesc[1:]
		fmt.Printf("After capitalization: %q\n", capitalized)
	}

	// Check if there's weird character replacement
	fmt.Printf("Hex bytes: % x\n", []byte(questionDesc))
}

func main() {
	testTheLogic()
}
