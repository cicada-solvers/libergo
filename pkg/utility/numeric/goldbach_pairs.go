package numeric

import "fmt"

// NewGoldbachPairs creates and returns a new instance of GoldbachPairs with an empty list of pairs.
func NewGoldbachPairs() *GoldbachPairs {
	return &GoldbachPairs{}
}

// GoldbachPairs represents a collection of GoldbachPair, storing pairs for computations related to Goldbach's conjecture.
type GoldbachPairs struct {
	GoldBachPairs []GoldbachPair
}

// GoldbachPair represents a pair of integers whose sum equals a given even number, based on Goldbach's conjecture.
// Number specifies the even number, while AddendOne and AddendTwo are the integers forming the sum.
type GoldbachPair struct {
	Number    int64
	AddendOne int64
	AddendTwo int64
}

// Add inserts a GoldbachPair into the GoldbachPairs collection if it does not already exist.
func (g *GoldbachPairs) Add(pair GoldbachPair) {
	if !g.checkIfPairExists(pair) {
		g.GoldBachPairs = append(g.GoldBachPairs, pair)
	}
}

// GetGoldbachPairs retrieves the collection of Goldbach pairs stored in the GoldbachPairs instance.
func (g *GoldbachPairs) GetGoldbachPairs() []GoldbachPair {
	return g.GoldBachPairs
}

// SolveForNumber computes all valid Goldbach pairs for an even number using a list of primes and adds them to the collection.
// Returns an error if the input number is not even.
func (g *GoldbachPairs) SolveForNumber(number int64, primeList *[]int64) error {
	if !IsNumberEven(number) {
		return fmt.Errorf("number %d is not even", number)
	}

	offset := 0
	numberTwo := int64(0)
	for _, prime := range *primeList {
		numberTwo, offset = g.getNumberTwo(number, prime, offset, primeList)
		if numberTwo > 0 {
			one, two := g.sortNumbersInPair(numberTwo, prime)
			g.Add(GoldbachPair{
				Number:    number,
				AddendOne: one,
				AddendTwo: two,
			})
		}
	}

	return nil
}

// getNumberTwo calculates an integer (numberTwo) such that addendOne + numberTwo equals a specified number using a prime list.
// Returns the calculated numberTwo and the updated offset value. If no match is found, returns 0 and the current offset.
func (g *GoldbachPairs) getNumberTwo(number, addendOne int64, offset int, primeList *[]int64) (int64, int) {
	newOffset := offset
	endIndex := len(*primeList) - 1 - offset

	if endIndex < 0 {
		return 0, newOffset
	}

	counter := 0
	for i := endIndex; i >= 0; i-- {
		counter++
		prime := (*primeList)[i]
		numberTwo := addendOne + prime
		if numberTwo == number {
			newOffset = newOffset + counter
			return numberTwo, newOffset
		}
	}

	return 0, newOffset
}

// sortNumbersInPair takes two integers and returns them in ascending order.
func (g *GoldbachPairs) sortNumbersInPair(one, two int64) (int64, int64) {
	if one > two {
		return two, one
	}

	return one, two
}

// checkIfPairExists checks if a given GoldbachPair already exists in the GoldbachPairs collection. Returns true if it does.
func (g *GoldbachPairs) checkIfPairExists(pair GoldbachPair) bool {
	for _, existingPair := range g.GoldBachPairs {
		if existingPair.AddendTwo == pair.AddendOne && existingPair.AddendTwo == pair.AddendTwo {
			return true
		}
	}

	return false
}
