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
		nums           []int
		expectedSubstr string
	}{
		{
			name:           "NoPrimes",
			nums:           []int{4, 6, 8},
			expectedSubstr: "No prime numbers found",
		},
		{
			name:           "SomePrimes",
			nums:           []int{4, 5, 7},
			expectedSubstr: "5, 7 are prime numbers",
		},
		{
			name:           "Mixed",
			nums:           []int{2, 10, 11},
			expectedSubstr: "2, 11 are prime numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// checkPrimeTool doesn't use tool.Context, so we can pass nil
			var tc tool.Context = nil
			args := checkPrimeToolArgs{Nums: tt.nums}

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
