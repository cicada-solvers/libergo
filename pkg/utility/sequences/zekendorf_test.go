package sequences

import (
	"math/big"
	"testing"
)

func TestGetFibonacciSequence(t *testing.T) {
	tests := []struct {
		maxNumber *big.Int
		expected  []*big.Int
	}{
		{
			big.NewInt(10),
			[]*big.Int{
				big.NewInt(1),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(5),
				big.NewInt(8),
			},
		},
		{
			big.NewInt(1),
			[]*big.Int{
				big.NewInt(1),
				big.NewInt(1),
			},
		},
	}

	for _, test := range tests {
		result, err := GetFibonacciSequence(test.maxNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Sequence) != len(test.expected) {
			t.Fatalf("expected length %d, got %d", len(test.expected), len(result.Sequence))
		}
		for i, v := range result.Sequence {
			if v.Cmp(test.expected[i]) != 0 {
				t.Errorf("expected %v, got %v", test.expected[i], v)
			}
		}
	}
}

func TestGetZekendorfRepresentationSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(10),
			false,
			[]*big.Int{
				big.NewInt(8),
				big.NewInt(2),
			},
		},
		{
			big.NewInt(10),
			true,
			[]*big.Int{
				big.NewInt(8),
				big.NewInt(2),
			},
		},
	}

	for _, test := range tests {
		result, err := GetZekendorfRepresentationSequence(test.maxNumber, test.isPositional)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Sequence) != len(test.expected) {
			t.Fatalf("expected %v, got %v", test.expected, result.Sequence)
		}
		for i, v := range result.Sequence {
			if v.Cmp(test.expected[i]) != 0 {
				t.Errorf("expected %v, got %v", test.expected[i], v)
			}
		}
	}
}
