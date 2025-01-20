package runelib

import (
	"testing"
)

func TestCharacterRepo(t *testing.T) {
	repo := NewCharacterRepo()

	if char := repo.GetANSICharFromDec(65, true); char != "A" {
		t.Errorf("Expected 'A', got '%s'", char)
	}

	if char := repo.GetASCIICharFromDec(65, true); char != "A" {
		t.Errorf("Expected 'A', got '%s'", char)
	}
}
