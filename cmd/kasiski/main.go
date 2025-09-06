package main

import (
	runelib "characterrepo"
	"errors"
	"flag"
	"fmt"
	"sort"
	"strings"
)

var charRepo *runelib.CharacterRepo

func main() {
	charRepo = runelib.NewCharacterRepo()
	textFlag := flag.String("text", "", "Text to analyze")
	flag.Parse()

	if *textFlag == "" {
		flag.Usage()
		return
	}

	text := strings.Split(strings.ToUpper(*textFlag), "")

	length, _ := determineKeyLength(text)
	fmt.Printf("Key length: %d\n", length)
}

// Kasiski examination to estimate Vigenère key length.
// 1) Normalize text (letters only, uppercase)
// 2) Find repeated n-grams (n in [3..5])
// 3) Record distances between repeated n-gram positions
// 4) Use GCD of distances; if GCD < 2, fallback to factor frequency
func determineKeyLength(text []string) (int, error) {
	norm := normalizeText(text)
	if len(norm) < 6 {
		return 0, errors.New("text too short for Kasiski examination")
	}

	// Collect distances between repeated n-grams
	var distances []int
	for n := 5; n >= 3; n-- { // prefer longer n-grams first
		pos := make(map[string][]int)
		for i := 0; i+n <= len(norm); i++ {
			sub := norm[i : i+n]
			pos[sub] = append(pos[sub], i)
		}
		for _, indices := range pos {
			if len(indices) < 2 {
				continue
			}
			for i := 1; i < len(indices); i++ {
				d := indices[i] - indices[i-1]
				if d > 0 {
					distances = append(distances, d)
				}
			}
		}
	}
	if len(distances) == 0 {
		return 0, errors.New("no repeated n-grams found; cannot estimate key length")
	}

	// First attempt: GCD of all distances
	if g := gcdSlice(distances); g >= 2 {
		return g, nil
	}

	// Fallback: factor frequency analysis
	// Count factors (>=2) of all distances and choose the most frequent
	freq := map[int]int{}
	for _, d := range distances {
		for _, f := range properFactors(d) {
			if f >= 2 {
				freq[f]++
			}
		}
	}
	if len(freq) == 0 {
		return 0, errors.New("unable to infer key length")
	}

	// Select factor with highest frequency; tie-breaker: smaller factor
	type kv struct {
		f int
		c int
	}
	var arr []kv
	for f, c := range freq {
		arr = append(arr, kv{f, c})
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].c == arr[j].c {
			return arr[i].f < arr[j].f
		}
		return arr[i].c > arr[j].c
	})
	if arr[0].f < 2 {
		return 0, errors.New("inconclusive factor analysis")
	}
	return arr[0].f, nil
}

// ... existing code ...
func normalizeText(s []string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if charRepo.IsRune(r, false) || charRepo.IsLetterInAlphabet(r) {
			b.WriteString(r)
		}
	}
	return b.String()
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func gcdSlice(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	g := nums[0]
	for i := 1; i < len(nums); i++ {
		g = gcd(g, nums[i])
		if g == 1 {
			return 1
		}
	}
	return g
}

func properFactors(n int) []int {
	if n <= 1 {
		return nil
	}
	// Collect all factors excluding 1 and n
	var res []int
	for f := 2; f*f <= n; f++ {
		if n%f == 0 {
			res = append(res, f)
			if f != n/f {
				res = append(res, n/f)
			}
		}
	}
	// Optional: include n itself, as some keys can equal distance,
	// but usually key length << distance; keep it excluded to reduce noise.
	sort.Ints(res)
	return res
}
