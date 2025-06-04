package cipher

import "strings"

// DecodeRailFence decodes a rail fence cipher with the given number of rails
func DecodeRailFence(ciphertext []string, rails int) string {
	if rails <= 1 || len(ciphertext) == 0 {
		return strings.Join(ciphertext, "")
	}

	// Create a 2D matrix of runes to represent the rails
	matrix := make([][]string, rails)
	for i := range matrix {
		matrix[i] = make([]string, len(ciphertext))
		for j := range matrix[i] {
			matrix[i][j] = "0" // Fill with placeholder
		}
	}

	// Mark the positions where characters will be with a placeholder
	row, direction := 0, -1
	for i := 0; i < len(ciphertext); i++ {
		matrix[row][i] = "*" // Mark position

		// Change direction at the boundaries
		if row == 0 || row == rails-1 {
			direction = -direction
		}
		row += direction
	}

	// Fill the matrix with the ciphertext
	idx := 0
	for r := 0; r < rails; r++ {
		for c := 0; c < len(ciphertext); c++ {
			if matrix[r][c] == "*" && idx < len(ciphertext) {
				matrix[r][c] = ciphertext[idx]
				idx++
			}
		}
	}

	// Read off the plaintext in zigzag order
	result := make([]string, len(ciphertext))
	idx = 0
	row, direction = 0, -1
	for i := 0; i < len(ciphertext); i++ {
		result[i] = matrix[row][i]

		// Change direction at the boundaries
		if row == 0 || row == rails-1 {
			direction = -direction
		}
		row += direction
	}

	return strings.Join(result, "")
}
