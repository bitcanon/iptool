package utils_test

import (
	"testing"

	"github.com/bitcanon/iptool/utils"
)

// TestClosestLargerPowerOfTwo tests the ClosestLargerPowerOfTwo function
// using various input values
func TestClosestLargerPowerOfTwo2(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    int
		expected uint32
	}{
		{name: "Zero", input: 0, expected: 1},
		{name: "One", input: 1, expected: 1},
		{name: "Two", input: 2, expected: 2},
		{name: "Three", input: 3, expected: 4},
		{name: "Four", input: 4, expected: 4},
		{name: "Five", input: 5, expected: 8},
		{name: "Six", input: 6, expected: 8},
		{name: "Seven", input: 7, expected: 8},
		{name: "Eight", input: 8, expected: 8},
		{name: "Nine", input: 9, expected: 16},
		{name: "Ten", input: 10, expected: 16},
		{name: "Eleven", input: 11, expected: 16},
		{name: "Twelve", input: 12, expected: 16},
		{name: "Thirteen", input: 13, expected: 16},
		{name: "Fourteen", input: 14, expected: 16},
		{name: "Fifteen", input: 15, expected: 16},
		{name: "Sixteen", input: 16, expected: 16},
		{name: "Seventeen", input: 17, expected: 32},
		{name: "Eighteen", input: 18, expected: 32},
		{name: "ThirtyTwo", input: 32, expected: 32},
		{name: "ThirtyThree", input: 33, expected: 64},
		{name: "SixtyFour", input: 64, expected: 64},
		{name: "SixtyFive", input: 65, expected: 128},
		{name: "OneHundred", input: 100, expected: 128},
		{name: "OneHundredTwentyEight", input: 128, expected: 128},
		{name: "OneHundredTwentyNine", input: 129, expected: 256},
		{name: "TwoHundredFiftyFive", input: 255, expected: 256},
		{name: "TwoHundredFiftySix", input: 256, expected: 256},
	}

	// Loop through test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Calculate the closest larger power of two
			result := utils.ClosestLargerPowerOfTwo(testCase.input)

			// Check if the result matches the expected value
			if result != testCase.expected {
				t.Errorf("expected: %d, got: %d", testCase.expected, result)
			}
		})
	}
}
