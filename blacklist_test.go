package main

import (
	"testing"
)

func TestIsBlacklisted(t *testing.T) {
	// Test with a name that is in the blacklist
	if !isBlacklisted("John") {
		t.Errorf("Expected 'John' to be blacklisted, but it was not")
	}

	// Test with a name that is not in the blacklist
	if isBlacklisted("Jane") {
		t.Errorf("Expected 'Jane' to not be blacklisted, but it was")
	}
}
