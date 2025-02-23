package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"titler"
)

// CurrentPosition struct.
type CurrentPosition struct {
	X int
	Y int
}

// replaceInvalidChars replaces all characters that are not base60 or base10 with a comma.
func replaceInvalidChars(baseString string) string {
	re := regexp.MustCompile(`[^0-9A-Za-z]`)
	return re.ReplaceAllString(baseString, ",")
}

// readStack function. Reads the stack in the file.
func readStack(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	var stack [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = replaceInvalidChars(line)
		stack = append(stack, strings.Split(line, ","))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stack, nil
}

// startStackProcessing function. Starts the stack processing.
func startStackProcessing(stack [][]string, position, processing string, alternating bool) string {
	currentPosition := CurrentPosition{0, 0}
	direction := "right"
	isXPlus := true
	isYPlus := true

	// Getting the initial position.
	switch position {
	case "top_left":
		currentPosition.X = 0
		currentPosition.Y = 0

		if processing == "column" {
			direction = "down"
		} else {
			direction = "right"
		}

		isXPlus = true
		isYPlus = true
	case "top_right":
		currentPosition.X = len(stack[0]) - 1
		currentPosition.Y = 0

		if processing == "column" {
			direction = "down"
		} else {
			direction = "left"
		}

		isXPlus = false
		isYPlus = true
	case "bottom_left":
		currentPosition.X = 0
		currentPosition.Y = len(stack) - 1

		if processing == "column" {
			direction = "up"
		} else {
			direction = "right"
		}

		isXPlus = true
		isYPlus = false
	case "bottom_right":
		currentPosition.X = len(stack) - 1
		currentPosition.Y = len(stack[0]) - 1

		if processing == "column" {
			direction = "up"
		} else {
			direction = "left"
		}

		isXPlus = false
		isYPlus = false
	}

	// Processing the stack.
	return traverseStack(stack, currentPosition, processing, direction, alternating, isXPlus, isYPlus)
}

// traverseStack function. Traverses the stack.
func traverseStack(stack [][]string, currentPosition CurrentPosition, processing, direction string, alternating, isXPlus, isYPlus bool) string {
	var result []string
	rows := len(stack)
	cols := len(stack[0])

	for {
		// Append the current position value to the result
		result = append(result, stack[currentPosition.Y][currentPosition.X])

		// Determine the next position
		if processing == "row" {
			if direction == "right" {
				if currentPosition.X+1 < cols && isXPlus {
					currentPosition.X++
				} else if currentPosition.X-1 >= 0 && !isXPlus {
					currentPosition.X--
				} else {
					if alternating {
						direction = "left"
						isXPlus = !isXPlus
					} else {
						currentPosition.X = 0
					}
					if currentPosition.Y+1 < rows && isYPlus {
						currentPosition.Y++
					} else if currentPosition.Y-1 >= 0 && !isYPlus {
						currentPosition.Y--
					} else {
						break
					}
				}
			} else if direction == "left" {
				if currentPosition.X-1 >= 0 && !isXPlus {
					currentPosition.X--
				} else if currentPosition.X+1 < cols && isXPlus {
					currentPosition.X++
				} else {
					if alternating {
						direction = "right"
						isXPlus = !isXPlus
					} else {
						currentPosition.X = cols - 1
					}
					if currentPosition.Y+1 < rows && isYPlus {
						currentPosition.Y++
					} else if currentPosition.Y-1 >= 0 && !isYPlus {
						currentPosition.Y--
					} else {
						break
					}
				}
			}
		} else if processing == "column" {
			if direction == "down" {
				if currentPosition.Y+1 < rows && isYPlus {
					currentPosition.Y++
				} else if currentPosition.Y-1 >= 0 && !isYPlus {
					currentPosition.Y--
				} else {
					if alternating {
						direction = "up"
						isYPlus = !isYPlus
					} else {
						currentPosition.Y = 0
					}
					if currentPosition.X+1 < cols && isXPlus {
						currentPosition.X++
					} else if currentPosition.X-1 >= 0 && !isXPlus {
						currentPosition.X--
					} else {
						break
					}
				}
			} else if direction == "up" {
				if currentPosition.Y-1 >= 0 && !isYPlus {
					currentPosition.Y--
				} else if currentPosition.Y+1 < rows && isYPlus {
					currentPosition.Y++
				} else {
					if alternating {
						direction = "down"
						isYPlus = !isYPlus
					} else {
						currentPosition.Y = rows - 1
					}
					if currentPosition.X+1 < cols && isXPlus {
						currentPosition.X++
					} else if currentPosition.X-1 >= 0 && !isXPlus {
						currentPosition.X--
					} else {
						break
					}
				}
			}
		}
	}

	return strings.Join(result, ",")
}

// main function.
func main() {
	titler.PrintTitle("Column Processing Program")

	fileFlag := flag.String("file", "", "The file containing the stack")
	positionFlag := flag.String("position", "top_left", "The initial position in the stack (top_left, top_right, bottom_left, bottom_right)")
	processingFlag := flag.String("processing", "row", "The processing type (row, column)")
	alternatingFlag := flag.Bool("alternating", false, "Alternate between directions")

	flag.Parse()

	if *fileFlag == "" || *positionFlag == "" || *processingFlag == "" {
		flag.Usage()
		return
	}

	if *processingFlag != "row" && *processingFlag != "column" {
		fmt.Println("Invalid processing type. Valid types are: row, column")
		return
	}

	stack, err := readStack(*fileFlag)
	if err != nil {
		fmt.Println("Error reading stack:", err)
		return
	}

	result := startStackProcessing(stack, *positionFlag, *processingFlag, *alternatingFlag)
	fmt.Println(result)
}
