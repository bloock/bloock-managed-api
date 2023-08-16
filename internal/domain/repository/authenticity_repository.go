package repository

import (
	"bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type AuthenticityRepository interface {
	CreateLocalKey(keyType key.KeyType) (domain.LocalKey, error)
	LoadLocalKey(keyType key.KeyType, publicKey string, privateKey string) (key.LocalKey, error)
	Sign(localKey key.LocalKey, keyType key.KeyType, commonName *string, data []byte) (record.Record, error)
	CreateManagedKey(name string, keyType key.KeyType, expiration int, level key.KeyProtectionLevel) (key.ManagedKey, error)
	Verify(record record.Record) (bool, error)
}
