package utils

import "golang.org/x/crypto/sha3"

func MerkleHashFunc(data []byte) ([]byte, error) {
	sha256Func := sha3.NewLegacyKeccak256()
	sha256Func.Write(data)
	return sha256Func.Sum(nil), nil
}
