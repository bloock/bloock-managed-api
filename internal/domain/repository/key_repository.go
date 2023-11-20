package repository

import (
	"context"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
)

type KeyRepository interface {
	LoadLocalKey(ctx context.Context, kty key.KeyType, publicKey string, privateKey *string) (*key.LocalKey, error)
	LoadManagedKey(ctx context.Context, kid string) (*key.ManagedKey, error)
	LoadLocalCertificate(ctx context.Context, pkcs12 []byte, pkcs12Password string) (*key.LocalCertificate, error)
	LoadManagedCertificate(ctx context.Context, kid string) (*key.ManagedCertificate, error)
}
