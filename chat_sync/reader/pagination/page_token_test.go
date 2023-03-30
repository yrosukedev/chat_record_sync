package pagination

import "testing"

func TestPageTokenEquality_compareByPointers(t *testing.T) {
	// Given
	token1 := NewPageToken(123)
	token2 := NewPageToken(123)

	// When
	if token1 == token2 {
		t.Error("pointers shouldn't be matched")
	}
}

func TestPageTokenEquality_compareByValues(t *testing.T) {
	// Given
	token1 := NewPageToken(123)
	token2 := NewPageToken(123)

	// When
	if *token1 != *token2 {
		t.Error("values should be matched")
	}
}
