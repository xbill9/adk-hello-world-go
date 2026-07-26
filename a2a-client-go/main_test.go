package main

import (
	"testing"
)

func TestRollDieTool(t *testing.T) {
	tests := []struct {
		name    string
		sides   int
		wantErr bool
	}{
		{"6 sides", 6, false},
		{"20 sides", 20, false},
		{"1 side", 1, false},
		{"100 sides", 100, false},
		{"0 sides", 0, true},
		{"negative sides", -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := rollDieToolArgs{Sides: tt.sides}
			// Passing nil for tool.Context as it is not used in the function
			got, err := rollDieTool(nil, args)
			if (err != nil) != tt.wantErr {
				t.Errorf("rollDieTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got < 1 || got > tt.sides) {
				t.Errorf("rollDieTool() = %v, want between 1 and %v", got, tt.sides)
			}
		})
	}
}
