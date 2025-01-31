package main

// Program represents the program
type Program struct {
	tasks chan []byte
}

// NewProgram creates a new Program
func NewProgram() *Program {
	return &Program{
		tasks: make(chan []byte, 10000),
	}
}

// generateAllByteArrays generates all possible byte arrays of a given length
func (p *Program) generateAllByteArrays(maxArrayLength int, startArray, stopArray []byte) {
	currentArray := make([]byte, len(startArray))
	copy(currentArray, startArray)
	p.generateByteArrays(maxArrayLength, 1, currentArray, stopArray)
	close(p.tasks)
}

// generateByteArrays generates all possible byte arrays of a given length
func (p *Program) generateByteArrays(maxArrayLength, currentArrayLevel int, passedArray, stopArray []byte) bool {
	startForValue := int(passedArray[currentArrayLevel-1])

	if currentArrayLevel == maxArrayLength {
		currentArray := make([]byte, maxArrayLength)
		if passedArray != nil {
			copy(currentArray, passedArray)
		}

		for i := startForValue; i < 256; i++ {
			currentArray[currentArrayLevel-1] = byte(i)
			p.tasks <- append([]byte(nil), currentArray...)
			if compareArrays(currentArray, stopArray) == 0 {
				return false
			}
		}
	} else {
		currentArray := make([]byte, maxArrayLength)
		if passedArray != nil {
			copy(currentArray, passedArray)
		}
		for i := startForValue; i < 256; i++ {
			currentArray[currentArrayLevel-1] = byte(i)
			shouldContinue := p.generateByteArrays(maxArrayLength, currentArrayLevel+1, currentArray, stopArray)
			if !shouldContinue {
				return false
			}
			currentArray[currentArrayLevel] = 0
		}
	}

	return true
}

// compareArrays compares two byte arrays and returns -1 if a < b, 0 if a == b, and 1 if a > b
func compareArrays(a, b []byte) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	if len(a) < len(b) {
		return -1
	} else if len(a) > len(b) {
		return 1
	}
	return 0
}
