package runer

import (
	"testing"
)

func TestTransposeLatinToRune(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A", "ᚪ"},
		{"AE", "ᚫ"},
		{"EA", "ᛠ"},
		{"EO", "ᛇ"},
		{"O", "ᚩ"},
		{"OE", "ᛟ"},
		{"T", "ᛏ"},
		{"TH", "ᚦ"},
		{"I", "ᛁ"},
		{"IO", "ᛡ"},
		{"ING", "ᛝ"},
		{"IA", "ᛡ"},
		{"N", "ᚾ"},
		{"HELLO", "ᚻᛖᛚᛚᚩ"},
	}

	for _, test := range tests {
		result := TransposeLatinToRune(test.input)
		if result != test.expected {
			t.Errorf("TransposeLatinToRune(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
