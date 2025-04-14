package cipher

type AlphabetType int

const (
	Latin AlphabetType = iota
	Rune
)

// indexOf returns the index of the target string in the slice, or -1 if not found.
func indexOf(slice []string, target string) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

// isLetter checks if a rune is a letter.
func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isUpper checks if a rune is an uppercase letter.
func isUpper(c rune) bool {
	return c >= 'A' && c <= 'Z'
}
