package numeric

import "fmt"

func NewGoldbachPairs() *GoldbachPairs {
	return &GoldbachPairs{}
}

type GoldbachPairs struct {
	GoldBachPairs []GoldbachPair
}

type GoldbachPair struct {
	Number    int64
	AddendOne int64
	AddendTwo int64
}

func (g *GoldbachPairs) Add(pair GoldbachPair) {
	if !g.checkIfPairExists(pair) {
		g.GoldBachPairs = append(g.GoldBachPairs, pair)
	}
}

func (g *GoldbachPairs) GetGoldbachPairs() []GoldbachPair {
	return g.GoldBachPairs
}

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

func (g *GoldbachPairs) sortNumbersInPair(one, two int64) (int64, int64) {
	if one > two {
		return two, one
	}

	return one, two
}

func (g *GoldbachPairs) checkIfPairExists(pair GoldbachPair) bool {
	for _, existingPair := range g.GoldBachPairs {
		if existingPair.AddendTwo == pair.AddendOne && existingPair.AddendTwo == pair.AddendTwo {
			return true
		}
	}

	return false
}
