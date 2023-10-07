package repository

import (
	"context"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type EncryptionRepository interface {
	EncryptRSAWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, publicKey string, privateKey *string) (*record.Record, error)
	EncryptRSAWithManagedKey(ctx context.Context, data []byte, kid string) (*record.Record, error)
	EncryptAESWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, key string) (*record.Record, error)
	EncryptAESWithManagedKey(ctx context.Context, data []byte, kid string) (*record.Record, error)
}
