package runer

import (
	"testing"
)

func TestPrepLatinToRune(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"quick", "CWICC"},
		{"zebra", "SEBRA"},
		{"king", "CING"},
		{"queen", "CWEEN"},
		{"victory", "UICTORY"},
		{"io", "IO"},
		{"ia", "IO"},
		{"hello", "HELLO"},
	}

	for _, test := range tests {
		result := PrepLatinToRune(test.input)
		if result != test.expected {
			t.Errorf("PrepLatinToRune(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
