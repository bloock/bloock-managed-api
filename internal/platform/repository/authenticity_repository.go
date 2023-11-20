package repository

import (
	"context"

	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/pkg"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/authenticity"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockAuthenticityRepository struct {
	client client.BloockClient
	logger zerolog.Logger
}

func NewBloockAuthenticityRepository(ctx context.Context, logger zerolog.Logger) repository.AuthenticityRepository {
	logger.With().Caller().Str("component", "authenticity-repository").Logger()

	c := client.NewBloockClient(pkg.GetApiKeyFromContext(ctx), "", pkg.GetEnvFromContext(ctx))

	return &BloockAuthenticityRepository{
		client: c,
		logger: logger,
	}
}

func (b BloockAuthenticityRepository) SignWithLocalKey(ctx context.Context, data []byte, localKey *key.LocalKey) (string, *record.Record, error) {
	signerArgs := authenticity.SignerArgs{
		LocalKey: localKey,
	}

	signer := authenticity.NewSigner(signerArgs)
	rec, err := b.client.RecordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}
	signatures, err := b.client.AuthenticityClient.GetSignatures(rec)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}

	return signatures[0].Signature, &rec, nil
}

func (b BloockAuthenticityRepository) SignWithManagedKey(ctx context.Context, data []byte, managedKey *key.ManagedKey) (string, *record.Record, error) {
	signerArgs := authenticity.SignerArgs{
		ManagedKey: managedKey,
	}

	signer := authenticity.NewSigner(signerArgs)
	rec, err := b.client.RecordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}
	signatures, err := b.client.AuthenticityClient.GetSignatures(rec)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}

	return signatures[0].Signature, &rec, nil
}

func (b BloockAuthenticityRepository) SignWithLocalCertificate(ctx context.Context, data []byte, localCertificate *key.LocalCertificate) (string, *record.Record, error) {
	signerArgs := authenticity.SignerArgs{
		LocalCertificate: localCertificate,
	}

	signer := authenticity.NewSigner(signerArgs)
	rec, err := b.client.RecordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}
	signatures, err := b.client.AuthenticityClient.GetSignatures(rec)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}

	return signatures[0].Signature, &rec, nil
}

func (b BloockAuthenticityRepository) SignWithManagedCertificate(ctx context.Context, data []byte, managedCertificate *key.ManagedCertificate) (string, *record.Record, error) {
	signerArgs := authenticity.SignerArgs{
		ManagedCertificate: managedCertificate,
	}

	signer := authenticity.NewSigner(signerArgs)
	rec, err := b.client.RecordClient.FromBytes(data).WithSigner(signer).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}
	signatures, err := b.client.AuthenticityClient.GetSignatures(rec)
	if err != nil {
		b.logger.Error().Err(err).Msg("")
		return "", nil, err
	}

	return signatures[0].Signature, &rec, nil
}
