package repository

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type AuthenticityRepository interface {
	SignECWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, publicKey string, privateKey *string) (string, record.Record, error)
	SignECWithLocalKeyEns(ctx context.Context, data []byte, kty key.KeyType, publicKey string, privateKey *string) (string, record.Record, error)
	SignECWithManagedKey(ctx context.Context, data []byte, kid string) (string, record.Record, error)
	SignECWithManagedKeyEns(ctx context.Context, data []byte, kid string) (string, record.Record, error)
}
