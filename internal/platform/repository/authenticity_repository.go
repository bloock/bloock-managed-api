package repository

import (
	"context"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/authenticity"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockAuthenticityRepository struct {
	keyClient          client.KeyClient
	authenticityClient client.AuthenticityClient
	recordClient       client.RecordClient
	logger             zerolog.Logger
}

func NewBloockAuthenticityRepository(logger zerolog.Logger) *BloockAuthenticityRepository {
	logger.With().Caller().Str("component", "authenticity-repository").Logger()

	return &BloockAuthenticityRepository{
		keyClient:          client.NewKeyClient(),
		authenticityClient: client.NewAuthenticityClient(),
		recordClient:       client.NewRecordClient(),
		logger:             logger,
	}
}

func (b BloockAuthenticityRepository) SignECWithLocalKey(ctx context.Context, data []byte, kty key.KeyType, publicKey string, privateKey *string) (string, record.Record, error) {
	signerArgs := authenticity.SignerArgs{}
	localKey, err := b.keyClient.LoadLocalKey(kty, publicKey, privateKey)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signerArgs.LocalKey = &localKey
	signer := authenticity.NewEcdsaSigner(signerArgs)
	rec, err := b.recordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signature, err := b.authenticityClient.Sign(rec, signer)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	return signature.Signature, rec, nil
}

func (b BloockAuthenticityRepository) SignECWithLocalKeyEns(ctx context.Context, data []byte, kty key.KeyType, publicKey string, privateKey *string) (string, record.Record, error) {
	signerArgs := authenticity.SignerArgs{}
	signer := authenticity.NewEnsSigner(signerArgs)

	localKey, err := b.keyClient.LoadLocalKey(kty, publicKey, privateKey)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signerArgs.LocalKey = &localKey
	rec, err := b.recordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signature, err := b.authenticityClient.Sign(rec, signer)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	return signature.Signature, rec, nil
}

func (b BloockAuthenticityRepository) SignECWithManagedKey(ctx context.Context, data []byte, kid string) (string, record.Record, error) {
	signerArgs := authenticity.SignerArgs{}

	managedKey, err := b.keyClient.LoadManagedKey(kid)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signerArgs.ManagedKey = &managedKey

	signer := authenticity.NewEcdsaSigner(signerArgs)
	rec, err := b.recordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signature, err := b.authenticityClient.Sign(rec, signer)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	return signature.Signature, rec, nil
}

func (b BloockAuthenticityRepository) SignECWithManagedKeyEns(ctx context.Context, data []byte, kid string) (string, record.Record, error) {
	signerArgs := authenticity.SignerArgs{}

	managedKey, err := b.keyClient.LoadManagedKey(kid)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signerArgs.ManagedKey = &managedKey

	signer := authenticity.NewEnsSigner(signerArgs)
	rec, err := b.recordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	signature, err := b.authenticityClient.Sign(rec, signer)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", record.Record{}, err
	}
	return signature.Signature, rec, nil
}
