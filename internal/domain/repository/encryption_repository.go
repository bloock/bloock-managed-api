package repository

import (
	"context"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type EncryptionRepository interface {
	EncryptWithLocalKey(ctx context.Context, data []byte, localKey key.LocalKey) (*record.Record, error)
	EncryptWithManagedKey(ctx context.Context, data []byte, managedKey key.ManagedKey, accessControl *key.AccessControl) (*record.Record, error)
	EncryptWithManagedCertificate(ctx context.Context, data []byte, managedCertificate key.ManagedCertificate, accessControl *key.AccessControl) (*record.Record, error)
}
