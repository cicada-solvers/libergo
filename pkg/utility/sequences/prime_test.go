package sequences

import (
	"math/big"
	"testing"
)

func TestIsPrime(t *testing.T) {
	tests := []struct {
		number   *big.Int
		expected bool
	}{
		{big.NewInt(1), false},
		{big.NewInt(2), true},
		{big.NewInt(3), true},
		{big.NewInt(4), false},
		{big.NewInt(17), true},
		{big.NewInt(18), false},
	}

	for _, test := range tests {
		result := IsPrime(test.number)
		if result != test.expected {
			t.Errorf("IsPrime(%v) = %v; want %v", test.number, result, test.expected)
		}
	}
}

func TestGetPrimeSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(10),
			false,
			[]*big.Int{
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(5),
				big.NewInt(7),
			},
		},
		{
			big.NewInt(3),
			true,
			[]*big.Int{
				big.NewInt(7),
			},
		},
	}

	for _, test := range tests {
		result, err := GetPrimeSequence(test.maxNumber, test.isPositional)
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

func TestGetFibonacciPrimeSequence(t *testing.T) {
	tests := []struct {
		maxNumber    *big.Int
		isPositional bool
		expected     []*big.Int
	}{
		{
			big.NewInt(10),
			false,
			[]*big.Int{
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(5),
			},
		},
		{
			big.NewInt(2),
			true,
			[]*big.Int{
				big.NewInt(5),
			},
		},
	}

	for _, test := range tests {
		result, err := GetFibonacciPrimeSequence(test.maxNumber, test.isPositional)
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
