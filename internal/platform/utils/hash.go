package utils

import (
	"golang.org/x/crypto/sha3"
	"regexp"
)

const lengthSHA256 = "64"

func MerkleHashFunc(data []byte) ([]byte, error) {
	sha256Func := sha3.NewLegacyKeccak256()
	sha256Func.Write(data)
	return sha256Func.Sum(nil), nil
}

func IsSHA256(record string) bool {
	matched, err := regexp.MatchString("^[a-f0-9]{"+lengthSHA256+"}$", record)
	if err != nil {
		return false
	}
	return matched
}
