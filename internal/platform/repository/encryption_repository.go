package repository

import (
	"context"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/encryption"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockEncryptionRepository struct {
	keyClient        client.KeyClient
	encryptionClient client.EncryptionClient
	recordClient     client.RecordClient
	logger           zerolog.Logger
}

func NewBloockEncryptionRepository(logger zerolog.Logger) *BloockEncryptionRepository {
	logger.With().Caller().Str("component", "encryption-repository").Logger()

	return &BloockEncryptionRepository{
		keyClient:        client.NewKeyClient(),
		encryptionClient: client.NewEncryptionClient(),
		recordClient:     client.NewRecordClient(),
		logger:           logger,
	}
}

func (b BloockEncryptionRepository) EncryptRSAWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, publicKey string, privateKey *string) (*record.Record, error) {
	encrypterArgs := encryption.EncrypterArgs{}
	localKey, err := b.keyClient.LoadLocalKey(kty, publicKey, privateKey)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}
	encrypterArgs.LocalKey = &localKey
	encrypter := encryption.NewRsaEncrypter(encrypterArgs)
	rec, err := b.recordClient.FromBytes(data).WithEncrypter(encrypter).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &rec, nil
}

func (b BloockEncryptionRepository) EncryptRSAWithManagedKey(ctx context.Context, data []byte, kid string) (*record.Record, error) {
	encrypterArgs := encryption.EncrypterArgs{}
	managedKey, err := b.keyClient.LoadManagedKey(kid)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}
	encrypterArgs.ManagedKey = &managedKey

	encrypter := encryption.NewRsaEncrypter(encrypterArgs)
	rec, err := b.recordClient.FromBytes(data).WithEncrypter(encrypter).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &rec, nil
}

func (b BloockEncryptionRepository) EncryptAESWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, key string) (*record.Record, error) {
	encrypterArgs := encryption.EncrypterArgs{}
	localKey, err := b.keyClient.LoadLocalKey(kty, key, nil)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}
	encrypterArgs.LocalKey = &localKey
	encrypter := encryption.NewAesEncrypter(encrypterArgs)
	rec, err := b.recordClient.FromBytes(data).WithEncrypter(encrypter).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &rec, nil
}

func (b BloockEncryptionRepository) EncryptAESWithManagedKey(ctx context.Context, data []byte, kid string) (*record.Record, error) {
	encrypterArgs := encryption.EncrypterArgs{}
	managedKey, err := b.keyClient.LoadManagedKey(kid)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}
	encrypterArgs.ManagedKey = &managedKey

	encrypter := encryption.NewAesEncrypter(encrypterArgs)
	rec, err := b.recordClient.FromBytes(data).WithEncrypter(encrypter).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &rec, nil
}
