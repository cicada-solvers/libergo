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
	g.GoldBachPairs = append(g.GoldBachPairs, pair)
}

func (g *GoldbachPairs) GetGoldbachPairs() []GoldbachPair {
	return g.GoldBachPairs
}

func (g *GoldbachPairs) SolveForNumber(number int64, primeList *[]int64) error {
	if !IsNumberEven(number) {
		return fmt.Errorf("number %d is not even", number)
	}

	for _, prime := range *primeList {
		numberTwo := g.getNumberTwo(number, prime, primeList)
		if numberTwo > 0 {
			g.Add(GoldbachPair{
				Number:    number,
				AddendOne: prime,
				AddendTwo: numberTwo,
			})
		}
	}

	return nil
}

func (g *GoldbachPairs) getNumberTwo(number, addendOne int64, primeList *[]int64) int64 {
	for _, prime := range *primeList {
		numberTwo := addendOne + prime
		if numberTwo == number {
			return numberTwo
		}
	}

	return 0
}
