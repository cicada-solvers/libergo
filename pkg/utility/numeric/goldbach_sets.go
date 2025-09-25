package numeric

import (
	"sequences"
)

// NewGoldbachSets creates and returns a new instance of GoldbachPairs with an empty list of pairs.
func NewGoldbachSets() *GoldbachSets {
	return &GoldbachSets{}
}

// GoldbachSets represents a collection of GoldbachPair, storing pairs for computations related to Goldbach's conjecture.
type GoldbachSets struct {
	GoldBachSets []GoldbachSet
}

// GoldbachSet represents a pair of integers whose sum equals a given even number, based on Goldbach's conjecture.
// Number specifies the even number, while AddendOne and AddendTwo are the integers forming the sum.
type GoldbachSet struct {
	Number      int64
	AddendOne   int64
	AddendTwo   int64
	AddendThree int64
}

// Solve finds all Goldbach sets for the given number.
func (g *GoldbachSets) Solve(number int64) {
	var primeNumbers []int64

	// Get all prime numbers up to the number
	for i := int64(2); i < number; i++ {
		if sequences.IsPrime64(i) {
			primeNumbers = append(primeNumbers, i)
		}
	}

	g.GetGoldBachSets(number, primeNumbers)
}

// GetGoldBachSets finds all Goldbach sets for the given number and prime set.
func (g *GoldbachSets) GetGoldBachSets(number int64, primeSet []int64) {
	for _, prime := range primeSet {
		firstSet := []int64{prime}
		g.GetNextPrime(number, firstSet, primeSet)
	}
}

// GetNextPrime recursively finds the next prime number in the given prime set.
func (g *GoldbachSets) GetNextPrime(number int64, currentSet, primeSet []int64) {
	if len(currentSet) < 3 {
		for _, prime := range primeSet {
			currentSet = append(currentSet, prime)
			g.GetNextPrime(number, currentSet, primeSet)
			currentSet = currentSet[:len(currentSet)-1]
		}
	} else {
		if currentSet[0]+currentSet[1]+currentSet[2] == number {
			currentSet = g.SortAddendValues(currentSet[0], currentSet[1], currentSet[2])
			completeSet := GoldbachSet{
				Number:      number,
				AddendOne:   currentSet[0],
				AddendTwo:   currentSet[1],
				AddendThree: currentSet[2],
			}

			if !g.ContainsSetAlready(completeSet) {
				g.GoldBachSets = append(g.GoldBachSets, completeSet)
			}
		}
	}
}

// SortAddendValues sorts the given addend values.
func (g *GoldbachSets) SortAddendValues(a, b, c int64) []int64 {
	result := []int64{a, b, c}

	// Simple bubble sort for just 3 elements
	for i := 0; i < 2; i++ {
		for j := 0; j < 2-i; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	return result
}

// ContainsSetAlready checks if the given GoldbachSet is already contained in the list of GoldbachSets.
func (g *GoldbachSets) ContainsSetAlready(set GoldbachSet) bool {
	for _, goldBachSet := range g.GoldBachSets {
		if goldBachSet.Number == set.Number && goldBachSet.AddendOne == set.AddendOne && goldBachSet.AddendTwo == set.AddendTwo && goldBachSet.AddendThree == set.AddendThree {
			return true
		}
	}

	return false
}
