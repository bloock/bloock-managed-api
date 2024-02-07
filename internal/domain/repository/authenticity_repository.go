package repository

import (
	"context"

	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type AuthenticityRepository interface {
	SignWithLocalKey(ctx context.Context, data []byte, localKey key.LocalKey) (string, *record.Record, error)
	SignWithManagedKey(ctx context.Context, data []byte, managedKey key.ManagedKey, accessControl *key.AccessControl) (string, *record.Record, error)
	SignWithLocalCertificate(ctx context.Context, data []byte, localCertificate key.LocalCertificate) (string, *record.Record, error)
	SignWithManagedCertificate(ctx context.Context, data []byte, managedCertificate key.ManagedCertificate, accessControl *key.AccessControl) (string, *record.Record, error)
}
