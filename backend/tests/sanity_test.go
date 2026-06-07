package tests

import "testing"

func TestAddSanity(t *testing.T) {
	observed := 2 + 3
	expected := 5

	if observed != expected {
		t.Fatalf("Sanity check failed: 2 + 3 = %d; want %d", observed, expected)
	}
}
