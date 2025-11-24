package main

import (
	"testing"
)

func TestRollDieTool(t *testing.T) {
	tests := []struct {
		name  string
		sides int
	}{
		{"6 sides", 6},
		{"20 sides", 20},
		{"1 side", 1},
		{"100 sides", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := rollDieToolArgs{Sides: tt.sides}
			// Passing nil for tool.Context as it is not used in the function
			got, err := rollDieTool(nil, args)
			if err != nil {
				t.Errorf("rollDieTool() error = %v", err)
				return
			}
			if got < 1 || got > tt.sides {
				t.Errorf("rollDieTool() = %v, want between 1 and %v", got, tt.sides)
			}
		})
	}
}
