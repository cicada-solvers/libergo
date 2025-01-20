package sequences

import (
	"math/big"
	"testing"
)

func TestGetCatalanSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(5),
			false,
			[]*big.Int{
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(5),
			},
		},
		{
			big.NewInt(3),
			true,
			[]*big.Int{
				big.NewInt(5),
			},
		},
	}

	for _, test := range tests {
		result, err := GetCatalanSequence(test.maxNumber, test.isPositional)
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

func TestGetCakeSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(5),
			false,
			[]*big.Int{
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(4),
				big.NewInt(8),
				big.NewInt(15),
				big.NewInt(26),
			},
		},
		{
			big.NewInt(3),
			true,
			[]*big.Int{
				big.NewInt(8),
			},
		},
	}

	for _, test := range tests {
		result, err := GetCakeSequence(test.maxNumber, test.isPositional)
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
