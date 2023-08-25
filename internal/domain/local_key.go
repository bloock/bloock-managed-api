package domain

import (
	"errors"
	"fmt"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

type LocalKey struct {
	id       uuid.UUID
	localKey key.LocalKey
	keyType  key.KeyType
}

func NewLocalKey(localKey key.LocalKey, keyType key.KeyType, id uuid.UUID) *LocalKey {
	return &LocalKey{localKey: localKey, keyType: keyType, id: id}
}

func NewLocalKeyID(id uuid.UUID, localKey key.LocalKey, keyType key.KeyType) *LocalKey {
	return &LocalKey{localKey: localKey, keyType: keyType, id: id}
}

func (l LocalKey) Id() uuid.UUID {
	return l.id
}

func (l LocalKey) LocalKey() key.LocalKey {
	return l.localKey
}

func (l LocalKey) KeyType() key.KeyType {
	return l.keyType
}

func (l LocalKey) KeyTypeStr() string {
	switch l.keyType {
	case key.EcP256k:
		return "EcP256k"
	case key.Rsa2048:
		return "Rsa2048"
	case key.Rsa3072:
		return "Rsa3072"
	case key.Rsa4096:
		return "Rsa4096"
	default:
		return ""
	}
}

var ErrInvalidKeyType = func(kty string) error {
	return errors.New(fmt.Sprintf("invalid key type: %s", kty))
}
