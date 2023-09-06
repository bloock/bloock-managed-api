package domain

import "github.com/bloock/bloock-sdk-go/v2/entity/key"

func ValidateKeyType(kty string) (key.KeyType, error) {
	switch kty {
	case "EcP256k":
		return key.EcP256k, nil
	case "Rsa2048":
		return key.Rsa2048, nil
	case "Rsa3072":
		return key.Rsa3072, nil
	case "Rsa4096":
		return key.Rsa4096, nil
	default:
		return -1, ErrInvalidKeyType(kty)
	}
}
