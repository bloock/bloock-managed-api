package repository

import (
	"context"

	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/pkg"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/rs/zerolog"
)

type BloockKeyRepository struct {
	client client.BloockClient
	logger zerolog.Logger
}

func NewBloockKeyRepository(ctx context.Context, logger zerolog.Logger) repository.KeyRepository {
	logger.With().Caller().Str("component", "key-repository").Logger()

	c := client.NewBloockClient(pkg.GetApiKeyFromContext(ctx), "", pkg.GetEnvFromContext(ctx))

	return &BloockEncryptionRepository{
		client: c,
		logger: logger,
	}
}

func (b BloockEncryptionRepository) LoadLocalKey(ctx context.Context, kty key.KeyType, publicKey string, privateKey *string) (*key.LocalKey, error) {
	localKey, err := b.client.KeyClient.LoadLocalKey(kty, publicKey, privateKey)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &localKey, nil
}

func (b BloockEncryptionRepository) LoadManagedKey(ctx context.Context, kid string) (*key.ManagedKey, error) {
	managedKey, err := b.client.KeyClient.LoadManagedKey(kid)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &managedKey, nil
}

func (b BloockEncryptionRepository) LoadLocalCertificate(ctx context.Context, pkcs12 []byte, pkcs12Password string) (*key.LocalCertificate, error) {
	localCertificate, err := b.client.KeyClient.LoadLocalCertificate(pkcs12, pkcs12Password)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &localCertificate, nil
}

func (b BloockEncryptionRepository) LoadManagedCertificate(ctx context.Context, kid string) (*key.ManagedCertificate, error) {
	managedCertificate, err := b.client.KeyClient.LoadManagedCertificate(kid)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &managedCertificate, nil
}
