package main

import (
	"os"
	"testing"
)

func TestIsOverAppLimit(t *testing.T) {
	// Create a temporary file for testing
	f, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer f.Close()

	// Write some test data to the file
	testData := []string{
		"1, AppA, John, 1.0, 2022-03-13T00:00:00Z",
		"2, AppB, Jane, 1.0, 2022-03-13T00:00:00Z",
		"3, AppA, John, 1.0, 2022-03-12T23:59:59Z",
		"4, AppA, John, 1.0, 2022-03-12T23:59:58Z",
		"5, AppA, John, 1.0, 2022-03-12T23:59:57Z",
		"6, AppA, John, 1.0, 2022-03-12T23:59:56Z",
		"7, AppA, John, 1.0, 2022-03-12T23:59:55Z",
		"8, AppA, John, 1.0, 2022-03-12T23:59:54Z",
	}
	for _, line := range testData {
		if _, err := f.WriteString(line + "\n"); err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
	}

	// Test with a name that is over the limit
	if isOverAppLimit("John", 4) {
		t.Errorf("Expected 'John' to be under the limit, but it was not")
	}

	// Test with a name that is under the limit
	if isOverAppLimit("Jane", 5) {
		t.Errorf("Expected 'Jane' to not be over the limit, but it was")
	}

	// Test with a name that is exactly at the limit
	if isOverAppLimit("John", 6) {
		t.Errorf("Expected 'John' to be at the limit, but it was not")
	}

	// Test with a limit of 0, which should always return false
	if isOverAppLimit("John", 0) {
		t.Errorf("Expected 'John' to be under the limit, but it was not")
	}
}
