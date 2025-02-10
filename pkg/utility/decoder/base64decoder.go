package decoder

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// DecodeCommand represents the input for the Base64 decoding.
type DecodeCommand struct {
	Input    string
	Encoding string
}

// DecodeBase64String decodes the Base64 input based on the specified encoding.
func DecodeBase64String(cmd *DecodeCommand) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(cmd.Input)
	if err != nil {
		return "", err
	}

	var decoded string
	if strings.ToUpper(cmd.Encoding) == "HEX" {
		decoded = hex.EncodeToString(decodedBytes)
	} else {
		decoded = string(decodedBytes)
	}

	return decoded, nil
}
