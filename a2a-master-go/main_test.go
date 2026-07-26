package main

import (
	"strconv"
	"testing"
)

func TestRollDieTool(t *testing.T) {
	// Test valid sides
	sides := 6
	args := rollDieToolArgs{Sides: sides}

	for i := 0; i < 50; i++ {
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

	// Test invalid sides
	invalidArgs := []rollDieToolArgs{
		{Sides: 0},
		{Sides: -3},
	}
	for _, arg := range invalidArgs {
		if _, err := rollDieTool(nil, arg); err == nil {
			t.Errorf("rollDieTool(%d sides) expected error, got nil", arg.Sides)
		}
	}
}
