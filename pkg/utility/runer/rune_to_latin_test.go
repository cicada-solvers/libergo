package runer

import (
	"testing"
)

func TestTransposeRuneToLatin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ᚪ", "A"},
		{"ᚫ", "AE"},
		{"ᛠ", "EA"},
		{"ᛇ", "EO"},
		{"ᚩ", "O"},
		{"ᛟ", "OE"},
		{"ᛏ", "T"},
		{"ᚦ", "TH"},
		{"ᛁ", "I"},
		{"ᛡ", "IO"},
		{"ᛝ", "ING"},
		{"ᚾ", "N"},
		{"ᚻᛖᛚᛚᚩ", "HELLO"},
	}

	for _, test := range tests {
		result := TransposeRuneToLatin(test.input)
		if result != test.expected {
			t.Errorf("TransposeRuneToLatin(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
