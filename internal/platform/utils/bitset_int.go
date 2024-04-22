package utils

import (
	"encoding/hex"
	"fmt"
)

type BitsetInt []int

func BitsetIntFromString(input string) (BitsetInt, error) {
	var groupSize = 4
	var groups []string
	for i := 0; i < len(input); i += groupSize {
		end := i + groupSize
		if end > len(input) {
			end = len(input)
		}
		groups = append(groups, input[i:end])
	}

	var output BitsetInt
	for _, group := range groups {
		decoded, err := hex.DecodeString(group)
		if err != nil {
			return nil, err
		}
		var intValue int
		for _, b := range decoded {
			intValue = intValue<<8 | int(b)
		}
		output = append(output, intValue)
	}

	return output, nil
}

func (b BitsetInt) ToString() (string, error) {
	var depthU16 []uint16
	for _, p := range b {
		if p > 65535 {
			// Handle the error case where the uint32 value is too large for uint16
			return "", fmt.Errorf("error: value too large for uint16: %d", p)
		}
		depthU16 = append(depthU16, uint16(p))
	}

	// Convert uint16 to []byte and concatenate
	var depthU8 []byte
	for _, x := range depthU16 {
		depthU8 = append(depthU8, Uint16ToBytes(x)...)
	}

	return hex.EncodeToString(depthU8), nil
}
