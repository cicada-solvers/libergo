package conversion

import (
	"fmt"
	"strconv"
	"strings"
)

// ConvertToByteArray converts a string to a byte array
func ConvertToByteArray(part string) ([]byte, error) {
	subParts := strings.Split(part, ",")
	var array []byte
	for _, subPart := range subParts {
		val, err := strconv.Atoi(subPart)
		if err != nil {
			return nil, fmt.Errorf("error converting string to byte: %v", err)
		}
		array = append(array, byte(val))
	}
	return array, nil
}
