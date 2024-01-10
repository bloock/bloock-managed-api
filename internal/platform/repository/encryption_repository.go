package repository

import (
	"context"

	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/pkg"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/encryption"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockEncryptionRepository struct {
	client client.BloockClient
	logger zerolog.Logger
}

func NewBloockEncryptionRepository(ctx context.Context, logger zerolog.Logger) repository.EncryptionRepository {
	logger.With().Caller().Str("component", "encryption-repository").Logger()

	c := client.NewBloockClient(pkg.GetApiKeyFromContext(ctx), nil, pkg.GetEnvFromContext(ctx))

	return &BloockEncryptionRepository{
		client: c,
		logger: logger,
	}
}

func (b BloockEncryptionRepository) EncryptWithLocalKey(ctx context.Context, data []byte, localKey key.LocalKey) (*record.Record, error) {
	encrypter := encryption.NewEncrypterWithLocalKey(localKey)
	rec, err := b.client.RecordClient.FromBytes(data).WithEncrypter(encrypter).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &rec, nil
}

func (b BloockEncryptionRepository) EncryptWithManagedKey(ctx context.Context, data []byte, managedKey key.ManagedKey) (*record.Record, error) {
	encrypter := encryption.NewEncrypterWithManagedKey(managedKey)
	rec, err := b.client.RecordClient.FromBytes(data).WithEncrypter(encrypter).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &rec, nil
}

func (b BloockEncryptionRepository) EncryptWithManagedCertificate(ctx context.Context, data []byte, managedCertificate key.ManagedCertificate) (*record.Record, error) {
	encrypter := encryption.NewEncrypterWithManagedCertificate(managedCertificate)
	rec, err := b.client.RecordClient.FromBytes(data).WithEncrypter(encrypter).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &rec, nil
}
