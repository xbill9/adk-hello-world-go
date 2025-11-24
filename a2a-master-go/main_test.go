package main

import (
	"strconv"
	"testing"
)

func TestRollDieTool(t *testing.T) {
	// Since rollDieTool uses rand.Intn, the output is random.
	// We can test if the output is a valid integer and within the range.

	sides := 6
	args := rollDieToolArgs{Sides: sides}

	// Run multiple times to ensure it stays within bounds
	for i := 0; i < 100; i++ {
		resultStr, err := rollDieTool(nil, args)
		if err != nil {
			t.Fatalf("rollDieTool returned error: %v", err)
		}

		result, err := strconv.Atoi(resultStr)
		if err != nil {
			t.Fatalf("rollDieTool returned non-integer string: %s", resultStr)
		}

		if result < 1 || result > sides {
			t.Errorf("rollDieTool result %d out of bounds [1, %d]", result, sides)
		}
	}
}
