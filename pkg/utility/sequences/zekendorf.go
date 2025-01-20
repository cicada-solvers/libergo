package sequences

import (
	"errors"
	"math/big"
)

// GetFibonacciSequence generates the Fibonacci sequence up to maxNumber.
func GetFibonacciSequence(maxNumber *big.Int) (*NumericSequence, error) {
	if maxNumber.Cmp(big.NewInt(1)) < 0 {
		return nil, errors.New("maxNumber must be greater than 0")
	}

	sequence := []*big.Int{big.NewInt(1), big.NewInt(1)}
	for {
		next := new(big.Int).Add(sequence[len(sequence)-1], sequence[len(sequence)-2])
		if next.Cmp(maxNumber) > 0 {
			break
		}
		sequence = append(sequence, next)
	}

	return &NumericSequence{
		Name:     "Fibonacci",
		Number:   maxNumber,
		Sequence: sequence,
	}, nil
}

// GetZekendorfRepresentationSequence generates the Zekendorf Representation sequence.
func GetZekendorfRepresentationSequence(maxNumber *big.Int, isPositional bool) (*NumericSequence, error) {
	retval := &NumericSequence{Name: "Zekendorf Representation", Number: new(big.Int).Set(maxNumber)}
	remainder := new(big.Int).Set(maxNumber)

	for remainder.Cmp(big.NewInt(0)) > 0 {
		fibSequence, err := GetFibonacciSequence(remainder)
		if err != nil {
			return nil, err
		}
		lastFib := fibSequence.Sequence[len(fibSequence.Sequence)-1]
		retval.Sequence = append(retval.Sequence, lastFib)
		remainder.Sub(remainder, lastFib)

		if isPositional && new(big.Int).SetUint64(uint64(len(retval.Sequence))).Cmp(maxNumber) > 0 {
			retval.Sequence = []*big.Int{retval.Sequence[maxNumber.Uint64()]}
			break
		}
	}

	return retval, nil
}
