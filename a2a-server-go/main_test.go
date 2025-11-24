package main

import (
	"strings"
	"testing"

	"google.golang.org/adk/tool"
)

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Negative", -5, false},
		{"Zero", 0, false},
		{"One", 1, false},
		{"Two", 2, true},
		{"Three", 3, true},
		{"Four", 4, false},
		{"Eleven", 11, true},
		{"Twenty", 20, false},
		{"LargePrime", 97, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrime(tt.input); got != tt.expected {
				t.Errorf("isPrime(%d) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCheckPrimeTool(t *testing.T) {
	tests := []struct {
		name           string
		num            int
		expectedSubstr string
	}{
		{
			name:           "NoPrimes",
			num:            4,
			expectedSubstr: "4 is not a prime number",
		},
		{
			name:           "Prime",
			num:            5,
			expectedSubstr: "5 is a prime number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// checkPrimeTool doesn't use tool.Context, so we can pass nil
			var tc tool.Context = nil
			args := checkPrimeToolArgs{Num: tt.num}

			got, err := checkPrimeTool(tc, args)
			if err != nil {
				t.Fatalf("checkPrimeTool returned unexpected error: %v", err)
			}

			if !strings.Contains(got, tt.expectedSubstr) {
				t.Errorf("checkPrimeTool output %q does not contain %q", got, tt.expectedSubstr)
			}
		})
	}
}
