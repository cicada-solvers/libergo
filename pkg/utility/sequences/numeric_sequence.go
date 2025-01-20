package sequences

import "math/big"

// NumericSequence represents a sequence of numbers.
type NumericSequence struct {
	Name     string
	Number   *big.Int
	Sequence []*big.Int
	Result   *big.Int
}
