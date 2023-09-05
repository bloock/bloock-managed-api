package repository

import (
	"bloock-managed-api/internal/domain"
	"context"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
	"github.com/rs/zerolog"
)

type BloockIntegrityRepository struct {
	integrityClient client.IntegrityClient
	logger          zerolog.Logger
}

func NewBloockIntegrityRepository(logger zerolog.Logger) *BloockIntegrityRepository {
	logger.With().Caller().Str("component", "integrity-repository").Logger()
	return &BloockIntegrityRepository{
		integrityClient: client.NewIntegrityClient(),
		logger: logger,
	}
}

func (b BloockIntegrityRepository) Certify(ctx context.Context, file []byte) (domain.Certification, error) {
	rec, err := client.NewRecordClient().FromFile(file).Build()
	if err != nil {
		b.logger.Error().Err(err).Msg("error certifying data")
		return domain.Certification{}, err
	}

	receipt, err := b.integrityClient.SendRecords([]record.Record{rec})
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
	}, nil
}
