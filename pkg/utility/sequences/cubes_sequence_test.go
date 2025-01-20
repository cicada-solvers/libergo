package sequences

import (
	"math/big"
	"testing"
)

func TestGetCubesSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(5),
			false,
			[]*big.Int{
				big.NewInt(0),
				big.NewInt(1),
				big.NewInt(8),
				big.NewInt(27),
				big.NewInt(64),
				big.NewInt(125),
			},
		},
		{
			big.NewInt(3),
			true,
			[]*big.Int{big.NewInt(27)},
		},
	}

	for _, test := range tests {
		result, err := GetCubesSequence(test.maxNumber, test.isPositional)
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
