package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	textFlag := flag.String("text", "", "Text to shift")
	shiftFlag := flag.Int("shift", 0, "Number of positions to shift to the left")
	textDirectionFlag := flag.String("direction", "left", "Direction to shift text (left or right)")
	flag.Parse()

	if *textFlag == "" || *shiftFlag == 0 {
		flag.Usage()
		return
	}

	letters := strings.Split(*textFlag, "")
	if *textDirectionFlag == "right" {
		result := shiftLettersRight(letters, *shiftFlag)
		fmt.Println(result)
		return
	} else {
		result := shiftLettersLeft(letters, *shiftFlag)
		fmt.Println(result)
		return
	}
}

func shiftLettersLeft(text []string, shift int) string {
	n := len(text)
	if n == 0 {
		return ""
	}
	// Normalize shift to [0, n)
	s := shift % n
	if s < 0 {
		s += n
	}
	rotated := append(append([]string{}, text[s:]...), text[:s]...)
	return strings.Join(rotated, "")
}

func shiftLettersRight(text []string, shift int) string {
	n := len(text)
	if n == 0 {
		return ""
	}
	// Normalize shift to [0, n)
	s := shift % n
	if s < 0 {
		s += n
	}
	if s == 0 {
		return strings.Join(text, "")
	}
	// Right rotation by s == left rotation by n - s
	left := n - s
	rotated := append(append([]string{}, text[left:]...), text[:left]...)
	return strings.Join(rotated, "")
}
