package repository

import (
	"context"

	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	"github.com/bloock/bloock-managed-api/internal/pkg"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockIntegrityRepository struct {
	client client.BloockClient
	logger zerolog.Logger
}

func NewBloockIntegrityRepository(ctx context.Context, logger zerolog.Logger) repository.IntegrityRepository {
	logger.With().Caller().Str("component", "integrity-repository").Logger()

	c := client.NewBloockClient(pkg.GetApiKeyFromContext(ctx), nil, pkg.GetEnvFromContext(ctx))

	return &BloockIntegrityRepository{
		client: c,
		logger: logger,
	}
}

func (b BloockIntegrityRepository) Certify(ctx context.Context, file []byte) (domain.Certification, error) {
	rec, err := client.NewRecordClient().FromFile(file).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("error certifying data")
		return domain.Certification{}, err
	}

	receipt, err := b.client.IntegrityClient.SendRecords([]record.Record{rec})
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return domain.Certification{}, err
	}
	dataHash, err := rec.GetHash()
	if err != nil {
		b.logger.Error().Err(err).Msg(err.Error())
		return domain.Certification{}, err
	}

	return domain.Certification{
		AnchorID: int(receipt[0].Anchor),
		Data:     file,
		Hash:     dataHash,
		Record:   &rec,
	}, nil
}
