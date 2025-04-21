package sequences

import (
	"fmt"
	"math"
	"math/big"
)

func GetCollatzSequence(n int64, isPosition bool) (*NumericSequence, error) {
	sequence := &NumericSequence{Name: "Collatz", Number: big.NewInt(n)}
	longestSequence := int64(0)

	if isPosition {
		for i := int64(1); i <= math.MaxInt32; i++ {
			sequence = &NumericSequence{Name: "Collatz", Number: big.NewInt(n)}
			sequence.Sequence = append(sequence.Sequence, big.NewInt(i))
			sequence, _ = getCollatzSequenceInternal(i, sequence)

			if int64(len(sequence.Sequence)) > longestSequence {
				fmt.Printf("Sequence %d - %d\n", i, int64(len(sequence.Sequence)))
				longestSequence = int64(len(sequence.Sequence))
			}

			if n == int64(len(sequence.Sequence)) {
				return sequence, nil
			}
		}

		fmt.Printf("Length not found for %d\n", n)
		return sequence, nil
	} else {
		sequence.Sequence = append(sequence.Sequence, big.NewInt(n))
		sequence, _ = getCollatzSequenceInternal(n, sequence)
	}

	return sequence, nil
}

func getCollatzSequenceInternal(n int64, sequence *NumericSequence) (*NumericSequence, error) {
	// Generate the Collatz sequence
	if n > 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		sequence.Sequence = append(sequence.Sequence, big.NewInt(n))
		return getCollatzSequenceInternal(n, sequence)
	} else if n < 1 {
		// Stop when n reaches 1
		err := fmt.Errorf("number must be greater than 1")
		return nil, err
	} else {
		// If n is 1, we stop the recursion
		return sequence, nil
	}
}
