package sequences

import (
	"math/big"
	"testing"
)

func TestGetTotientSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(9),
			false,
			[]*big.Int{
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(4),
				big.NewInt(5),
				big.NewInt(7),
				big.NewInt(8),
			},
		},
		{
			big.NewInt(5),
			true,
			[]*big.Int{
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(4),
			},
		},
	}

	for _, test := range tests {
		result, err := GetTotientSequence(test.maxNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Sequence) != len(test.expected) {
			t.Fatalf("expected %v, got %v", test.expected, result.Sequence)
		}
	}
}

func TestGetTotientPrimeSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(9),
			false,
			[]*big.Int{
				big.NewInt(2),
				big.NewInt(5),
				big.NewInt(7),
			},
		},
		{
			big.NewInt(5),
			true,
			[]*big.Int{
				big.NewInt(2),
				big.NewInt(3),
			},
		},
	}

	for _, test := range tests {
		result, err := GetTotientPrimeSequence(test.maxNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Sequence) != len(test.expected) {
			t.Fatalf("expected length %v, got %v", test.expected, result.Sequence)
		}
		for i, v := range result.Sequence {
			if v.Cmp(test.expected[i]) != 0 {
				t.Errorf("expected %v, got %v", test.expected[i], v)
			}
		}
	}
}
